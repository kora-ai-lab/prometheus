package capabilities

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus-dev/prometheus/internal/executor"
)

type mockExecutor struct {
	result *executor.ExecResult
}

func (m *mockExecutor) Execute(ctx context.Context, command string, opts executor.ExecOptions) *executor.ExecResult {
	return m.result
}

func TestTester_RunSyntaxCheck_ValidPython(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   "",
			Duration: 10 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunSyntaxCheck("python", "print('hello')")

	if !result.Ok {
		t.Errorf("expected ok=true for valid python code, got false")
	}
	if result.Type != "syntax" {
		t.Errorf("expected type=syntax, got %s", result.Type)
	}
}

func TestTester_RunSyntaxCheck_InvalidPython(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 1,
			Stdout:   "",
			Stderr:   "SyntaxError",
			Duration: 10 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunSyntaxCheck("python", "print('hello")

	if result.Ok {
		t.Errorf("expected ok=false for invalid python code, got true")
	}
}

func TestTester_RunSyntaxCheck_Bash(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 0,
			Stdout:   "",
			Stderr:   "",
			Duration: 10 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunSyntaxCheck("bash", "echo hello")

	if !result.Ok {
		t.Errorf("expected ok=true for valid bash, got false")
	}
	if result.Type != "syntax" {
		t.Errorf("expected type=syntax, got %s", result.Type)
	}
}

func TestTester_RunExecutionCheck_Success(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 0,
			Stdout:   "output",
			Stderr:   "",
			Duration: 100 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunExecutionCheck("python", "main.py", "hello")

	if !result.Ok {
		t.Errorf("expected ok=true, got false: %s", result.Output)
	}
	if result.Type != "execution" {
		t.Errorf("expected type=execution, got %s", result.Type)
	}
}

func TestTester_RunExecutionCheck_Failure(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 1,
			Stdout:   "",
			Stderr:   "error",
			Duration: 100 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunExecutionCheck("python", "main.py", "hello")

	if result.Ok {
		t.Errorf("expected ok=false for failed execution, got true")
	}
}

func TestTester_RunFunctionalCheck_Python(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 0,
			Stdout:   "passed",
			Stderr:   "",
			Duration: 500 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunFunctionalCheck("python", "test_example.py")

	if result.Type != "functional" {
		t.Errorf("expected type=functional, got %s", result.Type)
	}
	if !result.Ok {
		t.Errorf("expected ok=true for passing test, got false")
	}
}

func TestTester_RunFunctionalCheck_Failure(t *testing.T) {
	mock := &mockExecutor{
		result: &executor.ExecResult{
			ExitCode: 1,
			Stdout:   "",
			Stderr:   "test failed",
			Duration: 500 * time.Millisecond,
		},
	}
	tester := NewTester(mock)

	result := tester.RunFunctionalCheck("pytest", "test_example.py")

	if result.Ok {
		t.Errorf("expected ok=false for failing test, got true")
	}
}