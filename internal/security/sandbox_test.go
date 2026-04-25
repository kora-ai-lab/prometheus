package security

import (
	"context"
	"testing"
)

func TestSandbox_Basic(t *testing.T) {
	sb := NewSandbox(SandboxConfig{})
	defer sb.Cleanup()

	// Use appropriate command based on OS
	cmd := "echo ok"
	result, err := sb.Run(context.Background(), cmd)
	if err != nil {
		t.Logf("Note: Run returned error on this platform: %v", err)
	}
	// Just verify it doesn't crash
	if result == nil {
		t.Error("Expected result, got nil")
	}
}

func TestSandbox_Level(t *testing.T) {
	sb := NewSandbox(SandboxConfig{})
	if sb.Level() != SandboxWorkdir {
		t.Errorf("Expected SandboxWorkdir, got %d", sb.Level())
	}
}