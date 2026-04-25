package tracing

import (
	"sync"
	"time"
)

type Span struct {
	TraceID    TraceID
	SpanID     SpanID
	ParentID   *SpanID
	Operation string
	StartTime  time.Time
	EndTime    time.Time
	Attributes map[string]interface{}
	mu         sync.RWMutex
}

func StartSpan(operation string, attrs ...interface{}) *Span {
	ctx := CurrentContext()

	var parentID *SpanID
	if ctx != nil {
		parentID = &ctx.SpanID
	}

	span := &Span{
		TraceID:    GenerateTraceID(),
		SpanID:     GenerateSpanID(),
		ParentID:   parentID,
		Operation: operation,
		StartTime:  time.Now(),
		Attributes: make(map[string]interface{}),
	}

	for i := 0; i < len(attrs); i += 2 {
		if i+1 < len(attrs) {
			span.Attributes[attrs[i].(string)] = attrs[i+1]
		}
	}

	return span
}

func (s *Span) End() {
	s.mu.Lock()
	s.EndTime = time.Now()
	s.mu.Unlock()
}

func (s *Span) SetAttribute(key string, value interface{}) {
	s.mu.Lock()
	s.Attributes[key] = value
	s.mu.Unlock()
}

func (s *Span) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}