package domain

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"time"
)

const (
	dateLayout    = "2006-01-02"
	summaryMonths = 6
)

// MonthlySummary is the read-model for a user's monthly spending overview. Its
// list fields are returned in full, ranked order; the presentation layer applies
// any truncation (e.g. top 3 changes, most-recent 5) and renders display labels.
type MonthlySummary struct {
	TotalIncome        ValueAndChange[int]
	TotalSpent         ValueAndChange[int]
	NetSavings         ValueAndChange[int]
	SavingsRate        float64
	LastMonthSpent     int
	SummaryStatement   string
	TopCategories      []CategorySpending
	MonthlyTrend       []MonthSpending
	DailyTrend         []DaySpending
	BiggestChanges     []CategorySpendingChange
	TopMerchants       []MerchantSpending
	RecentTransactions []Transaction
}

// SummaryInput is the raw, repository-sourced data the calculator needs to build
// a MonthlySummary. The application layer gathers it; the calculator owns the math.
type SummaryInput struct {
	Month, Year          string
	CurrentMonth         []Transaction
	TotalIncome          ValueAndChange[int]
	TotalSpent           ValueAndChange[int]
	NetSavings           ValueAndChange[int]
	LastMonthSpent       int
	TopMerchants         []MerchantSpending
	TopCategories        []CategorySpending
	MonthlyExpenses      map[string]int // keyed "2006-01"
	CurCategorySpending  map[Category]int
	PrevCategorySpending map[Category]int
}

// SummaryCalculator is a stateless domain service that turns raw monthly data
// into a MonthlySummary. All methods are pure functions of their inputs.
type SummaryCalculator struct{}

func NewSummaryCalculator() SummaryCalculator { return SummaryCalculator{} }

// Calculate assembles the full monthly summary from the supplied input.
func (c SummaryCalculator) Calculate(in SummaryInput) (*MonthlySummary, error) {
	savingsRate := roundTo2Dp(safeDivide(in.NetSavings.Value, in.TotalIncome.Value) * 100)

	monthlyTrend, err := c.MonthlyTrend(in.Month, in.Year, in.MonthlyExpenses)
	if err != nil {
		return nil, err
	}

	changeInSpending := float64(in.TotalSpent.Value - in.LastMonthSpent)
	top := c.TopOutflow(in.CurrentMonth)
	narrative := c.Narrative(changeInSpending, top, in.LastMonthSpent > 0, in.NetSavings.Value, savingsRate)

	return &MonthlySummary{
		TotalIncome:        in.TotalIncome,
		TotalSpent:         in.TotalSpent,
		NetSavings:         in.NetSavings,
		SavingsRate:        savingsRate,
		LastMonthSpent:     in.LastMonthSpent,
		SummaryStatement:   narrative,
		TopCategories:      in.TopCategories,
		MonthlyTrend:       monthlyTrend,
		DailyTrend:         c.DailySpendingTrend(in.CurrentMonth),
		BiggestChanges:     c.CategorySpendingChanges(in.CurCategorySpending, in.PrevCategorySpending),
		TopMerchants:       in.TopMerchants,
		RecentTransactions: c.RecentTransactions(in.CurrentMonth),
	}, nil
}

// Narrative builds the one-line spending summary sentence. Returns "" when there
// is no outflow to describe.
func (SummaryCalculator) Narrative(
	difference float64,
	topTransaction *Transaction,
	hasPrevMonth bool,
	netSavings int,
	savingsRate float64,
) string {
	if topTransaction == nil {
		return ""
	}

	suffix := fmt.Sprintf(
		"Your biggest expense was %s at $%.2f, and you saved $%.2f (%.2f%% of income).",
		topTransaction.Category,
		float64(topTransaction.Amount)/100,
		float64(netSavings)/100,
		savingsRate,
	)
	if !hasPrevMonth {
		return suffix
	}
	direction := "more"
	if difference < 0 {
		direction = "less"
	}
	return fmt.Sprintf("You spent %.0f%% %s than last month. ", math.Abs(roundTo2Dp(difference/100)), direction) + suffix
}

