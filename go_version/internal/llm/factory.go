package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/kora-ai-lab/prometheus/internal/config"
)

type ErrNoLLMAvailable struct {
	Hint string
}

func (e ErrNoLLMAvailable) Error() string {
	return "no llm available: " + e.Hint
}

func AutoDetect(cfg *config.LLMConfig, serverPath string) (ModelProvider, error) {
	if cfg == nil {
		return nil, ErrNoLLMAvailable{Hint: "missing llm config"}
	}

	// Cloud providers (fastest first)
	if os.Getenv("GROQ_API_KEY") != "" {
		return NewGroqProvider(os.Getenv("GROQ_API_KEY"), cfg.ModelName), nil
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return NewAnthropicProvider(cfg.ModelName), nil
	}
	if os.Getenv("GOOGLE_API_KEY") != "" {
		return NewGoogleProvider(cfg.ModelName), nil
	}

	// Try Ollama (local, slower but offline capable)
	ollama := NewOllamaProvider(cfg.Endpoint, cfg.ModelName)
	if ollama.IsAvailable() {
		return ollama, nil
	}

	// Then try local GGUF (offline fallback, slower on CPU)
	resolvedServerPath := resolveServerPath(cfg, serverPath)
	if cfg.Provider == "local" && cfg.ModelPath != "" {
		if _, err := os.Stat(cfg.ModelPath); err == nil {
			if provider, err := NewLocalLlamaProvider(resolvedServerPath, cfg.ModelPath, cfg.VisionModelPath); err == nil {
				return provider, nil
			}
		}
	}

	return nil, ErrNoLLMAvailable{Hint: "start Ollama or configure a model"}
}

func postJSON(ctx context.Context, client *http.Client, url string, payload any) (io.ReadCloser, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(data))
	}
	return resp.Body, nil
}

func decodeJSON(r io.Reader, out any) error {
	return json.NewDecoder(r).Decode(out)
}

func resolveServerPath(cfg *config.LLMConfig, embeddedPath string) string {
	if cfg != nil && cfg.ServerPath != "" {
		return cfg.ServerPath
	}
	if v := os.Getenv("PROMETHEUS_LLM_SERVER_PATH"); v != "" {
		return v
	}
	return embeddedPath
}
