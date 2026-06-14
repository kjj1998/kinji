package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// MockRepository is a function-backed test double for TransactionRepository.
type MockRepository struct {
	GetMonthlyTransactionsFn func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	SaveTransactionsFn       func(ctx context.Context, userId string, transactions []shared.Transaction) error
	GetTransactionPeriodsFn  func(ctx context.Context, userId string) ([]domain.Period, error)
}

// compile-time check that MockRepository satisfies the interface.
var _ TransactionRepository = (*MockRepository)(nil)

func (m *MockRepository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockRepository) SaveTransactions(ctx context.Context, userId string, transactions []shared.Transaction) error {
	return m.SaveTransactionsFn(ctx, userId, transactions)
}

func (m *MockRepository) GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error) {
	return m.GetTransactionPeriodsFn(ctx, userId)
}
