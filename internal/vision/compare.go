package vision

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

type Comparator struct {
	provider llm.ModelProvider
}

func NewComparator(provider llm.ModelProvider) *Comparator {
	return &Comparator{provider: provider}
}

func (c *Comparator) Diff(img1, img2 string) (string, error) {
	abs1, err := filepath.Abs(img1)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for img1: %w", err)
	}
	abs2, err := filepath.Abs(img2)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for img2: %w", err)
	}

	data, err := os.ReadFile(abs1)
	if err != nil {
		return "", fmt.Errorf("failed to read img1: %w", err)
	}
	data2, err := os.ReadFile(abs2)
	if err != nil {
		return "", fmt.Errorf("failed to read img2: %w", err)
	}

	cmd := exec.Command("magick", "compare", "-metric", "AE", abs1, abs2, os.DevNull)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		_ = cmd
	}

	diffCount := 0.0
	minLen := len(data)
	if len(data2) < minLen {
		minLen = len(data2)
	}
	for i := 0; i < minLen; i++ {
		diffCount += math.Abs(float64(data[i]) - float64(data2[i]))
	}

	diffCount += math.Abs(float64(len(data) - len(data2)))

	return fmt.Sprintf("pixel diff: %d different bytes", int(diffCount)), nil
}

func (c *Comparator) CompareWithLLM(ctx context.Context, mockupPath, resultPath string) (string, error) {
	if c.provider == nil {
		return "", fmt.Errorf(" no provider configured")
	}

	mockupData, err := os.ReadFile(mockupPath)
	if err != nil {
		return "", fmt.Errorf("failed to read mockup: %w", err)
	}

	resultData, err := os.ReadFile(resultPath)
	if err != nil {
		return "", fmt.Errorf("failed to read result: %w", err)
	}

	messages := []llm.Message{
		{
			Role: "user",
			Content: `Compare these two images and describe the visual differences. 
The first image is the mockup/design reference and the second is the rendered result.
Please analyze layout, colors, spacing, and any missing or incorrect elements.`,
			Images: [][]byte{mockupData, resultData},
		},
	}

	resp, err := c.provider.Complete(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("LLM analysis failed: %w", err)
	}

	return resp.Content, nil
}