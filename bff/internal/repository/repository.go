package repository

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/dto"
	"github.com/kjj1998/kinji/bff/internal/models"
)

type Repository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]models.Transaction, error)
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]dto.Merchant, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]dto.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(
		ctx context.Context,
		userId, month, year string,
	) (
		dto.ValueAndChange[int],
		dto.ValueAndChange[int],
		dto.ValueAndChange[int],
		int,
		error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(
		ctx context.Context,
		userId, month, year string,
	) (map[models.Category]int, map[models.Category]int, error)
	SaveTransactions(
		ctx context.Context,
		userId string,
		transactions []models.Transaction,
	) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]models.Period, error)
}
