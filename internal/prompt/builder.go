package prompt

import (
	"fmt"
	"strings"

	"github.com/prometheus-dev/prometheus/internal/discovery"
	"github.com/prometheus-dev/prometheus/internal/llm"
)

type Builder struct {
	systemPromptBase string
	contextWindow    int
	historyRatio     float64
	env              *discovery.EnvironmentProfile
	capabilities     []string
	patterns         []string
}

func NewBuilder(info *llm.ModelInfo, env *discovery.EnvironmentProfile, capabilities []string) *Builder {
	contextWindow := 4096
	if info != nil && info.ContextWindow > 0 {
		contextWindow = info.ContextWindow
	}
	return &Builder{
		systemPromptBase: LoadSystemPrompt(),
		contextWindow:    contextWindow,
		historyRatio:     0.60,
		env:              env,
		capabilities:     capabilities,
	}
}

func (b *Builder) Build() string {
	budget := int(float64(b.contextWindow) * (1.0 - b.historyRatio))
	out := b.systemPromptBase
	if budget <= EstimateTokens(out) {
		return out
	}

	if envBlock := b.buildEnvBlock(); envBlock != "" {
		out += "\n\n" + envBlock
	}
	if len(b.capabilities) > 0 {
		out += "\n\nCAPABILITIES: " + strings.Join(b.capabilities, ", ")
	}
	if len(b.patterns) > 0 {
		out += "\n\nPATTERNS: " + strings.Join(b.patterns, " | ")
	}
	return out
}

func (b *Builder) BuildMessages(history []llm.Message) []llm.Message {
	msgs := []llm.Message{{Role: "system", Content: b.Build()}}
	return append(msgs, history...)
}

func (b *Builder) buildEnvBlock() string {
	if b.env == nil {
		return ""
	}

	tools := b.env.AvailableTools
	if len(tools) > 8 {
		tools = tools[:8]
	}

	return fmt.Sprintf(
		"ENV:%s/%s RAM:%dMB CPU:%d NET:%t PM:%s TOOLS:%s",
		b.env.OS,
		b.env.Arch,
		b.env.RAMMb,
		b.env.CPUCores,
		b.env.Internet,
		b.env.PackageManager,
		strings.Join(tools, ","),
	)
}
