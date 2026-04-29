package context

import (
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

func BenchmarkContextManager_Add(b *testing.B) {
	m := New(nil)
	msg := llm.Message{Role: "user", Content: "benchmark message content"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Add(msg)
	}
}

func BenchmarkContextManager_BuildMessages(b *testing.B) {
	m := New(nil)
	for i := 0; i < 10; i++ {
		m.Add(llm.Message{Role: "user", Content: "test message"})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.BuildMessages("system prompt")
	}
}

func BenchmarkContextManager_Snapshot(b *testing.B) {
	m := New(nil)
	for i := 0; i < 10; i++ {
		m.Add(llm.Message{Role: "user", Content: "test message"})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Snapshot()
	}
}
