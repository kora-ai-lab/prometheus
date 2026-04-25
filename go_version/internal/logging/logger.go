package logging

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/prometheus-dev/prometheus/internal/executor"
	"github.com/prometheus-dev/prometheus/internal/llm"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func New(home string) (*Logger, error) {
	logDir := filepath.Join(home, "logs")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(logDir, time.Now().Format("2006-01-02")+".jsonl")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *Logger) write(kind string, payload map[string]any) {
	if l == nil || l.file == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	record := map[string]any{
		"time": time.Now().UTC().Format(time.RFC3339),
		"kind": kind,
	}
	for k, v := range payload {
		record[k] = v
	}
	data, _ := json.Marshal(record)
	_, _ = l.file.Write(append(data, '\n'))
}

func (l *Logger) LogTaskStart(id, goal string) {
	l.write("task_start", map[string]any{
		"task_id": id,
		"goal":    RedactSecrets(goal),
	})
}

func (l *Logger) LogTaskEnd(id string, status any) {
	l.write("task_end", map[string]any{
		"task_id": id,
		"status":  status,
	})
}

func (l *Logger) LogLLMCall(id string, resp *llm.Response, duration time.Duration) {
	payload := map[string]any{
		"task_id":      id,
		"duration_ms":  duration.Milliseconds(),
		"input_tokens": 0,
		"output":       "",
	}
	if resp != nil {
		payload["input_tokens"] = resp.InputTokens
		payload["output_tokens"] = resp.OutputTokens
		payload["output"] = RedactSecrets(resp.Content)
	}
	l.write("llm_call", payload)
}

func (l *Logger) LogExec(id string, result *executor.ExecResult) {
	l.write("exec", map[string]any{
		"task_id":     id,
		"command":     RedactSecrets(result.Command),
		"stdout":      RedactSecrets(result.Stdout),
		"stderr":      RedactSecrets(result.Stderr),
		"exit_code":   result.ExitCode,
		"duration_ms": result.Duration.Milliseconds(),
		"timed_out":   result.TimedOut,
	})
}

func (l *Logger) LogFileCreated(id, path string) {
	l.write("file_created", map[string]any{
		"task_id": id,
		"path":    path,
	})
}

func (l *Logger) LogBrowserAction(id, action string) {
	l.write("browser_action", map[string]any{
		"task_id": id,
		"action":  action,
	})
}

func (l *Logger) LogVisionCapture(id, target string) {
	l.write("vision_capture", map[string]any{
		"task_id": id,
		"target":  target,
	})
}
