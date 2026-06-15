package shared

import "testing"

func TestParseMonth_Valid(t *testing.T) {
	m, err := ParseMonth("06", "2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Start().Format("2006-01-02"); got != "2026-06-01" {
		t.Errorf("Start = %s, want 2026-06-01", got)
	}
}

func TestParseMonth_Invalid(t *testing.T) {
	for _, tc := range []struct{ month, year string }{
		{"13", "2026"},
		{"00", "2026"},
		{"ab", "2026"},
		{"06", "20x6"},
	} {
		if _, err := ParseMonth(tc.month, tc.year); err == nil {
			t.Errorf("ParseMonth(%q, %q): expected error", tc.month, tc.year)
		}
	}
}

func TestMonthEnd(t *testing.T) {
	cases := []struct {
		month, year string
		wantEnd     string
	}{
		{"01", "2026", "2026-01-31"},
		{"02", "2026", "2026-02-28"}, // non-leap
		{"02", "2024", "2024-02-29"}, // leap
		{"04", "2026", "2026-04-30"},
		{"12", "2026", "2026-12-31"},
	}
	for _, tc := range cases {
		m, err := ParseMonth(tc.month, tc.year)
		if err != nil {
			t.Fatalf("ParseMonth(%q,%q): %v", tc.month, tc.year, err)
		}
		if got := m.End().Format("2006-01-02"); got != tc.wantEnd {
			t.Errorf("End(%s-%s) = %s, want %s", tc.year, tc.month, got, tc.wantEnd)
		}
	}
}

func TestMonthRange(t *testing.T) {
	m, _ := ParseMonth("02", "2024")
	start, end := m.Range()
	if start.Format("2006-01-02") != "2024-02-01" || end.Format("2006-01-02") != "2024-02-29" {
		t.Errorf("Range = (%s, %s)", start.Format("2006-01-02"), end.Format("2006-01-02"))
	}
}

func TestAddMonths(t *testing.T) {
	m, _ := ParseMonth("06", "2026")
	if got := m.AddMonths(2).Key(); got != "2026-08" {
		t.Errorf("AddMonths(2).Key = %s, want 2026-08", got)
	}
	// negative crosses the year boundary
	jan, _ := ParseMonth("01", "2026")
	if got := jan.AddMonths(-1).Key(); got != "2025-12" {
		t.Errorf("AddMonths(-1).Key = %s, want 2025-12", got)
	}
}

func TestPrevious(t *testing.T) {
	jan, _ := ParseMonth("01", "2026")
	if got := jan.Previous().Key(); got != "2025-12" {
		t.Errorf("Previous().Key = %s, want 2025-12", got)
	}
}

func TestKey(t *testing.T) {
	m, _ := ParseMonth("06", "2026")
	if got := m.Key(); got != "2026-06" {
		t.Errorf("Key = %s, want 2026-06", got)
	}
}

func TestGetMonthRangeDateStrings(t *testing.T) {
	start, end := GetMonthRangeDateStrings("02", "2024")
	if start != "2024-02-01" || end != "2024-02-29" {
		t.Errorf("GetMonthRangeDateStrings = (%q, %q), want (2024-02-01, 2024-02-29)", start, end)
	}
}
