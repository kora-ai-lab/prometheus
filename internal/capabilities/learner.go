package capabilities

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

type LearnedPattern struct {
	ID           string    `json:"id"`
	Pattern      string    `json:"pattern"`
	Steps        []string `json:"steps"`
	Keywords    []string  `json:"keywords"`
	SuccessCount int       `json:"success_count"`
	LastUsed    time.Time `json:"last_used"`
	CreatedAt   time.Time `json:"created_at"`
}

type Learner struct {
	patterns []LearnedPattern
}

func NewLearner() *Learner {
	return &Learner{
		patterns: make([]LearnedPattern, 0),
	}
}

func (l *Learner) Learn(taskGoal string, successfulSteps []string) {
	keywords := extractKeywords(taskGoal)

	pattern := LearnedPattern{
		ID:           generateID(),
		Pattern:    keywordsToPattern(keywords),
		Steps:      successfulSteps,
		Keywords:   keywords,
		SuccessCount: 1,
		LastUsed:   time.Now(),
		CreatedAt:  time.Now(),
	}

	l.patterns = append(l.patterns, pattern)
}

func (l *Learner) FindSimilar(goal string) *LearnedPattern {
	keywords := extractKeywords(goal)

	var best *LearnedPattern
	var bestScore int

	for i := range l.patterns {
		p := &l.patterns[i]
		score := keywordOverlap(keywords, p.Keywords)

		if score > bestScore {
			bestScore = score
			best = p
		}
	}

	if bestScore >= 2 {
		best.SuccessCount++
		best.LastUsed = time.Now()
		return best
	}

	return nil
}

func (l *Learner) SuggestNextSteps(goal string) []string {
	similar := l.FindSimilar(goal)
	if similar != nil {
		return similar.Steps
	}
	return nil
}

func (l *Learner) GetPatterns() []LearnedPattern {
	return l.patterns
}

func (l *Learner) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(l.patterns, "", "  ")
}

func (l *Learner) ImportJSON(data []byte) error {
	return json.Unmarshal(data, &l.patterns)
}

func (l *Learner) GetTopPatterns(n int) []LearnedPattern {
	if n > len(l.patterns) {
		n = len(l.patterns)
	}
	sorted := make([]LearnedPattern, len(l.patterns))
	copy(sorted, l.patterns)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SuccessCount > sorted[j].SuccessCount
	})
	return sorted[:n]
}

var stopwords = []string{"the", "a", "an", "buy", "purchase", "get", "me", "this", "on", "to", "for", "can", "you", "please", "i", "want", "need", "have", "from"}

func extractKeywords(goal string) []string {
	var keywords []string
	words := strings.Fields(goal)

	for _, w := range words {
		w = strings.ToLower(w)
		if !contains(stopwords, w) && len(w) > 2 {
			keywords = append(keywords, w)
		}
	}

	return keywords
}

func keywordsToPattern(keywords []string) string {
	if len(keywords) >= 3 {
		return keywords[0] + "_" + keywords[1] + "_" + keywords[2]
	}
	return strings.Join(keywords, "_")
}

func keywordOverlap(a, b []string) int {
	count := 0
	for _, wa := range a {
		for _, wb := range b {
			if wa == wb {
				count++
			}
		}
	}
	return count
}

func generateID() string {
	return time.Now().Format("20060102150405")
}