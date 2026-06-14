package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// TransactionRepository is the read port the summary feature needs.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]domain.MerchantSpending, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]domain.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (
		totalIncome domain.ValueAndChange[int],
		totalSpent domain.ValueAndChange[int],
		netSavings domain.ValueAndChange[int],
		lastMonthSpent int,
		err error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (
		current map[shared.Category]int,
		previous map[shared.Category]int,
		err error,
	)
}
