package vision

import "context"

type NoOpVisionProvider struct{}

func (n *NoOpVisionProvider) Analyze(_ context.Context, _ []byte, _ string) (string, error) {
	return "[vision unavailable - continuing without analysis]", nil
}

func (n *NoOpVisionProvider) HasVision() bool {
	return false
}
