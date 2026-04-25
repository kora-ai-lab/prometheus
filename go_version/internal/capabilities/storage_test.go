package capabilities

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_SaveTool(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	spec := &Spec{
		Name:        "test_tool",
		Language:    "go",
		Description: "A test tool",
		MainFile:    "main.go",
		TestFile:    "main_test.go",
	}

	code := `package main

func main() {
    println("Hello")
}`
	testCode := `package main

import "testing"

func TestMain(t *testing.T) {
    t.Run("test", func(t *testing.T) {
        t.Log("test")
    })
}`

	err := storage.SaveTool(spec, code, testCode)
	require.NoError(t, err)

	toolPath := filepath.Join(tmpDir, "forged", "test_tool")
	_, err = os.Stat(toolPath)
	require.NoError(t, err, "tool directory should exist")

	mainFile := filepath.Join(toolPath, "main.go")
	content, err := os.ReadFile(mainFile)
	require.NoError(t, err, "main.go should exist")
	require.Equal(t, code, string(content))

	testFile := filepath.Join(toolPath, "main_test.go")
	content, err = os.ReadFile(testFile)
	require.NoError(t, err, "main_test.go should exist")
	require.Equal(t, testCode, string(content))

	metaFile := filepath.Join(toolPath, "meta.toml")
	content, err = os.ReadFile(metaFile)
	require.NoError(t, err, "meta.toml should exist")
	require.Contains(t, string(content), "type = 'forged'")
	require.Contains(t, string(content), "name = 'test_tool'")
	require.Contains(t, string(content), "language = 'go'")
}

func TestStorage_LoadTool(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	spec := &Spec{
		Name:        "load_test",
		Language:    "python",
		Description: "A tool for loading",
		MainFile:    "main.py",
		TestFile:    "test_main.py",
	}

	code := `print("hello")`
	testCode := `def test_hello(): assert True`

	err := storage.SaveTool(spec, code, testCode)
	require.NoError(t, err)

	loadedCode, loadedTestCode, meta, err := storage.LoadTool("load_test")
	require.NoError(t, err)
	require.Equal(t, code, loadedCode)
	require.Equal(t, testCode, loadedTestCode)
	require.Contains(t, meta, "name = 'load_test'")
	require.Contains(t, meta, "language = 'python'")
}

func TestStorage_DeleteTool(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	spec := &Spec{
		Name:     "delete_me",
		Language: "go",
		MainFile: "main.go",
		TestFile: "main_test.go",
	}
	err := storage.SaveTool(spec, "code", "test")
	require.NoError(t, err)

	err = storage.DeleteTool("delete_me")
	require.NoError(t, err)

	toolPath := filepath.Join(tmpDir, "forged", "delete_me")
	_, err = os.Stat(toolPath)
	require.True(t, os.IsNotExist(err), "tool directory should be deleted")
}

func TestStorage_ListForged(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	err := storage.SaveTool(&Spec{Name: "tool1", Language: "go", MainFile: "main.go", TestFile: "main_test.go"}, "code1", "test1")
	require.NoError(t, err)
	err = storage.SaveTool(&Spec{Name: "tool2", Language: "python", MainFile: "main.py", TestFile: "test_main.py"}, "code2", "test2")
	require.NoError(t, err)

	tools, err := storage.ListForged()
	require.NoError(t, err)
	require.Len(t, tools, 2)
	require.Contains(t, tools, "tool1")
	require.Contains(t, tools, "tool2")
}

func TestStorage_GetPath(t *testing.T) {
	tmpDir := t.TempDir()
	storage := NewStorage(tmpDir)

	path := storage.GetPath("my_tool")
	expected := filepath.Join(tmpDir, "forged", "my_tool")
	require.Equal(t, expected, path)
}