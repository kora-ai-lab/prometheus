package browser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestADB_Screenshot(t *testing.T) {
	if !hasADB() {
		t.Skip("adb not available")
	}

	adb := NewADB(nil)
	defer adb.Close()

	path, err := adb.Screenshot()
	if err != nil {
		t.Fatalf("Screenshot() error = %v", err)
	}
	if path == "" {
		t.Error("Screenshot() returned empty path")
	}
	defer os.Remove(path)

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Screenshot file does not exist at %s", path)
	}

	// Verify it's a PNG file
	ext := filepath.Ext(path)
	if ext != ".png" {
		t.Errorf("Expected .png file, got %s", ext)
	}
}

func TestADB_Tap(t *testing.T) {
	if !hasADB() {
		t.Skip("adb not available")
	}

	adb := NewADB(nil)
	defer adb.Close()

	// Test Tap with valid coordinates
	err := adb.Tap(500, 500)
	if err != nil {
		t.Fatalf("Tap(500, 500) error = %v", err)
	}

	// Test Tap with edge coordinates
	err = adb.Tap(0, 0)
	if err != nil {
		t.Fatalf("Tap(0, 0) error = %v", err)
	}

	err = adb.Tap(1080, 1920)
	if err != nil {
		t.Fatalf("Tap(1080, 1920) error = %v", err)
	}
}

func TestADB_InputText(t *testing.T) {
	if !hasADB() {
		t.Skip("adb not available")
	}

	adb := NewADB(nil)
	defer adb.Close()

	// Test InputText with simple text
	err := adb.InputText("hello")
	if err != nil {
		t.Fatalf("InputText('hello') error = %v", err)
	}

	// Test InputText with spaces
	err = adb.InputText("hello world")
	if err != nil {
		t.Fatalf("InputText('hello world') error = %v", err)
	}

	// Test InputText with numbers
	err = adb.InputText("1234567890")
	if err != nil {
		t.Fatalf("InputText('1234567890') error = %v", err)
	}
}
