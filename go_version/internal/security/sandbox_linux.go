//go:build linux
// +build linux

package security

import (
	"context"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func canUseNamespaces() bool {
	cmd := exec.Command("unshare", "--user", "--map-root-user", "echo", "ok")
	_ = cmd.Run()
	return cmd.ProcessState != nil && cmd.ProcessState.Success()
}

func newNamespaceSandbox(cfg SandboxConfig) *NamespaceSandbox {
	return &NamespaceSandbox{cfg: cfg}
}

type NamespaceSandbox struct {
	cfg SandboxConfig
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