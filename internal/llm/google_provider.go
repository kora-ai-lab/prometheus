package llm

import (
	"context"
	"errors"
)

type GoogleProvider struct {
	info *ModelInfo
}

func NewGoogleProvider(model string) *GoogleProvider {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GoogleProvider{
		info: &ModelInfo{Name: model, ContextWindow: 1_000_000, Provider: "google", HasVision: true},
	}
}

func (p *GoogleProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
	_ = ctx
	_ = messages
	return nil, errors.New("google provider scaffold is not wired yet")
}

func (p *GoogleProvider) Stream(ctx context.Context, messages []Message, tokens chan<- string) error {
	defer close(tokens)
	return errors.New("google provider scaffold is not wired yet")
}

func (p *GoogleProvider) ModelInfo() *ModelInfo { return p.info }
func (p *GoogleProvider) IsAvailable() bool     { return false }
func (p *GoogleProvider) HasVision() bool       { return true }
func (p *GoogleProvider) Close() error          { return nil }
