package update

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadAndVerify(binaryURL, expectedSHA256 string) (tmpPath string, err error) {
	expectedSHA256 = strings.ToLower(strings.ReplaceAll(expectedSHA256, " ", ""))
	if len(expectedSHA256) != 64 {
		return "", fmt.Errorf("invalid SHA256 length: %d", len(expectedSHA256))
	}

	resp, err := http.Get(binaryURL)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	tmpDir := os.TempDir()
	tmpPath = filepath.Join(tmpDir, "prometheus-update")

	f, err := os.Create(tmpPath)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	writer := io.MultiWriter(f, h)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("write failed: %w", err)
	}

	actualHash := fmt.Sprintf("%x", h.Sum(nil))
	if actualHash != expectedSHA256 {
		os.Remove(tmpPath)
		return "", fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA256, actualHash)
	}

	return tmpPath, nil
}