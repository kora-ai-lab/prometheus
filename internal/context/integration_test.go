//go:build integration

package context

import (
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

func TestSnapshotRestoreRoundtrip(t *testing.T) {
	t.Parallel()
	m := New(nil)

	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
		{Role: "user", Content: "How are you?"},
	}

	for _, msg := range messages {
		m.Add(msg)
	}

	snapshot := m.Snapshot()
	m.Restore(snapshot)

	restored := m.BuildMessages("system")
	if len(restored) < 2 {
		t.Errorf("expected at least system + messages, got %d", len(restored))
	}
}

func TestContextCompactionFlow(t *testing.T) {
	t.Parallel()
	m := New(nil)

	for i := 0; i < 20; i++ {
		m.Add(llm.Message{
			Role:    "user",
			Content: "This is a test message with some content to fill up the context window gradually",
		})
	}

	messages := m.BuildMessages("system prompt")
	if len(messages) < 1 {
		t.Error("expected messages after compaction")
	}
}

func TestMultiSessionResume(t *testing.T) {
	t.Parallel()
	m := New(nil)

	m.Add(llm.Message{Role: "user", Content: "Session 1 message"})
	snapshot1 := m.Snapshot()

	m2 := New(nil)
	m2.Restore(snapshot1)
	m2.Add(llm.Message{Role: "user", Content: "Session 2 message"})

	messages := m2.BuildMessages("system")
	found := false
	for _, msg := range messages {
		if msg.Content == "Session 2 message" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Session 2 message in built messages")
	}
}

func TestBuildMessagesWithWarmSummary(t *testing.T) {
	t.Parallel()
	m := New(nil)
	m.warmSummary = "Previous conversation summary"

	messages := m.BuildMessages("system prompt")
	if len(messages) != 2 {
		t.Errorf("expected 2 messages (system + warm summary), got %d", len(messages))
	}
	if messages[1].Content != "[PREVIOUS_CONTEXT]\nPrevious conversation summary" {
		t.Error("expected warm summary in second message")
	}
}
