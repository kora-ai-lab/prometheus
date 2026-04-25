package main

import (
	"bytes"
	"os/exec"
)

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func ExecuteCommand(command string) (*ExecutionResult, error) {
	cmd := exec.Command("cmd", "/C", command) // Use cmd /C for Windows
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			return nil, err
		}
	}

	return &ExecutionResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}
