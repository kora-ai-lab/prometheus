package tracing

import (
	"crypto/rand"
	"encoding/hex"
)

type TraceID [16]byte

type SpanID [8]byte

func GenerateTraceID() TraceID {
	var id TraceID
	rand.Read(id[:])
	return id
}

func GenerateSpanID() SpanID {
	var id SpanID
	rand.Read(id[:])
	return id
}

func (t TraceID) String() string {
	return hex.EncodeToString(t[:])
}

func (s SpanID) String() string {
	return hex.EncodeToString(s[:])
}