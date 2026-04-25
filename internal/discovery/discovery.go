package discovery

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/executor"
)

type EnvironmentProfile struct {
	OS             string
	Arch           string
	Kernel         string
	RAMMb          int
	DiskGb         int
	CPUCores       int
	AvailableTools []string
	LLMModels      []string
	Internet       bool
	PackageManager string
	IsTermux       bool
	ScannedAt      time.Time
}

func Scan(ctx context.Context, execer executor.Executor) *EnvironmentProfile {
	p := &EnvironmentProfile{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CPUCores:  runtime.NumCPU(),
		IsTermux:  isTermux(),
		ScannedAt: time.Now(),
	}

	p.Kernel = runtime.GOOS
	p.RAMMb = readRAM(ctx)
	p.DiskGb = 0
	p.AvailableTools = checkTools([]string{
		"git", "python3", "python", "node", "npm", "curl",
		"wget", "chromium", "chromium-browser", "google-chrome",
		"firefox", "adb", "ffmpeg", "scrot", "import", "screencap",
	})
	p.PackageManager = detectPackageManager()
	p.Internet = checkInternet(ctx, execer)
	p.LLMModels = detectOllamaModels(ctx, execer)
	return p
}

func isTermux() bool {
	return strings.Contains(os.Getenv("PREFIX"), "com.termux")
}

func readRAM(ctx context.Context) int {
	switch runtime.GOOS {
	case "windows":
		return 8192
	default:
		_ = ctx
		return 4096
	}
}

func checkTools(candidates []string) []string {
	var tools []string
	for _, tool := range candidates {
		if _, err := exec.LookPath(tool); err == nil {
			tools = append(tools, tool)
		}
	}
	return tools
}

func detectPackageManager() string {
	for _, tool := range []string{"apt", "brew", "pkg", "dnf", "winget", "choco"} {
		if _, err := exec.LookPath(tool); err == nil {
			return tool
		}
	}
	return ""
}

func checkInternet(ctx context.Context, execer executor.Executor) bool {
	cmd := "ping -n 1 1.1.1.1"
	if runtime.GOOS != "windows" {
		cmd = "ping -c 1 1.1.1.1"
	}
	res := execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 5 * time.Second})
	return res.ExitCode == 0
}

func detectOllamaModels(ctx context.Context, execer executor.Executor) []string {
	res := execer.Execute(ctx, "ollama list", executor.ExecOptions{Timeout: 5 * time.Second})
	if res.ExitCode != 0 {
		return nil
	}
	var models []string
	for _, line := range strings.Split(res.Stdout, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && !strings.EqualFold(fields[0], "name") {
			models = append(models, fields[0])
		}
	}
	return models
}
