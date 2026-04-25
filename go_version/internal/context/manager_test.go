package context

import (
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

func TestManager_Snapshot(t *testing.T) {
	m := &Manager{
		hotBuffer: []llm.Message{
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there"},
		},
		warmSummary:   `{"goal":"test","done":[],"state":"active"}`,
		contextWindow: 4096,
	}

	snapshot := m.Snapshot()

	if snapshot == nil {
		t.Fatal("Snapshot returned nil")
	}

	if snapshot["warm_summary"] != m.warmSummary {
		t.Errorf("expected warm_summary %q, got %q", m.warmSummary, snapshot["warm_summary"])
	}
}

func TestManager_Restore(t *testing.T) {
	m := &Manager{
		contextWindow: 4096,
	}

	data := map[string]any{
		"warm_summary": `{"goal":"test","done":["task1"]}`,
		"hot_buffer":   `[{"role":"user","content":"test"}]`,
	}

	err := m.Restore(data)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	if m.warmSummary != data["warm_summary"] {
		t.Errorf("expected warm_summary %q, got %q", data["warm_summary"], m.warmSummary)
	}

	if len(m.hotBuffer) != 1 {
		t.Errorf("expected 1 message, got %d", len(m.hotBuffer))
	}
}