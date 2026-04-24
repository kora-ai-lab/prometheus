package task

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus-dev/prometheus/internal/executor"
	"github.com/prometheus-dev/prometheus/internal/llm"
	"github.com/prometheus-dev/prometheus/internal/prompt"
)

func (t *Task) Run(ctx context.Context, deps *TaskDeps) error {
	deps.Logger.LogTaskStart(t.ID, t.Goal)
	defer deps.Logger.LogTaskEnd(t.ID, t.Status)

	for t.Status == StatusRunning {
		messages := deps.PromptBuilder.BuildMessages(t.Context)
		start := time.Now()
		resp, err := deps.Provider.Complete(ctx, messages)
		deps.Logger.LogLLMCall(t.ID, resp, time.Since(start))
		if err != nil {
			t.Retries++
			if t.Retries >= t.MaxRetries {
				t.Status = StatusFailed
				return deps.TaskStore.Save(t)
			}
			t.Context = append(t.Context, llm.Message{Role: "user", Content: "LLM error: " + err.Error()})
			continue
		}

		t.Context = append(t.Context, llm.Message{Role: "assistant", Content: resp.Content})
		action, parseErr := prompt.ParseAction(resp.Content)
		if parseErr != nil {
			t.ParseErrors++
			if t.ParseErrors >= t.MaxParseErrors {
				t.Status = StatusFailed
				return deps.TaskStore.Save(t)
			}
			t.Context = append(t.Context, llm.Message{
				Role:    "user",
				Content: prompt.FormatParseRepair(parseErr, resp.Content),
			})
			continue
		}
		t.ParseErrors = 0

		var observation string
		switch action.Action {
		case "exec":
			if allowed, secErr := deps.Security.Allow(action.Command); !allowed {
				observation = "BLOCKED BY SECURITY: " + secErr.Error()
			} else {
				result := deps.Executor.Execute(ctx, action.Command, executor.ExecOptions{})
				deps.Logger.LogExec(t.ID, result)
				observation = formatObservation(result)
				if isCommandNotFound(result) && deps.CapEngine != nil {
					tool := extractToolName(action.Command)
					if ensureErr := deps.CapEngine.Ensure(ctx, tool); ensureErr == nil {
						result = deps.Executor.Execute(ctx, action.Command, executor.ExecOptions{})
						deps.Logger.LogExec(t.ID, result)
						observation = "[TOOL INSTALLED]\n" + formatObservation(result)
					}
				}
			}
		case "create":
			cf := action.CreateFile
			if err := os.MkdirAll(filepath.Dir(cf.Path), 0o755); err != nil {
				observation = "ERROR mkdir: " + err.Error()
			} else if err := os.WriteFile(cf.Path, []byte(cf.Content), 0o644); err != nil {
				observation = "ERROR write: " + err.Error()
			} else {
				observation = fmt.Sprintf("FILE_CREATED: %s (%d bytes)", cf.Path, len(cf.Content))
				deps.Logger.LogFileCreated(t.ID, cf.Path)
			}
		case "browser":
			observation = deps.Browser.Do(ctx, action)
			deps.Logger.LogBrowserAction(t.ID, action.BrowserAction)
		case "vision":
			observation = deps.Browser.VisionResult(ctx, action)
			deps.Logger.LogVisionCapture(t.ID, action.VisionTarget)
		case "ask":
			t.Status = StatusBlocked
			t.BlockedReason = action.Question
			t.UpdatedAt = time.Now()
			return deps.TaskStore.Save(t)
		case "done":
			t.Status = StatusDone
			t.UpdatedAt = time.Now()
			return deps.TaskStore.Save(t)
		case "error":
			t.Retries++
			observation = "Agent reported an error and must retry with another approach."
			if t.Retries >= t.MaxRetries {
				t.Status = StatusFailed
				t.UpdatedAt = time.Now()
				return deps.TaskStore.Save(t)
			}
		}

		t.Context = append(t.Context, llm.Message{Role: "user", Content: observation})
		t.UpdatedAt = time.Now()
		if err := deps.TaskStore.Save(t); err != nil {
			return err
		}
	}
	return nil
}

func formatObservation(result *executor.ExecResult) string {
	return fmt.Sprintf(
		"COMMAND: %s\nEXIT: %d\nSTDOUT:\n%s\nSTDERR:\n%s",
		result.Command,
		result.ExitCode,
		result.Stdout,
		result.Stderr,
	)
}

func isCommandNotFound(result *executor.ExecResult) bool {
	text := strings.ToLower(result.Stderr + "\n" + result.Stdout)
	return strings.Contains(text, "not found") || strings.Contains(text, "is not recognized")
}

func extractToolName(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
