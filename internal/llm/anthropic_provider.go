package llm

import (
	"context"
	"errors"
)

type AnthropicProvider struct {
	info *ModelInfo
}

func NewAnthropicProvider(model string) *AnthropicProvider {
	if model == "" {
		model = "claude-sonnet"
	}
	return &AnthropicProvider{
		info: &ModelInfo{Name: model, ContextWindow: 200000, Provider: "anthropic", HasVision: true},
	}
}

func (p *AnthropicProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
	_ = ctx
	_ = messages
	return nil, errors.New("anthropic provider scaffold is not wired yet")
}

func (p *AnthropicProvider) Stream(ctx context.Context, messages []Message, tokens chan<- string) error {
	defer close(tokens)
	return errors.New("anthropic provider scaffold is not wired yet")
}

func (p *AnthropicProvider) ModelInfo() *ModelInfo { return p.info }
func (p *AnthropicProvider) IsAvailable() bool     { return false }
func (p *AnthropicProvider) HasVision() bool       { return true }
func (p *AnthropicProvider) Close() error          { return nil }
