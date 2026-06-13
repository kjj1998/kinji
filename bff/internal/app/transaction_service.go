package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/kjj1998/kinji/bff/internal/model"
)

// TransactionService is the use-case API for working with a user's transactions.
type TransactionService interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]model.Transaction, error)
	ImportStatement(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]model.Transaction, error)
	SaveTransactions(ctx context.Context, userId string, transactions []model.Transaction) ([]model.Transaction, error)
	GetPeriods(ctx context.Context, userId string) ([]model.Period, error)
}

type transactionService struct {
	repo   TransactionRepository
	parser StatementParser
}

func NewTransactionService(repo TransactionRepository, parser StatementParser) TransactionService {
	return &transactionService{repo: repo, parser: parser}
}

func (t *transactionService) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]model.Transaction, error) {
	slog.DebugContext(ctx, "get monthly transactions for", "userId", userId, "month", month, "year", year)
	return t.repo.GetMonthlyTransactions(ctx, userId, month, year)
}

// ImportStatement reads an uploaded statement PDF, extracts its rows via the
// parser, validates the running-balance invariant, and stamps each resulting
// transaction with the user id and a fresh id. PDF decryption/validation lives in
// the parser adapter; balance reconciliation lives in model.Statement.
func (t *transactionService) ImportStatement(
	ctx context.Context,
	userId string,
	statement multipart.File,
	password string,
	onProgress func(stage string),
) ([]model.Transaction, error) {
	pdfBytes, err := io.ReadAll(statement)
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}
	onProgress("uploaded")

	lines, err := t.parser.Extract(ctx, pdfBytes, password, onProgress)
	if err != nil {
		slog.ErrorContext(ctx, "statement extraction failed", "err", err)
		return nil, fmt.Errorf("extract statement: %w", err)
	}

	onProgress("checking_balances")
	statementAgg := model.NewStatement(lines)
	if err := statementAgg.Validate(); err != nil {
		return nil, fmt.Errorf("validate statement: %w", err)
	}

	txns := statementAgg.Transactions()
	for i := range txns {
		txns[i].UserID = userId
		txns[i].ID = uuid.Must(uuid.NewV7()).String()
	}
	return txns, nil
}

func (t *transactionService) SaveTransactions(
	ctx context.Context,
	userId string,
	transactions []model.Transaction,
) ([]model.Transaction, error) {
	slog.Info(fmt.Sprintf("saving %d transactions into the database", len(transactions)))

	if err := t.repo.SaveTransactions(ctx, userId, transactions); err != nil {
		return nil, fmt.Errorf("saving transactions: %w", err)
	}

	slog.Info(fmt.Sprintf("saved %d transactions into the database", len(transactions)))
	return transactions, nil
}

func (t *transactionService) GetPeriods(ctx context.Context, userId string) ([]model.Period, error) {
	slog.InfoContext(ctx, "getting transaction periods for userId", "userId", userId)
	return t.repo.GetTransactionPeriods(ctx, userId)
}
