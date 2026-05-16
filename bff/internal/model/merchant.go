package model

type Merchant struct {
	Name     string   `json:"name"`
	Amount   int      `json:"amount"`
	Category Category `json:"category"`
}
