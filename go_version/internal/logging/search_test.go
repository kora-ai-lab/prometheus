package logging

import (
	"context"
	"testing"
)

func TestTokenize(t *testing.T) {
	engine := NewSearchEngine("", nil)

	tokens := engine.tokenize("Hello World 123")
	if len(tokens) != 3 {
		t.Errorf("Expected 3 tokens, got %d", len(tokens))
	}
}

func TestSearchEngine_SearchBM25(t *testing.T) {
	engine := NewSearchEngine("", nil)
	engine.docStore = map[docID]LogEntry{
		0: {Level: "task_start", TaskID: "task1", Event: map[string]any{"goal": "test"}},
		1: {Level: "task_end", TaskID: "task2", Event: map[string]any{"status": "done"}},
	}
	engine.index = map[string][]docID{
		"task_start": {0},
		"task_end":   {1},
		"task":       {0, 1},
	}

	results, err := engine.Search(context.Background(), "task start", 5)
	if err != nil {
		t.Errorf("Search failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected results, got none")
	}
}

func TestCosineSimilarity(t *testing.T) {
	a := []float64{1, 0, 0}
	b := []float64{1, 0, 0}

	score := cosineSimilarity(a, b)
	if score != 1.0 {
		t.Errorf("Expected 1.0, got %f", score)
	}

	a = []float64{1, 0, 0}
	b = []float64{0, 1, 0}

	score = cosineSimilarity(a, b)
	if score != 0 {
		t.Errorf("Expected 0, got %f", score)
	}
}