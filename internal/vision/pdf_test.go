package vision

import (
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/executor"
)

func TestNewPDFConverter(t *testing.T) {
	pc := NewPDFConverter(nil)
	if pc == nil {
		t.Fatal("NewPDFConverter returned nil")
	}
}

func TestNewPDFConverter_WithCapEngine(t *testing.T) {
	exec := executor.NewShellExecutor()
	capEng := capabilities.NewEngine(exec, nil, nil)
	pc := NewPDFConverter(capEng)
	if pc == nil {
		t.Fatal("NewPDFConverter returned nil")
	}
	if pc.capEngine != capEng {
		t.Error("capEngine not set correctly")
	}
}

func TestPDFConverter_Convert_MissingPDF(t *testing.T) {
	pc := NewPDFConverter(nil)
	_, err := pc.Convert("", "")
	if err == nil {
		t.Error("expected error for empty pdfPath")
	}
}