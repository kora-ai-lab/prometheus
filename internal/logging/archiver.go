package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kora-ai-lab/prometheus/internal/llm"
)

var nowFunc = func() time.Time {
	return time.Now().UTC()
}

type Archiver struct {
	logsDir    string
	summaryDir string
	archiveDir string
	provider   llm.ModelProvider
}

func NewArchiver(logsDir, summaryDir, archiveDir string, provider llm.ModelProvider) *Archiver {
	return &Archiver{
		logsDir:    logsDir,
		summaryDir: summaryDir,
		archiveDir: archiveDir,
		provider:   provider,
	}
}

func (a *Archiver) ArchivePreviousMonth(ctx context.Context) error {
	now := nowFunc()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	prevMonth := currentMonth.AddDate(0, -1, 0)
	monthStr := prevMonth.Format("2006-01")

	destDir := filepath.Join(a.archiveDir, monthStr)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	if err := a.moveMonthLogs(monthStr, destDir); err != nil {
		return fmt.Errorf("failed to move month logs: %w", err)
	}

	if err := a.generateMonthlySummary(ctx, monthStr, destDir); err != nil {
		return fmt.Errorf("failed to generate monthly summary: %w", err)
	}

	return nil
}

func (a *Archiver) moveMonthLogs(monthStr, destDir string) error {
	pattern := monthStr + "-*.jsonl.zst"
	matches, err := filepath.Glob(filepath.Join(a.logsDir, pattern))
	if err != nil {
		return fmt.Errorf("failed to glob logs: %w", err)
	}

	for _, srcPath := range matches {
		filename := filepath.Base(srcPath)
		dstPath := filepath.Join(destDir, filename)

		if err := os.Rename(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to move %s: %w", filename, err)
		}
	}

	return nil
}

func (a *Archiver) generateMonthlySummary(ctx context.Context, monthStr, archiveDir string) error {
	pattern := monthStr + "-*.md"
	matches, err := filepath.Glob(filepath.Join(a.summaryDir, pattern))
	if err != nil {
		return fmt.Errorf("failed to glob summaries: %w", err)
	}

	sort.Strings(matches)

	var dailySummaries []string
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		filename := filepath.Base(path)
		date := strings.TrimSuffix(filename, ".md")
		dailySummaries = append(dailySummaries, fmt.Sprintf("## %s\n\n%s", date, string(data)))
	}

	if len(dailySummaries) == 0 {
		return fmt.Errorf("no daily summaries found for %s", monthStr)
	}

	combined := strings.Join(dailySummaries, "\n\n---\n\n")
	summary, err := a.createMonthlySummary(ctx, monthStr, combined)
	if err != nil {
		return err
	}

	summaryPath := filepath.Join(archiveDir, "summary.md")
	return os.WriteFile(summaryPath, []byte(summary), 0644)
}

func (a *Archiver) createMonthlySummary(ctx context.Context, monthStr, combinedDaily string) (string, error) {
	generated := time.Now().UTC().Format(time.RFC3339)

	if a.provider != nil {
		messages := []llm.Message{
			{
				Role: "system",
				Content: `You are a helpful assistant that creates monthly journal summaries.
Create a markdown summary with these sections:
- ## Overview (2-3 sentences)
- ## Highlights (key achievements)
- ## Statistics (aggregate stats)`,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Create a monthly journal summary for %s based on these daily summaries:\n\n%s", monthStr, combinedDaily),
			},
		}

		resp, err := a.provider.Complete(ctx, messages)
		if err == nil && resp != nil && resp.Content != "" {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("<!-- prometheus_monthly_summary_version:1 month:%s generated:%s -->\n", monthStr, generated))
			sb.WriteString(resp.Content)
			return sb.String(), nil
		}
	}

	return a.fallbackMonthlySummary(monthStr, combinedDaily, generated)
}

func (a *Archiver) fallbackMonthlySummary(monthStr, combinedDaily, generated string) (string, error) {
	var totalTasks, completedTasks, llmCalls int
	var totalInputTokens, totalOutputTokens int

	lines := strings.Split(combinedDaily, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Tasks started:") {
			var count int
			fmt.Sscanf(line, "Tasks started: %d", &count)
			totalTasks += count
		}
		if strings.Contains(line, "Tasks completed:") {
			var count int
			fmt.Sscanf(line, "Tasks completed: %d", &count)
			completedTasks += count
		}
		if strings.Contains(line, "LLM calls:") {
			var count int
			fmt.Sscanf(line, "LLM calls: %d", &count)
			llmCalls += count
		}
		if strings.Contains(line, "Total input tokens:") {
			var count int
			fmt.Sscanf(line, "Total input tokens: %d", &count)
			totalInputTokens += count
		}
		if strings.Contains(line, "Total output tokens:") {
			var count int
			fmt.Sscanf(line, "Total output tokens: %d", &count)
			totalOutputTokens += count
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<!-- prometheus_monthly_summary_version:1 month:%s generated:%s -->\n", monthStr, generated))
	sb.WriteString(fmt.Sprintf("# Monthly Journal - %s\n\n", monthStr))

	sb.WriteString("## Overview\n")
	if completedTasks > 0 {
		sb.WriteString(fmt.Sprintf("A productive month with %d tasks completed.\n\n", completedTasks))
	} else if totalTasks > 0 {
		sb.WriteString(fmt.Sprintf("A busy month with %d tasks started.\n\n", totalTasks))
	} else {
		sb.WriteString("A month with recorded activity.\n\n")
	}

	sb.WriteString("## Statistics\n")
	sb.WriteString(fmt.Sprintf("- Total tasks started: %d\n", totalTasks))
	sb.WriteString(fmt.Sprintf("- Tasks completed: %d\n", completedTasks))
	sb.WriteString(fmt.Sprintf("- LLM calls: %d\n", llmCalls))
	sb.WriteString(fmt.Sprintf("- Total input tokens: %d\n", totalInputTokens))
	sb.WriteString(fmt.Sprintf("- Total output tokens: %d\n", totalOutputTokens))

	return sb.String(), nil
}