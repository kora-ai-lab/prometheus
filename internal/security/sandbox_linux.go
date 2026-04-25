package security

import (
	"context"
	"syscall"
	"os/exec"
	"strconv"
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
	if canUseNamespaces() {
		return newNamespaceSandbox(cfg)
	}
	return &WorkdirSandbox{cfg: cfg}
}

func canUseNamespaces() bool {
	cmd := exec.Command("unshare", "--user", "--map-root-user", "echo", "ok")
	_ = cmd.Run()
	return cmd.ProcessState != nil && cmd.ProcessState.Success()
}

type NamespaceSandbox struct {
	cfg SandboxConfig
}

func newNamespaceSandbox(cfg SandboxConfig) *NamespaceSandbox {
	return &NamespaceSandbox{cfg: cfg}
}

func (s *NamespaceSandbox) Level() SandboxLevel { return SandboxNamespace }

func (s *NamespaceSandbox) Run(ctx context.Context, cmd string) (*SandboxResult, error) {
	timeout := s.cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	workDir := s.cfg.WorkDir
	if workDir == "" {
		workDir = "/tmp"
	}
	maxRAM := s.cfg.MaxRAMMb
	if maxRAM == 0 {
		maxRAM = 512
	}
	maxCPUSec := s.cfg.MaxCPUSec
	if maxCPUSec == 0 {
		maxCPUSec = 30
	}
	shellCmd := exec.CommandContext(ctx, "unshare",
		"--user",
		"--map-root-user",
		"--mount",
		"--pid",
		"--fork",
		"--mount-proc",
		"--UTS",
		"--IPC",
		"--cgroup",
		"--net",
		"bash", "-c",
		"ulimit -v "+strconv.Itoa(maxRAM*1024)+
		" && ulimit -t "+strconv.Itoa(maxCPUSec)+
		" && cd "+workDir+
		" && PATH=/usr/local/bin:/usr/bin:/bin && "+cmd,
	)
	shellCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER,
	}
	out, err := shellCmd.CombinedOutput()
	result := &SandboxResult{Stdout: string(out), ExitCode: 0}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}
	return result, nil
}

func (s *NamespaceSandbox) Cleanup() error {
	return nil
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
	result := &SandboxResult{Stdout: string(out), ExitCode: 0}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}
	return result, nil
}

func (s *WorkdirSandbox) Cleanup() error {
	return nil
}