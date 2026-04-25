package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoggerAsyncWrites(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer logger.Close()

	time.Sleep(50 * time.Millisecond)

	entries := 10
	for i := 0; i < entries; i++ {
		logger.LogTaskStart("test-"+string(rune('0'+i)), "test goal")
	}

	logger.Close()

	time.Sleep(50 * time.Millisecond)

	logFiles, err := os.ReadDir(filepath.Join(tmpDir, "logs"))
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(logFiles) != 1 {
		t.Errorf("expected 1 log file, got %d", len(logFiles))
	}

	content, _ := os.ReadFile(filepath.Join(tmpDir, "logs", logFiles[0].Name()))
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "task_start") {
			count++
		}
	}
	if count != entries {
		t.Errorf("expected %d entries, got %d", entries, count)
	}
}

func TestLoggerDailyRotation(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	logger.Close()

	time.Sleep(50 * time.Millisecond)

	logFiles, err := os.ReadDir(filepath.Join(tmpDir, "logs"))
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(logFiles) != 1 {
		t.Errorf("expected 1 log file, got %d", len(logFiles))
	}

	today := time.Now().Format("2006-01-02")
	if !strings.Contains(logFiles[0].Name(), today) {
		t.Errorf("expected filename to contain %s, got %s", today, logFiles[0].Name())
	}
	if !strings.HasSuffix(logFiles[0].Name(), ".jsonl") {
		t.Errorf("expected .jsonl suffix, got %s", logFiles[0].Name())
	}
}

func TestLoggerFlushOnClose(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	logger.LogTaskStart("test-1", "goal1")

	logger.Close()

	logFiles, _ := os.ReadDir(filepath.Join(tmpDir, "logs"))
	content, _ := os.ReadFile(filepath.Join(tmpDir, "logs", logFiles[0].Name()))
	if !strings.Contains(string(content), "goal1") {
		t.Error("expected content to contain logged goal")
	}
}

func TestLoggerRedaction(t *testing.T) {
	tmpDir := t.TempDir()
	logger, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	logger.LogTaskStart("test", "password=secret123")

	logger.Close()

	logFiles, _ := os.ReadDir(filepath.Join(tmpDir, "logs"))
	content, _ := os.ReadFile(filepath.Join(tmpDir, "logs", logFiles[0].Name()))
	if strings.Contains(string(content), "password=secret123") {
		t.Error("password should be redacted")
	}
	if !strings.Contains(string(content), "[REDACTED]") {
		t.Error("expected [REDACTED] in output")
	}
}