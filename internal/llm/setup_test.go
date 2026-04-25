package llm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/config"
	"github.com/kora-ai-lab/prometheus/internal/discovery"
)

func TestFirstRunSetupUsesExistingModel(t *testing.T) {
	home := t.TempDir()
	modelsDir := filepath.Join(home, "models")
	if err := os.MkdirAll(modelsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	modelPath := filepath.Join(modelsDir, "a-model.gguf")
	if err := os.WriteFile(modelPath, []byte("fake"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := config.Save(home, config.DefaultConfig()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	env := &discovery.EnvironmentProfile{RAMMb: 4096}
	if err := FirstRunSetup(home, env, ioDiscard{}); err != nil {
		t.Fatalf("FirstRunSetup() error = %v", err)
	}

	cfg, err := config.Load(home)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LLM.ModelPath != modelPath {
		t.Fatalf("ModelPath = %q, want %q", cfg.LLM.ModelPath, modelPath)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
