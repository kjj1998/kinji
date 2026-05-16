package model

type ValueAndChange[T int | float64] struct {
	Value  T       `json:"value"`
	Change float64 `json:"change"`
}

type TransactionSummary struct {
	TotalIncome        int                      `json:"totalIncome"`
	TotalSpent         ValueAndChange[int]      `json:"totalSpent"`
	NetSavings         ValueAndChange[int]      `json:"netSavings"`
	SavingsRate        ValueAndChange[float64]  `json:"savingsRate"`
	LastMonthSpent     int                      `json:"lastMonthSpent"`
	TopCategory        *CategorySpending        `json:"topCategory"`
	MonthlySummary     string                   `json:"monthlySummary"`
	SpendingByCategory []CategorySpending       `json:"spendingByCategory"`
	MonthlyTrend       []DateSpending           `json:"monthlyTrend"`
	DailyTrend         []DateSpending           `json:"dailyTrend"`
	BiggestChanges     []CategorySpendingChange `json:"biggestChanges"`
	TopMerchants       []Merchant               `json:"topMerchants"`
	RecentTransactions []Transaction            `json:"recentTransactions"`
}
