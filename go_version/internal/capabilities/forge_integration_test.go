package capabilities

import (
	"context"
	"os"
	"testing"

	"github.com/prometheus-dev/prometheus/internal/llm"
)

func TestForge_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test (run without -short to execute)")
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		t.Skip("GROQ_API_KEY not set, skipping integration test")
	}

	llm := llm.NewGroqProvider(apiKey, "")
	if !llm.IsAvailable() {
		t.Skip("Groq provider not available, skipping integration test")
	}

	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)
	tester := NewTester("", tmpDir)

	forge := NewForge(llm, storage, tester, nil)

	result, err := forge.Forge(context.Background(), "a tool that prints hello world")
	if err != nil {
		t.Fatalf("Forge() error = %v", err)
	}

	if result.Path == "" {
		t.Error("result.Path is empty")
	}
}