package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

func main() {
	fmt.Println("Testing Ollama...")
	provider := llm.NewOllamaProvider("http://127.0.0.1:11434", "phi3:mini")
	
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	
	resp, err := provider.Complete(ctx, []llm.Message{{Role: "user", Content: "Say hi in 3 words"}})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Println("Response:", resp.Content)
}