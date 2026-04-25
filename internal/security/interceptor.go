package security

import (
	"errors"
	"strings"

	"github.com/kora-ai-lab/prometheus/internal/config"
)

type Interceptor struct {
	cfg config.SecurityConfig
}

func New(cfg config.SecurityConfig) *Interceptor {
	return &Interceptor{cfg: cfg}
}

func (i *Interceptor) Allow(command string) (bool, error) {
	score := 0
	lower := strings.ToLower(command)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lower, strings.ToLower(pattern.Needle)) {
			score += pattern.Score
		}
	}

	switch {
	case score >= 91:
		return false, errors.New("command blocked by security policy")
	case score >= 71 && i.cfg.DangerousOpsConfirm:
		return false, errors.New("command requires explicit confirmation")
	default:
		return true, nil
	}
}
