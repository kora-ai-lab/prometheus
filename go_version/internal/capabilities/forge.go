package capabilities

import (
	"context"
	"fmt"

	ctxmgr "github.com/prometheus-dev/prometheus/internal/context"
)

type ForgeResult struct {
	Name     string
	Path     string
	Verified bool
	Language string
}

type LLM interface {
	GenerateSpec(ctx context.Context, task string) (*Spec, error)
	GenerateCode(ctx context.Context, spec *Spec, task string) (code, testCode string, err error)
}

type TesterInterface interface {
	RunSyntaxCheck(lang, code string) TestResult
	RunExecutionCheck(lang, mainFile, testInput string) TestResult
	RunFunctionalCheck(lang, testFile string) TestResult
}

type Forge struct {
	llm            LLM
	storage        *Storage
	tester         TesterInterface
	contextManager *ctxmgr.Manager
}

func NewForge(llm LLM, storage *Storage, tester TesterInterface, contextManager *ctxmgr.Manager) *Forge {
	return &Forge{
		llm:            llm,
		storage:        storage,
		tester:         tester,
		contextManager: contextManager,
	}
}

func (f *Forge) Forge(ctx context.Context, task string) (*ForgeResult, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		spec, err := f.llm.GenerateSpec(ctx, task)
		if err != nil {
			return nil, fmt.Errorf("spec generation failed: %w", err)
		}

		if err := spec.Validate(); err != nil {
			if attempt < maxRetries-1 {
				task = fmt.Sprintf("%s\n\nError: %v", task, err)
				continue
			}
			return nil, fmt.Errorf("spec validation failed: %w", err)
		}

		code, testCode, err := f.llm.GenerateCode(ctx, spec, task)
		if err != nil {
			return nil, fmt.Errorf("code generation failed: %w", err)
		}

		err = f.storage.SaveTool(spec, code, testCode)
		if err != nil {
			return nil, fmt.Errorf("failed to save tool: %w", err)
		}

		if !f.runTests(spec, code, testCode) {
			if attempt < maxRetries-1 {
				task = fmt.Sprintf("%s\n\nTests failed, please fix the code", task)
				continue
			}
			return nil, fmt.Errorf("tests failed after %d attempts", maxRetries)
		}

		result := &ForgeResult{
			Name:     spec.Name,
			Path:     f.storage.GetPath(spec.Name),
			Verified: true,
			Language: spec.Language,
		}

		Registry[result.Name] = Capability{
			Name:        result.Name,
			Type:        "forged",
			InstallCmds: map[string]string{},
			Description: "forged capability",
		}

		if f.contextManager != nil {
			f.contextManager.AppendBlockC(fmt.Sprintf("- %s: %s | %s | %s/",
				result.Name, spec.Description, result.Language, result.Path))
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func (f *Forge) runTests(spec *Spec, code, testCode string) bool {
	syntaxResult := f.tester.RunSyntaxCheck(spec.Language, code)
	if !syntaxResult.Ok {
		return false
	}

	executionResult := f.tester.RunExecutionCheck(spec.Language, spec.MainFile, "")
	if !executionResult.Ok {
		return false
	}

	functionalResult := f.tester.RunFunctionalCheck(spec.Language, spec.TestFile)
	return functionalResult.Ok
}