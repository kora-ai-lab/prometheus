package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type GroqProvider struct {
	apiKey   string
	model    string
	endpoint string
	client  *http.Client
	info    *ModelInfo
}

func NewGroqProvider(apiKey, model string) *GroqProvider {
	endpoint := "https://api.groq.com/openai/v1"
	if model == "" {
		model = "llama-3.3-70b-versatile"
	}
	return &GroqProvider{
		apiKey:   apiKey,
		model:    model,
		endpoint: endpoint,
		client:  &http.Client{Timeout: 60 * time.Second},
		info: &ModelInfo{
			Name:          model,
			ContextWindow: 8192,
			Provider:      "groq",
		},
	}
}

type openAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

func (p *GroqProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
	reqBody := openAIRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	}

	body, err := p.postWithAuth(ctx, p.endpoint+"/chat/completions", reqBody)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var parsed groqResponse
	if err := decodeJSON(body, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Choices) == 0 {
		return &Response{Content: ""}, nil
	}
	return &Response{Content: parsed.Choices[0].Message.Content}, nil
}

func (p *GroqProvider) postWithAuth(ctx context.Context, url string, payload any) (io.ReadCloser, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(data))
	}
	return resp.Body, nil
}

func (p *GroqProvider) Stream(ctx context.Context, messages []Message, tokens chan<- string) error {
	defer close(tokens)
	resp, err := p.Complete(ctx, messages)
	if err != nil {
		return err
	}
	tokens <- resp.Content
	return nil
}

func (p *GroqProvider) ModelInfo() *ModelInfo { return p.info }
func (p *GroqProvider) IsAvailable() bool {
	return p.apiKey != ""
}
func (p *GroqProvider) HasVision() bool { return false }
func (p *GroqProvider) Close() error    { return nil }