package logging

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/prometheus-dev/prometheus/internal/llm"
)

type LogEntry struct {
	Ts      string         `json:"ts"`
	Session string        `json:"session"`
	TaskID  string        `json:"task_id,omitempty"`
	Level   string        `json:"level"`
	Event   map[string]any `json:"event"`
}

type Summarizer struct {
	logsDir    string
	summaryDir string
	provider  llm.ModelProvider
}

func NewSummarizer(logsDir, summaryDir string, provider llm.ModelProvider) *Summarizer {
	return &Summarizer{
		logsDir:    logsDir,
		summaryDir: summaryDir,
		provider:  provider,
	}
}

func (s *Summarizer) SummarizeDay(ctx context.Context, date string) error {
	if err := os.MkdirAll(s.summaryDir, 0755); err != nil {
		return err
	}

	events, err := s.loadDayEvents(date)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return fmt.Errorf("no events found for date %s", date)
	}

	digest := s.buildDigest(events)
	summary, err := s.generateSummary(ctx, date, digest)
	if err != nil {
		return err
	}

	summaryFile := filepath.Join(s.summaryDir, date+".md")
	return os.WriteFile(summaryFile, []byte(summary), 0644)
}

func (s *Summarizer) loadDayEvents(date string) ([]LogEntry, error) {
	var entries []LogEntry

	uncompressedPath := filepath.Join(s.logsDir, date+".jsonl")
	if _, err := os.Stat(uncompressedPath); err == nil {
		entries = s.readLogFile(uncompressedPath, entries)
	}

	compressedPath := filepath.Join(s.logsDir, date+".jsonl.zst")
	if _, err := os.Stat(compressedPath); err == nil {
		entries = s.readCompressedFile(compressedPath, entries)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Ts < entries[j].Ts
	})

	return entries, nil
}

func (s *Summarizer) readLogFile(path string, entries []LogEntry) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return entries
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry LogEntry
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &entry); err == nil {
			if entry.Ts == "" {
				var raw map[string]json.RawMessage
				if err := json.Unmarshal(line, &raw); err == nil {
					if v, ok := raw["ts"]; ok {
						json.Unmarshal(v, &entry.Ts)
					}
					if v, ok := raw["session"]; ok {
						json.Unmarshal(v, &entry.Session)
					}
					if v, ok := raw["task_id"]; ok {
						json.Unmarshal(v, &entry.TaskID)
					}
					if v, ok := raw["kind"]; ok {
						json.Unmarshal(v, &entry.Level)
					}
					if v, ok := raw["level"]; ok {
						json.Unmarshal(v, &entry.Level)
					}
					if v, ok := raw["event"]; ok {
						json.Unmarshal(v, &entry.Event)
					}
				}
			}
			entries = append(entries, entry)
		}
	}
	return entries
}

func (s *Summarizer) readCompressedFile(path string, entries []LogEntry) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return entries
	}
	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return s.readLogFile(path, entries)
	}
	defer decoder.Close()

	scanner := bufio.NewScanner(decoder)
	for scanner.Scan() {
		var entry LogEntry
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	if scanner.Err() != nil && len(entries) == 0 {
		return s.readLogFile(path, entries)
	}

	if len(entries) == 0 {
		return s.readLogFile(path, entries)
	}
	return entries
}

