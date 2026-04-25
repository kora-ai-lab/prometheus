package llm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ProgressWriter struct {
	Writer     io.Writer
	OnProgress func(downloaded, total int64)
	Total      int64
	Downloaded int64
}

func (w *ProgressWriter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	w.Downloaded += int64(n)
	if w.OnProgress != nil {
		w.OnProgress(w.Downloaded, w.Total)
	}
	return n, err
}

func Download(ctx context.Context, entry *ModelEntry, dest string, progress func(downloaded, total int64)) error {
	if entry == nil || entry.URL == "" {
		return fmt.Errorf("model entry missing download URL")
	}

	tmpPath := dest + ".tmp"
	var startByte int64
	if fi, err := os.Stat(tmpPath); err == nil {
		startByte = fi.Size()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, entry.URL, nil)
	if err != nil {
		return err
	}
	if startByte > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startByte))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	hash := sha256.New()
	if startByte > 0 {
		existing, err := os.ReadFile(tmpPath)
		if err != nil {
			return err
		}
		if _, err := hash.Write(existing); err != nil {
			return err
		}
	}

	counter := &ProgressWriter{
		Writer:     io.MultiWriter(f, hash),
		OnProgress: progress,
		Total:      entry.SizeBytes,
		Downloaded: startByte,
	}
	if _, err := io.Copy(counter, resp.Body); err != nil {
		return err
	}

	got := hex.EncodeToString(hash.Sum(nil))
	if entry.SHA256 != "" && got != entry.SHA256 {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("sha256 mismatch: got %s want %s", got, entry.SHA256)
	}

	return os.Rename(tmpPath, dest)
}
