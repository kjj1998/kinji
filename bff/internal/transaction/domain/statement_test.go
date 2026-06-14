package domain_test

import (
	"errors"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// lines builds a small statement: an opening row, an outflow, and an inflow,
// with balances that reconcile by default.
func validLines() []domain.StatementLine {
	return []domain.StatementLine{
		{Txn: shared.Transaction{Merchant: "OPENING", Amount: 0, Direction: shared.Inflow}, Balance: 10000},
		{Txn: shared.Transaction{Merchant: "COFFEE", Amount: 500, Direction: shared.Outflow}, Balance: 9500},
		{Txn: shared.Transaction{Merchant: "SALARY", Amount: 200000, Direction: shared.Inflow}, Balance: 209500},
	}
}

func TestStatementValidate(t *testing.T) {
	t.Run("reconciling balances pass", func(t *testing.T) {
		if err := domain.NewStatement(validLines()).Validate(); err != nil {
			t.Fatalf("expected valid statement, got %v", err)
		}
	})

	t.Run("single row always valid", func(t *testing.T) {
		lines := []domain.StatementLine{
			{Txn: shared.Transaction{Merchant: "ONLY", Amount: 500, Direction: shared.Outflow}, Balance: 9500},
		}
		if err := domain.NewStatement(lines).Validate(); err != nil {
			t.Fatalf("expected valid statement, got %v", err)
		}
	})

	t.Run("mismatched balance returns ErrBalanceMismatch", func(t *testing.T) {
		lines := validLines()
		lines[2].Balance = 999999 // tamper

		err := domain.NewStatement(lines).Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrBalanceMismatch) {
			t.Errorf("expected ErrBalanceMismatch, got %v", err)
		}
	})
}

func TestStatementTransactions(t *testing.T) {
	txns := domain.NewStatement(validLines()).Transactions()
	if len(txns) != 3 {
		t.Fatalf("expected 3 transactions, got %d", len(txns))
	}
	if txns[1].Merchant != "COFFEE" {
		t.Errorf("expected COFFEE at index 1, got %q", txns[1].Merchant)
	}
}
