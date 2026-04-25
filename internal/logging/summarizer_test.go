package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

type mockProvider struct {
	resp *llm.Response
}

func (m *mockProvider) Complete(ctx context.Context, messages []llm.Message) (*llm.Response, error) {
	if m.resp != nil {
		return m.resp, nil
	}
	return nil, fmt.Errorf("no response")
}

func (m *mockProvider) Stream(ctx context.Context, messages []llm.Message, tokens chan<- string) error {
	return nil
}

func (m *mockProvider) ModelInfo() *llm.ModelInfo {
	return &llm.ModelInfo{Name: "test"}
}

func (m *mockProvider) IsAvailable() bool {
	return m.resp != nil
}

func (m *mockProvider) HasVision() bool {
	return false
}

func (m *mockProvider) Close() error {
	return nil
}

type mockProviderUnavailable struct{}

func (m *mockProviderUnavailable) Complete(ctx context.Context, messages []llm.Message) (*llm.Response, error) {
	return nil, fmt.Errorf("provider unavailable")
}

func (m *mockProviderUnavailable) Stream(ctx context.Context, messages []llm.Message, tokens chan<- string) error {
	return fmt.Errorf("provider unavailable")
}

func (m *mockProviderUnavailable) ModelInfo() *llm.ModelInfo {
	return nil
}

func (m *mockProviderUnavailable) IsAvailable() bool {
	return false
}

func (m *mockProviderUnavailable) HasVision() bool {
	return false
}

func (m *mockProviderUnavailable) Close() error {
	return nil
}

func TestSummarizer_New(t *testing.T) {
	logsDir := t.TempDir()
	summaryDir := t.TempDir()
	provider := &mockProvider{}

	s := NewSummarizer(logsDir, summaryDir, provider)
	if s == nil {
		t.Fatal("NewSummarizer returned nil")
	}
	if s.logsDir != logsDir {
		t.Errorf("logsDir = %s, want %s", s.logsDir, logsDir)
	}
	if s.summaryDir != summaryDir {
		t.Errorf("summaryDir = %s, want %s", s.summaryDir, summaryDir)
	}
}

func TestSummarizer_SummarizeDay(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)

	date := "2025-01-15"
	logFile := filepath.Join(logsDir, date+".jsonl")
	logData := `{"ts":"2025-01-15T09:00:00Z","session":"s1","task_id":"task-1","level":"task_start","event":{"goal":"test goal"}}
{"ts":"2025-01-15T09:05:00Z","session":"s1","task_id":"task-1","level":"task_end","event":{"status":"done"}}
{"ts":"2025-01-15T10:00:00Z","session":"s1","task_id":"task-2","level":"llm_call","event":{"input_tokens":100,"output_tokens":50}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	provider := &mockProvider{
		resp: &llm.Response{
			Content:      "# Journal 2025-01-15\n\n## En bref\nTest summary.\n\n## Accompli\n- Completed task-1\n- Made progress on task-2\n\n## Stats\n- 2 tasks\n- 150 total tokens",
			InputTokens:  100,
			OutputTokens: 50,
		},
	}

	s := NewSummarizer(logsDir, summaryDir, provider)
	ctx := context.Background()

	if err := s.SummarizeDay(ctx, date); err != nil {
		t.Fatalf("SummarizeDay failed: %v", err)
	}

	summaryFile := filepath.Join(summaryDir, date+".md")
	data, err := os.ReadFile(summaryFile)
	if err != nil {
		t.Fatalf("summary file not created: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("summary content is empty")
	}
}

func TestSummarizer_SummarizeDay_Compressed(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)

	date := "2025-01-16"
	logFile := filepath.Join(logsDir, date+".jsonl.zst")
	logData := `{"ts":"2025-01-16T09:00:00Z","session":"s1","task_id":"task-3","level":"task_start","event":{"goal":"another test"}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write compressed log file: %v", err)
	}

	provider := &mockProvider{
		resp: &llm.Response{
			Content:      "# Journal 2025-01-16\n\n## En bref\nAnother day.\n\n## Accompli\n- Started task-3\n\n## Stats\n- 1 task",
			InputTokens:  50,
			OutputTokens: 25,
		},
	}

	s := NewSummarizer(logsDir, summaryDir, provider)
	ctx := context.Background()

	if err := s.SummarizeDay(ctx, date); err != nil {
		t.Fatalf("SummarizeDay failed: %v", err)
	}

	summaryFile := filepath.Join(summaryDir, date+".md")
	data, err := os.ReadFile(summaryFile)
	if err != nil {
		t.Fatalf("summary file not created: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("summary content is empty")
	}
}

