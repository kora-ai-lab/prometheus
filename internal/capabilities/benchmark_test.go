package capabilities

import (
	"context"
	"testing"
)

func BenchmarkEngine_Ensure(b *testing.B) {
	engine := NewEngine(nil, nil, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Ensure(context.Background(), "python3")
	}
}

func BenchmarkRegistry_Lookup(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Registry["python3"]
		_ = Registry["node"]
		_ = Registry["git"]
		_ = Registry["docker"]
	}
}

func BenchmarkRegistry_Lookup_Miss(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Registry["nonexistent"]
	}
}
