package embedded

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/prometheus-dev/prometheus/internal/config"
	"github.com/prometheus-dev/prometheus/internal/llm"
)

var ErrPlatformNotSupported = errors.New("embedded llama-server is not available on this platform")
var ErrPlaceholderArtifact = errors.New("embedded llama-server artifact is a scaffold placeholder")

func ExtractServer() (string, error) {
	if len(ServerBinary) > 0 && !bytes.Contains(ServerBinary, []byte("PROMETHEUS_PLACEHOLDER_ARTIFACT")) {
		return extractEmbedded()
	}
	return downloadLlamaServer()
}

func extractEmbedded() (string, error) {
	home := config.PrometheusHome()
	dest := filepath.Join(home, "runtime", "llama-server")
	if runtime.GOOS == "windows" {
		dest += ".exe"
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return "", err
	}
	if ok, err := existingOK(dest, sha256OfEmbedded()); err == nil && ok {
		return dest, nil
	}
	if err := os.WriteFile(dest, ServerBinary, 0o700); err != nil {
		return "", err
	}
	return dest, nil
}

func downloadLlamaServer() (string, error) {
	entry := llm.LlamaServerEntry
	platform := runtime.GOOS + "/" + runtime.GOARCH
	plat, ok := entry.Platforms[platform]
	if !ok {
		return "", ErrPlatformNotSupported
	}

	home := config.PrometheusHome()
	runtimeDir := filepath.Join(home, "runtime")
	if err := os.MkdirAll(runtimeDir, 0o700); err != nil {
		return "", err
	}

	zipPath := filepath.Join(runtimeDir, plat.ZIPFilename)
	exePath := filepath.Join(runtimeDir, plat.ExeName)

	if _, err := os.Stat(exePath); err == nil {
		return exePath, nil
	}

	fmt.Printf("⬇️  Downloading llama-server %s...\n", entry.Version)
	resp, err := http.Get(entry.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	f, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	f.Close()

	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer zf.Close()

	for _, file := range zf.File {
		if strings.HasPrefix(file.Name, "llama-server") {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

		 dst, err := os.OpenFile(exePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(dst, rc); err != nil {
				dst.Close()
				return "", err
			}
			dst.Close()
			break
		}
	}

	os.Remove(zipPath)
	return exePath, nil
}

func sha256OfEmbedded() string {
	sum := sha256.Sum256(ServerBinary)
	return hex.EncodeToString(sum[:])
}

func existingOK(path, expected string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == expected, nil
}
