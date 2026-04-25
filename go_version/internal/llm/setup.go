package llm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kora-ai-lab/prometheus/internal/config"
	"github.com/kora-ai-lab/prometheus/internal/discovery"
)

var ErrSetupRequired = errors.New("local model not configured")

func FirstRunSetup(home string, env *discovery.EnvironmentProfile, ui io.Writer) error {
	model := SelectModel(env.RAMMb)
	modelsDir := filepath.Join(home, "models")

	existing, err := DetectLocalModels(modelsDir)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		if err := config.UpdateLLM(home, func(cfg *config.LLMConfig) {
			cfg.Provider = "local"
			cfg.ModelPath = existing[0]
		}); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(ui, "PROMETHEUS INITIAL SETUP\nUsing existing local model: %s\n", existing[0])
		return nil
	}

	expectedPath := filepath.Join(modelsDir, model.Filename)
	_, _ = fmt.Fprintf(
		ui,
		"PROMETHEUS INITIAL SETUP\nRecommended model: %s\nExpected path: %s\nPlace a GGUF there, set PROMETHEUS_LLM_MODEL_PATH, or configure Ollama.\n",
		model.Name,
		expectedPath,
	)
	if err := config.UpdateLLM(home, func(cfg *config.LLMConfig) {
		cfg.Provider = "local"
		cfg.ModelPath = expectedPath
	}); err != nil {
		return err
	}
	return ErrSetupRequired
}

func DetectLocalModels(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var models []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".gguf") {
			models = append(models, filepath.Join(dir, entry.Name()))
		}
	}
	sort.Strings(models)
	return models, nil
}
