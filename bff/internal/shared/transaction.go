package shared

import "fmt"

type Category string
type Direction string

const (
	Inflow  Direction = "INFLOW"
	Outflow Direction = "OUTFLOW"
)

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
	CategoryTransfer      Category = "Transfer"
)

type Transaction struct {
	ID        string
	UserID    string
	Date      string
	Merchant  string
	Category  Category
	Amount    int
	Direction Direction
	Notes     string
	Split     int
}

// IsInflow reports whether the transaction is money coming in.
func (t Transaction) IsInflow() bool { return t.Direction == Inflow }

// IsOutflow reports whether the transaction is money going out.
func (t Transaction) IsOutflow() bool { return t.Direction == Outflow }

func (d Direction) IsValid() bool {
	return d == Inflow || d == Outflow
}

func (c Category) IsValid() bool {
	switch c {
	case CategoryEntertainment, CategoryFood, CategoryGroceries,
		CategoryHealth, CategoryIncome, CategoryShopping,
		CategorySubscriptions, CategoryTransport, CategoryUtilities, CategoryCredit, CategoryTransfer:
		return true
	}
	return false
}

// ParseCategory converts a raw string into a Category, returning
// ErrInvalidCategory if it is not one of the known categories.
func ParseCategory(s string) (Category, error) {
	c := Category(s)
	if !c.IsValid() {
		return "", fmt.Errorf("%q: %w", s, ErrInvalidCategory)
	}
	return c, nil
}

// ParseDirection converts a raw string into a Direction, returning
// ErrInvalidDirection if it is neither INFLOW nor OUTFLOW.
func ParseDirection(s string) (Direction, error) {
	d := Direction(s)
	if !d.IsValid() {
		return "", fmt.Errorf("%q: %w", s, ErrInvalidDirection)
	}
	return d, nil
}
