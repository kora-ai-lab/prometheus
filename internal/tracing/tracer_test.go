package tracing

import (
	"testing"
)

func TestGenerateTraceID(t *testing.T) {
	id1 := GenerateTraceID()
	id2 := GenerateTraceID()

	if id1.String() == id2.String() {
		t.Error("Expected unique trace IDs")
	}

	if len(id1.String()) != 32 {
		t.Errorf("Expected 32 hex chars, got %d", len(id1.String()))
	}
}

func TestGenerateSpanID(t *testing.T) {
	id1 := GenerateSpanID()
	id2 := GenerateSpanID()

	if id1.String() == id2.String() {
		t.Error("Expected unique span IDs")
	}

	if len(id1.String()) != 16 {
		t.Errorf("Expected 16 hex chars, got %d", len(id1.String()))
	}
}