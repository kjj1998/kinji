package dto

import "github.com/kjj1998/kinji/bff/internal/model"

// Transaction is the wire representation of a model.Transaction.
type Transaction struct {
	ID        string          `json:"id"`
	UserID    string          `json:"userId"`
	Date      string          `json:"date"`
	Merchant  string          `json:"merchant"`
	Category  model.Category  `json:"category"`
	Amount    int             `json:"amount"`
	Direction model.Direction `json:"direction"`
	Notes     string          `json:"notes,omitempty"`
	Split     int             `json:"split,omitempty"`
}

// ToTransaction maps a domain transaction to its wire representation.
func ToTransaction(t model.Transaction) Transaction {
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
func (t Transaction) Domain() model.Transaction {
	return model.Transaction{
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
func ToTransactions(txns []model.Transaction) []Transaction {
	out := make([]Transaction, len(txns))
	for i, t := range txns {
		out[i] = ToTransaction(t)
	}
	return out
}

// DomainTransactions maps a slice of wire transactions back to domain entities.
func DomainTransactions(in []Transaction) []model.Transaction {
	out := make([]model.Transaction, len(in))
	for i, t := range in {
		out[i] = t.Domain()
	}
	return out
}
