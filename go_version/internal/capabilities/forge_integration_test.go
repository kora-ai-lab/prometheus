package capabilities

import (
	"testing"
)

func TestForge_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test (run without -short to execute)")
	}

	t.Skip("Integration test requires LLM wrapper implementation - skipping")
}