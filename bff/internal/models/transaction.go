package models

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
	CategoryCredit        Category = "Credit"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryEntertainment, CategoryFood, CategoryGroceries,
		CategoryHealth, CategoryIncome, CategoryShopping,
		CategorySubscriptions, CategoryTransport, CategoryUtilities, CategoryCredit:
		return true
	}
	return false
}

type Direction string

const (
	Inflow  Direction = "INFLOW"
	Outflow Direction = "OUTFLOW"
)

func (d Direction) IsValid() bool {
	return d == Inflow || d == Outflow
}

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

type Transactions struct {
	Transactions   []Transaction              `json:"transactions"`
	Availabilities []TransactionsAvailability `json:"availabilities"`
}
