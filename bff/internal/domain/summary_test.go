package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestNarrative(t *testing.T) {
	var c SummaryCalculator
	tests := []struct {
		name           string
		difference     float64
		topTransaction *Transaction
		hasPrevMonth   bool
		netSavings     int
		savingsRate    float64
		expected       string
	}{
		{
			name:           "top transaction nil",
			topTransaction: nil,
			expected:       "",
		},
		{
			name:           "hasPrevMonth is false",
			topTransaction: &Transaction{Category: CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   false,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference > 0",
			difference:     96.45,
			topTransaction: &Transaction{Category: CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   true,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "You spent 1% more than last month. Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference < 0",
			difference:     -96.45,
			topTransaction: &Transaction{Category: CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   true,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "You spent 1% less than last month. Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := c.Narrative(tc.difference, tc.topTransaction, tc.hasPrevMonth, tc.netSavings, tc.savingsRate)
			if actual != tc.expected {
				t.Errorf("Narrative failed: expected %s, got %s", tc.expected, actual)
			}
		})
	}
}

func TestDailySpendingTrend(t *testing.T) {
	var c SummaryCalculator
	zeroWeek := func(set map[time.Weekday]int) []DaySpending {
		week := []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday,
			time.Thursday, time.Friday, time.Saturday, time.Sunday,
		}
		out := make([]DaySpending, len(week))
		for i, d := range week {
			out[i] = DaySpending{Weekday: d, Amount: set[d]}
		}
		return out
	}

	tests := []struct {
		name     string
		txs      []Transaction
		expected []DaySpending
	}{
		{
			name: "skips INFLOW transactions",
			txs: []Transaction{
				{Date: "2026-06-08", Amount: 1000, Direction: Inflow},  // Monday
				{Date: "2026-06-08", Amount: 500, Direction: Outflow},  // Monday
			},
			expected: zeroWeek(map[time.Weekday]int{time.Monday: 500}),
		},
		{
			name: "skips unparseable dates",
			txs: []Transaction{
				{Date: "not-a-date", Amount: 999, Direction: Outflow},
				{Date: "", Amount: 999, Direction: Outflow},
				{Date: "2026-06-10", Amount: 300, Direction: Outflow}, // Wednesday
			},
			expected: zeroWeek(map[time.Weekday]int{time.Wednesday: 300}),
		},
		{
			name: "groups by weekday Mon to Sun and sums amounts",
			txs: []Transaction{
				{Date: "2026-06-08", Amount: 100, Direction: Outflow}, // Monday
				{Date: "2026-06-15", Amount: 200, Direction: Outflow}, // Monday (next week)
				{Date: "2026-06-09", Amount: 50, Direction: Outflow},  // Tuesday
				{Date: "2026-06-10", Amount: 60, Direction: Outflow},  // Wednesday
				{Date: "2026-06-11", Amount: 70, Direction: Outflow},  // Thursday
				{Date: "2026-06-12", Amount: 80, Direction: Outflow},  // Friday
				{Date: "2026-06-13", Amount: 90, Direction: Outflow},  // Saturday
				{Date: "2026-06-14", Amount: 110, Direction: Outflow}, // Sunday
			},
			expected: zeroWeek(map[time.Weekday]int{
				time.Monday: 300, time.Tuesday: 50, time.Wednesday: 60,
				time.Thursday: 70, time.Friday: 80, time.Saturday: 90, time.Sunday: 110,
			}),
		},
		{
			name:     "always returns 7 buckets for empty input",
			txs:      nil,
			expected: zeroWeek(map[time.Weekday]int{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := c.DailySpendingTrend(tc.txs)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("DailySpendingTrend failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestCategorySpendingChanges(t *testing.T) {
	var c SummaryCalculator
	tests := []struct {
		name     string
		cur      map[Category]int
		prev     map[Category]int
		expected []CategorySpendingChange
	}{
		{
			name: "IsNew when prev is zero",
			cur: map[Category]int{
				CategoryFood:      100,
				CategoryTransport: 200,
			},
			prev: map[Category]int{
				CategoryFood: 50,
			},
			expected: []CategorySpendingChange{
				{Category: CategoryFood, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: CategoryTransport, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "no baseline falls back to sort by amount",
			cur: map[Category]int{
				CategoryFood:      100,
				CategoryTransport: 300,
				CategoryShopping:  200,
			},
			prev: map[Category]int{},
			expected: []CategorySpendingChange{
				{Category: CategoryTransport, Amount: 300, Change: 300, PercentageChange: 0, IsNew: true},
				{Category: CategoryShopping, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
				{Category: CategoryFood, Amount: 100, Change: 100, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "with baseline sorts by absolute percentage change",
			cur: map[Category]int{
				CategoryFood:      10,
				CategoryTransport: 100,
				CategoryShopping:  100,
			},
			prev: map[Category]int{
				CategoryFood:      100,
				CategoryTransport: 50,
				CategoryShopping:  90,
			},
			expected: []CategorySpendingChange{
				{Category: CategoryTransport, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: CategoryFood, Amount: 10, Change: -90, PercentageChange: -90, IsNew: false},
				{Category: CategoryShopping, Amount: 100, Change: 10, PercentageChange: 11, IsNew: false},
			},
		},
		{
			name: "returns full ranked list (no truncation in domain)",
			cur: map[Category]int{
				CategoryFood:      200,
				CategoryTransport: 150,
				CategoryShopping:  130,
				CategoryHealth:    110,
			},
			prev: map[Category]int{
				CategoryFood:      100,
				CategoryTransport: 100,
				CategoryShopping:  100,
				CategoryHealth:    100,
			},
			expected: []CategorySpendingChange{
				{Category: CategoryFood, Amount: 200, Change: 100, PercentageChange: 100, IsNew: false},
				{Category: CategoryTransport, Amount: 150, Change: 50, PercentageChange: 50, IsNew: false},
				{Category: CategoryShopping, Amount: 130, Change: 30, PercentageChange: 30, IsNew: false},
				{Category: CategoryHealth, Amount: 110, Change: 10, PercentageChange: 10, IsNew: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := c.CategorySpendingChanges(tc.cur, tc.prev)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("CategorySpendingChanges failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestRecentTransactions(t *testing.T) {
	var c SummaryCalculator

	t.Run("nil returns empty slice not nil", func(t *testing.T) {
		actual := c.RecentTransactions(nil)
		if actual == nil {
			t.Fatal("expected non-nil empty slice, got nil")
		}
		if len(actual) != 0 {
			t.Errorf("expected empty slice, got %v", actual)
		}
	})

	t.Run("sorts by date descending", func(t *testing.T) {
		txs := []Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		expected := []Transaction{
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
			{ID: "a", Date: "2026-06-01"},
		}
		actual := c.RecentTransactions(txs)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("does not mutate original slice", func(t *testing.T) {
		txs := []Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		original := []Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		c.RecentTransactions(txs)
		if !reflect.DeepEqual(txs, original) {
			t.Errorf("original slice was mutated: expected %v, got %v", original, txs)
		}
	})
}

func TestMonthlyTrend(t *testing.T) {
	var c SummaryCalculator

	t.Run("returns 6 month buckets ending at target month, missing months default to 0", func(t *testing.T) {
		monthlyExpenses := map[string]int{
			"2026-04": 500,
			"2026-06": 1000,
		}
		firstOf := func(m time.Month) time.Time {
			return time.Date(2026, m, 1, 0, 0, 0, 0, time.UTC)
		}
		expected := []MonthSpending{
			{Month: firstOf(time.January), Amount: 0},
			{Month: firstOf(time.February), Amount: 0},
			{Month: firstOf(time.March), Amount: 0},
			{Month: firstOf(time.April), Amount: 500},
			{Month: firstOf(time.May), Amount: 0},
			{Month: firstOf(time.June), Amount: 1000},
		}
		actual, err := c.MonthlyTrend("06", "2026", monthlyExpenses)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("invalid month string returns error", func(t *testing.T) {
		actual, err := c.MonthlyTrend("13", "2026", map[string]int{})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if actual != nil {
			t.Errorf("expected nil result on error, got %v", actual)
		}
	})

	t.Run("invalid year string returns error", func(t *testing.T) {
		actual, err := c.MonthlyTrend("06", "abcd", map[string]int{})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if actual != nil {
			t.Errorf("expected nil result on error, got %v", actual)
		}
	})
}

func TestTopOutflow(t *testing.T) {
	var c SummaryCalculator
	tests := []struct {
		name     string
		txs      []Transaction
		expected *Transaction
	}{
		{
			name:     "empty returns nil",
			txs:      []Transaction{},
			expected: nil,
		},
		{
			name: "ignores INFLOW transactions",
			txs: []Transaction{
				{ID: "a", Amount: 1000, Direction: Inflow},
				{ID: "b", Amount: 200, Direction: Outflow},
			},
			expected: &Transaction{ID: "b", Amount: 200, Direction: Outflow},
		},
		{
			name: "all INFLOW returns nil",
			txs: []Transaction{
				{ID: "a", Amount: 1000, Direction: Inflow},
				{ID: "b", Amount: 500, Direction: Inflow},
			},
			expected: nil,
		},
		{
			name: "picks max amount OUTFLOW",
			txs: []Transaction{
				{ID: "a", Amount: 100, Direction: Outflow},
				{ID: "b", Amount: 900, Direction: Outflow},
				{ID: "c", Amount: 500, Direction: Outflow},
				{ID: "d", Amount: 2000, Direction: Inflow},
			},
			expected: &Transaction{ID: "b", Amount: 900, Direction: Outflow},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := c.TopOutflow(tc.txs)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("TopOutflow failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}