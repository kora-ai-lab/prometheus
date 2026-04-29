package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/klauspost/compress/zstd"
)

type SearchResult struct {
	Entry  LogEntry
	Score  float64
	Method string
}

type SearchEngine struct {
	logsDir  string
	index    map[string][]docID
	docStore map[docID]LogEntry
	provider interface{}
}

type docID int

func NewSearchEngine(logsDir string, provider interface{}) *SearchEngine {
	return &SearchEngine{
		logsDir:  logsDir,
		index:    make(map[string][]docID),
		docStore: make(map[docID]LogEntry),
		provider: provider,
	}
}

func (s *SearchEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	tokens := s.tokenize(query)

	if s.provider != nil {
		return s.searchEmbedding(ctx, query, limit)
	}
	return s.searchBM25(tokens, limit), nil
}

func (s *SearchEngine) SearchTemporal(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	tq := ParseTemporalQuery(query)

	var dates []string

	if tq.StartDate != "" && tq.EndDate != "" {
		start, err1 := time.Parse("2006-01-02", tq.StartDate)
		end, err2 := time.Parse("2006-01-02", tq.EndDate)

		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid date format: start=%v, end=%v", err1, err2)
		}

		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			dates = append(dates, d.Format("2006-01-02"))
		}
	} else {
		dates = []string{time.Now().Format("2006-01-02")}
	}

	for _, date := range dates {
		dayEntries, err := loadDayEvents(date, s.logsDir)
		if err != nil {
			continue
		}

		for i, entry := range dayEntries {
			id := docID(len(s.docStore) + i)
			s.docStore[id] = entry

			tokens := s.tokenize(s.entryToText(entry))
			for _, token := range tokens {
				s.index[token] = append(s.index[token], id)
			}
		}
	}

	return s.Search(ctx, query, limit)
}

func (s *SearchEngine) tokenize(text string) []string {
	text = strings.ToLower(text)
	var tokens []string
	var buf strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else if buf.Len() > 0 {
			tokens = append(tokens, buf.String())
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return tokens
}

func (s *SearchEngine) entryToText(entry LogEntry) string {
	var parts []string
	parts = append(parts, entry.Level)
	if entry.TaskID != "" {
		parts = append(parts, entry.TaskID)
	}
	for _, v := range entry.Event {
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, " ")
}

func (s *SearchEngine) searchBM25(queryTokens []string, limit int) []SearchResult {
	if len(s.index) == 0 {
		return nil
	}

	n := float64(len(s.docStore))
	avgDL := n
	if n > 0 {
		var total int
		for _, docs := range s.index {
			total += len(docs)
		}
		avgDL = float64(total) / n
	}

	scored := make(map[docID]float64)

	for _, token := range queryTokens {
		docIDs := s.index[token]
		df := float64(len(docIDs))
		if df == 0 {
			continue
		}

		idf := math.Log((n - df + 0.5) / (df + 0.5))

		for _, id := range docIDs {
			tf := float64(s.termFreq(id, token))
			k1, b := 1.5, 0.75
			score := idf * (tf * (k1 + 1)) / (tf + k1*(1-b+b*avgDL/avgDL))
			scored[id] += score
		}
	}

	var results []SearchResult
	for id, score := range scored {
		results = append(results, SearchResult{
			Entry:  s.docStore[id],
			Score:  score,
			Method: "bm25",
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func (s *SearchEngine) termFreq(docID docID, term string) int {
	entry := s.docStore[docID]
	text := s.entryToText(entry)
	tokens := s.tokenize(text)

	var count int
	for _, t := range tokens {
		if t == term {
			count++
		}
	}
	return count
}

func (s *SearchEngine) searchEmbedding(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	embProvider, ok := s.provider.(interface {
		GenerateEmbedding(ctx context.Context, text string) ([]float64, error)
	})
	if !ok || embProvider == nil {
		return s.searchBM25(s.tokenize(query), limit), nil
	}

	embedding, err := embProvider.GenerateEmbedding(ctx, query)
	if err != nil || len(embedding) == 0 {
		return s.searchBM25(s.tokenize(query), limit), nil
	}

	var results []SearchResult
	for _, entry := range s.docStore {
		entryText := s.entryToText(entry)
		if entryText == "" {
			continue
		}

		docEmb, err := embProvider.GenerateEmbedding(ctx, entryText)
		if err != nil || len(docEmb) == 0 {
			continue
		}

		score := cosineSimilarity(embedding, docEmb)
		results = append(results, SearchResult{
			Entry:  entry,
			Score:  score,
			Method: "embedding",
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (s *SearchEngine) Index(ctx context.Context, date string) error {
	entries, err := loadDayEvents(date, s.logsDir)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		id := docID(i)
		s.docStore[id] = entry

		tokens := s.tokenize(s.entryToText(entry))
		for _, token := range tokens {
			s.index[token] = append(s.index[token], id)
		}
	}

	return nil
}

func loadDayEvents(date, logsDir string) ([]LogEntry, error) {
	var entries []LogEntry

	uncompressedPath := filepath.Join(logsDir, date+".jsonl")
	if _, err := os.Stat(uncompressedPath); err == nil {
		entries = readLogFile(uncompressedPath, entries)
	}

	compressedPath := filepath.Join(logsDir, date+".jsonl.zst")
	if _, err := os.Stat(compressedPath); err == nil {
		entries = readCompressedFile(compressedPath, entries)
	}

	return entries, nil
}

func readLogFile(path string, entries []LogEntry) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return entries
	}
	defer file.Close()

	for {
		var entry LogEntry
		dec := json.NewDecoder(file)
		if err := dec.Decode(&entry); err != nil {
			break
		}
		entries = append(entries, entry)
	}

	return entries
}

func readCompressedFile(path string, entries []LogEntry) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return entries
	}
	defer file.Close()

	zstdReader, err := zstd.NewReader(file)
	if err != nil {
		return entries
	}
	defer zstdReader.Close()

	for {
		var entry LogEntry
		dec := json.NewDecoder(zstdReader)
		if err := dec.Decode(&entry); err != nil {
			break
		}
		entries = append(entries, entry)
	}

	return entries
}