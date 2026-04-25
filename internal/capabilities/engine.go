package capabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/discovery"
	"github.com/kora-ai-lab/prometheus/internal/executor"
	"github.com/kora-ai-lab/prometheus/internal/logging"
)

type Engine struct {
	execer    executor.Executor
	env       *discovery.EnvironmentProfile
	logger    *logging.Logger
	installed map[string]bool
	mu        sync.Mutex
	forge     *Forge
}

func NewEngine(execer executor.Executor, env *discovery.EnvironmentProfile, logger *logging.Logger) *Engine {
	return &Engine{
		execer:    execer,
		env:       env,
		logger:    logger,
		installed: map[string]bool{},
	}
}

func (e *Engine) SetForge(forge *Forge) {
	e.forge = forge
}

func (e *Engine) Execer() executor.Executor {
	return e.execer
}

func (e *Engine) Ensure(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("empty capability name")
	}

	// 1. Check if already installed
	if _, err := exec.LookPath(name); err == nil {
		return nil
	}

	e.mu.Lock()
	if e.installed[name] {
		e.mu.Unlock()
		return nil
	}
	e.mu.Unlock()

	// 2. Check registry
	capability, ok := Registry[name]
	if ok {
		cmd := capability.InstallCmds[e.env.PackageManager]
		if cmd == "" {
			cmd = capability.InstallCmds["pip"]
		}
		if cmd == "" {
			return fmt.Errorf("capability %q has no install recipe", name)
		}
		res := e.execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 15 * time.Minute})
		if res.ExitCode != 0 {
			return fmt.Errorf("install %q failed: %s", name, res.Stderr)
		}
		e.mu.Lock()
		e.installed[name] = true
		e.mu.Unlock()
		return nil
	}

	// 3. Try package managers
	if cap, pm, err := e.TryPackageManager(ctx, name); err == nil {
		cmd := cap.InstallCmds[pm]
		res := e.execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 10 * time.Minute})
		if res.ExitCode == 0 {
			e.mu.Lock()
			e.installed[name] = true
			e.mu.Unlock()
			return nil
		}
	}

	// 4. Try internet
	if _, source, err := e.TryInternet(ctx, name); err == nil {
		fmt.Printf("Found %s online via %s\n", name, source)
	}

	// 5. Try Forge
	if e.forge != nil {
		result, err := e.forge.Forge(ctx, "create capability: "+name)
		if err != nil {
			return fmt.Errorf("forge failed for %q: %w", name, err)
		}
		if result != nil && result.Verified {
			e.mu.Lock()
			e.installed[name] = true
			e.mu.Unlock()
			return nil
		}
	}

	// 6. Return not found
	return fmt.Errorf("capability %q not found in registry, package managers, or online", name)
}

// TryPackageManager searches for a tool in package managers
func (e *Engine) TryPackageManager(ctx context.Context, name string) (*Capability, string, error) {
	searches := map[string]string{
		"pip":    "pip show " + name,
		"npm":    "npm view " + name + " 2>/dev/null",
		"apt":    "apt-cache show " + name + " 2>/dev/null | head -5",
		"brew":   "brew search " + name + " 2>/dev/null",
		"cargo":  "cargo search " + name + " --limit 1",
	}

	pm := e.env.PackageManager
	if searchCmd, ok := searches[pm]; ok {
		res := e.execer.Execute(ctx, searchCmd, executor.ExecOptions{Timeout: 30 * time.Second})
		if res.ExitCode == 0 && len(res.Stdout) > 10 {
			return &Capability{
				Name:         name,
				InstallCmds:  map[string]string{pm: searchCmd},
				Type:         pm,
				SizeMb:       10,
				Description:  "installed via " + pm,
			}, pm, nil
		}
	}
	return nil, "", fmt.Errorf("not found in package managers")
}

// TryInternet searches PyPI, npm, crates.io for a tool
func (e *Engine) TryInternet(ctx context.Context, name string) (*Capability, string, error) {
	urls := []string{
		"https://pypi.org/pypi/" + name + "/json",
		"https://registry.npmjs.org/" + name,
		"https://crates.io/api/v1/crates/" + name,
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		var data map[string]any
		if json.NewDecoder(resp.Body).Decode(&data) == nil {
			return &Capability{
				Name:         name,
				Type:         "internet",
				SizeMb:       10,
				Description:  "found via internet search",
			}, "internet", nil
		}
	}
	return nil, "", fmt.Errorf("not found online")
}
