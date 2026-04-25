package main

import (
	"fmt"
	"strings"
)

type ContextManager struct {
	SystemPrompt string
	History      string
	Goal         string
	MaxTokens    int // Simplified: max lines/chars for this POC
}

func NewContextManager(systemPrompt, goal string) *ContextManager {
	return &ContextManager{
		SystemPrompt: systemPrompt,
		Goal:         goal,
		MaxTokens:    50, // Limit to 50 entries before compaction
	}
}

func (cm *ContextManager) AddEntry(command, observation string) {
	cm.History += fmt.Sprintf("\nCommand: %s\nObservation: %s", command, observation)
}

func (cm *ContextManager) GetFullContext() string {
	return fmt.Sprintf("%s\n\nGoal: %s\n%s\n\nWhat is your next action?", cm.SystemPrompt, cm.Goal, cm.History)
}

func (cm *ContextManager) Compact(llm ModelProvider) (string, error) {
	if len(strings.Split(cm.History, "\nCommand: ")) < 10 {
		return cm.History, nil // Not enough history to compact yet
	}

	prompt := fmt.Sprintf(
		"Summarize the following execution history. Keep critical decisions, errors, and the current state. Be concise.\n\n%s",
		cm.History,
	)

	summary, err := llm.Generate(prompt)
	if err != nil {
		return "", err
	}

	cm.History = fmt.Sprintf("--- COMPACTED HISTORY ---\n%s\n--- END SUMMARY ---", summary)
	return cm.History, nil
}