func TestSummarizer_LoadDayEvents(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logsDir, 0755)

	date := "2025-01-17"
	logFile := filepath.Join(logsDir, date+".jsonl")
	logData := `{"ts":"2025-01-17T09:00:00Z","session":"s1","task_id":"task-1","level":"task_start","event":{"goal":"test"}}
{"ts":"2025-01-17T09:01:00Z","session":"s1","task_id":"task-1","level":"task_end","event":{"status":"done"}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	provider := &mockProvider{}
	s := NewSummarizer(logsDir, t.TempDir(), provider)

	events, err := s.loadDayEvents(date)
	if err != nil {
		t.Fatalf("loadDayEvents failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("got %d events, want 2", len(events))
	}
}

func TestSummarizer_BuildDigest(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logsDir, 0755)

	events := []LogEntry{
		{
			Ts:      "2025-01-18T09:00:00Z",
			Session: "s1",
			TaskID:  "task-1",
			Level:   "task_start",
			Event:   map[string]any{"goal": "test goal"},
		},
		{
			Ts:      "2025-01-18T09:05:00Z",
			Session: "s1",
			TaskID:  "task-1",
			Level:   "task_end",
			Event:   map[string]any{"status": "done"},
		},
		{
			Ts:      "2025-01-18T10:00:00Z",
			Session: "s1",
			TaskID:  "task-2",
			Level:   "llm_call",
			Event:   map[string]any{"input_tokens": 100, "output_tokens": 50},
		},
	}

	provider := &mockProvider{}
	s := NewSummarizer(logsDir, t.TempDir(), provider)

	digest := s.buildDigest(events)
	if len(digest) == 0 {
		t.Error("digest is empty")
	}
}

func TestSummarizer_FallbackSummary(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	summaryDir := filepath.Join(tmpDir, "summaries")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(summaryDir, 0755)

	date := "2025-01-19"
	logFile := filepath.Join(logsDir, date+".jsonl")
	logData := `{"ts":"2025-01-19T09:00:00Z","session":"s1","task_id":"task-1","level":"task_start","event":{"goal":"test"}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	provider := &mockProviderUnavailable{}
	s := NewSummarizer(logsDir, summaryDir, provider)
	ctx := context.Background()

	err := s.SummarizeDay(ctx, date)
	if err != nil {
		t.Fatalf("SummarizeDay with fallback failed: %v", err)
	}

	summaryFile := filepath.Join(summaryDir, date+".md")
	data, err := os.ReadFile(summaryFile)
	if err != nil {
		t.Fatalf("summary file not created: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("fallback summary content is empty")
	}
}

func TestLogEntry_JSONFields(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logsDir, 0755)

	date := "2025-01-20"
	logFile := filepath.Join(logsDir, date+".jsonl")
	logData := `{"time":"2025-01-20T09:00:00Z","session":"s1","task_id":"task-x","kind":"task_start"}
{"ts":"2025-01-20T09:01:00Z","session":"s1","level":"info","event":{"msg":"hello"}}
`
	if err := os.WriteFile(logFile, []byte(logData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	provider := &mockProvider{}
	s := NewSummarizer(logsDir, t.TempDir(), provider)

	events, err := s.loadDayEvents(date)
	if err != nil {
		t.Fatalf("loadDayEvents failed: %v", err)
	}

	if len(events) != 2 {
		t.Errorf("got %d events, want 2", len(events))
	}
}