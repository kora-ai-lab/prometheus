package logging

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SecurityLogger struct {
	file *os.File
	mu   sync.Mutex
}

func NewSecurityLogger(home string) (*SecurityLogger, error) {
	secDir := filepath.Join(home, "security")
	if err := os.MkdirAll(secDir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(secDir, "events.jsonl")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}
	return &SecurityLogger{file: file}, nil
}

func (sl *SecurityLogger) Close() error {
	if sl == nil || sl.file == nil {
		return nil
	}
	return sl.file.Close()
}

func (sl *SecurityLogger) write(event string, payload map[string]any) {
	if sl == nil || sl.file == nil {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()

	record := map[string]any{
		"time": time.Now().UTC().Format(time.RFC3339),
		"event": event,
	}
	for k, v := range payload {
		record[k] = v
	}
	data, _ := json.Marshal(record)
	_, _ = sl.file.Write(append(data, '\n'))
}

func (sl *SecurityLogger) LogBlockedCommand(cmd string, reasons []string) {
	sl.write("command_blocked", map[string]any{
		"command": RedactSecrets(cmd),
		"reasons": reasons,
	})
}

func (sl *SecurityLogger) LogConfirmedCommand(cmd string) {
	sl.write("command_confirmed", map[string]any{
		"command": RedactSecrets(cmd),
	})
}

func (sl *SecurityLogger) LogSASTFinding(cmd, finding string) {
	sl.write("sast_finding", map[string]any{
		"command": RedactSecrets(cmd),
		"finding": finding,
	})
}

func (sl *SecurityLogger) LogPortDetected(port int) {
	sl.write("port_detected", map[string]any{
		"port": port,
	})
}