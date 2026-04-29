package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/api"
	"github.com/kora-ai-lab/prometheus/internal/browser"
	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/config"
	"github.com/kora-ai-lab/prometheus/internal/discovery"
	"github.com/kora-ai-lab/prometheus/internal/executor"
	"github.com/kora-ai-lab/prometheus/internal/llm"
	"github.com/kora-ai-lab/prometheus/internal/llm/embedded"
	"github.com/kora-ai-lab/prometheus/internal/logging"
	"github.com/kora-ai-lab/prometheus/internal/metrics"
	"github.com/kora-ai-lab/prometheus/internal/prompt"
	"github.com/kora-ai-lab/prometheus/internal/security"
	"github.com/kora-ai-lab/prometheus/internal/storage"
	"github.com/kora-ai-lab/prometheus/internal/task"
	"github.com/kora-ai-lab/prometheus/internal/ui"
	"github.com/kora-ai-lab/prometheus/internal/update"
	"github.com/kora-ai-lab/prometheus/internal/vault"
	"github.com/kora-ai-lab/prometheus/internal/vision"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	home, err := config.EnsureHome()
	exitOnError(err, "initializing prometheus home")

	cfg, err := config.Load(home)
	exitOnError(err, "loading config")

	logger, err := logging.New(home)
	exitOnError(err, "creating logger")
	defer logger.Close()

	store, err := storage.Open(home)
	exitOnError(err, "opening storage")
	defer store.Close()

	execer := executor.NewRateLimitedExecutor(cfg.Security.RateLimitExecPerSec)
	env := discovery.Scan(ctx, execer)

	if handled := handleCLI(home, cfg, env); handled {
		return
	}

	serverPath, err := embedded.ExtractServer()
	if err != nil && !errors.Is(err, embedded.ErrPlaceholderArtifact) && !errors.Is(err, embedded.ErrPlatformNotSupported) {
		exitOnError(err, "extracting embedded llama-server")
	}

	if cfg.LLM.ModelPath == "" {
		logger.Info("model path not set, using auto-detection or Ollama fallback")
	}

	provider, err := llm.AutoDetect(&cfg.LLM, serverPath)
	exitOnError(err, "selecting llm provider")
	defer provider.Close()

	visionProvider := vision.VisionProvider(&vision.NoOpVisionProvider{})
	if provider.HasVision() {
		visionProvider = vision.NewLocalProvider()
	}

	capEngine := capabilities.NewEngine(execer, env, logger)
	promptBuilder := prompt.NewBuilder(provider.ModelInfo(), env, nil)
	browserMgr := browser.NewManager(capEngine, visionProvider)
	defer browserMgr.Close()
	securityInterceptor := security.New(cfg.Security)
	_ = vault.New(filepath.Join(home, "vault.enc"))

	if len(os.Args) > 1 && os.Args[1] == "--web" {
		taskMgr := api.NewTaskManager(func() *task.TaskDeps {
			return &task.TaskDeps{
				Provider:      provider,
				Executor:      execer,
				Vision:        visionProvider,
				PromptBuilder: promptBuilder,
				CapEngine:     capEngine,
				Security:      securityInterceptor,
				Logger:        logger,
				TaskStore:     store,
			}
		})
		server := ui.NewWebServer(cfg.UI.WebHost, cfg.UI.WebPort, taskMgr, cfg)
		fmt.Printf("Prometheus web UI listening on http://%s:%d\n", cfg.UI.WebHost, cfg.UI.WebPort)
		fmt.Printf("API Token: %s\n", server.AuthToken())
		exitOnError(server.Start(), "starting web ui")
		return
	}

	goal := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if goal == "" {
		if !ui.IsInteractive() || isCI() {
			showBanner()
			fmt.Println("Usage:")
			fmt.Println("  prometheus <goal>           Run with a goal")
			fmt.Println("  prometheus --web            Start web UI at http://localhost:8080")
			fmt.Println("  prometheus --help           Show help")
			fmt.Println("")
			fmt.Println("Examples:")
			fmt.Println("  prometheus 'list files in current folder'")
			fmt.Println("  prometheus 'search for info about AI'")
			fmt.Println("  prometheus --web")
			return
		}
		showBanner()
		fmt.Print("> ")
		line, _ := ui.ReadLine("")
		goal = line
		if goal == "" {
			return
		}
	}
	if goal == "" {
		fmt.Println("No goal provided. Usage: prometheus <goal> or prometheus --web")
		return
	}

	t := task.New(goal)
	exitOnError(store.Save(t), "saving task")

	deps := &task.TaskDeps{
		Provider:      provider,
		Executor:      execer,
		Vision:        visionProvider,
		Browser:       browserMgr,
		PromptBuilder: promptBuilder,
		CapEngine:     capEngine,
		Security:      securityInterceptor,
		Logger:        logger,
		TaskStore:     store,
	}

	for t.Status == task.StatusRunning || t.Status == task.StatusBlocked {
		if t.Status == task.StatusBlocked {
			answer, err := ui.ReadLine("\n⊙ PROMETHEUS A BESOIN D'UNE INFO:\n" + t.BlockedReason + "\n> ")
			exitOnError(err, "reading blocked response")
			t.Resume(answer)
		}
		exitOnError(t.Run(ctx, deps), "running task")
	}

	switch t.Status {
	case task.StatusDone:
		fmt.Println("\n✓ Terminé")
	case task.StatusFailed:
		fmt.Println("\n✗ Échec — voir les logs pour plus de détails")
		os.Exit(1)
	}
}

