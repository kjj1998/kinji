package handler

import "github.com/kjj1998/kinji/bff/internal/shared"

// Transaction is the wire representation of a shared.Transaction.
type Transaction struct {
	ID        string           `json:"id"`
	UserID    string           `json:"userId"`
	Date      string           `json:"date"`
	Merchant  string           `json:"merchant"`
	Category  shared.Category  `json:"category"`
	Amount    int              `json:"amount"`
	Direction shared.Direction `json:"direction"`
	Notes     string           `json:"notes,omitempty"`
	Split     int              `json:"split,omitempty"`
}

// ToTransaction maps a domain transaction to its wire representation.
func ToTransaction(t shared.Transaction) Transaction {
	return Transaction{
		ID:        t.ID,
		UserID:    t.UserID,
		Date:      t.Date,
		Merchant:  t.Merchant,
		Category:  t.Category,
		Amount:    t.Amount,
		Direction: t.Direction,
		Notes:     t.Notes,
		Split:     t.Split,
	}
}

// Domain maps a wire transaction back to the domain entity.
func (t Transaction) Domain() shared.Transaction {
	return shared.Transaction{
		ID:        t.ID,
		UserID:    t.UserID,
		Date:      t.Date,
		Merchant:  t.Merchant,
		Category:  t.Category,
		Amount:    t.Amount,
		Direction: t.Direction,
		Notes:     t.Notes,
		Split:     t.Split,
	}
}

// ToTransactions maps a slice of domain transactions to wire representations.
// The result is never nil so it marshals as [] rather than null.
func ToTransactions(txns []shared.Transaction) []Transaction {
	out := make([]Transaction, len(txns))
	for i, t := range txns {
		out[i] = ToTransaction(t)
	}
	return out
}

// DomainTransactions maps a slice of wire transactions back to domain entities.
func DomainTransactions(in []Transaction) []shared.Transaction {
	out := make([]shared.Transaction, len(in))
	for i, t := range in {
		out[i] = t.Domain()
	}
	return out
}
