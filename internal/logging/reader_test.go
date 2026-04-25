package logging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReader_ReadDay_Jsonl(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	testData := `{"level":"info","msg":"test1"}
{"level":"warn","msg":"test2"}
`
	logFile := filepath.Join(logsDir, "2024-01-15.jsonl")
	if err := os.WriteFile(logFile, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	r := NewReader(logsDir, archiveDir)
	entries, err := r.ReadDay("2024-01-15")
	if err != nil {
		t.Fatalf("ReadDay failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Level != "info" {
		t.Errorf("expected level 'info', got %v", entries[0].Level)
	}
}

func TestReader_ReadDay_Zstd(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	testData := `{"level":"info","msg":"test1"}
{"level":"warn","msg":"test2"}
`
	srcFile := filepath.Join(logsDir, "2024-01-15.jsonl")
	if err := os.WriteFile(srcFile, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	c := NewCompressor(logsDir)
	zstFile := filepath.Join(logsDir, "2024-01-15.jsonl.zst")
	if err := c.CompressFile(srcFile, zstFile); err != nil {
		t.Fatalf("failed to compress file: %v", err)
	}
	os.Remove(srcFile)

	r := NewReader(logsDir, archiveDir)
	entries, err := r.ReadDay("2024-01-15")
	if err != nil {
		t.Fatalf("ReadDay failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestReader_ReadDay_Archived(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	testData := `{"level":"info","msg":"archived"}`
	archiveMonthDir := filepath.Join(archiveDir, "2024-01")
	os.MkdirAll(archiveMonthDir, 0755)

	srcFile := filepath.Join(archiveMonthDir, "2024-01-15.jsonl")
	if err := os.WriteFile(srcFile, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	c := NewCompressor(archiveDir)
	zstFile := srcFile + ".zst"
	if err := c.CompressFile(srcFile, zstFile); err != nil {
		t.Fatalf("failed to compress file: %v", err)
	}
	os.Remove(srcFile)

	r := NewReader(logsDir, archiveDir)
	entries, err := r.ReadDay("2024-01-15")
	if err != nil {
		t.Fatalf("ReadDay failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestReader_ReadDay_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	logsDir := filepath.Join(tmpDir, "logs")
	archiveDir := filepath.Join(tmpDir, "archive")
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	r := NewReader(logsDir, archiveDir)
	_, err := r.ReadDay("2024-01-15")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}