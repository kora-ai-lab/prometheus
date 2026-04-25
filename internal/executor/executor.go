package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus-dev/prometheus/internal/config"
	"github.com/prometheus-dev/prometheus/internal/security"
)

type ExecOptions struct {
	Timeout time.Duration
	WorkDir string
	Env     []string
}

type ExecResult struct {
	Command  string
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	TimedOut bool
}

type Executor interface {
	Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult
}

type ShellExecutor struct {
	securityInterceptor *security.Interceptor
}

func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{
		securityInterceptor: security.New(config.SecurityConfig{
			DangerousOpsConfirm: true,
		}),
	}
}

func (s *ShellExecutor) Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult {
	if s.securityInterceptor != nil {
		if allowed, err := s.securityInterceptor.Allow(command); err != nil || !allowed {
			stderr := err.Error()
			if strings.Contains(stderr, "requires explicit confirmation") {
				stderr = "SECURITY_CONFIRMATION_REQUIRED: " + stderr
			} else {
				stderr = "SECURITY_BLOCKED: " + stderr
			}
			return &ExecResult{
				Command:  command,
				ExitCode: -1,
				Stderr:   stderr,
			}
		}
	}

	return executeImpl(ctx, command, opts)
}

func executeImpl(ctx context.Context, command string, opts ExecOptions) *ExecResult {
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	default:
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	if opts.WorkDir != "" {
		cmd.Dir = opts.WorkDir
	}
	cmd.Env = append(os.Environ(), opts.Env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := &ExecResult{
		Command:  command,
		Stdout:   TruncateMid(stdout.String(), 50_000),
		Stderr:   TruncateMid(stderr.String(), 20_000),
		Duration: duration,
		TimedOut: ctx.Err() == context.DeadlineExceeded,
	}

	if err == nil {
		return result
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	} else {
		result.ExitCode = -1
		if result.Stderr == "" {
			result.Stderr = err.Error()
		}
	}

	return result
}

func Execute(ctx context.Context, command string, opts ExecOptions) *ExecResult {
	return executeImpl(ctx, command, opts)
}

func TruncateMid(s string, maxChars int) string {
	if len(s) <= maxChars {
		return s
	}

	half := maxChars * 2 / 5
	skipped := s[half : len(s)-half]
	return s[:half] + fmt.Sprintf("\n...[truncated %d lines]...\n", strings.Count(skipped, "\n")) + s[len(s)-half:]
}