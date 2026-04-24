package embedded

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"

	"github.com/prometheus-dev/prometheus/internal/config"
)

var ErrPlatformNotSupported = errors.New("embedded llama-server is not available on this platform")
var ErrPlaceholderArtifact = errors.New("embedded llama-server artifact is a scaffold placeholder")

func ExtractServer() (string, error) {
	if len(ServerBinary) == 0 || ServerName == "" {
		return "", ErrPlatformNotSupported
	}
	if bytes.Contains(ServerBinary, []byte("PROMETHEUS_PLACEHOLDER_ARTIFACT")) {
		return "", ErrPlaceholderArtifact
	}

	home := config.PrometheusHome()
	dest := filepath.Join(home, "runtime", "llama-server")
	if filepath.Ext(ServerName) == ".exe" {
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
