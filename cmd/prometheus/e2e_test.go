//go:build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestCLIFullFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	binPath := buildBinary(t)
	defer os.Remove(binPath)

	port := "19099"
	cmd := exec.Command(binPath, "--web-port", port, "--web-host", "127.0.0.1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start prometheus: %v", err)
	}
	defer cmd.Process.Kill()

	baseURL := fmt.Sprintf("http://127.0.0.1:%s", port)
	waitForService(t, baseURL)

	t.Run("submit_task", func(t *testing.T) {
		resp, err := http.Post(baseURL+"/api/tasks", "application/json",
			bytes.NewBufferString(`{"goal":"test task"}`))
		if err != nil {
			t.Fatalf("failed to submit task: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("unexpected status %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("check_status", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/tasks")
		if err != nil {
			t.Fatalf("failed to get tasks: %v", err)
		}
		defer resp.Body.Close()

		var tasks []map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
	})
}

func TestServiceLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	if !isAdmin() {
		t.Skip("skipping service lifecycle test: not running as admin")
	}

	binPath := buildBinary(t)
	defer os.Remove(binPath)

	t.Run("install", func(t *testing.T) {
		cmd := exec.Command(binPath, "service", "install")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("install failed: %v\n%s", err, output)
		}
	})

	t.Run("start", func(t *testing.T) {
		cmd := exec.Command(binPath, "service", "start")
		if err := cmd.Run(); err != nil {
			t.Fatalf("start failed: %v", err)
		}
		time.Sleep(2 * time.Second)
	})

	t.Run("stop", func(t *testing.T) {
		cmd := exec.Command(binPath, "service", "stop")
		if err := cmd.Run(); err != nil {
			t.Fatalf("stop failed: %v", err)
		}
	})

	t.Run("uninstall", func(t *testing.T) {
		cmd := exec.Command(binPath, "service", "uninstall")
		if err := cmd.Run(); err != nil {
			t.Fatalf("uninstall failed: %v", err)
		}
	})
}

func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "prometheus-test.exe")

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "C:\\Users\\junio\\OneDrive\\AI AGENT HACKATON\\Prometheus\\go_version\\cmd\\prometheus"
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	return binPath
}

func waitForService(t *testing.T, baseURL string) {
	t.Helper()
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL)
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("service did not start in time")
}

func isAdmin() bool {
	cmd := exec.Command("net", "session")
	return cmd.Run() == nil
}
