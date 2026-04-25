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
	mu          sync.Mutex
	file        *os.File
	currentDate string
	home        string
	writeCh     chan writeRequest
	done        chan struct{}
	pending    int
}

type writeRequest struct {
	level  string
	taskID string
	event map[string]any
}

func New(home string) (*Logger, error) {
	logDir := filepath.Join(home, "logs")
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return nil, err
	}
	l := &Logger{
		home:    home,
		writeCh: make(chan writeRequest, 512),
		done:    make(chan struct{}),
	}
	if err := l.rotate(); err != nil {
		return nil, err
	}
	go l.writerLoop()
	return l, nil
}

func (l *Logger) rotate() error {
	today := time.Now().Format("2006-01-02")
	if l.file != nil && l.currentDate == today {
		return nil
	}
	if l.file != nil {
		l.file.Close()
	}
	path := filepath.Join(l.home, "logs", today+".jsonl")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	l.file = file
	l.currentDate = today
	return nil
}

func (l *Logger) writerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var batch []writeRequest
	flush := func() {
		if len(batch) == 0 {
			return
		}
		l.writeSync(batch)
		batch = nil
		l.pending = 0
	}
	for {
		select {
		case <-l.done:
			for len(l.writeCh) > 0 {
				batch = append(batch, <-l.writeCh)
			}
			flush()
			return
		case req := <-l.writeCh:
			batch = append(batch, req)
			l.pending++
			if l.pending >= 100 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}

func (l *Logger) writeSync(batch []writeRequest) {
	if l == nil || l.file == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.rotate(); err != nil {
		return
	}

	for _, req := range batch {
		record := map[string]any{
			"time":   time.Now().UTC().Format(time.RFC3339),
			"kind":   req.level,
			"task_id": req.taskID,
		}
		for k, v := range req.event {
			record[k] = v
		}
		data, _ := json.Marshal(record)
		_, _ = l.file.Write(append(data, '\n'))
	}
}

func (l *Logger) Log(level, taskID string, event map[string]any) {
	if l == nil {
		return
	}
	redacted := make(map[string]any)
	for k, v := range event {
		if s, ok := v.(string); ok {
			redacted[k] = RedactSecrets(s)
		} else {
			redacted[k] = v
		}
	}
	select {
	case l.writeCh <- writeRequest{level: level, taskID: taskID, event: redacted}:
	default:
	}
}

func (l *Logger) Close() error {
	if l == nil {
		return nil
	}
	select {
	case <-l.done:
	default:
		close(l.done)
	}
	time.Sleep(50 * time.Millisecond)
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) write(kind string, payload map[string]any) {
	if l == nil || l.writeCh == nil {
		return
	}
	redacted := make(map[string]any)
	for k, v := range payload {
		if s, ok := v.(string); ok {
			redacted[k] = RedactSecrets(s)
		} else {
			redacted[k] = v
		}
	}
	select {
	case l.writeCh <- writeRequest{level: kind, taskID: "", event: redacted}:
	default:
	}
}

func (l *Logger) LogTaskStart(id, goal string) {
	l.write("task_start", map[string]any{
		"task_id": id,
		"goal":    goal,
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
		"duration_ms": duration.Milliseconds(),
		"input_tokens": 0,
		"output":       "",
	}
	if resp != nil {
		payload["input_tokens"] = resp.InputTokens
		payload["output_tokens"] = resp.OutputTokens
		payload["output"] = resp.Content
	}
	l.write("llm_call", payload)
}

func (l *Logger) LogExec(id string, result *executor.ExecResult) {
	l.write("exec", map[string]any{
		"task_id":     id,
		"command":     result.Command,
		"stdout":      result.Stdout,
		"stderr":      result.Stderr,
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