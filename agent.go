package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	LLM          ModelProvider
	DB           *DB
	CurrentGoal  string
	State        AgentState
	CurrentTaskID string
	StepCount    int
	MaxRetries   int
}

func NewAgent(llm ModelProvider, db *DB, goal string) *Agent {
	taskID := uuid.New().String()
	return &Agent{
		LLM:         llm,
		DB:          db,
		CurrentGoal: goal,
		State:       StateThinking,
		CurrentTaskID: taskID,
		StepCount:   0,
		MaxRetries:  3,
	}
}

func (a *Agent) Step() (AgentState, bool) {
	a.StepCount++

	switch a.State {
	case StateThinking:
		return a.think()

	case StateExecuting:
		return a.execute()

	case StateObserving:
		return a.observe()

	case StateBlocked:
		return a.State, false

	case StateCompleted, StateError:
		return a.State, true
	}

	return a.State, false
}

func (a *Agent) think() (AgentState, bool) {
	dbStatus := TaskStatusRunning
	if err := a.DB.UpdateTaskStatus(a.CurrentTaskID, dbStatus); err != nil {
		fmt.Printf("DB update error: %v\n", err)
	}

	prompt := fmt.Sprintf(
		"Goal: %s\n\nThink about what command to execute. Respond with COMMAND: <yourcommand> or DONE if complete.",
		a.CurrentGoal,
	)

	response, err := a.LLM.Generate(prompt)
	if err != nil {
		a.handleError(err)
		return StateError, true
	}

	a.DB.LogExecution(a.CurrentTaskID, a.StepCount, "think", response, "", 0)

	if response == "DONE" {
		a.State = StateCompleted
		a.DB.CompleteTask(a.CurrentTaskID)
		return a.State, true
	}

	a.State = StateExecuting
	return a.State, false
}

func (a *Agent) execute() (AgentState, bool) {
	prompt := fmt.Sprintf(
		"Goal: %s\n\nWhat is the command to execute? Just respond with the command, nothing else.",
		a.CurrentGoal,
	)

	response, err := a.LLM.Generate(prompt)
	if err != nil {
		a.handleError(err)
		return StateError, true
	}

	result, err := ExecuteCommand(response)
	if err != nil {
		a.State = StateBlocked
		a.DB.BlockTask(a.CurrentTaskID, err.Error())
		return a.State, false
	}

	a.DB.LogExecution(a.CurrentTaskID, a.StepCount, response, result.Stdout, result.Stderr, result.ExitCode)

	a.State = StateObserving
	return a.State, false
}

func (a *Agent) observe() (AgentState, bool) {
	prompt := fmt.Sprintf(
		"Previous command output:\nstdout: %s\nstderr: %s\nexit code: %d\n\nIs the goal complete? Respond DONE or continue.",
		"", "", 0,
	)

	response, err := a.LLM.Generate(prompt)
	if err != nil {
		a.handleError(err)
		return StateError, true
	}

	if response == "DONE" {
		a.State = StateCompleted
		a.DB.CompleteTask(a.CurrentTaskID)
		return a.State, true
	}

	a.State = StateThinking
	return a.State, false
}

func (a *Agent) handleError(err error) {
	a.DB.IncrementRetry(a.CurrentTaskID)

	if a.StepCount >= a.MaxRetries {
		a.State = StateError
		a.DB.FailTask(a.CurrentTaskID, err.Error())
	} else {
		a.State = StateBlocked
		a.DB.BlockTask(a.CurrentTaskID, err.Error())
	}
}

func (a *Agent) ResolveBlock(answer string) {
	a.DB.UnblockTask(a.CurrentTaskID, answer)
	a.State = StateThinking
}

func (a *Agent) GetTaskInfo() (id string, goal string, status string, step int) {
	return a.CurrentTaskID, a.CurrentGoal, a.State.String(), a.StepCount
}

func (a *Agent) String() string {
	return fmt.Sprintf(
		"Agent{id=%s, goal=%s, state=%s, step=%d}",
		a.CurrentTaskID[:8], a.CurrentGoal[:20], a.State.String(), a.StepCount,
	)
}

func TestModel(name string) bool {
	fmt.Printf("\n--- Testing Model: %s ---\n", name)

	folderName := fmt.Sprintf("test_%s", name)
	os.RemoveAll(folderName)

	llm := NewOllamaProvider("http://localhost:11434", name)
	dbPath := fmt.Sprintf("test_%s.db", name)
	db, err := NewDB(dbPath)
	if err != nil {
		fmt.Printf("DB error: %v\n", err)
		return false
	}
	defer os.Remove(dbPath)

	goal := fmt.Sprintf("Create a directory '%s' and a file inside it called 'success.txt' with the content 'Verified'", folderName)
	agent := NewAgent(llm, db, goal)

	state, done := agent.Step()
	if done && state == StateCompleted {
		fmt.Printf("Agent completed. Final state: %d\n", state)
		return true
	}

	fmt.Printf("Agent did not complete properly. Final state: %d\n", state)
	return false
}

func init() {
	_ = time.Now()
}