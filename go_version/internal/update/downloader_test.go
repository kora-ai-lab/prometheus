package update

import (
	"testing"
)

func TestDownloadAndVerify_BadChecksum(t *testing.T) {
	tmpPath, err := DownloadAndVerify(
		"https://example.com/prometheus",
		"invalid-checksum-not-sha256",
	)
	if err == nil {
		t.Error("Expected error with invalid checksum, got nil")
	}
	if tmpPath != "" {
		t.Error("Expected empty path on error")
	}
}

func TestDownloadAndVerify_InvalidHashLength(t *testing.T) {
	_, err := DownloadAndVerify("https://example.com/prometheus", "abc123")
	if err == nil {
		t.Error("Expected error with short hash")
	}
}