func (s *Summarizer) buildDigest(events []LogEntry) string {
	var tasksStarted, tasksCompleted, llmCalls int
	var totalInputTokens, totalOutputTokens int
	taskMap := make(map[string]bool)
	var eventsByTask []string

	for _, e := range events {
		switch e.Level {
		case "task_start":
			if e.TaskID != "" {
				taskMap[e.TaskID] = true
				tasksStarted++
				eventsByTask = append(eventsByTask, fmt.Sprintf("- Started task %s", e.TaskID))
			}
		case "task_end":
			if e.TaskID != "" && taskMap[e.TaskID] {
				tasksCompleted++
				status := "unknown"
				if s, ok := e.Event["status"]; ok {
					if str, ok := s.(string); ok {
						status = str
					}
				}
				eventsByTask = append(eventsByTask, fmt.Sprintf("- Completed task %s (%s)", e.TaskID, status))
			}
		case "llm_call":
			llmCalls++
			if inputTokens, ok := toFloat64(e.Event["input_tokens"]); ok {
				totalInputTokens += int(inputTokens)
			}
			if outputTokens, ok := toFloat64(e.Event["output_tokens"]); ok {
				totalOutputTokens += int(outputTokens)
			}
		}
	}

	if len(eventsByTask) > 10 {
		eventsByTask = eventsByTask[:10]
	}

	var sb strings.Builder
	sb.WriteString("Events:\n")
	for _, e := range eventsByTask {
		sb.WriteString(e)
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("\nTasks: started=%d, completed=%d\n", tasksStarted, tasksCompleted))
	sb.WriteString(fmt.Sprintf("LLM calls: %d (input: %d tokens, output: %d tokens)\n", llmCalls, totalInputTokens, totalOutputTokens))

	return sb.String()
}

func (s *Summarizer) generateSummary(ctx context.Context, date, digest string) (string, error) {
	generated := time.Now().UTC().Format(time.RFC3339)

	if s.provider != nil {
		messages := []llm.Message{
			{
				Role: "system",
				Content: `You are a helpful assistant that creates daily journal summaries. 
Create a markdown summary with these sections:
- ## En bref (2-3 sentences summary)
- ## Accompli (bullet list of key accomplishments)
- ## Stats (statistics)`,
			},
			{
				Role: "user",
				Content: fmt.Sprintf("Create a daily journal summary for %s based on these events:\n%s", date, digest),
			},
		}

		resp, err := s.provider.Complete(ctx, messages)
		if err == nil && resp != nil && resp.Content != "" {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("<!-- prometheus_summary_version:1 date:%s generated:%s -->\n", date, generated))
			sb.WriteString(resp.Content)
			return sb.String(), nil
		}
	}

	return s.fallbackSummary(date, digest, generated)
}

func (s *Summarizer) fallbackSummary(date, digest, generated string) (string, error) {
	var tasksStarted, tasksCompleted, llmCalls int
	var totalInputTokens, totalOutputTokens int

	lines := strings.Split(digest, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Started task") {
			tasksStarted++
		}
		if strings.Contains(line, "Completed task") {
			tasksCompleted++
		}
		if strings.Contains(line, "LLM calls:") {
			fmt.Sscanf(line, "LLM calls: %d (input: %d tokens, output: %d tokens)", &llmCalls, &totalInputTokens, &totalOutputTokens)
		}
	}

	var accomplishments []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			accomplishments = append(accomplishments, line)
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<!-- prometheus_summary_version:1 date:%s generated:%s -->\n", date, generated))
	sb.WriteString(fmt.Sprintf("# Journal %s\n\n", date))

	sb.WriteString("## En bref\n")
	if tasksCompleted > 0 {
		sb.WriteString(fmt.Sprintf("A productive day with %d tasks completed.\n\n", tasksCompleted))
	} else if tasksStarted > 0 {
		sb.WriteString(fmt.Sprintf("A busy day with %d tasks started.\n\n", tasksStarted))
	} else {
		sb.WriteString("A quiet day with activity.\n\n")
	}

	sb.WriteString("## Accompli\n")
	if len(accomplishments) > 0 {
		for _, a := range accomplishments {
			sb.WriteString(a)
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("- No specific accomplishments recorded\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Stats\n")
	sb.WriteString(fmt.Sprintf("- Tasks started: %d\n", tasksStarted))
	sb.WriteString(fmt.Sprintf("- Tasks completed: %d\n", tasksCompleted))
	sb.WriteString(fmt.Sprintf("- LLM calls: %d\n", llmCalls))
	sb.WriteString(fmt.Sprintf("- Total input tokens: %d\n", totalInputTokens))
	sb.WriteString(fmt.Sprintf("- Total output tokens: %d\n", totalOutputTokens))

	return sb.String(), nil
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}