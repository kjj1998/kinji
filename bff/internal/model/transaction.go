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
