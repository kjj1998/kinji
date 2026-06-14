package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// MockRepository is a function-backed test double for TransactionRepository.
type MockRepository struct {
	GetMonthlyTransactionsFn   func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	GetMonthlyTopMerchantsFn   func(ctx context.Context, userId, month, year string, limit int) ([]domain.MerchantSpending, error)
	GetMonthlyTopCategoriesFn  func(ctx context.Context, userId, month, year string, limit int) ([]domain.CategorySpending, error)
	GetTotalsFn                func(ctx context.Context, userId, month, year string) (domain.ValueAndChange[int], domain.ValueAndChange[int], domain.ValueAndChange[int], int, error)
	GetLastSixMonthsExpensesFn func(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingFn      func(ctx context.Context, userId, month, year string) (map[shared.Category]int, map[shared.Category]int, error)
}

// compile-time check that MockRepository satisfies the interface.
var _ TransactionRepository = (*MockRepository)(nil)

func (m *MockRepository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]domain.MerchantSpending, error) {
	return m.GetMonthlyTopMerchantsFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]domain.CategorySpending, error) {
	return m.GetMonthlyTopCategoriesFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (domain.ValueAndChange[int], domain.ValueAndChange[int], domain.ValueAndChange[int], int, error) {
	return m.GetTotalsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error) {
	return m.GetLastSixMonthsExpensesFn(ctx, userId, month, year)
}

func (m *MockRepository) GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (map[shared.Category]int, map[shared.Category]int, error) {
	return m.GetCategorySpendingFn(ctx, userId, month, year)
}