func exitOnError(err error, step string) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "prometheus: %s: %v\n", step, err)
	os.Exit(1)
}

func isCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

func showBanner() {
	fmt.Println(`
╔═══════════════════════════════════════╗
║     Prometheus v1.0.2                 ║
║     AI-first agent runtime            ║
╚═══════════════════════════════════════╝
`)
}

func handleCLI(home string, cfg *config.Config, env *discovery.EnvironmentProfile) bool {
	if len(os.Args) < 2 {
		return false
	}

	switch os.Args[1] {
	case "setup":
		exitOnError(llm.FirstRunSetup(home, env, os.Stdout), "running setup")
		return true
	case "metrics":
		data, _ := json.MarshalIndent(metrics.New().Snapshot(), "", "  ")
		fmt.Println(string(data))
		fmt.Println("Note: runtime metrics persistence is scaffolded but not yet connected.")
		return true
	case "logs":
		date := time.Now().Format("2006-01-02")
		if len(os.Args) > 2 {
			date = os.Args[2]
		}
		path := filepath.Join(home, "logs", date+".jsonl")
		raw, err := os.ReadFile(path)
		exitOnError(err, "reading logs")
		fmt.Print(string(raw))
		return true
	case "vault":
		if len(os.Args) > 2 && os.Args[2] == "list" {
			v := vault.New(filepath.Join(home, "vault.enc"))
			keys, err := v.List()
			exitOnError(err, "listing vault keys")
			for _, key := range keys {
				fmt.Println(key)
			}
			return true
		}
		fmt.Println("Usage: prometheus vault list")
		return true
	case "selftest":
		fmt.Printf("selftest ok\nprovider=%s\npackage_manager=%s\n", cfg.LLM.Provider, env.PackageManager)
		fmt.Println("Note: end-to-end selftest scenarios are scaffold placeholders at this stage.")
		return true
	case "update":
		hasUpdate, latestVersion, err := update.CheckForUpdate("anomalyco", "prometheus")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
			return true
		}

		if !hasUpdate {
			fmt.Println("You are running the latest version.")
			return true
		}

		fmt.Printf("Version %s is available. Upgrade? (y/n) ", latestVersion)

		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" {
			fmt.Println("Update cancelled.")
			return true
		}

		fmt.Println("Update download and apply is scaffolded.")
		fmt.Println("(Full implementation would download, verify checksum, and apply the update)")
		return true
	}
	return false
}
