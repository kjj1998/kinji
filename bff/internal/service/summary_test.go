package service

import (
	"reflect"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/dto"
	"github.com/kjj1998/kinji/bff/internal/models"
)

func TestGenerateMonthlySummary(t *testing.T) {
	tests := []struct {
		name           string // Descriptive name of the specific scenario
		difference     float64
		topTransaction *models.Transaction
		hasPrevMonth   bool
		netSavings     int
		savingsRate    float64
		expected       string // The output you expect to get
	}{
		{
			name:           "top transaction nil",
			topTransaction: nil,
			expected:       "",
		},
		{
			name:           "hasPrevMonth is false",
			topTransaction: &models.Transaction{Category: models.CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   false,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference > 0",
			difference:     96.45,
			topTransaction: &models.Transaction{Category: models.CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   true,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "You spent 1% more than last month. Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
		{
			name:           "difference < 0",
			difference:     -96.45,
			topTransaction: &models.Transaction{Category: models.CategoryEntertainment, Amount: 50000},
			hasPrevMonth:   true,
			netSavings:     145000,
			savingsRate:    43.4,
			expected:       "You spent 1% less than last month. Your biggest expense was Entertainment at $500.00, and you saved $1450.00 (43.40% of income).",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := generateMonthlySummary(
				tc.difference, tc.topTransaction, tc.hasPrevMonth, tc.netSavings, tc.savingsRate)
			if actual != tc.expected {
				t.Errorf("generateMonthlySummary failed: expected %s, got %s", tc.expected, actual)
			}
		})
	}
}

func TestComputeDailySpendingTrend(t *testing.T) {
	tests := []struct {
		name     string
		txs      []models.Transaction
		expected []dto.DateSpending
	}{
		{
			name: "skips INFLOW transactions",
			txs: []models.Transaction{
				{Date: "2026-06-08", Amount: 1000, Direction: models.Inflow}, // Monday
				{Date: "2026-06-08", Amount: 500, Direction: models.Outflow}, // Monday
			},
			expected: []dto.DateSpending{
				{Date: "Mon", Amount: 500},
				{Date: "Tue", Amount: 0},
				{Date: "Wed", Amount: 0},
				{Date: "Thu", Amount: 0},
				{Date: "Fri", Amount: 0},
				{Date: "Sat", Amount: 0},
				{Date: "Sun", Amount: 0},
			},
		},
		{
			name: "skips unparseable dates",
			txs: []models.Transaction{
				{Date: "not-a-date", Amount: 999, Direction: models.Outflow},
				{Date: "", Amount: 999, Direction: models.Outflow},
				{Date: "2026-06-10", Amount: 300, Direction: models.Outflow}, // Wednesday
			},
			expected: []dto.DateSpending{
				{Date: "Mon", Amount: 0},
				{Date: "Tue", Amount: 0},
				{Date: "Wed", Amount: 300},
				{Date: "Thu", Amount: 0},
				{Date: "Fri", Amount: 0},
				{Date: "Sat", Amount: 0},
				{Date: "Sun", Amount: 0},
			},
		},
		{
			name: "groups by weekday Mon to Sun and sums amounts",
			txs: []models.Transaction{
				{Date: "2026-06-08", Amount: 100, Direction: models.Outflow}, // Monday
				{Date: "2026-06-15", Amount: 200, Direction: models.Outflow}, // Monday (next week)
				{Date: "2026-06-09", Amount: 50, Direction: models.Outflow},  // Tuesday
				{Date: "2026-06-10", Amount: 60, Direction: models.Outflow},  // Wednesday
				{Date: "2026-06-11", Amount: 70, Direction: models.Outflow},  // Thursday
				{Date: "2026-06-12", Amount: 80, Direction: models.Outflow},  // Friday
				{Date: "2026-06-13", Amount: 90, Direction: models.Outflow},  // Saturday
				{Date: "2026-06-14", Amount: 110, Direction: models.Outflow}, // Sunday
			},
			expected: []dto.DateSpending{
				{Date: "Mon", Amount: 300},
				{Date: "Tue", Amount: 50},
				{Date: "Wed", Amount: 60},
				{Date: "Thu", Amount: 70},
				{Date: "Fri", Amount: 80},
				{Date: "Sat", Amount: 90},
				{Date: "Sun", Amount: 110},
			},
		},
		{
			name: "always returns 7 buckets with 3-letter day labels for empty input",
			txs:  nil,
			expected: []dto.DateSpending{
				{Date: "Mon", Amount: 0},
				{Date: "Tue", Amount: 0},
				{Date: "Wed", Amount: 0},
				{Date: "Thu", Amount: 0},
				{Date: "Fri", Amount: 0},
				{Date: "Sat", Amount: 0},
				{Date: "Sun", Amount: 0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := computeDailySpendingTrend(tc.txs)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("computeDailySpendingTrend failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestComputeCategoriesWithBiggestSpendingChange(t *testing.T) {
	tests := []struct {
		name     string
		cur      map[models.Category]int
		prev     map[models.Category]int
		expected []dto.CategorySpendingChange
	}{
		{
			name: "IsNew when prev is zero",
			cur: map[models.Category]int{
				models.CategoryFood:      100,
				models.CategoryTransport: 200,
			},
			prev: map[models.Category]int{
				models.CategoryFood: 50,
			},
			// len(prev) != 0 so sorted by abs(percentageChange): Food 100 > Transport 0.
			expected: []dto.CategorySpendingChange{
				{Category: models.CategoryFood, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: models.CategoryTransport, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "no baseline falls back to sort by amount",
			cur: map[models.Category]int{
				models.CategoryFood:      100,
				models.CategoryTransport: 300,
				models.CategoryShopping:  200,
			},
			prev: map[models.Category]int{},
			// len(prev) == 0: percentageChange is 0 for all, sorted by amount desc.
			expected: []dto.CategorySpendingChange{
				{Category: models.CategoryTransport, Amount: 300, Change: 300, PercentageChange: 0, IsNew: true},
				{Category: models.CategoryShopping, Amount: 200, Change: 200, PercentageChange: 0, IsNew: true},
				{Category: models.CategoryFood, Amount: 100, Change: 100, PercentageChange: 0, IsNew: true},
			},
		},
		{
			name: "with baseline sorts by absolute percentage change",
			cur: map[models.Category]int{
				models.CategoryFood:      10,
				models.CategoryTransport: 100,
				models.CategoryShopping:  100,
			},
			prev: map[models.Category]int{
				models.CategoryFood:      100,
				models.CategoryTransport: 50,
				models.CategoryShopping:  90,
			},
			// abs(pct): Transport 100 > Food 90 (from -90) > Shopping 11.
			expected: []dto.CategorySpendingChange{
				{Category: models.CategoryTransport, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: models.CategoryFood, Amount: 10, Change: -90, PercentageChange: -90, IsNew: false},
				{Category: models.CategoryShopping, Amount: 100, Change: 10, PercentageChange: 11, IsNew: false},
			},
		},
		{
			name: "caps at top 3",
			cur: map[models.Category]int{
				models.CategoryFood:      200,
				models.CategoryTransport: 150,
				models.CategoryShopping:  130,
				models.CategoryHealth:    110,
			},
			prev: map[models.Category]int{
				models.CategoryFood:      100,
				models.CategoryTransport: 100,
				models.CategoryShopping:  100,
				models.CategoryHealth:    100,
			},
			// abs(pct): Food 100, Transport 50, Shopping 30, Health 10 -> top 3 only.
			expected: []dto.CategorySpendingChange{
				{Category: models.CategoryFood, Amount: 200, Change: 100, PercentageChange: 100, IsNew: false},
				{Category: models.CategoryTransport, Amount: 150, Change: 50, PercentageChange: 50, IsNew: false},
				{Category: models.CategoryShopping, Amount: 130, Change: 30, PercentageChange: 30, IsNew: false},
			},
		},
		{
			name: "fewer than 3 returns all",
			cur: map[models.Category]int{
				models.CategoryFood:      100,
				models.CategoryTransport: 50,
			},
			prev: map[models.Category]int{
				models.CategoryFood:      50,
				models.CategoryTransport: 40,
			},
			// abs(pct): Food 100 > Transport 25.
			expected: []dto.CategorySpendingChange{
				{Category: models.CategoryFood, Amount: 100, Change: 50, PercentageChange: 100, IsNew: false},
				{Category: models.CategoryTransport, Amount: 50, Change: 10, PercentageChange: 25, IsNew: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := computeCategoriesWithBiggestSpendingChange(tc.cur, tc.prev)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("computeCategoriesWithBiggestSpendingChange failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestRecentTransactions(t *testing.T) {
	t.Run("nil returns empty slice not nil", func(t *testing.T) {
		actual := recentTransactions(nil, 5)
		if actual == nil {
			t.Fatal("expected non-nil empty slice, got nil")
		}
		if len(actual) != 0 {
			t.Errorf("expected empty slice, got %v", actual)
		}
	})

	t.Run("sorts by date descending", func(t *testing.T) {
		txs := []models.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		expected := []models.Transaction{
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
			{ID: "a", Date: "2026-06-01"},
		}
		actual := recentTransactions(txs, 3)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("n greater than length clamps to all", func(t *testing.T) {
		txs := []models.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
		}
		expected := []models.Transaction{
			{ID: "b", Date: "2026-06-10"},
			{ID: "a", Date: "2026-06-01"},
		}
		actual := recentTransactions(txs, 100)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("does not mutate original slice", func(t *testing.T) {
		txs := []models.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		original := []models.Transaction{
			{ID: "a", Date: "2026-06-01"},
			{ID: "b", Date: "2026-06-10"},
			{ID: "c", Date: "2026-06-05"},
		}
		recentTransactions(txs, 3)
		if !reflect.DeepEqual(txs, original) {
			t.Errorf("original slice was mutated: expected %v, got %v", original, txs)
		}
	})
}

func TestBuildMonthlyTrend(t *testing.T) {
	t.Run("returns 6 buckets ending at target month with Jan-style labels and missing months default to 0", func(t *testing.T) {
		monthlyExpenses := map[string]int{
			"2026-04": 500,
			"2026-06": 1000,
		}
		expected := []dto.DateSpending{
			{Date: "Jan", Amount: 0},
			{Date: "Feb", Amount: 0},
			{Date: "Mar", Amount: 0},
			{Date: "Apr", Amount: 500},
			{Date: "May", Amount: 0},
			{Date: "Jun", Amount: 1000},
		}
		actual, err := buildMonthlyTrend("06", "2026", monthlyExpenses)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("invalid month string returns error", func(t *testing.T) {
		actual, err := buildMonthlyTrend("13", "2026", map[string]int{})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if actual != nil {
			t.Errorf("expected nil result on error, got %v", actual)
		}
	})

	t.Run("invalid year string returns error", func(t *testing.T) {
		actual, err := buildMonthlyTrend("06", "abcd", map[string]int{})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if actual != nil {
			t.Errorf("expected nil result on error, got %v", actual)
		}
	})
}

func TestGetTopTransaction(t *testing.T) {
	tests := []struct {
		name     string
		txs      []models.Transaction
		expected *models.Transaction
	}{
		{
			name:     "empty returns nil",
			txs:      []models.Transaction{},
			expected: nil,
		},
		{
			name: "ignores INFLOW transactions",
			txs: []models.Transaction{
				{ID: "a", Amount: 1000, Direction: models.Inflow},
				{ID: "b", Amount: 200, Direction: models.Outflow},
			},
			expected: &models.Transaction{ID: "b", Amount: 200, Direction: models.Outflow},
		},
		{
			name: "all INFLOW returns nil",
			txs: []models.Transaction{
				{ID: "a", Amount: 1000, Direction: models.Inflow},
				{ID: "b", Amount: 500, Direction: models.Inflow},
			},
			expected: nil,
		},
		{
			name: "picks max amount OUTFLOW",
			txs: []models.Transaction{
				{ID: "a", Amount: 100, Direction: models.Outflow},
				{ID: "b", Amount: 900, Direction: models.Outflow},
				{ID: "c", Amount: 500, Direction: models.Outflow},
				{ID: "d", Amount: 2000, Direction: models.Inflow},
			},
			expected: &models.Transaction{ID: "b", Amount: 900, Direction: models.Outflow},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getTopTransaction(tc.txs)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("getTopTransaction failed: expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
