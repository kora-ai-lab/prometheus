package logging

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/prometheus-dev/prometheus/internal/llm"
)

func TestArchiver_New(t *testing.T) {
	logsDir := t.TempDir()
	summaryDir := t.TempDir()
	archiveDir := t.TempDir()
	provider := &mockProvider{}

	a := NewArchiver(logsDir, summaryDir, archiveDir, provider)
	if a == nil {
		t.Fatal("NewArchiver returned nil")
	}
	if a.logsDir != logsDir {
		t.Errorf("logsDir = %s, want %s", a.logsDir, logsDir)
	}
	if a.summaryDir != summaryDir {
		t.Errorf("summaryDir = %s, want %s", a.summaryDir, summaryDir)
	}
	if a.archiveDir != archiveDir {
		t.Errorf("archiveDir = %s, want %s", a.archiveDir, archiveDir)
	}
}

func TestArchiver_ArchivePreviousMonth(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	monthStr := "2025-03"
	logFile := filepath.Join(logsDir, monthStr+"-15.jsonl.zst")
	logData := `{"ts":"2025-03-15T09:00:00Z","session":"s1","task_id":"task-1","level":"task_start","event":{"goal":"test"}}
{"ts":"2025-03-15T09:05:00Z","session":"s1","task_id":"task-1","level":"task_end","event":{"status":"done"}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	summaryFile := filepath.Join(summaryDir, "2025-03-15.md")
	summaryData := `# Journal 2025-03-15

## En bref
Test day.

## Stats
- Tasks started: 1
- Tasks completed: 1
`
	if err := os.WriteFile(summaryFile, []byte(summaryData), 0644); err != nil {
		t.Fatalf("failed to write summary file: %v", err)
	}

	provider := &mockProvider{}
	a := NewArchiver(logsDir, summaryDir, archiveDir, provider)
	ctx := context.Background()

	originalNow := nowFunc
	nowFunc = func() time.Time {
		return time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC)
	}
	defer func() { nowFunc = originalNow }()

	if err := a.ArchivePreviousMonth(ctx); err != nil {
		t.Fatalf("ArchivePreviousMonth failed: %v", err)
	}

	expectedLogPath := filepath.Join(archiveDir, monthStr, "2025-03-15.jsonl.zst")
	if _, err := os.Stat(expectedLogPath); err != nil {
		t.Errorf("archived log file not found: %v", err)
	}

	expectedSummaryPath := filepath.Join(archiveDir, monthStr, "summary.md")
	if _, err := os.Stat(expectedSummaryPath); err != nil {
		t.Errorf("monthly summary file not found: %v", err)
	}
}

func TestArchiver_MoveMonthLogs(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	destDir := filepath.Join(tmpDir, "archive", "2025-03")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(destDir, 0755)

	monthStr := "2025-03"
	logFile := filepath.Join(logsDir, monthStr+"-01.jsonl.zst")
	if err := os.WriteFile(logFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	provider := &mockProvider{}
	a := NewArchiver(logsDir, t.TempDir(), t.TempDir(), provider)

	if err := a.moveMonthLogs(monthStr, destDir); err != nil {
		t.Fatalf("moveMonthLogs failed: %v", err)
	}

	expectedPath := filepath.Join(destDir, "2025-03-01.jsonl.zst")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("log file not moved: %v", err)
	}

	if _, err := os.Stat(logFile); err == nil {
		t.Error("original log file still exists")
	}
}

func TestArchiver_GenerateMonthlySummary(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	monthStr := "2025-03"
	summaryFile := filepath.Join(summaryDir, "2025-03-15.md")
	summaryData := `# Journal 2025-03-15

## En bref
Test day.

## Stats
- Tasks started: 2
- Tasks completed: 1
- LLM calls: 5
`
	if err := os.WriteFile(summaryFile, []byte(summaryData), 0644); err != nil {
		t.Fatalf("failed to write summary file: %v", err)
	}

	provider := &mockProvider{
		resp: &llm.Response{
			Content: "# Monthly Summary\n\n## Overview\nTest month.",
		},
	}
	a := NewArchiver(logsDir, summaryDir, archiveDir, provider)
	ctx := context.Background()

	destDir := filepath.Join(archiveDir, monthStr)
	os.MkdirAll(destDir, 0755)
	if err := a.generateMonthlySummary(ctx, monthStr, destDir); err != nil {
		t.Fatalf("generateMonthlySummary failed: %v", err)
	}

	summaryPath := filepath.Join(destDir, "summary.md")
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("monthly summary not created: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("monthly summary content is empty")
	}
}

func TestArchiver_GenerateMonthlySummary_NoSummaries(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	provider := &mockProvider{}
	a := NewArchiver(logsDir, summaryDir, archiveDir, provider)
	ctx := context.Background()

	monthStr := "2025-03"
	destDir := filepath.Join(archiveDir, monthStr)
	err := a.generateMonthlySummary(ctx, monthStr, destDir)
	if err == nil {
		t.Error("expected error when no summaries found")
	}
}

func TestArchiver_FallbackMonthlySummary(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)

	monthStr := "2025-03"
	summaryFile := filepath.Join(summaryDir, "2025-03-15.md")
	summaryData := `# Journal 2025-03-15

## Stats
- Tasks started: 2
- Tasks completed: 1
- LLM calls: 5
- Total input tokens: 1000
- Total output tokens: 500
`
	if err := os.WriteFile(summaryFile, []byte(summaryData), 0644); err != nil {
		t.Fatalf("failed to write summary file: %v", err)
	}

	provider := &mockProviderUnavailable{}
	a := NewArchiver(logsDir, summaryDir, archiveDir, provider)
	ctx := context.Background()

	destDir := filepath.Join(archiveDir, monthStr)
	os.MkdirAll(destDir, 0755)

	err := a.generateMonthlySummary(ctx, monthStr, destDir)
	if err != nil {
		t.Fatalf("generateMonthlySummary with fallback failed: %v", err)
	}

	summaryPath := filepath.Join(destDir, "summary.md")
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("monthly summary not created: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("fallback summary content is empty")
	}
	if !strings.Contains(content, "Monthly Journal") {
		t.Error("fallback summary missing header")
	}
}