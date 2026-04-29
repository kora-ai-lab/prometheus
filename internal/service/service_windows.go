//go:build windows

package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
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
	"github.com/kora-ai-lab/prometheus/internal/prompt"
	"github.com/kora-ai-lab/prometheus/internal/security"
	"github.com/kora-ai-lab/prometheus/internal/storage"
	"github.com/kora-ai-lab/prometheus/internal/task"
	"github.com/kora-ai-lab/prometheus/internal/ui"
	"github.com/kora-ai-lab/prometheus/internal/vault"
	"github.com/kora-ai-lab/prometheus/internal/vision"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type handler struct {
	stopCh chan struct{}
}

func (h *handler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-h.stopCh
		cancel()
	}()

	go func() {
		if err := runService(ctx); err != nil {
			cancel()
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
Loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				close(h.stopCh)
				break Loop
			case svc.Interrogate:
				changes <- c.CurrentStatus
			}
		}
	}
	changes <- svc.Status{State: svc.Stopped}
	return false, 0
}

func runService(ctx context.Context) error {
	home, err := config.EnsureHome()
	if err != nil {
		return fmt.Errorf("initializing prometheus home: %w", err)
	}

	cfg, err := config.Load(home)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger, err := logging.New(home)
	if err != nil {
		return fmt.Errorf("creating logger: %w", err)
	}

	store, err := storage.Open(home)
	if err != nil {
		return fmt.Errorf("opening storage: %w", err)
	}

	execer := executor.NewRateLimitedExecutor(cfg.Security.RateLimitExecPerSec)
	env := discovery.Scan(ctx, execer)

	serverPath, err := embedded.ExtractServer()
	if err != nil && !errors.Is(err, embedded.ErrPlaceholderArtifact) && !errors.Is(err, embedded.ErrPlatformNotSupported) {
		return fmt.Errorf("extracting embedded llama-server: %w", err)
	}

	provider, err := llm.AutoDetect(&cfg.LLM, serverPath)
	if err != nil {
		return fmt.Errorf("selecting llm provider: %w", err)
	}

	visionProvider := vision.VisionProvider(&vision.NoOpVisionProvider{})
	if provider.HasVision() {
		visionProvider = vision.NewLocalProvider()
	}

	capEngine := capabilities.NewEngine(execer, env, logger)
	promptBuilder := prompt.NewBuilder(provider.ModelInfo(), env, nil)
	browserMgr := browser.NewManager(capEngine, visionProvider)
	securityInterceptor := security.New(cfg.Security)
	_ = vault.New(filepath.Join(home, "vault.enc"))

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
	if err := writeToken(server.AuthToken()); err != nil {
		logger.Error("failed to write token", err)
	}

	fmt.Printf("Prometheus service running on http://%s:%d\n", cfg.UI.WebHost, cfg.UI.WebPort)
	fmt.Printf("API Token: %s\n", server.AuthToken())

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	select {
	case <-ctx.Done():
		server.Shutdown(context.Background())
		provider.Close()
		browserMgr.Close()
		store.Close()
		logger.Close()
		return nil
	case err := <-errCh:
		provider.Close()
		browserMgr.Close()
		store.Close()
		logger.Close()
		return err
	}
}

func Install() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting to service control manager: %w", err)
	}
	defer m.Disconnect()

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	s, err := m.CreateService(ServiceName, exePath, mgr.Config{
		DisplayName: DisplayName,
		Description: Description,
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return fmt.Errorf("creating service: %w", err)
	}
	defer s.Close()

	return nil
}

func Uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting to service control manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(ServiceName)
	if err != nil {
		return fmt.Errorf("opening service: %w", err)
	}
	defer s.Close()

	status, err := s.Query()
	if err == nil && status.State == svc.Running {
		_, err = s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("stopping service: %w", err)
		}
		for i := 0; i < 30; i++ {
			status, err = s.Query()
			if err != nil || status.State == svc.Stopped {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	err = s.Delete()
	if err != nil {
		return fmt.Errorf("deleting service: %w", err)
	}

	return nil
}

func Status() (string, error) {
	m, err := mgr.Connect()
	if err != nil {
		return "", fmt.Errorf("connecting to service control manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(ServiceName)
	if err != nil {
		return "not installed", nil
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return "", fmt.Errorf("querying service status: %w", err)
	}

	switch status.State {
	case svc.Running:
		return "running", nil
	case svc.Stopped:
		return "stopped", nil
	default:
		return "unknown", nil
	}
}

func Start() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting to service control manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(ServiceName)
	if err != nil {
		return fmt.Errorf("opening service: %w", err)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		return fmt.Errorf("starting service: %w", err)
	}

	return nil
}

func Stop() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connecting to service control manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(ServiceName)
	if err != nil {
		return fmt.Errorf("opening service: %w", err)
	}
	defer s.Close()

	_, err = s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("stopping service: %w", err)
	}

	return nil
}

func Run() error {
	inService, err := svc.IsWindowsService()
	if err != nil {
		return fmt.Errorf("checking if running as windows service: %w", err)
	}

	if inService {
		h := &handler{stopCh: make(chan struct{})}
		return svc.Run(ServiceName, h)
	}

	return runService(context.Background())
}

func writeToken(token string) error {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}
	dir := filepath.Join(localAppData, "Prometheus")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "token.txt"), []byte(token), 0600)
}

func getLocalAddr(host string, port int) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}
