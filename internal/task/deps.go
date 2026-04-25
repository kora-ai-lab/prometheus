package task

import (
	"github.com/prometheus-dev/prometheus/internal/browser"
	"github.com/prometheus-dev/prometheus/internal/capabilities"
	"github.com/prometheus-dev/prometheus/internal/executor"
	"github.com/prometheus-dev/prometheus/internal/llm"
	"github.com/prometheus-dev/prometheus/internal/logging"
	"github.com/prometheus-dev/prometheus/internal/prompt"
	"github.com/prometheus-dev/prometheus/internal/security"
	"github.com/prometheus-dev/prometheus/internal/vision"
)

type TaskStore interface {
	Save(*Task) error
}

type TaskDeps struct {
	Provider      llm.ModelProvider
	Executor      executor.Executor
	Vision        vision.VisionProvider
	Browser       *browser.Manager
	PromptBuilder *prompt.Builder
	CapEngine     *capabilities.Engine
	Security      *security.Interceptor
	Logger        *logging.Logger
	TaskStore     TaskStore
}
