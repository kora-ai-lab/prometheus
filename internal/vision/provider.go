package vision

import "context"

type VisionProvider interface {
	Analyze(ctx context.Context, imageBytes []byte, question string) (string, error)
	HasVision() bool
}
