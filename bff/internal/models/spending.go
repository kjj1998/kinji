package models

type CategorySpending struct {
	Category Category `json:"category"`
	Amount   int      `json:"amount"`
}

type DateSpending struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type MonthlyExpense struct {
	Month  string `json:"month"`
	Amount int    `json:"amount"`
}

type CategorySpendingChange struct {
	Category         Category `json:"category"`
	Amount           int      `json:"amount"`
	Change           int      `json:"change"`
	PercentageChange int      `json:"percentageChange"`
	IsNew            bool     `json:"isNew"`
}
