package tracing

import "sync"

type TraceContext struct {
	TraceID  TraceID
	SpanID   SpanID
	ParentID *SpanID
}

var currentContext struct {
	mu   sync.RWMutex
	ctx *TraceContext
}

func WithContext(ctx *TraceContext) func() {
	currentContext.mu.Lock()
	old := currentContext.ctx
	currentContext.ctx = ctx
	currentContext.mu.Unlock()

	return func() {
		currentContext.mu.Lock()
		currentContext.ctx = old
		currentContext.mu.Unlock()
	}
}

func CurrentContext() *TraceContext {
	currentContext.mu.RLock()
	defer currentContext.mu.RUnlock()
	return currentContext.ctx
}