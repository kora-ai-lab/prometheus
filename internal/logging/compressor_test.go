package logging

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
)

func TestCompressor_CompressFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "test.jsonl")
	dst := filepath.Join(tmpDir, "test.jsonl.zst")

	srcData := make([]byte, 5*1024*1024)
	for i := range srcData {
		srcData[i] = byte(i % 256)
	}
	if err := os.WriteFile(src, srcData, 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	c := NewCompressor(tmpDir)
	if err := c.CompressFile(src, dst); err != nil {
		t.Fatalf("CompressFile failed: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("compressed file not created: %v", err)
	}

	originalSize := int64(len(srcData))
	compressedSize := info.Size()
	ratio := float64(compressedSize) / float64(originalSize)

	t.Logf("Original: %d bytes, Compressed: %d bytes, Ratio: %.2f%%", originalSize, compressedSize, ratio*100)

	if ratio > 0.15 {
		t.Errorf("compression ratio too high: got %.2f%%, want < 15%%", ratio*100)
	}
}

func TestCompressor_VerifyIntegrity(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "test.jsonl")
	dst := filepath.Join(tmpDir, "test.jsonl.zst")

	testData := []byte(`{"level":"info","msg":"test log entry"}` + "\n")
	if err := os.WriteFile(src, testData, 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	c := NewCompressor(tmpDir)
	if err := c.CompressFile(src, dst); err != nil {
		t.Fatalf("CompressFile failed: %v", err)
	}

	if !c.VerifyIntegrity(dst) {
		t.Error("VerifyIntegrity returned false for valid compressed file")
	}
}

func TestCompressor_VerifyIntegrity_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.zst")

	if err := os.WriteFile(invalidFile, []byte("not zstd data"), 0644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	c := NewCompressor(tmpDir)
	if c.VerifyIntegrity(invalidFile) {
		t.Error("VerifyIntegrity returned true for invalid compressed file")
	}
}

func TestCompressor_CompressOldLogs(t *testing.T) {
	tmpDir := t.TempDir()

	for i := 0; i < 3; i++ {
		fname := filepath.Join(tmpDir, "app-2024-01-0"+string(rune('1'+i))+".jsonl")
		if err := os.WriteFile(fname, []byte(`{"level":"info"}`+"\n"), 0644); err != nil {
			t.Fatalf("failed to write log file: %v", err)
		}
	}

	c := NewCompressor(tmpDir)
	ctx := context.Background()
	if err := c.CompressOldLogs(ctx); err != nil {
		t.Fatalf("CompressOldLogs failed: %v", err)
	}

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	zstCount := 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".zst" {
			zstCount++
		}
	}

	if zstCount != 3 {
		t.Errorf("expected 3 .zst files, got %d", zstCount)
	}
}

func TestCompressor_CompressionLevel(t *testing.T) {
	c := NewCompressor("/tmp")
	if c.level != zstd.SpeedBestCompression {
		t.Errorf("expected SpeedBestCompression level, got %v", c.level)
	}
}
