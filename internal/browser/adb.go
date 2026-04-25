package browser

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ADB represents an Android Debug Bridge client for mobile automation
type ADB struct {
	device interface{} // capEngine is optional, stored as interface to avoid import
}

// NewADB creates a new ADB client and auto-installs adb if needed
func NewADB(capEngine interface{}) *ADB {
	// Try to call Ensure if capEngine implements the right interface
	if capEngine != nil {
		if ce, ok := capEngine.(interface{ Ensure(context.Context, string) error }); ok {
			ce.Ensure(context.Background(), "adb")
		}
	}

	// Get device ID from adb devices output
	device := detectDevice()

	return &ADB{
		device: device,
	}
}

// detectDevice finds the first connected Android device
func detectDevice() string {
	cmd := exec.Command("adb", "devices")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "default"
	}

	// Parse adb devices output to find first device
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == "device" {
			return parts[0]
		}
	}

	return "default"
}

// Screenshot captures a screenshot from the device and saves it as PNG to a temp file
func (a *ADB) Screenshot() (string, error) {
	device := "default"
	if d, ok := a.device.(string); ok {
		device = d
	}
	tmpDir := os.TempDir()
	localPath := filepath.Join(tmpDir, fmt.Sprintf("adb_screenshot_%d.png", time.Now().UnixNano()))

	// Execute adb exec-out screencap -p to get screenshot
	cmd := exec.Command("adb", "-s", device, "exec-out", "screencap", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("screencap failed: %w", err)
	}

	// Write to temporary file
	if err := os.WriteFile(localPath, out.Bytes(), 0600); err != nil {
		return "", fmt.Errorf("write screenshot file: %w", err)
	}

	return localPath, nil
}

// Tap touches the screen at the specified coordinates
func (a *ADB) Tap(x, y int) error {
	device := "default"
	if d, ok := a.device.(string); ok {
		device = d
	}
	cmd := exec.Command("adb", "-s", device, "shell", "input", "tap", fmt.Sprintf("%d", x), fmt.Sprintf("%d", y))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tap(%d, %d) failed: %w", x, y, err)
	}
	return nil
}

// InputText sends text input to the device
func (a *ADB) InputText(text string) error {
	device := "default"
	if d, ok := a.device.(string); ok {
		device = d
	}
	cmd := exec.Command("adb", "-s", device, "shell", "input", "text", text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("input text failed: %w", err)
	}
	return nil
}

// Close performs cleanup
func (a *ADB) Close() error {
	// No resources to cleanup for now
	return nil
}

// hasADB checks if adb command is available
func hasADB() bool {
	cmd := exec.Command("adb", "version")
	return cmd.Run() == nil
}
