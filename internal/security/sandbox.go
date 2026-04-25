package security

import (
	"context"
	"os/exec"
	"runtime"
	"time"
)

type SandboxLevel int

const (
	SandboxNone     SandboxLevel = 0
	SandboxWorkdir   SandboxLevel = 1
	SandboxNamespace SandboxLevel = 2
)

type SandboxConfig struct {
	MaxRAMMb      int
	MaxCPUSec     int
	MaxFileSizeMb int
	WorkDir       string
	Timeout       time.Duration
}

type SandboxResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	TimedOut bool
}

type Sandbox interface {
	Level() SandboxLevel
	Run(ctx context.Context, cmd string) (*SandboxResult, error)
	Cleanup() error
}

func NewSandbox(cfg SandboxConfig) Sandbox {
	if runtime.GOOS == "linux" && canUseNamespaces() {
		return newNamespaceSandbox(cfg)
	}
return &WorkdirSandbox{cfg: cfg}
}

type WorkdirSandbox struct {
	cfg SandboxConfig
}

func (s *WorkdirSandbox) Level() SandboxLevel { return SandboxWorkdir }

func (s *WorkdirSandbox) Run(ctx context.Context, cmd string) (*SandboxResult, error) {
	workDir := s.cfg.WorkDir
	
	timeout := 5 * time.Minute
	if s.cfg.Timeout > 0 {
		timeout = s.cfg.Timeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	execCmd := exec.CommandContext(ctx, "sh", "-c", cmd)
	if workDir != "" {
		execCmd.Dir = workDir
	}

	out, err := execCmd.CombinedOutput()
	result := &SandboxResult{
		Stdout:   string(out),
		ExitCode: 0,
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result, nil
}

func (s *WorkdirSandbox) Cleanup() error { return nil }
