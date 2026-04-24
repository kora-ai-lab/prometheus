package llm

import "context"

type Message struct {
	Role    string
	Content string
	Images  [][]byte
}

type Response struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

type ModelInfo struct {
	Name          string
	ContextWindow int
	Provider      string
	HasVision     bool
}

type ModelProvider interface {
	Complete(ctx context.Context, messages []Message) (*Response, error)
	Stream(ctx context.Context, messages []Message, tokens chan<- string) error
	ModelInfo() *ModelInfo
	IsAvailable() bool
	HasVision() bool
	Close() error
}
