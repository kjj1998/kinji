package dto

import (
	"github.com/kjj1998/kinji/bff/internal/shared"
	webdto "github.com/kjj1998/kinji/bff/internal/shared/webdto"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// View policy: how many ranked items the monthly summary screen shows.
const (
	maxBiggestChanges     = 3
	maxRecentTransactions = 5
)

type ValueAndChange[T int | float64] struct {
	Value  T `json:"value"`
	Change T `json:"change,omitempty"`
}

type Merchant struct {
	Name     string          `json:"name"`
	Amount   int             `json:"amount"`
	Category shared.Category `json:"category"`
}

type CategorySpending struct {
	Category shared.Category `json:"category"`
	Amount   int             `json:"amount"`
}

type DateSpending struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CategorySpendingChange struct {
	Category         shared.Category `json:"category"`
	Amount           int             `json:"amount"`
	Change           int             `json:"change"`
	PercentageChange int             `json:"percentageChange"`
	IsNew            bool            `json:"isNew"`
}

type TransactionSummary struct {
	TotalIncome        ValueAndChange[int]      `json:"totalIncome"`
	TotalSpent         ValueAndChange[int]      `json:"totalSpent"`
	NetSavings         ValueAndChange[int]      `json:"netSavings"`
	SavingsRate        float64                  `json:"savingsRate"`
	LastMonthSpent     int                      `json:"lastMonthSpent"`
	MonthlySummary     string                   `json:"monthlySummary"`
	TopCategories      []CategorySpending       `json:"topCategories"`
	MonthlyTrend       []DateSpending           `json:"monthlyExpenses"`
	DailyTrend         []DateSpending           `json:"dailyTrend"`
	BiggestChanges     []CategorySpendingChange `json:"biggestChanges"`
	TopMerchants       []Merchant               `json:"topMerchants"`
	RecentTransactions []webdto.Transaction     `json:"recentTransactions"`
}

// ToTransactionSummary maps the domain read-model to its wire representation,
// rendering display labels (weekday/month) and applying the view truncation
// (top biggest changes, most-recent transactions).
func ToTransactionSummary(s *domain.MonthlySummary) *TransactionSummary {
	if s == nil {
		return nil
	}
	return &TransactionSummary{
		TotalIncome:        toValueAndChange(s.TotalIncome),
		TotalSpent:         toValueAndChange(s.TotalSpent),
		NetSavings:         toValueAndChange(s.NetSavings),
		SavingsRate:        s.SavingsRate,
		LastMonthSpent:     s.LastMonthSpent,
		MonthlySummary:     s.SummaryStatement,
		TopCategories:      toCategorySpendings(s.TopCategories),
		MonthlyTrend:       toMonthlyTrend(s.MonthlyTrend),
		DailyTrend:         toDailyTrend(s.DailyTrend),
		BiggestChanges:     toCategorySpendingChanges(capSlice(s.BiggestChanges, maxBiggestChanges)),
		TopMerchants:       toMerchants(s.TopMerchants),
		RecentTransactions: webdto.ToTransactions(capSlice(s.RecentTransactions, maxRecentTransactions)),
	}
}

func toValueAndChange(v domain.ValueAndChange[int]) ValueAndChange[int] {
	return ValueAndChange[int]{Value: v.Value, Change: v.Change}
}

func toCategorySpendings(in []domain.CategorySpending) []CategorySpending {
	out := make([]CategorySpending, len(in))
	for i, c := range in {
		out[i] = CategorySpending{Category: c.Category, Amount: c.Amount}
	}
	return out
}

func toMerchants(in []domain.MerchantSpending) []Merchant {
	out := make([]Merchant, len(in))
	for i, m := range in {
		out[i] = Merchant{Name: m.Name, Amount: m.Amount, Category: m.Category}
	}
	return out
}

func toCategorySpendingChanges(in []domain.CategorySpendingChange) []CategorySpendingChange {
	out := make([]CategorySpendingChange, len(in))
	for i, c := range in {
		out[i] = CategorySpendingChange{
			Category:         c.Category,
			Amount:           c.Amount,
			Change:           c.Change,
			PercentageChange: c.PercentageChange,
			IsNew:            c.IsNew,
		}
	}
	return out
}

// toDailyTrend renders each weekday bucket as a three-letter label (e.g. "Mon").
func toDailyTrend(in []domain.DaySpending) []DateSpending {
	out := make([]DateSpending, len(in))
	for i, d := range in {
		out[i] = DateSpending{Date: d.Weekday.String()[:3], Amount: d.Amount}
	}
	return out
}

// toMonthlyTrend renders each month bucket as a short month label (e.g. "Jan").
func toMonthlyTrend(in []domain.MonthSpending) []DateSpending {
	out := make([]DateSpending, len(in))
	for i, m := range in {
		out[i] = DateSpending{Date: m.Month.Format("Jan"), Amount: m.Amount}
	}
	return out
}

// capSlice returns at most n leading elements of s.
func capSlice[T any](s []T, n int) []T {
	if len(s) > n {
		return s[:n]
	}
	return s
}
