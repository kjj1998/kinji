package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// TransactionRepository is the persistence port for the transaction feature.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	SaveTransactions(ctx context.Context, userId string, transactions []shared.Transaction) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error)
}

// StatementParser is the port for turning a bank-statement PDF into raw extracted
// rows. Implementations handle decryption and extraction; balance validation is a
// domain concern (see domain.Statement), so the parser returns rows with their
// running balances and does not reconcile them.
type StatementParser interface {
	Extract(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error)
}
