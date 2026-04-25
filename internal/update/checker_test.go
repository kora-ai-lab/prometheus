package update

import (
	"testing"
)

func TestCheckForUpdate_GitHubUnreachable(t *testing.T) {
	t.Skip("Requires network - run manually to test")
	hasUpdate, latest, err := CheckForUpdate("anomalyco", "prometheus")
	if err != nil {
		t.Logf("Expected network error or valid response: %v", err)
	}
	_ = hasUpdate
	_ = latest
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		local   string
		remote  string
		want    int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"2.0.0", "1.9.9", 1},
	}

	for _, tt := range tests {
		got, err := CompareVersions(tt.local, tt.remote)
		if err != nil {
			t.Fatalf("CompareVersions(%s, %s) error: %v", tt.local, tt.remote, err)
		}
		if got != tt.want {
			t.Errorf("CompareVersions(%s, %s) = %d, want %d", tt.local, tt.remote, got, tt.want)
		}
	}
}