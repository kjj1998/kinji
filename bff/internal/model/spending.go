package model

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
