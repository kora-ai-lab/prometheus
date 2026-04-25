package sast

import (
	"context"
	"os"
	"testing"
)

func TestScanner_ScanFile(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.py")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	
	f.WriteString(`api_key = "sk-test123456789"`)
	f.Close()

	scanner := NewScanner()
	set, err := scanner.ScanFile(context.Background(), f.Name())
	if err != nil {
		t.Fatal(err)
	}
	
	if len(set.Findings) == 0 {
		t.Error("Expected findings, got none")
	}
}