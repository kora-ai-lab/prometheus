package vision

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/prometheus-dev/prometheus/internal/capabilities"
	"github.com/prometheus-dev/prometheus/internal/executor"
)

type PDFConverter struct {
	capEngine *capabilities.Engine
	execer   executor.Executor
}

func NewPDFConverter(capEngine *capabilities.Engine) *PDFConverter {
	var execer executor.Executor
	if capEngine != nil {
		execer = capEngine.Execer()
	}
	if execer == nil {
		execer = executor.NewShellExecutor()
	}
	return &PDFConverter{
		capEngine: capEngine,
		execer:   execer,
	}
}

func (p *PDFConverter) Convert(pdfPath, outputDir string) ([]string, error) {
	if pdfPath == "" {
		return nil, fmt.Errorf("pdfPath is required")
	}
	if outputDir == "" {
		outputDir = filepath.Dir(pdfPath)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output dir: %w", err)
	}

	ctx := context.Background()

	if p.capEngine != nil {
		if err := p.capEngine.Ensure(ctx, "pdftoppm"); err != nil {
			return nil, fmt.Errorf("failed to ensure pdftoppm: %w", err)
		}
	}

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), ".pdf")
	outputPattern := filepath.Join(outputDir, baseName)

	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = fmt.Sprintf("pdftoppm -png -singlefile %s %s", pdfPath, outputPattern)
	default:
		cmd = fmt.Sprintf("pdftoppm -png -singlefile %s %s", pdfPath, outputPattern)
	}

	res := p.execer.Execute(ctx, cmd, executor.ExecOptions{Timeout: 2 * time.Minute})
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("pdftoppm failed: %s", res.Stderr)
	}

	pngPath := outputPattern + "-1.png"
	if _, err := os.Stat(pngPath); err != nil {
		pngPath = outputPattern + ".png"
	}

	var pngs []string
	if _, err := os.Stat(pngPath); err == nil {
		pngs = append(pngs, pngPath)
	}

	if len(pngs) == 0 {
		files, err := os.ReadDir(outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read output dir: %w", err)
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".png") {
				pngs = append(pngs, filepath.Join(outputDir, f.Name()))
			}
		}
	}

	return pngs, nil
}