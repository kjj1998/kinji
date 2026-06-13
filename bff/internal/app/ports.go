package app

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/domain"
)

// TransactionRepository is the persistence port owned by the application layer.
// Implementations live in the adapter ring and return domain types only.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]domain.Transaction, error)
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
		current map[domain.Category]int,
		previous map[domain.Category]int,
		err error,
	)
	SaveTransactions(ctx context.Context, userId string, transactions []domain.Transaction) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error)
}

// StatementParser is the port for turning a bank-statement PDF into raw extracted
// rows. Implementations handle decryption and extraction; balance validation is a
// domain concern (see domain.Statement), so the parser returns rows with their
// running balances and does not reconcile them.
type StatementParser interface {
	Extract(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error)
}
