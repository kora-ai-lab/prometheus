package sast

import (
	"regexp"
)

func getAuthRules() []*Rule {
	return []*Rule{
		{
			ID:      "AUTH001",
			Pattern: regexp.MustCompile(`(?i)@app\.route.*\ndef\s+\w+\s*\([^)]*\):\s*$`),
		},
		{
			ID:      "AUTH002",
			Pattern: regexp.MustCompile(`(?i)(\/admin|\/superuser|\/root|\/management)`),
		},
		{
			ID:      "AUTH003",
			Pattern: regexp.MustCompile(`(?i)get_object_or_404.*request\.`),
		},
		{
			ID:      "AUTH010",
			Pattern: regexp.MustCompile(`(?i)(password|token|secret)\s*==\s*(?:request|input)`),
		},
		{
			ID:      "AUTH011",
			Pattern: regexp.MustCompile(`(?i)hashlib\.(md5|sha1)\(.*password`),
		},
		{
			ID:      "AUTH012",
			Pattern: regexp.MustCompile(`(?i)jwt\.decode.*verify\s*=\s*False`),
		},
		{
			ID:      "AUTH013",
			Pattern: regexp.MustCompile(`(?i)jwt.*secret\s*[=:]\s*["'](?:secret|mysecret)`),
		},
		{
			ID:      "AUTH020",
			Pattern: regexp.MustCompile(`(?i)pickle\.(loads|load)\s*\(\s*.*request`),
		},
		{
			ID:      "AUTH021",
			Pattern: regexp.MustCompile(`(?i)__proto__\s*:`),
		},
		{
			ID:      "AUTH022",
			Pattern: regexp.MustCompile(`yaml\.load\s*\([^)]*\)`),
		},
	}
}