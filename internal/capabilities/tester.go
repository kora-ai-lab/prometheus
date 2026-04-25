package capabilities

import (
	"context"
	"time"

	"github.com/prometheus-dev/prometheus/internal/executor"
)

type TestResult struct {
	Ok       bool
	Type     string
	Output   string
	Duration time.Duration
}

type Tester struct {
	execer  executor.Executor
	timeout time.Duration
}

func NewTester(execer executor.Executor) *Tester {
	return &Tester{
		execer:  execer,
		timeout: 30 * time.Second,
	}
}

func (t *Tester) RunSyntaxCheck(lang, code string) TestResult {
	var cmd string

	switch lang {
	case "python":
		cmd = "python -m py_compile -c \"" + code + "\""
	case "bash":
		cmd = "bash -n -c \"" + code + "\""
	default:
		cmd = lang + " -m py_compile -c \"" + code + "\""
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	result := t.execer.Execute(ctx, cmd, executor.ExecOptions{
		Timeout: t.timeout,
	})

	return TestResult{
		Ok:       result.ExitCode == 0,
		Type:     "syntax",
		Output:   result.Stdout + result.Stderr,
		Duration: result.Duration,
	}
}

func (t *Tester) RunExecutionCheck(lang, mainFile, testInput string) TestResult {
	var cmd string

	switch lang {
	case "python":
		cmd = "python " + mainFile
	case "bash":
		cmd = "bash " + mainFile
	default:
		cmd = lang + " " + mainFile
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	envVars := []string{}
	if testInput != "" {
		envVars = append(envVars, "INPUT="+testInput)
	}

	result := t.execer.Execute(ctx, cmd, executor.ExecOptions{
		Timeout: t.timeout,
		Env:     envVars,
	})

	return TestResult{
		Ok:       result.ExitCode == 0,
		Type:     "execution",
		Output:   result.Stdout + result.Stderr,
		Duration: result.Duration,
	}
}

func (t *Tester) RunFunctionalCheck(lang, testFile string) TestResult {
	var cmd string

	switch lang {
	case "python":
		cmd = "pytest " + testFile
	case "bash":
		cmd = "bats " + testFile
	default:
		cmd = lang + " " + testFile
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	result := t.execer.Execute(ctx, cmd, executor.ExecOptions{
		Timeout: t.timeout,
	})

	return TestResult{
		Ok:       result.ExitCode == 0,
		Type:     "functional",
		Output:   result.Stdout + result.Stderr,
		Duration: result.Duration,
	}
}