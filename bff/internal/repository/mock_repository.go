package repository

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/dto"
	"github.com/kjj1998/kinji/bff/internal/models"
)

type MockRepository struct {
	GetMonthlyTransactionsFn   func(ctx context.Context, userId, month, year string) ([]models.Transaction, error)
	GetMonthlyTopMerchantsFn   func(ctx context.Context, userId, month, year string, limit int) ([]dto.Merchant, error)
	GetMonthlyTopCategoriesFn  func(ctx context.Context, userId, month, year string, limit int) ([]dto.CategorySpending, error)
	GetTotalsFn                func(ctx context.Context, userId, month, year string) (dto.ValueAndChange[int], dto.ValueAndChange[int], dto.ValueAndChange[int], int, error)
	GetLastSixMonthsExpensesFn func(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingFn      func(ctx context.Context, userId, month, year string) (map[models.Category]int, map[models.Category]int, error)
	SaveTransactionsFn         func(ctx context.Context, userId string, transactions []models.Transaction) error
	GetTransactionPeriodsFn    func(ctx context.Context, userId string) ([]models.Period, error)
}

// compile-time check that MockRepository satisfies the interface
var _ Repository = (*MockRepository)(nil)

func (m *MockRepository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]models.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]dto.Merchant, error) {
	return m.GetMonthlyTopMerchantsFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]dto.CategorySpending, error) {
	return m.GetMonthlyTopCategoriesFn(ctx, userId, month, year, limit)
}

func (m *MockRepository) GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (dto.ValueAndChange[int], dto.ValueAndChange[int], dto.ValueAndChange[int], int, error) {
	return m.GetTotalsFn(ctx, userId, month, year)
}

func (m *MockRepository) GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error) {
	return m.GetLastSixMonthsExpensesFn(ctx, userId, month, year)
}

func (m *MockRepository) GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (map[models.Category]int, map[models.Category]int, error) {
	return m.GetCategorySpendingFn(ctx, userId, month, year)
}

func (m *MockRepository) SaveTransactions(ctx context.Context, userId string, transactions []models.Transaction) error {
	return m.SaveTransactionsFn(ctx, userId, transactions)
}

func (m *MockRepository) GetTransactionPeriods(ctx context.Context, userId string) ([]models.Period, error) {
	return m.GetTransactionPeriodsFn(ctx, userId)
}
