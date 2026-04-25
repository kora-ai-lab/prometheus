package prompt

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus-dev/prometheus/internal/config"
)

//go:embed system_default.md
var embeddedSystemPrompt string

func LoadSystemPrompt() string {
	userPrompt := filepath.Join(config.PrometheusHome(), "prompts", "system_v1.md")
	if data, err := os.ReadFile(userPrompt); err == nil && isCompatiblePrompt(data) {
		return string(data)
	}
	if data, err := os.ReadFile(filepath.Join("assets", "prompts", "system_v1.md")); err == nil && isCompatiblePrompt(data) {
		return string(data)
	}
	return embeddedSystemPrompt
}

func isCompatiblePrompt(data []byte) bool {
	text := string(data)
	return strings.Contains(text, "prometheus_prompt_version:") &&
		strings.Contains(text, `"action": "exec|ask|browser|vision|create|done|error"`)
}
