package model

import (
	"errors"
	"testing"
)

// lines builds a small statement: an opening row, an outflow, and an inflow,
// with balances that reconcile by default.
func validLines() []StatementLine {
	return []StatementLine{
		{Txn: Transaction{Merchant: "OPENING", Amount: 0, Direction: Inflow}, Balance: 10000},
		{Txn: Transaction{Merchant: "COFFEE", Amount: 500, Direction: Outflow}, Balance: 9500},
		{Txn: Transaction{Merchant: "SALARY", Amount: 200000, Direction: Inflow}, Balance: 209500},
	}
}

func TestStatementValidate(t *testing.T) {
	t.Run("reconciling balances pass", func(t *testing.T) {
		if err := NewStatement(validLines()).Validate(); err != nil {
			t.Fatalf("expected valid statement, got %v", err)
		}
	})

	t.Run("single row always valid", func(t *testing.T) {
		lines := []StatementLine{
			{Txn: Transaction{Merchant: "ONLY", Amount: 500, Direction: Outflow}, Balance: 9500},
		}
		if err := NewStatement(lines).Validate(); err != nil {
			t.Fatalf("expected valid statement, got %v", err)
		}
	})

	t.Run("mismatched balance returns ErrBalanceMismatch", func(t *testing.T) {
		lines := validLines()
		lines[2].Balance = 999999 // tamper

		err := NewStatement(lines).Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrBalanceMismatch) {
			t.Errorf("expected ErrBalanceMismatch, got %v", err)
		}
	})
}

func TestStatementTransactions(t *testing.T) {
	txns := NewStatement(validLines()).Transactions()
	if len(txns) != 3 {
		t.Fatalf("expected 3 transactions, got %d", len(txns))
	}
	if txns[1].Merchant != "COFFEE" {
		t.Errorf("expected COFFEE at index 1, got %q", txns[1].Merchant)
	}
}
