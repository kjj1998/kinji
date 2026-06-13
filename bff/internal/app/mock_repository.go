package app

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/model"
)

// MockRepository is a function-backed test double for TransactionRepository.
type MockRepository struct {
	GetMonthlyTransactionsFn   func(ctx context.Context, userId, month, year string) ([]model.Transaction, error)
	GetMonthlyTopMerchantsFn   func(ctx context.Context, userId, month, year string, limit int) ([]model.MerchantSpending, error)
	GetMonthlyTopCategoriesFn  func(ctx context.Context, userId, month, year string, limit int) ([]model.CategorySpending, error)
	GetTotalsFn                func(ctx context.Context, userId, month, year string) (model.ValueAndChange[int], model.ValueAndChange[int], model.ValueAndChange[int], int, error)
	GetLastSixMonthsExpensesFn func(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingFn      func(ctx context.Context, userId, month, year string) (map[model.Category]int, map[model.Category]int, error)
	SaveTransactionsFn         func(ctx context.Context, userId string, transactions []model.Transaction) error
	GetTransactionPeriodsFn    func(ctx context.Context, userId string) ([]model.Period, error)
}

// compile-time check that MockRepository satisfies the interface.
var _ TransactionRepository = (*MockRepository)(nil)

func (m *MockRepository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]model.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]model.MerchantSpending, error) {
	return m.GetMonthlyTopMerchantsFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]model.CategorySpending, error) {
	return m.GetMonthlyTopCategoriesFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (model.ValueAndChange[int], model.ValueAndChange[int], model.ValueAndChange[int], int, error) {
	return m.GetTotalsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error) {
	return m.GetLastSixMonthsExpensesFn(ctx, userId, month, year)
}

func (m *MockRepository) GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (map[model.Category]int, map[model.Category]int, error) {
	return m.GetCategorySpendingFn(ctx, userId, month, year)
}

func (m *MockRepository) SaveTransactions(ctx context.Context, userId string, transactions []model.Transaction) error {
	return m.SaveTransactionsFn(ctx, userId, transactions)
}

func (m *MockRepository) GetTransactionPeriods(ctx context.Context, userId string) ([]model.Period, error) {
	return m.GetTransactionPeriodsFn(ctx, userId)
}
