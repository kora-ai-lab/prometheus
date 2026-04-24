package vision

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus-dev/prometheus/internal/executor"
)

func CaptureScreen(ctx context.Context, execer executor.Executor) ([]byte, error) {
	tmpFile := filepath.Join(os.TempDir(), "prometheus-screen.png")
	var cmd string

	switch runtime.GOOS {
	case "windows":
		cmd = `powershell -Command "Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $bounds=[System.Windows.Forms.Screen]::PrimaryScreen.Bounds; $bmp=New-Object System.Drawing.Bitmap $bounds.Width,$bounds.Height; $g=[System.Drawing.Graphics]::FromImage($bmp); $g.CopyFromScreen($bounds.Location,[System.Drawing.Point]::Empty,$bounds.Size); $bmp.Save('` + strings.ReplaceAll(tmpFile, `\`, `\\`) + `',[System.Drawing.Imaging.ImageFormat]::Png)"`
	case "darwin":
		cmd = "screencapture -x " + tmpFile
	default:
		cmd = "sh -c 'if command -v scrot >/dev/null 2>&1; then scrot " + tmpFile + "; elif command -v import >/dev/null 2>&1; then import -window root " + tmpFile + "; else exit 1; fi'"
	}

	res := execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 10 * time.Second})
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("capture failed: %s", res.Stderr)
	}
	defer os.Remove(tmpFile)
	return os.ReadFile(tmpFile)
}
