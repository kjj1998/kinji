package model

type Merchant struct {
	Name     string   `json:"name"`
	Amount   float64  `json:"amount"`
	Category Category `json:"category"`
}
