package vision

import (
	"context"
	"errors"
)

type LocalVisionProvider struct{}

func NewLocalProvider() *LocalVisionProvider {
	return &LocalVisionProvider{}
}

func (p *LocalVisionProvider) Analyze(ctx context.Context, imageBytes []byte, question string) (string, error) {
	_, _, _ = ctx, imageBytes, question
	return "", errors.New("local vision provider is scaffolded but not implemented")
}

func (p *LocalVisionProvider) HasVision() bool {
	return false
}
