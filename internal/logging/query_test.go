package logging

import (
	"testing"
	"time"
)

func TestParseTemporalQuery_Today(t *testing.T) {
	result := ParseTemporalQuery("aujourd'hui")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	today := time.Now().Format("2006-01-02")
	if result.StartDate != today {
		t.Errorf("StartDate = %s, want %s", result.StartDate, today)
	}
	if result.EndDate != today {
		t.Errorf("EndDate = %s, want %s", result.EndDate, today)
	}
}

func TestParseTemporalQuery_Yesterday(t *testing.T) {
	result := ParseTemporalQuery("hier")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if result.StartDate != yesterday {
		t.Errorf("StartDate = %s, want %s", result.StartDate, yesterday)
	}
	if result.EndDate != yesterday {
		t.Errorf("EndDate = %s, want %s", result.EndDate, yesterday)
	}
}

func TestParseTemporalQuery_DaysAgo(t *testing.T) {
	result := ParseTemporalQuery("il y a 7 jours")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	expected := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	if result.StartDate != expected {
		t.Errorf("StartDate = %s, want %s", result.StartDate, expected)
	}
	if result.EndDate != expected {
		t.Errorf("EndDate = %s, want %s", result.EndDate, expected)
	}
}

func TestParseTemporalQuery_ThreeDaysAgo(t *testing.T) {
	result := ParseTemporalQuery("il y a 3 jours")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	expected := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	if result.StartDate != expected {
		t.Errorf("StartDate = %s, want %s", result.StartDate, expected)
	}
	if result.EndDate != expected {
		t.Errorf("EndDate = %s, want %s", result.EndDate, expected)
	}
}

func TestParseTemporalQuery_LastMonday(t *testing.T) {
	result := ParseTemporalQuery("lundi dernier")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	now := time.Now()
	weekday := int(now.Weekday())
	mondayOffset := weekday + 7
	if weekday == 0 {
		mondayOffset = 6
	} else {
		mondayOffset = weekday - 1 + 7
	}
	lastMonday := now.AddDate(0, 0, -mondayOffset).Format("2006-01-02")

	if result.StartDate != lastMonday {
		t.Errorf("StartDate = %s, want %s", result.StartDate, lastMonday)
	}
	if result.EndDate != lastMonday {
		t.Errorf("EndDate = %s, want %s", result.EndDate, lastMonday)
	}
}

func TestParseTemporalQuery_LastWeek(t *testing.T) {
	result := ParseTemporalQuery("semaine dernière")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

now := time.Now()
	weekday := int(now.Weekday())
	mondayOffset := weekday + 7
	if weekday == 0 {
		mondayOffset = 6
	} else {
		mondayOffset = weekday - 1 + 7
	}
	lastMonday := now.AddDate(0, 0, -mondayOffset).Format("2006-01-02")

	sundayOffset := mondayOffset - 6
	lastSunday := now.AddDate(0, 0, -sundayOffset).Format("2006-01-02")

	if result.StartDate != lastMonday {
		t.Errorf("StartDate = %s, want %s", result.StartDate, lastMonday)
	}
	if result.EndDate != lastSunday {
		t.Errorf("EndDate = %s, want %s", result.EndDate, lastSunday)
	}
}

func TestParseTemporalQuery_LastMonth(t *testing.T) {
	result := ParseTemporalQuery("mois dernier")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	startDate := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	if result.StartDate != startDate {
		t.Errorf("StartDate = %s, want %s", result.StartDate, startDate)
	}
	if result.EndDate != endDate {
		t.Errorf("EndDate = %s, want %s", result.EndDate, endDate)
	}
}

func TestParseTemporalQuery_WithKeywords(t *testing.T) {
	result := ParseTemporalQuery("aujourd'hui code review")
	if result == nil {
		t.Fatal("ParseTemporalQuery returned nil")
	}

	if len(result.Keywords) != 2 {
		t.Errorf("Keywords = %v, want 2 keywords", result.Keywords)
	}
}