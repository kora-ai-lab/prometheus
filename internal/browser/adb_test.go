package browser

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewADB_NoCapEngine(t *testing.T) {
	adb := NewADB(nil)
	if adb == nil {
		t.Fatal("NewADB returned nil")
	}
	if adb.device == nil {
		t.Error("device should not be nil")
	}
}

func TestADB_Screenshot_NoDevice(t *testing.T) {
	// Without adb installed, Screenshot should fail gracefully
	adb := &ADB{device: "default"}
	path, err := adb.Screenshot()

	// If adb is not installed, we expect an error
	if err == nil {
		// adb is installed - verify file was created
		if path == "" {
			t.Error("expected non-empty path")
		} else {
			defer os.Remove(path)
			if !strings.HasSuffix(path, ".png") {
				t.Errorf("expected .png extension, got %s", path)
			}
		}
	}
}

func TestADB_Tap_NoDevice(t *testing.T) {
	adb := &ADB{device: "default"}
	err := adb.Tap(100, 200)

	// Without adb, this should fail
	if err == nil {
		// adb is installed and device connected - acceptable
		t.Log("Tap succeeded (device connected)")
	} else if !strings.Contains(err.Error(), "failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestADB_InputText_NoDevice(t *testing.T) {
	adb := &ADB{device: "default"}
	err := adb.InputText("hello world")

	if err == nil {
		t.Log("InputText succeeded (device connected)")
	} else if !strings.Contains(err.Error(), "failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestADB_Close(t *testing.T) {
	adb := &ADB{device: "default"}
	err := adb.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestDetectDevice_NoADB(t *testing.T) {
	// This test depends on whether adb is installed
	// Just verify it doesn't panic
	device := detectDevice()
	if device == "" {
		t.Error("detectDevice returned empty string")
	}
}

func TestHasADB(t *testing.T) {
	// Just verify it doesn't panic
	result := hasADB()
	t.Logf("adb available: %v", result)
}

func TestADB_Screenshot_FilePermissions(t *testing.T) {
	// Skip if adb not available
	if !hasADB() {
		t.Skip("adb not available")
	}

	adb := &ADB{device: "default"}
	path, err := adb.Screenshot()
	if err != nil {
		t.Fatalf("Screenshot failed: %v", err)
	}
	defer os.Remove(path)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Cannot stat screenshot: %v", err)
	}

	// Verify file permissions (0600 = owner read/write only)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestADB_Tap_InvalidCoords(t *testing.T) {
	adb := &ADB{device: "default"}

	// Negative coordinates should still execute (device may reject)
	_ = adb.Tap(-1, -1)
}

func TestADB_InputText_Empty(t *testing.T) {
	adb := &ADB{device: "default"}

	// Empty text should still execute
	_ = adb.InputText("")
}

func TestADB_Screenshot_ConcurrentPaths(t *testing.T) {
	// Verify that multiple screenshot calls produce unique paths
	adb := &ADB{device: "default"}

	path1, err1 := adb.Screenshot()
	if err1 != nil {
		t.Skip("adb not available")
	}
	defer os.Remove(path1)

	path2, err2 := adb.Screenshot()
	if err2 != nil {
		t.Fatalf("Second screenshot failed: %v", err2)
	}
	defer os.Remove(path2)

	if path1 == path2 {
		t.Error("consecutive screenshots should have unique paths")
	}
}

func TestADB_DeviceID(t *testing.T) {
	// Test with explicit device ID
	adb := &ADB{device: "emulator-5554"}

	if d, ok := adb.device.(string); !ok || d != "emulator-5554" {
		t.Errorf("expected device 'emulator-5554', got %v", adb.device)
	}
}

func TestNewADB_WithCapEngine(t *testing.T) {
	// Test with a mock capEngine that implements Ensure
	mockEngine := &mockCapEngine{ensureCalled: false}
	adb := NewADB(mockEngine)

	if adb == nil {
		t.Fatal("NewADB returned nil")
	}
	if !mockEngine.ensureCalled {
		t.Error("Ensure should have been called on capEngine")
	}
}

type mockCapEngine struct {
	ensureCalled bool
}

func (m *mockCapEngine) Ensure(ctx context.Context, name string) error {
	m.ensureCalled = true
	return nil
}
