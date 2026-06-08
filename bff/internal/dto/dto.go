package dto

import "github.com/kjj1998/kinji/bff/internal/models"

type ValueAndChange[T int | float64] struct {
	Value  T `json:"value"`
	Change T `json:"change,omitempty"`
}

func NewValueAndChange[T int | float64](values []T) ValueAndChange[T] {
	if len(values) == 0 {
		return ValueAndChange[T]{}
	}
	v := ValueAndChange[T]{Value: values[0]}
	if len(values) > 1 {
		v.Change = values[0] - values[1]
	}

	return v
}

type Merchant struct {
	Name     string          `json:"name"`
	Amount   int             `json:"amount"`
	Category models.Category `json:"category"`
}

type CategorySpending struct {
	Category models.Category `json:"category"`
	Amount   int             `json:"amount"`
}

type DateSpending struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CategorySpendingChange struct {
	Category         models.Category `json:"category"`
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
	RecentTransactions []models.Transaction     `json:"recentTransactions"`
}