// DailySpendingTrend sums outflows into seven Monday-to-Sunday weekday buckets.
// Unparseable dates and inflows are ignored.
func (SummaryCalculator) DailySpendingTrend(txns []Transaction) []DaySpending {
	totals := make(map[time.Weekday]int)
	for _, t := range txns {
		if t.IsInflow() {
			continue
		}
		date, err := time.Parse(dateLayout, t.Date)
		if err != nil {
			continue
		}
		totals[date.Weekday()] += t.Amount
	}

	week := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday,
		time.Thursday, time.Friday, time.Saturday, time.Sunday,
	}
	trend := make([]DaySpending, len(week))
	for i, day := range week {
		trend[i] = DaySpending{Weekday: day, Amount: totals[day]}
	}
	return trend
}

// CategorySpendingChanges ranks every category by how much its spending moved
// between the previous and current period. With a prior baseline it sorts by the
// magnitude of percentage change; with no baseline it falls back to raw amount.
// The full ranked list is returned; the caller decides how many to show.
func (SummaryCalculator) CategorySpendingChanges(cur, prev map[Category]int) []CategorySpendingChange {
	categories := make(map[Category]struct{}, len(cur)+len(prev))
	for cat := range cur {
		categories[cat] = struct{}{}
	}
	for cat := range prev {
		categories[cat] = struct{}{}
	}

	result := make([]CategorySpendingChange, 0, len(categories))
	for cat := range categories {
		curAmount := cur[cat]
		prevAmount := prev[cat]
		result = append(result, CategorySpendingChange{
			Category:         cat,
			Amount:           curAmount,
			Change:           curAmount - prevAmount,
			PercentageChange: int(percentageChange(curAmount, prevAmount)),
			IsNew:            prevAmount == 0,
		})
	}

	if len(prev) == 0 {
		// No baseline: "biggest movers" is undefined, fall back to biggest spenders.
		sortByAmountDesc(result, func(c CategorySpendingChange) int {
			return c.Amount
		})
	} else {
		sortByAmountDesc(result, func(c CategorySpendingChange) int {
			if c.PercentageChange < 0 {
				return -c.PercentageChange
			}
			return c.PercentageChange
		})
	}

	return result
}

// RecentTransactions returns the transactions sorted by date, most recent first.
// The full slice is returned (never nil); the caller decides how many to show.
func (SummaryCalculator) RecentTransactions(txns []Transaction) []Transaction {
	sorted := slices.Clone(txns)
	if sorted == nil {
		sorted = []Transaction{}
	}
	slices.SortFunc(sorted, func(a, b Transaction) int {
		return cmp.Compare(b.Date, a.Date)
	})
	return sorted
}

// MonthlyTrend builds the trailing six-month outflow series ending at the given
// month, keyed by the first instant of each month. Missing months default to 0.
func (SummaryCalculator) MonthlyTrend(month, year string, monthlyExpenses map[string]int) ([]MonthSpending, error) {
	to, err := ParseMonth(month, year)
	if err != nil {
		return nil, err
	}

	trend := make([]MonthSpending, summaryMonths)
	for i := range summaryMonths {
		m := to.AddMonths(-(summaryMonths - 1 - i))
		trend[i] = MonthSpending{
			Month:  m.Start(),
			Amount: monthlyExpenses[m.Key()],
		}
	}
	return trend, nil
}

// TopOutflow returns the single largest outflow, or nil when there are none.
func (SummaryCalculator) TopOutflow(txns []Transaction) *Transaction {
	var top *Transaction
	for i := range txns {
		t := &txns[i]
		if !t.IsOutflow() {
			continue
		}
		if top == nil || t.Amount > top.Amount {
			top = t
		}
	}
	return top
}

// Number constrains the numeric types the spending math operates on.
type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

func percentageChange[T Number](current, previous T) float64 {
	return roundTo2Dp(safeDivide(current-previous, previous) * 100)
}

func roundTo2Dp(value float64) float64 {
	return math.Round(value*100) / 100
}

func safeDivide[T Number](a, b T) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func sortByAmountDesc[T any](s []T, amount func(T) int) {
	slices.SortFunc(s, func(a, b T) int {
		if amount(b) > amount(a) {
			return 1
		}
		if amount(b) < amount(a) {
			return -1
		}
		return 0
	})
}