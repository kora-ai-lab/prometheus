package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type LocalLlamaProvider struct {
	serverPath string
	modelPath  string
	baseURL    string
	cmd        *exec.Cmd
	client     *http.Client
	info       *ModelInfo
	mu         sync.Mutex
}

func NewLocalLlamaProvider(serverPath, modelPath, visionPath string) (*LocalLlamaProvider, error) {
	if serverPath == "" {
		return nil, errors.New("embedded llama-server is not available")
	}
	if _, err := os.Stat(serverPath); err != nil {
		return nil, fmt.Errorf("llama-server binary unavailable: %w", err)
	}
	if _, err := os.Stat(modelPath); err != nil {
		return nil, fmt.Errorf("gguf model unavailable: %w", err)
	}

	p := &LocalLlamaProvider{
		serverPath: serverPath,
		modelPath:  modelPath,
		client:     &http.Client{Timeout: 120 * time.Second},
		info: &ModelInfo{
			Name:          filepathBase(modelPath),
			ContextWindow: 4096,
			Provider:      "local",
			HasVision:     visionPath != "",
		},
	}
	if err := p.start(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *LocalLlamaProvider) start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil {
		return nil
	}

	port, err := freePort()
	if err != nil {
		return err
	}
	p.baseURL = "http://127.0.0.1:" + strconv.Itoa(port)

	args := []string{
		"--model", p.modelPath,
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(port),
		"--ctx-size", strconv.Itoa(p.info.ContextWindow),
		"--no-webui",
	}
	cmd := exec.Command(p.serverPath, args...)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return err
	}
	p.cmd = cmd

	if err := p.waitReady(30 * time.Second); err != nil {
		if p.cmd != nil && p.cmd.Process != nil {
			p.cmd.Process.Kill()
		}
		p.cmd = nil
		return err
	}
	if info, err := p.fetchModelInfo(); err == nil && info != nil {
		p.info = info
	}
	return nil
}

func (p *LocalLlamaProvider) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := p.client.Get(p.baseURL + "/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("llama-server not ready after %s", timeout)
}

type openAIChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

type openAIModelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Details struct {
			ContextLength int `json:"context_length"`
		} `json:"details"`
	} `json:"data"`
}

func (p *LocalLlamaProvider) fetchModelInfo() (*ModelInfo, error) {
	resp, err := p.client.Get(p.baseURL + "/v1/models")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed openAIModelsResponse
	if err := decodeJSON(resp.Body, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Data) == 0 {
		return p.info, nil
	}
	return &ModelInfo{
		Name:          parsed.Data[0].ID,
		ContextWindow: maxInt(parsed.Data[0].Details.ContextLength, p.info.ContextWindow),
		Provider:      "local",
		HasVision:     p.info.HasVision,
	}, nil
}

func (p *LocalLlamaProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
	body, err := postJSON(ctx, p.client, p.baseURL+"/v1/chat/completions", openAIChatRequest{
		Model:    p.info.Name,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var parsed openAIChatResponse
	if err := decodeJSON(body, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, errors.New("local llama response had no choices")
	}
	return &Response{
		Content:      parsed.Choices[0].Message.Content,
		InputTokens:  parsed.Usage.PromptTokens,
		OutputTokens: parsed.Usage.CompletionTokens,
	}, nil
}

func (p *LocalLlamaProvider) Stream(ctx context.Context, messages []Message, tokens chan<- string) error {
	defer close(tokens)
	resp, err := p.Complete(ctx, messages)
	if err != nil {
		return err
	}
	tokens <- resp.Content
	return nil
}

func (p *LocalLlamaProvider) ModelInfo() *ModelInfo { return p.info }
func (p *LocalLlamaProvider) IsAvailable() bool     { return p != nil && p.cmd != nil }
func (p *LocalLlamaProvider) HasVision() bool       { return p.info != nil && p.info.HasVision }

func (p *LocalLlamaProvider) Close() error {
	if p == nil || p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	_ = p.cmd.Process.Kill()
	_, _ = p.cmd.Process.Wait()
	p.cmd = nil
	return nil
}

type OllamaProvider struct {
	endpoint string
	model    string
	client   *http.Client
	info     *ModelInfo
}

func NewOllamaProvider(endpoint, model string) *OllamaProvider {
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11434"
	}
	if model == "" {
		model = "phi3:mini"
	}
	return &OllamaProvider{
		endpoint: endpoint,
		model:    model,
		client:   &http.Client{Timeout: 300 * time.Second},
		info: &ModelInfo{
			Name:          model,
			ContextWindow: 8192,
			Provider:      "ollama",
		},
	}
}

type ollamaRequest struct {
	Model           string    `json:"model"`
	Messages        []Message `json:"messages"`
	Stream          bool      `json:"stream"`
	ResponseFormat  any       `json:"response_format,omitempty"`
	Temperature     float64   `json:"temperature,omitempty"`
}

type ollamaChatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

func (p *OllamaProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
	body, err := postJSON(ctx, p.client, p.endpoint+"/api/chat", ollamaRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var parsed ollamaChatResponse
	if err := decodeJSON(body, &parsed); err != nil {
		return nil, err
	}
	return &Response{Content: parsed.Message.Content}, nil
}

func (p *OllamaProvider) Stream(ctx context.Context, messages []Message, tokens chan<- string) error {
	defer close(tokens)
	resp, err := p.Complete(ctx, messages)
	if err != nil {
		return err
	}
	tokens <- resp.Content
	return nil
}

func (p *OllamaProvider) ModelInfo() *ModelInfo { return p.info }
func (p *OllamaProvider) IsAvailable() bool {
	resp, err := p.client.Get(p.endpoint + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode < 500
}
func (p *OllamaProvider) HasVision() bool { return false }
func (p *OllamaProvider) Close() error    { return nil }

func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()
	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("unexpected listener address")
	}
	return addr.Port, nil
}

func prometheusActionSchema_() map[string]any {
	return map[string]any{
		"name": "PrometheusAction",
		"strict": true,
		"schema": map[string]any{
			"type": "object",
			"required": []string{"thinking", "action", "why"},
			"properties": map[string]any{
				"thinking":   map[string]any{"type": "string"},
				"action":     map[string]any{"type": "string", "enum": []string{"exec", "ask", "browser", "vision", "create", "done", "error"}},
				"command":    map[string]any{"type": "string"},
				"dangerous":  map[string]any{"type": "boolean"},
				"why":        map[string]any{"type": "string"},
				"question":   map[string]any{"type": "string"},
				"create_file": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path":    map[string]any{"type": "string"},
						"content": map[string]any{"type": "string"},
					},
					"required": []string{"path", "content"},
				},
			},
		},
	}
}

func filepathBase(path string) string {
	if path == "" {
		return "local-llama"
	}
	return filepath.Base(path)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
