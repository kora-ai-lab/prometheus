package context

import (
	stdcontext "context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus-dev/prometheus/internal/llm"
	"github.com/prometheus-dev/prometheus/internal/prompt"
)

type Manager struct {
	hotBuffer     []llm.Message
	warmSummary   string
	contextWindow int
	threshold     float64
	keepLast      int
	provider      llm.ModelProvider
}

func New(provider llm.ModelProvider) *Manager {
	m := &Manager{
		contextWindow: 4096,
		threshold:     0.60,
		keepLast:      5,
		provider:      provider,
	}
	if provider != nil && provider.ModelInfo() != nil && provider.ModelInfo().ContextWindow > 0 {
		m.contextWindow = provider.ModelInfo().ContextWindow
	}
	m.configure()
	return m
}

func (m *Manager) configure() {
	switch {
	case m.contextWindow < 4000:
		m.threshold = 0.60
		m.keepLast = 5
	case m.contextWindow < 8000:
		m.threshold = 0.60
		m.keepLast = 5
	case m.contextWindow < 32000:
		m.threshold = 0.65
		m.keepLast = 10
	default:
		m.threshold = 0.70
		m.keepLast = 20
	}
}

func (m *Manager) Add(msg llm.Message) {
	m.hotBuffer = append(m.hotBuffer, msg)
	if m.usageRatio() > m.threshold {
		m.compact()
	}
}

func (m *Manager) BuildMessages(systemPrompt string) []llm.Message {
	msgs := []llm.Message{{Role: "system", Content: systemPrompt}}
	if m.warmSummary != "" {
		msgs = append(msgs, llm.Message{
			Role:    "system",
			Content: "[PREVIOUS_CONTEXT]\n" + m.warmSummary,
		})
	}
	return append(msgs, m.hotBuffer...)
}

func (m *Manager) usageRatio() float64 {
	total := prompt.EstimateTokens(m.warmSummary)
	for _, msg := range m.hotBuffer {
		total += prompt.EstimateTokens(msg.Content)
	}
	return float64(total) / float64(m.contextWindow)
}

func (m *Manager) compact() {
	if m.provider == nil || len(m.hotBuffer) <= m.keepLast {
		return
	}
	toCompact := m.hotBuffer[:len(m.hotBuffer)-m.keepLast]
	var sb strings.Builder
	for _, msg := range toCompact {
		sb.WriteString(msg.Role)
		sb.WriteString(": ")
		sb.WriteString(msg.Content)
		sb.WriteByte('\n')
	}

	resp, err := m.provider.Complete(stdcontext.Background(), []llm.Message{
		{
			Role: "user",
			Content: fmt.Sprintf(
				"Summarize in compact JSON with goal, done, decisions, state, next. Max 250 tokens.\n%s",
				sb.String(),
			),
		},
	})
	if err != nil {
		return
	}

	m.warmSummary = resp.Content
	m.hotBuffer = append([]llm.Message{}, m.hotBuffer[len(m.hotBuffer)-m.keepLast:]...)
}

func (m *Manager) CompactWithTimeout(timeout time.Duration) {
	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), timeout)
	defer cancel()
	_ = ctx
	m.compact()
}

func (m *Manager) Snapshot() map[string]any {
	return map[string]any{
		"warm_summary": m.warmSummary,
		"hot_buffer":   m.hotBuffer,
		"context_used": int(m.usageRatio() * float64(m.contextWindow)),
	}
}

func (m *Manager) Restore(data map[string]any) error {
	if data == nil {
		return nil
	}
	if ws, ok := data["warm_summary"].(string); ok {
		m.warmSummary = ws
	}
	if hb, ok := data["hot_buffer"].(string); ok {
		var msgs []llm.Message
		if err := json.Unmarshal([]byte(hb), &msgs); err == nil {
			m.hotBuffer = msgs
		}
	}
	return nil
}
