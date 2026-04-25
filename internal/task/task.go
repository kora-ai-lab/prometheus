package task

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/prometheus-dev/prometheus/internal/llm"
)

type TaskStatus string

const (
	StatusRunning   TaskStatus = "running"
	StatusBlocked   TaskStatus = "blocked"
	StatusDone      TaskStatus = "done"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

type Task struct {
	ID             string
	Goal           string
	Status         TaskStatus
	Context        []llm.Message
	Memory         map[string]any
	BlockedReason  string
	Retries        int
	MaxRetries     int
	ParseErrors    int
	MaxParseErrors int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func New(goal string) *Task {
	now := time.Now()
	return &Task{
		ID:             newID(),
		Goal:           goal,
		Status:         StatusRunning,
		Context:        []llm.Message{{Role: "user", Content: goal}},
		Memory:         map[string]any{},
		MaxRetries:     5,
		MaxParseErrors: 3,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}

func (t *Task) Resume(answer string) {
	t.Context = append(t.Context, llm.Message{
		Role:    "user",
		Content: "Réponse à ta question: " + answer,
	})
	t.Status = StatusRunning
	t.BlockedReason = ""
	t.UpdatedAt = time.Now()
}
