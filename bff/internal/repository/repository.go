package repository

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/models"
)

type Repository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) (models.Transactions, error)
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]models.Merchant, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]models.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(
		ctx context.Context,
		userId, month, year string,
	) (
		models.ValueAndChange[int],
		models.ValueAndChange[int],
		models.ValueAndChange[int],
		int,
		error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(
		ctx context.Context,
		userId, month, year string,
	) (map[models.Category]int, map[models.Category]int, error)
}
