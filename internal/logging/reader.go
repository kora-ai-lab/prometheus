package logging

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

type Reader struct {
	logsDir    string
	archiveDir string
}

func NewReader(logsDir, archiveDir string) *Reader {
	return &Reader{
		logsDir:    logsDir,
		archiveDir: archiveDir,
	}
}

func (r *Reader) ReadDay(date string) ([]LogEntry, error) {
	candidates := []string{
		filepath.Join(r.logsDir, date+".jsonl"),
		filepath.Join(r.logsDir, date+".jsonl.zst"),
		filepath.Join(r.archiveDir, date[0:7], date+".jsonl.zst"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return r.readFile(path)
		}
	}

	return nil, errors.New("log file not found")
}

func (r *Reader) readFile(path string) ([]LogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reader io.Reader = file

	if strings.HasSuffix(path, ".zst") {
		zstdReader, err := zstd.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer zstdReader.Close()
		reader = zstdReader
	}

	var entries []LogEntry
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}