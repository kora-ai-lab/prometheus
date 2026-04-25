package capabilities

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

type mockLLM struct {
	spec       *Spec
	code       string
	testCode   string
	specErr    error
	codeErr    error
	callCount  int
}

func (m *mockLLM) GenerateSpec(ctx context.Context, task string) (*Spec, error) {
	m.callCount++
	if m.specErr != nil {
		return nil, m.specErr
	}
	return m.spec, nil
}

func (m *mockLLM) GenerateCode(ctx context.Context, spec *Spec, task string) (string, string, error) {
	m.callCount++
	if m.codeErr != nil {
		return "", "", m.codeErr
	}
	return m.code, m.testCode, nil
}

type mockTester struct {
	syntaxResult  TestResult
	executionResult TestResult
	functionalResult TestResult
}

func (m *mockTester) RunSyntaxCheck(lang, code string) TestResult {
	return m.syntaxResult
}

func (m *mockTester) RunExecutionCheck(lang, mainFile, testInput string) TestResult {
	return m.executionResult
}

func (m *mockTester) RunFunctionalCheck(lang, testFile string) TestResult {
	return m.functionalResult
}

func TestForge_Forge_Success(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	tester := &mockTester{
		syntaxResult:     TestResult{Ok: true, Type: "syntax"},
		executionResult:  TestResult{Ok: true, Type: "execution"},
		functionalResult: TestResult{Ok: true, Type: "functional"},
	}
	llm := &mockLLM{
		spec: &Spec{
			Name:        "test-tool",
			Language:    "python",
			Description: "A test tool",
			MainFile:    "main.py",
			TestFile:    "test_main.py",
		},
		code:     "print('hello')",
		testCode: "def test_hello(): assert True",
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a test tool")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Name != "test-tool" {
		t.Errorf("expected name=test-tool, got %s", result.Name)
	}
	if result.Language != "python" {
		t.Errorf("expected language=python, got %s", result.Language)
	}
	if !result.Verified {
		t.Error("expected verified=true")
	}
	expectedPath := filepath.Join(tmpDir, "forged", "test-tool")
	if result.Path != expectedPath {
		t.Errorf("expected path=%s, got %s", expectedPath, result.Path)
	}
}

func TestForge_Forge_InvalidSpec(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	tester := &mockTester{}
	llm := &mockLLM{
		spec: &Spec{
			Name:     "",
			Language: "python",
		},
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a tool")

	if err == nil {
		t.Fatal("expected error for invalid spec, got nil")
	}
	if result != nil {
		t.Error("expected nil result for invalid spec")
	}
}

func TestForge_Forge_RetryOnTestFailure(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	callCount := 0
	tester := &mockTester{
		functionalResult: TestResult{Ok: false, Type: "functional", Output: "test failed"},
	}
	llm := &mockLLM{
		spec: &Spec{
			Name:        "retry-tool",
			Language:    "python",
			Description: "A tool that retries",
			MainFile:    "main.py",
			TestFile:    "test_main.py",
		},
		code:     "print('hello')",
		testCode: "def test_hello(): assert True",
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a retry tool")

	if err == nil {
		t.Fatal("expected error after test failures, got nil")
	}
	if result != nil {
		t.Error("expected nil result after failures")
	}
	if callCount > 3 {
		t.Error("should retry up to 3 times")
	}
	_ = callCount
}

func TestForge_Forge_MaxRetries(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	tester := &mockTester{
		functionalResult: TestResult{Ok: false, Type: "functional", Output: "test failed"},
	}
	llm := &mockLLM{
		spec: &Spec{
			Name:        "max-retry-tool",
			Language:    "python",
			Description: "A tool that always fails",
			MainFile:    "main.py",
			TestFile:    "test_main.py",
		},
		code:     "print('hello')",
		testCode: "def test_hello(): assert True",
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a tool that always fails")

	if err == nil {
		t.Fatal("expected error after max retries, got nil")
	}
	if result != nil {
		t.Error("expected nil result after max retries")
	}
}

func TestForge_SpecGenerationError(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	tester := &mockTester{}
	llm := &mockLLM{
		specErr: errors.New("spec generation failed"),
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a tool")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestForge_CodeGenerationError(t *testing.T) {
	tmpDir := t.TempDir()

	storage := NewStorage(tmpDir)
	tester := &mockTester{}
	llm := &mockLLM{
		spec: &Spec{
			Name:        "test-tool",
			Language:    "python",
			Description: "A test tool",
			MainFile:    "main.py",
			TestFile:    "test_main.py",
		},
		codeErr: errors.New("code generation failed"),
	}

	forge := NewForge(llm, storage, tester)

	result, err := forge.Forge(context.Background(), "create a tool")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}