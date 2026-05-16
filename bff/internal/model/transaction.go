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

type Direction string

const (
	Inflow  Direction = "INFLOW"
	Outflow Direction = "OUTFLOW"
)

type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Date      string    `json:"date"`
	Merchant  string    `json:"merchant"`
	Category  Category  `json:"category"`
	Amount    int       `json:"amount"`
	Direction Direction `json:"direction"`
	Notes     string    `json:"notes,omitempty"`
	Split     int       `json:"split,omitempty"`
}

type TransactionsAvailability struct {
	Year   int   `json:"year"`
	Months []int `json:"months"`
}
