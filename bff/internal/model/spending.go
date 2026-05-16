package model

type CategorySpending struct {
	Category Category `json:"category"`
	Amount   int      `json:"amount"`
}

type DateSpending struct {
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CategorySpendingChange struct {
	Category         Category `json:"category"`
	Amount           int      `json:"amount"`
	Change           int      `json:"change"`
	PercentageChange int      `json:"percentageChange"`
}
