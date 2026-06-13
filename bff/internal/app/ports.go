package app

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/model"
)

// TransactionRepository is the persistence port owned by the application layer.
// Implementations live in the adapter ring and return domain types only.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]model.Transaction, error)
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
	SaveTransactions(ctx context.Context, userId string, transactions []model.Transaction) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]model.Period, error)
}

// StatementParser is the port for turning a bank-statement PDF into raw extracted
// rows. Implementations handle decryption and extraction; balance validation is a
// domain concern (see model.Statement), so the parser returns rows with their
// running balances and does not reconcile them.
type StatementParser interface {
	Extract(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]model.StatementLine, error)
}
