package tracing

import (
	"testing"
)

func TestSpan_Basic(t *testing.T) {
	span := StartSpan("test.operation", "key", "value")
	defer span.End()

	span.SetAttribute("extra", "attr")

	if span.Operation != "test.operation" {
		t.Errorf("Expected operation, got %s", span.Operation)
	}

	if span.Attributes["key"] != "value" {
		t.Error("Expected key=value in attributes")
	}

	dur := span.Duration()
	if dur < 0 {
		t.Error("Duration should be positive")
	}
}