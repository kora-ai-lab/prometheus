package vision

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/executor"
)

type Capture struct {
	outputDir  string
	capEngine  *capabilities.Engine
	execer     executor.Executor
	triggers   map[string]bool
}

func NewCapture(outputDir string, capEngine *capabilities.Engine) *Capture {
	var execer executor.Executor
	if capEngine != nil {
		execer = capEngine.Execer()
	}
	if execer == nil {
		execer = executor.NewShellExecutor()
	}
	if outputDir == "" {
		outputDir = "/tmp/screenshots"
	}
	return &Capture{
		outputDir: outputDir,
		capEngine: capEngine,
		execer:    execer,
		triggers: map[string]bool{
			"server_started": true,
			"html_created":   true,
			"explicit":      true,
		},
	}
}

func (c *Capture) DetectOS() string {
	return runtime.GOOS
}

func (c *Capture) Capture(filename string) (string, error) {
	if c.outputDir == "" {
		c.outputDir = "/tmp/screenshots"
	}
	if err := os.MkdirAll(c.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output dir: %w", err)
	}

	outputPath := filepath.Join(c.outputDir, filename)
	if err := c.captureToFile(outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (c *Capture) captureToFile(outputPath string) error {
	ctx := context.Background()
	tmpFile := filepath.Join(os.TempDir(), "prometheus-screen.png")
	var cmd string

	switch runtime.GOOS {
	case "windows":
		cmd = `powershell -Command "Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $bounds=[System.Windows.Forms.Screen]::PrimaryScreen.Bounds; $bmp=New-Object System.Drawing.Bitmap $bounds.Width,$bounds.Height; $g=[System.Drawing.Graphics]::FromImage($bmp); $g.CopyFromScreen($bounds.Location,[System.Drawing.Point]::Empty,$bounds.Size); $bmp.Save('` + strings.ReplaceAll(tmpFile, `\`, `\\`) + `',[System.Drawing.Imaging.ImageFormat]::Png)"`
	case "darwin":
		cmd = "screencapture -x " + tmpFile
	default:
		if c.capEngine != nil {
			_ = c.capEngine.Ensure(ctx, "scrot")
		}
		cmd = "sh -c 'if command -v scrot >/dev/null 2>&1; then scrot " + tmpFile + "; elif command -v import >/dev/null 2>&1; then import -window root " + tmpFile + "; else exit 1; fi'"
	}

	res := c.execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 10 * time.Second})
	if res.ExitCode != 0 {
		return fmt.Errorf("capture failed: %s", res.Stderr)
	}

	return c.copyToOutput(tmpFile, outputPath)
}

func (c *Capture) copyToOutput(tmpFile, outputPath string) error {
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to read capture: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	os.Remove(tmpFile)
	return nil
}

func (c *Capture) Trigger(event string, data string) error {
	if !c.triggers[event] {
		return fmt.Errorf("unknown event: %s", event)
	}
	filename := fmt.Sprintf("capture_%s_%d.png", event, time.Now().Unix())
	_, err := c.Capture(filename)
	return err
}