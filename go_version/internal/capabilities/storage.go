package capabilities

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

type Storage struct {
	baseDir string
}

func NewStorage(baseDir string) *Storage {
	return &Storage{baseDir: baseDir}
}

type metaData struct {
	Type       string    `toml:"type"`
	Language   string    `toml:"language"`
	Name       string    `toml:"name"`
	Verified   bool      `toml:"verified"`
	ForgedAt   time.Time `toml:"forged_at"`
	RetryCount int       `toml:"retry_count"`
	MainFile   string    `toml:"main_file"`
	TestFile   string    `toml:"test_file"`
}

func (s *Storage) SaveTool(spec *Spec, code, testCode string) error {
	toolDir := filepath.Join(s.baseDir, "forged", spec.Name)
	err := os.MkdirAll(toolDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	mainFile := filepath.Join(toolDir, spec.MainFile)
	err = os.WriteFile(mainFile, []byte(code), 0644)
	if err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	testFile := filepath.Join(toolDir, spec.TestFile)
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	if err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}

	meta := metaData{
		Type:       "tool",
		Language:   spec.Language,
		Name:       spec.Name,
		Verified:   false,
		ForgedAt:   time.Now(),
		RetryCount: 0,
		MainFile:   spec.MainFile,
		TestFile:   spec.TestFile,
	}

	metaFile := filepath.Join(toolDir, "meta.toml")
	metaBytes, err := toml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	err = os.WriteFile(metaFile, metaBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func (s *Storage) LoadTool(name string) (code, testCode, meta string, err error) {
	toolDir := filepath.Join(s.baseDir, "forged", name)

	metaBytes, err := os.ReadFile(filepath.Join(toolDir, "meta.toml"))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read metadata: %w", err)
	}
	meta = string(metaBytes)

	var m metaData
	err = toml.Unmarshal(metaBytes, &m)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse metadata: %w", err)
	}

	codeBytes, err := os.ReadFile(filepath.Join(toolDir, m.MainFile))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read main file: %w", err)
	}
	code = string(codeBytes)

	if m.TestFile != "" {
		testBytes, err := os.ReadFile(filepath.Join(toolDir, m.TestFile))
		if err != nil {
			return "", "", "", fmt.Errorf("failed to read test file: %w", err)
		}
		testCode = string(testBytes)
	}

	return code, testCode, meta, nil
}

func (s *Storage) DeleteTool(name string) error {
	toolDir := filepath.Join(s.baseDir, "forged", name)
	err := os.RemoveAll(toolDir)
	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}
	return nil
}

func (s *Storage) ListForged() ([]string, error) {
	forgedDir := filepath.Join(s.baseDir, "forged")
	entries, err := os.ReadDir(forgedDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read forged directory: %w", err)
	}

	var tools []string
	for _, e := range entries {
		if e.IsDir() {
			tools = append(tools, e.Name())
		}
	}
	return tools, nil
}

func (s *Storage) GetPath(name string) string {
	return filepath.Join(s.baseDir, "forged", name)
}