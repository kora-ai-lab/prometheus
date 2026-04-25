package sast

import (
	"regexp"
)

func getInjectionRules() []*Rule {
	return []*Rule{
		{
			ID:      "INJ001",
			Pattern: regexp.MustCompile(`(?i)(execute|executemany)\s*\(\s*["'].*\+`),
		},
		{
			ID:      "INJ002",
			Pattern: regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*\.format\(`),
		},
		{
			ID:      "INJ003",
			Pattern: regexp.MustCompile(`(?i)fmt\.Sprintf.*(SELECT|INSERT|UPDATE|DELETE)`),
		},
		{
			ID:      "INJ004",
			Pattern: regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE)[^"']*["']\s*%\s*[\(\w]`),
		},
		{
			ID:      "INJ010",
			Pattern: regexp.MustCompile(`(?i)os\.(system|popen)\s*\(\s*.*request`),
		},
		{
			ID:      "INJ011",
			Pattern: regexp.MustCompile(`subprocess.*shell\s*=\s*True.*request`),
		},
		{
			ID:      "INJ012",
			Pattern: regexp.MustCompile(`(?i)eval\s*\(\s*.*request`),
		},
		{
			ID:      "INJ013",
			Pattern: regexp.MustCompile(`(?i)\beval\s*\(\s*.*req\.`),
		},
		{
			ID:      "INJ020",
			Pattern: regexp.MustCompile(`\.innerHTML\s*[+=]\s*[^"']`),
		},
		{
			ID:      "INJ021",
			Pattern: regexp.MustCompile(`document\.(write|writeln)\s*\(.*req\.`),
		},
		{
			ID:      "INJ030",
			Pattern: regexp.MustCompile(`(?i)open\s*\(\s*.*request`),
		},
		{
			ID:      "INJ040",
			Pattern: regexp.MustCompile(`(?i)(xml\.etree|lxml\.etree).*parse`),
		},
		{
			ID:      "INJ050",
			Pattern: regexp.MustCompile(`requests\.(get|post).*request\.`),
		},
		{
			ID:      "INJ051",
			Pattern: regexp.MustCompile(`fetch\s*\(\s*(?:req\.|request\.)\)`),
		},
	}
}