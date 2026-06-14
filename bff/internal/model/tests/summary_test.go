package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/kjj1998/kinji/bff/internal/model"
)

func TestNarrative(t *testing.T) {
	var c model.SummaryCalculator
	tests := []struct {
		name           string
		difference     float64
		topTransaction *model.Transaction
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
			topTransaction: &model.Transaction{Category: model.CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   false,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference > 0",
			difference:     96.45,
			topTransaction: &model.Transaction{Category: model.CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   true,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "You spent 1% more than last month. Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference < 0",
			difference:     -96.45,
			topTransaction: &model.Transaction{Category: model.CategoryEntertainment, Amount: 50000},
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
	var c model.SummaryCalculator
	zeroWeek := func(set map[time.Weekday]int) []model.DaySpending {
		week := []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday,
			time.Thursday, time.Friday, time.Saturday, time.Sunday,
		}
		out := make([]model.DaySpending, len(week))
		for i, d := range week {
			out[i] = model.DaySpending{Weekday: d, Amount: set[d]}
		}
		return out
	}

	tests := []struct {
		name     string
		txs      []model.Transaction
		expected []model.DaySpending
	}{
		{
			name: "skips INFLOW transactions",
			txs: []model.Transaction{
				{Date: "2026-06-08", Amount: 1000, Direction: model.Inflow}, // Monday
				{Date: "2026-06-08", Amount: 500, Direction: model.Outflow}, // Monday
			},
			expected: zeroWeek(map[time.Weekday]int{time.Monday: 500}),
		},
		{
			name: "skips unparseable dates",
			txs: []model.Transaction{
				{Date: "not-a-date", Amount: 999, Direction: model.Outflow},
				{Date: "", Amount: 999, Direction: model.Outflow},
				{Date: "2026-06-10", Amount: 300, Direction: model.Outflow}, // Wednesday
			},
			expected: zeroWeek(map[time.Weekday]int{time.Wednesday: 300}),
		},
		{
			name: "groups by weekday Mon to Sun and sums amounts",
			txs: []model.Transaction{
				{Date: "2026-06-08", Amount: 100, Direction: model.Outflow}, // Monday
				{Date: "2026-06-15", Amount: 200, Direction: model.Outflow}, // Monday (next week)
				{Date: "2026-06-09", Amount: 50, Direction: model.Outflow},  // Tuesday
				{Date: "2026-06-10", Amount: 60, Direction: model.Outflow},  // Wednesday
				{Date: "2026-06-11", Amount: 70, Direction: model.Outflow},  // Thursday
				{Date: "2026-06-12", Amount: 80, Direction: model.Outflow},  // Friday
				{Date: "2026-06-13", Amount: 90, Direction: model.Outflow},  // Saturday
				{Date: "2026-06-14", Amount: 110, Direction: model.Outflow}, // Sunday
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
	var c model.SummaryCalculator
	tests := []struct {
		name     string
		cur      map[model.Category]int
		prev     map[model.Category]int
		expected []model.CategorySpendingChange
	}{
		{
			name: "IsNew when prev is zero",
			cur: map[model.Category]int{
				model.CategoryFood:      100,
				model.CategoryTransport: 200,
			},
			prev: map[model.Category]int{
				model.CategoryFood: 50,
			},
			expected: []model.CategorySpendingChange{
				{Category: model.CategoryFood, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: model.CategoryTransport, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "no baseline falls back to sort by amount",
			cur: map[model.Category]int{
				model.CategoryFood:      100,
				model.CategoryTransport: 300,
				model.CategoryShopping:  200,
			},
			prev: map[model.Category]int{},
			expected: []model.CategorySpendingChange{
				{Category: model.CategoryTransport, Amount: 300, Change: 300, PercentageChange: 0, IsNew: true},
				{Category: model.CategoryShopping, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
				{Category: model.CategoryFood, Amount: 100, Change: 100, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "with baseline sorts by absolute percentage change",
			cur: map[model.Category]int{
				model.CategoryFood:      10,
				model.CategoryTransport: 100,
				model.CategoryShopping:  100,
			},
			prev: map[model.Category]int{
				model.CategoryFood:      100,
				model.CategoryTransport: 50,
				model.CategoryShopping:  90,
			},
			expected: []model.CategorySpendingChange{
				{Category: model.CategoryTransport, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: model.CategoryFood, Amount: 10, Change: -90, PercentageChange: -90, IsNew: false},
				{Category: model.CategoryShopping, Amount: 100, Change: 10, PercentageChange: 11, IsNew: false},
			},
		},
		{
			name: "returns full ranked list (no truncation in domain)",
			cur: map[model.Category]int{
				model.CategoryFood:      200,
				model.CategoryTransport: 150,
				model.CategoryShopping:  130,
				model.CategoryHealth:    110,
			},
			prev: map[model.Category]int{
				model.CategoryFood:      100,
				model.CategoryTransport: 100,
				model.CategoryShopping:  100,
				model.CategoryHealth:    100,
			},
			expected: []model.CategorySpendingChange{
				{Category: model.CategoryFood, Amount: 200, Change: 100, PercentageChange: 100, IsNew: false},
				{Category: model.CategoryTransport, Amount: 150, Change: 50, PercentageChange: 50, IsNew: false},
				{Category: model.CategoryShopping, Amount: 130, Change: 30, PercentageChange: 30, IsNew: false},
				{Category: model.CategoryHealth, Amount: 110, Change: 10, PercentageChange: 10, IsNew: false},
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
	var c model.SummaryCalculator

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
		txs := []model.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		expected := []model.Transaction{
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
		txs := []model.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		original := []model.Transaction{
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
	var c model.SummaryCalculator

	t.Run("returns 6 month buckets ending at target month, missing months default to 0", func(t *testing.T) {
		monthlyExpenses := map[string]int{
			"2026-04": 500,
			"2026-06": 1000,
		}
		firstOf := func(m time.Month) time.Time {
			return time.Date(2026, m, 1, 0, 0, 0, 0, time.UTC)
		}
		expected := []model.MonthSpending{
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
	var c model.SummaryCalculator
	tests := []struct {
		name     string
		txs      []model.Transaction
		expected *model.Transaction
	}{
		{
			name:     "empty returns nil",
			txs:      []model.Transaction{},
			expected: nil,
		},
		{
			name: "ignores INFLOW transactions",
			txs: []model.Transaction{
				{ID: "a", Amount: 1000, Direction: model.Inflow},
				{ID: "b", Amount: 200, Direction: model.Outflow},
			},
			expected: &model.Transaction{ID: "b", Amount: 200, Direction: model.Outflow},
		},
		{
			name: "all INFLOW returns nil",
			txs: []model.Transaction{
				{ID: "a", Amount: 1000, Direction: model.Inflow},
				{ID: "b", Amount: 500, Direction: model.Inflow},
			},
			expected: nil,
		},
		{
			name: "picks max amount OUTFLOW",
			txs: []model.Transaction{
				{ID: "a", Amount: 100, Direction: model.Outflow},
				{ID: "b", Amount: 900, Direction: model.Outflow},
				{ID: "c", Amount: 500, Direction: model.Outflow},
				{ID: "d", Amount: 2000, Direction: model.Inflow},
			},
			expected: &model.Transaction{ID: "b", Amount: 900, Direction: model.Outflow},
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

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		a, b     int     // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "normal division", a: 4, b: 2, expected: 2.0},
		{name: "divide by 0", a: 4, b: 0, expected: 0.0},
		{name: "negative operands", a: -10, b: 5, expected: -2.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := model.SafeDivide(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("SafeDivide(%d, %d) failed: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
			}
		})
	}
}

func TestPercentageChange(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		a, b     int     // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "increase", a: 45, b: 40, expected: 12.5},
		{name: "decrease", a: 35, b: 40, expected: -12.5},
		{name: "previous equals 0", a: 10, b: 0, expected: 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := model.PercentageChange(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("PercentageChange(%d, %d) failed: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
			}
		})
	}
}

func TestRoundTo2Dp(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		input    float64 // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "half up rounding", input: 45.555, expected: 45.56},
		{name: "values with one dp", input: 45.4, expected: 45.40},
		{name: "already rounded values", input: 45.40, expected: 45.40},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := model.RoundTo2Dp(tc.input)
			if actual != tc.expected {
				t.Errorf("RoundTo2Dp(%f) failed: expected %f, got %f", tc.input, tc.expected, actual)
			}
		})
	}
}

func TestSortByAmountDesc_Struct(t *testing.T) {
	type item struct {
		name   string
		amount int
	}
	got := []item{
		{"a", 10}, {"b", 50}, {"c", 30}, {"d", 50},
	}
	model.SortByAmountDesc(got, func(i item) int { return i.amount })

	// assert non-strictly descending by amount (don't assert order within the two 50s)
	for i := 1; i < len(got); i++ {
		if got[i-1].amount < got[i].amount {
			t.Errorf("not descending at %d: %v", i, got)
		}
	}
}
