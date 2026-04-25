package vision

import (
	"context"
	"testing"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

type mockProvider struct {
	output string
	err    error
}

func (m *mockProvider) Complete(ctx context.Context, messages []llm.Message) (*llm.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &llm.Response{Content: m.output}, nil
}

func (m *mockProvider) Stream(ctx context.Context, messages []llm.Message, tokens chan<- string) error {
	return m.err
}

func (m *mockProvider) ModelInfo() *llm.ModelInfo {
	return &llm.ModelInfo{HasVision: true}
}

func (m *mockProvider) IsAvailable() bool {
	return m.err == nil
}

func (m *mockProvider) HasVision() bool {
	return true
}

func (m *mockProvider) Close() error {
	return nil
}

func TestNewComparator(t *testing.T) {
	provider := &mockProvider{}
	c := NewComparator(provider)
	if c == nil {
		t.Fatal("NewComparator returned nil")
	}
	if c.provider != provider {
		t.Error(" Comparator did not store provider")
	}
}

func TestComparator_Diff(t *testing.T) {
	provider := &mockProvider{}
	c := NewComparator(provider)

	result, err := c.Diff("testdata/mockup.png", "testdata/result.png")
	if err != nil {
		t.Errorf("Diff failed: %v", err)
	}
	if result == "" {
		t.Error("Diff returned empty result")
	}
}

func TestComparator_CompareWithLLM(t *testing.T) {
	provider := &mockProvider{output: "Images are visually similar with minor differences"}
	c := NewComparator(provider)

	result, err := c.CompareWithLLM(context.Background(), "testdata/mockup.png", "testdata/result.png")
	if err != nil {
		t.Errorf("CompareWithLLM failed: %v", err)
	}
	if result == "" {
		t.Error("CompareWithLLM returned empty result")
	}
}