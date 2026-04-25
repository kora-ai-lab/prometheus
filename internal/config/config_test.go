package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	home := t.TempDir()
	cfg := DefaultConfig()
	cfg.LLM.ModelPath = filepath.Join(home, "models", "phi.gguf")
	cfg.LLM.ServerPath = filepath.Join(home, "runtime", "llama-server")
	cfg.UI.WebEnabled = true

	if err := Save(home, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := Load(home)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got.LLM.ModelPath != cfg.LLM.ModelPath {
		t.Fatalf("ModelPath = %q, want %q", got.LLM.ModelPath, cfg.LLM.ModelPath)
	}
	if got.LLM.ServerPath != cfg.LLM.ServerPath {
		t.Fatalf("ServerPath = %q, want %q", got.LLM.ServerPath, cfg.LLM.ServerPath)
	}
	if !got.UI.WebEnabled {
		t.Fatalf("WebEnabled = false, want true")
	}
}
