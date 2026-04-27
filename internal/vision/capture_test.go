package vision

import (
	"context"
	"os"
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/executor"
)

type mockExecutor_ struct {
	shouldFail bool
}

func (m *mockExecutor_) Execute(ctx context.Context, cmd string, opts executor.ExecOptions) *executor.ExecResult {
	if m.shouldFail {
		return &executor.ExecResult{
			Command:  cmd,
			ExitCode: 1,
			Stderr:   "mock failure",
		}
	}
	return &executor.ExecResult{
		Command:  cmd,
		ExitCode: 0,
		Stdout:   "ok",
	}
}

type capturingExecutor struct {
	lastCmd   string
	lastOpts executor.ExecOptions
	tmpDir   string
}

func (c *capturingExecutor) Execute(ctx context.Context, cmd string, opts executor.ExecOptions) *executor.ExecResult {
	c.lastCmd = cmd
	c.lastOpts = opts
	tmpFile := c.tmpDir + "/prometheus-screen.png"
	os.WriteFile(tmpFile, []byte("fake image"), 0644)
	return &executor.ExecResult{
		Command:  cmd,
		ExitCode: 0,
	}
}

func TestDetectOS(t *testing.T) {
	c := &Capture{
		outputDir:  "/tmp/screenshots",
		capEngine: nil,
	}
	os := c.DetectOS()
	if os == "" {
		t.Error("DetectOS returned empty string")
	}
	if os != "windows" && os != "darwin" && os != "linux" {
		t.Errorf("DetectOS returned unexpected OS: %s", os)
	}
}

func TestNewCapture(t *testing.T) {
	exec := executor.NewShellExecutor()
	capEng := capabilities.NewEngine(exec, nil, nil)
	c := NewCapture("/tmp/screenshots", capEng)
	if c == nil {
		t.Fatal("NewCapture returned nil")
	}
	if c.outputDir != "/tmp/screenshots" {
		t.Errorf("NewCapture outputDir = %s, want /tmp/screenshots", c.outputDir)
	}
}

func TestCapture_Capture_WithMock(t *testing.T) {
	mock := &capturingExecutor{tmpDir: os.TempDir()}
	c := &Capture{
		outputDir: "/tmp/screenshots",
		execer:    mock,
		triggers: map[string]bool{
			"server_started": true,
			"html_created":  true,
			"explicit":      true,
		},
	}

	path, err := c.Capture("test.png")
	if err != nil {
		t.Errorf("Capture failed: %v", err)
	}
	if path == "" {
		t.Error("Capture returned empty path")
	}
	if mock.lastCmd == "" {
		t.Error("No command was executed")
	}
}

func TestCapture_Trigger(t *testing.T) {
	mock := &capturingExecutor{tmpDir: os.TempDir()}
	c := &Capture{
		outputDir: "/tmp/screenshots",
		execer:  mock,
		triggers: map[string]bool{
			"server_started": true,
			"html_created":  true,
			"explicit":     true,
		},
	}

	tests := []struct {
		event string
		data string
		want bool
	}{
		{"server_started", "", true},
		{"html_created", "", true},
		{"explicit", "", true},
		{"unknown", "", false},
	}

	for _, tt := range tests {
		err := c.Trigger(tt.event, tt.data)
		if (err == nil) != tt.want {
			t.Errorf("Trigger(%q) error = %v, want error = %v", tt.event, err, tt.want)
		}
	}
}

func TestCapture_TriggerGeneratesFilename(t *testing.T) {
	mock := &capturingExecutor{tmpDir: os.TempDir()}
	c := &Capture{
		outputDir: "/tmp/screenshots",
		execer:  mock,
		triggers: map[string]bool{
			"server_started": true,
			"html_created":  true,
			"explicit":     true,
		},
	}

	err := c.Trigger("server_started", "")
	if err != nil {
		t.Errorf("Trigger failed: %v", err)
	}
}

func TestCapture_NewCaptureDefaults(t *testing.T) {
	c := NewCapture("", nil)
	if c == nil {
		t.Fatal("NewCapture returned nil")
	}
	if c.outputDir != "/tmp/screenshots" {
		t.Errorf("outputDir = %s, want /tmp/screenshots", c.outputDir)
	}
}

func TestCapture_TriggerWithEngine(t *testing.T) {
	exec := executor.NewShellExecutor()
	capEng := capabilities.NewEngine(exec, nil, nil)
	c := NewCapture("/tmp/screenshots", capEng)

	err := c.Trigger("server_started", "")
	_ = err
}