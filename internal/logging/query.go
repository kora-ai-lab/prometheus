package logging

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type TimeQuery struct {
	StartDate string
	EndDate  string
	Keywords []string
}

func ParseTemporalQuery(q string) *TimeQuery {
	if q == "" {
		return &TimeQuery{}
	}

	q = strings.ToLower(strings.TrimSpace(q))

	today := time.Now().UTC()
	result := &TimeQuery{}

	switch {
	case strings.Contains(q, "aujourd'hui"):
		result.StartDate = today.Format("2006-01-02")
		result.EndDate = result.StartDate

	case strings.Contains(q, "hier"):
		yesterday := today.AddDate(0, 0, -1)
		result.StartDate = yesterday.Format("2006-01-02")
		result.EndDate = result.StartDate

	case matchesLastMonday(q):
		result = parseLastMonday(today)

	case strings.Contains(q, "semaine dernière"):
		result = parseLastWeek(today)

	case strings.Contains(q, "mois dernier"):
		result = parseLastMonth(today)

	case matchesDaysAgo(q):
		result = parseDaysAgo(q, today)
	}

	result.Keywords = extractKeywords(q, result.StartDate, result.EndDate)

	return result
}

func matchesLastMonday(q string) bool {
	days := []string{"lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi", "dimanche"}
	for _, day := range days {
		if strings.Contains(q, day+" dernier") || strings.Contains(q, "dernier "+day) {
			return true
		}
	}
	return false
}

func parseLastMonday(now time.Time) *TimeQuery {
	weekday := int(now.Weekday())
	var mondayOffset int
	if weekday == 0 {
		mondayOffset = 6
	} else {
		mondayOffset = weekday - 1
	}
	mondayOffset += 7
	lastMonday := now.AddDate(0, 0, -mondayOffset)
	return &TimeQuery{
		StartDate: lastMonday.Format("2006-01-02"),
		EndDate:   lastMonday.Format("2006-01-02"),
	}
}

func parseLastWeek(now time.Time) *TimeQuery {
	weekday := int(now.Weekday())
	var mondayOffset int
	if weekday == 0 {
		mondayOffset = 6
	} else {
		mondayOffset = weekday - 1
	}
	mondayOffset += 7
	lastMonday := now.AddDate(0, 0, -mondayOffset)
	lastSunday := lastMonday.AddDate(0, 0, 6)
	return &TimeQuery{
		StartDate: lastMonday.Format("2006-01-02"),
		EndDate:   lastSunday.Format("2006-01-02"),
	}
}

func parseLastMonth(now time.Time) *TimeQuery {
	lastMonth := now.AddDate(0, -1, 0)
	startDate := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC)
	return &TimeQuery{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}
}

var daysAgoRegex = regexp.MustCompile(`il y a (\d+) jours`)

func matchesDaysAgo(q string) bool {
	return daysAgoRegex.MatchString(q)
}

func parseDaysAgo(q string, now time.Time) *TimeQuery {
	matches := daysAgoRegex.FindStringSubmatch(q)
	if len(matches) < 2 {
		return &TimeQuery{}
	}
	var n int
	if _, err := fmt.Sscanf(matches[1], "%d", &n); err != nil {
		return &TimeQuery{}
	}
	date := now.AddDate(0, 0, -n)
	dateStr := date.Format("2006-01-02")
	return &TimeQuery{
		StartDate: dateStr,
		EndDate:   dateStr,
	}
}

func extractKeywords(q, startDate, endDate string) []string {
	var keywords []string
	patterns := []string{
		"aujourd'hui",
		"hier",
		"il y a",
		"jours",
		"semaine dernière",
		"mois dernier",
		"lundi dernier",
		"mardi dernier",
		"mercredi dernier",
		"jeudi dernier",
		"vendredi dernier",
		"samedi dernier",
		"dimanche dernier",
	}

	testQ := q
	for _, p := range patterns {
		testQ = strings.ReplaceAll(testQ, p, "")
	}

	testQ = strings.ReplaceAll(testQ, "dernier", "")
	testQ = strings.ReplaceAll(testQ, "il", "")
	testQ = strings.ReplaceAll(testQ, "y", "")
	testQ = strings.ReplaceAll(testQ, "a", "")

	words := strings.Fields(testQ)
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w != "" {
			keywords = append(keywords, w)
		}
	}

	return keywords
}