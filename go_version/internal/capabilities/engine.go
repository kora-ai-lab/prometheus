package capabilities

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/prometheus-dev/prometheus/internal/discovery"
	"github.com/prometheus-dev/prometheus/internal/executor"
	"github.com/prometheus-dev/prometheus/internal/logging"
)

type Engine struct {
	execer    executor.Executor
	env       *discovery.EnvironmentProfile
	logger    *logging.Logger
	installed map[string]bool
	mu        sync.Mutex
}

func NewEngine(execer executor.Executor, env *discovery.EnvironmentProfile, logger *logging.Logger) *Engine {
	return &Engine{
		execer:    execer,
		env:       env,
		logger:    logger,
		installed: map[string]bool{},
	}
}

func (e *Engine) Ensure(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("empty capability name")
	}
	if _, err := exec.LookPath(name); err == nil {
		return nil
	}

	e.mu.Lock()
	if e.installed[name] {
		e.mu.Unlock()
		return nil
	}
	e.mu.Unlock()

	capability, ok := Registry[name]
	if !ok {
		return fmt.Errorf("capability %q not found in registry", name)
	}
	cmd := capability.InstallCmds[e.env.PackageManager]
	if cmd == "" {
		return fmt.Errorf("capability %q has no install recipe for %s", name, e.env.PackageManager)
	}

	res := e.execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 15 * time.Minute})
	e.logger.LogExec("capability:"+name, res)
	if res.ExitCode != 0 {
		return fmt.Errorf("install %q failed: %s", name, res.Stderr)
	}

	e.mu.Lock()
	e.installed[name] = true
	e.mu.Unlock()
	return nil
}
