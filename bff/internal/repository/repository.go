package repository

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/model"
)

type Repository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) (model.Transactions, error)
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]model.Merchant, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]model.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(
		ctx context.Context,
		userId, month, year string,
	) (
		model.ValueAndChange[int],
		model.ValueAndChange[int],
		model.ValueAndChange[int],
		int,
		error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(
		ctx context.Context,
		userId, month, year string,
	) (map[model.Category]int, map[model.Category]int, error)
}
