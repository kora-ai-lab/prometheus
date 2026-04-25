package llm

import "context"

type Message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  [][]byte `json:"images,omitempty"`
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
