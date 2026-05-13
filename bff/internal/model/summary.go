package model

type ValueAndChange struct {
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type TransactionSummary struct {
	TotalIncome        float64                  `json:"totalIncome"`
	TotalSpent         ValueAndChange           `json:"totalSpent"`
	NetSavings         ValueAndChange           `json:"netSavings"`
	SavingsRate        ValueAndChange           `json:"savingsRate"`
	LastMonthSpent     float64                  `json:"lastMonthSpent"`
	TopCategory        *CategorySpending        `json:"topCategory"`
	MonthlySummary     string                   `json:"monthlySummary"`
	SpendingByCategory []CategorySpending       `json:"spendingByCategory"`
	MonthlyTrend       []DateSpending           `json:"monthlyTrend"`
	DailyTrend         []DateSpending           `json:"dailyTrend"`
	BiggestChanges     []CategorySpendingChange `json:"biggestChanges"`
	TopMerchants       []Merchant               `json:"topMerchants"`
	RecentTransactions []Transaction            `json:"recentTransactions"`
}
