package security

import (
	"os"
	"runtime"
	"testing"
)

func TestCheckPermissions(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	os.Chmod(f.Name(), 0777)
	
	findings := CheckPermissions([]string{f.Name()})
	
	if len(findings) == 0 {
		t.Error("Expected to find world-writable")
	}
	
	if findings[0].Severity != "high" {
		t.Errorf("Expected high severity, got %s", findings[0].Severity)
	}
}

func TestCheckPermissions_Private(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("os.Chmod doesn't work on Windows")
	}
	
	f, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	os.Chmod(f.Name(), 0600)
	
	findings := CheckPermissions([]string{f.Name()})
	
	if len(findings) != 0 {
		t.Error("Private file should not have findings")
	}
}