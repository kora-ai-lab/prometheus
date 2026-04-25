package task

import (
	"github.com/kora-ai-lab/prometheus/internal/browser"
	"github.com/kora-ai-lab/prometheus/internal/capabilities"
	"github.com/kora-ai-lab/prometheus/internal/executor"
	"github.com/kora-ai-lab/prometheus/internal/llm"
	"github.com/kora-ai-lab/prometheus/internal/logging"
	"github.com/kora-ai-lab/prometheus/internal/prompt"
	"github.com/kora-ai-lab/prometheus/internal/security"
	"github.com/kora-ai-lab/prometheus/internal/vision"
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
