package model

type Category string

const (
	CategoryEntertainment Category = "Entertainment"
	CategoryFood          Category = "Food"
	CategoryGroceries     Category = "Groceries"
	CategoryHealth        Category = "Health"
	CategoryIncome        Category = "Income"
	CategoryShopping      Category = "Shopping"
	CategorySubscriptions Category = "Subscriptions"
	CategoryTransport     Category = "Transport"
	CategoryUtilities     Category = "Utilities"
)

type Transaction struct {
	ID       string   `json:"id"`
	UserID   string   `json:"userId"`
	Date     string   `json:"date"`
	Merchant string   `json:"merchant"`
	Category Category `json:"category"`
	Amount   float64  `json:"amount"`
	Notes    *string  `json:"notes"`
	Split    *float64 `json:"split,omitempty"`
}

type ValueAndChange struct {
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type CategorySpending struct {
	Category Category `json:"category"`
	Amount   float64  `json:"amount"`
}

type DateSpending struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type CategorySpendingChange struct {
	Category         Category `json:"category"`
	Amount           float64  `json:"amount"`
	Change           float64  `json:"change"`
	PercentageChange int      `json:"percentageChange"`
}

type Merchant struct {
	Name     string   `json:"name"`
	Amount   float64  `json:"amount"`
	Category Category `json:"category"`
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
