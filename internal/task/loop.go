package task

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/executor"
	"github.com/kora-ai-lab/prometheus/internal/llm"
	"github.com/kora-ai-lab/prometheus/internal/prompt"
)

func (t *Task) Run(ctx context.Context, deps *TaskDeps) error {
	deps.Logger.LogTaskStart(t.ID, t.Goal)
	defer deps.Logger.LogTaskEnd(t.ID, t.Status)

	t.SetProgress("Initializing...")
	if err := deps.TaskStore.Save(t); err != nil {
		return err
	}

	for t.Status == StatusRunning {
		t.SetProgress("Thinking...")
		if err := deps.TaskStore.Save(t); err != nil {
			return err
		}

		messages := deps.PromptBuilder.BuildMessages(t.Context)
		start := time.Now()
		resp, err := deps.Provider.Complete(ctx, messages)
		deps.Logger.LogLLMCall(t.ID, resp, time.Since(start))
		if err != nil {
			errMsg := err.Error()
			isRetryable := strings.Contains(errMsg, "context deadline exceeded") ||
				strings.Contains(errMsg, "i/o timeout") ||
				strings.Contains(errMsg, "Server closed connection")
			if !isRetryable || t.Retries >= t.MaxRetries {
				t.Status = StatusFailed
				t.Error = errMsg
				t.SetProgress("Failed: " + errMsg)
				return deps.TaskStore.Save(t)
			}
			t.Retries++
			t.Context = append(t.Context, llm.Message{Role: "user", Content: "LLM timeout, retrying..."})
			continue
		}

		t.Context = append(t.Context, llm.Message{Role: "assistant", Content: resp.Content})
		action, parseErr := prompt.ParseAction(resp.Content)
		if parseErr != nil {
			t.ParseErrors++
			if t.ParseErrors >= t.MaxParseErrors {
				t.Status = StatusFailed
				t.Error = parseErr.Error()
				t.SetProgress("Failed: " + parseErr.Error())
				return deps.TaskStore.Save(t)
			}
			t.Context = append(t.Context, llm.Message{
				Role:    "user",
				Content: prompt.FormatParseRepair(parseErr, resp.Content),
			})
			continue
		}
		t.ParseErrors = 0

		t.SetProgress("Executing: " + action.Action)
		if err := deps.TaskStore.Save(t); err != nil {
			return err
		}

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
		case "create", "create_file", "crAcer_fichier":
			cf := action.CreateFile
			if cf == nil && action.Command != "" {
				cf = &prompt.CreateFileAction{
					Path:    extractFilePath(action.Command),
					Content: extractFileContent(action.Command),
				}
			}
			if cf == nil || cf.GetPath() == "" {
				observation = "ERROR: create action missing file info"
			} else {
				abs, _ := filepath.Abs(cf.GetPath())
				dir := filepath.Dir(abs)
				if dir != "." && dir != "" {
					if err := os.MkdirAll(dir, 0o755); err != nil {
						observation = "ERROR mkdir: " + err.Error()
					} else if err := os.WriteFile(abs, []byte(cf.GetContent()), 0o644); err != nil {
						observation = "ERROR write: " + err.Error()
					} else {
						observation = fmt.Sprintf("FILE_CREATED: %s (%d bytes)", abs, len(cf.GetContent()))
						deps.Logger.LogFileCreated(t.ID, abs)
					}
				} else {
					if err := os.WriteFile(abs, []byte(cf.GetContent()), 0o644); err != nil {
						observation = "ERROR write: " + err.Error()
					} else {
						observation = fmt.Sprintf("FILE_CREATED: %s (%d bytes)", abs, len(cf.GetContent()))
						deps.Logger.LogFileCreated(t.ID, abs)
					}
				}
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
			t.SetProgress("Waiting for input: " + action.Question)
			return deps.TaskStore.Save(t)
		case "done":
			t.Status = StatusDone
			if action.Command != "" {
				t.Result = action.Command
			} else {
				t.Result = action.Thinking
			}
			t.SetProgress("Done")
			return deps.TaskStore.Save(t)
		case "error":
			t.Retries++
			observation = "Agent reported an error and must retry with another approach."
			if t.Retries >= t.MaxRetries {
				t.Status = StatusFailed
				t.Error = "max retries exceeded"
				t.SetProgress("Failed: max retries exceeded")
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

func extractFilePath(command string) string {
	fields := strings.Fields(command)
	for i, f := range fields {
		if f == ">" && i+1 < len(fields) {
			return strings.TrimSpace(fields[i+1])
		}
		if f == "echo" && i+2 < len(fields) {
			return strings.TrimSpace(fields[len(fields)-1])
		}
		if f == "touch" || f == "cat" || f == "new-item" {
			if i+1 < len(fields) {
				return strings.TrimSpace(fields[i+1])
			}
		}
	}
	if len(fields) >= 2 {
		return strings.TrimSpace(fields[len(fields)-1])
	}
	return ""
}

func extractFileContent(command string) string {
	fields := strings.Fields(command)
	for i, f := range fields {
		if f == ">" && i > 0 {
			content := strings.Join(fields[:i], " ")
			return strings.TrimSpace(content)
		}
		if f == "-Content" && i+1 < len(fields) {
			return strings.TrimSpace(fields[i+1])
		}
	}
	return ""
}

func totalChars(messages []llm.Message) int {
	n := 0
	for _, m := range messages {
		n += len(m.Content)
	}
	return n
}
