package sast

import (
	"regexp"
)

func getSecretsRules() []*Rule {
	return []*Rule{
		{
			ID:      "SEC001",
			Pattern: regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[=:]\s*["']([\w-]{16,})["']`),
		},
		{
			ID:      "SEC002",
			Pattern: regexp.MustCompile(`(?i)(password|passwd|pwd|secret)\s*[=:]\s*["'][^"']{3,}["']`),
		},
		{
			ID:      "SEC003",
			Pattern: regexp.MustCompile(`(?i)(jwt|bearer|auth[_-]?token)\s*[=:]\s*["'](\w{20,})["']`),
		},
		{
			ID:      "SEC004",
			Pattern: regexp.MustCompile(`-----BEGIN (RSA|EC|OPENSSH|DSA) PRIVATE KEY-----`),
		},
		{
			ID:      "SEC005",
			Pattern: regexp.MustCompile(`(ghp_\w{36}|glpat-\w{20}|xoxb-\d+-\w+|sk-\w{48})`),
		},
		{
			ID:      "SEC006",
			Pattern: regexp.MustCompile(`(?i)\b(md5|sha1|des|3des|rc4)\s*\(`),
		},
		{
			ID:      "SEC007",
			Pattern: regexp.MustCompile(`(?i)(verify\s*=\s*False|NODE_TLS_REJECT_UNAUTHORIZED\s*=\s*0)`),
		},
	}
}