package update

import (
	"testing"
)

func TestApplyUpdate_DryRun(t *testing.T) {
	err := ApplyUpdate("/nonexistent/prometheus")
	if err == nil {
		t.Error("Expected error with non-existent binary")
	}
}

func TestCopyFile_NonExistent(t *testing.T) {
	err := copyFile("/nonexistent/src", "/nonexistent/dst")
	if err == nil {
		t.Error("Expected error copying nonexistent file")
	}
}