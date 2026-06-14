package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/model"
)

// TransactionRepository is the persistence port owned by the application layer.
// Implementations live in the adapter ring and return domain types only.
type TransactionRepository interface {
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]model.MerchantSpending, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]model.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (
		totalIncome model.ValueAndChange[int],
		totalSpent model.ValueAndChange[int],
		netSavings model.ValueAndChange[int],
		lastMonthSpent int,
		err error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (
		current map[model.Category]int,
		previous map[model.Category]int,
		err error,
	)
}
