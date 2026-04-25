package sast

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Severity string

const (
	Critical Severity = "CRITICAL"
	High     Severity = "HIGH"
	Medium   Severity = "MEDIUM"
	Low      Severity = "LOW"
	Info     Severity = "INFO"
)

type Finding struct {
	RuleID   string
	File    string
	Line    int
	Code    string
	Severity Severity
}

type FindingSet struct {
	File     string
	Findings []*Finding
}

type Rule struct {
	ID      string
	Pattern *regexp.Regexp
}

type Scanner struct {
	rules []*Rule
}

func NewScanner() *Scanner {
	return &Scanner{
		rules: loadAllRules(),
	}
}

func loadAllRules() []*Rule {
	var rules []*Rule
	rules = append(rules, getSecretsRules()...)
	rules = append(rules, getInjectionRules()...)
	rules = append(rules, getAuthRules()...)
	return rules
}

func (s *Scanner) ScanFile(ctx context.Context, path string) (*FindingSet, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	set := &FindingSet{File: path}
	lines := strings.Split(string(content), "\n")

	for lineNum, line := range lines {
		for _, rule := range s.rules {
			if rule.Pattern.MatchString(line) {
				set.Findings = append(set.Findings, &Finding{
					RuleID:   rule.ID,
					File:    path,
					Line:    lineNum + 1,
					Code:    strings.TrimSpace(line),
					Severity: High,
				})
			}
		}
	}

	return set, nil
}

func (s *Scanner) ScanDir(ctx context.Context, dir string) ([]*FindingSet, error) {
	var results []*FindingSet

	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil { return err }
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == ".git" || name == "__pycache__" ||
			   name == "vendor" || name == ".venv" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".py" || ext == ".js" || ext == ".ts" || ext == ".go" ||
		   ext == ".sh" || ext == ".sql" {
			set, err := s.ScanFile(ctx, path)
			if err == nil && len(set.Findings) > 0 {
				results = append(results, set)
			}
		}
		return nil
	})

	return results, nil
}

func getInjectionRules() []*Rule { return nil }
func getAuthRules() []*Rule { return nil }