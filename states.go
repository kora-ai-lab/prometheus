package main

import "fmt"

type AgentState int

const (
	StateThinking AgentState = iota // 0
	StateExecuting                 // 1
	StateObserving                  // 2
	StateBlocked                   // 3
	StateCompleted                 // 4
	StateError                    // 5
)

func (s AgentState) String() string {
	switch s {
	case StateThinking:
		return "thinking"
	case StateExecuting:
		return "executing"
	case StateObserving:
		return "observing"
	case StateBlocked:
		return "blocked"
	case StateCompleted:
		return "completed"
	case StateError:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}