package logging

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

type Compressor struct {
	logsDir string
	level   zstd.EncoderLevel
}

func NewCompressor(logsDir string) *Compressor {
	return &Compressor{
		logsDir: logsDir,
		level:   zstd.SpeedBestCompression,
	}
}

func (c *Compressor) CompressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	encoder, err := zstd.NewWriter(dstFile, zstd.WithEncoderLevel(c.level))
	if err != nil {
		return err
	}
	defer encoder.Close()

	if _, err := io.Copy(encoder, srcFile); err != nil {
		return err
	}

	return encoder.Close()
}

func (c *Compressor) VerifyIntegrity(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return false
	}
	defer decoder.Close()

	_, err = io.Copy(io.Discard, decoder)
	return err == nil
}

func (c *Compressor) CompressOldLogs(ctx context.Context) error {
	entries, err := os.ReadDir(c.logsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		name := entry.Name()
		if filepath.Ext(name) != ".jsonl" {
			continue
		}

		src := filepath.Join(c.logsDir, name)
		dst := src + ".zst"

		if _, err := os.Stat(dst); err == nil {
			continue
		}

		if err := c.CompressFile(src, dst); err != nil {
			return err
		}

		if err := os.Remove(src); err != nil {
			return err
		}
	}

	return nil
}

var ErrInvalidZstd = errors.New("invalid zstd data")
