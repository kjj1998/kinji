package service

import (
	"context"
	"log/slog"

	"github.com/kohjunjie/kinji/bff/internal/model"
	"github.com/kohjunjie/kinji/bff/internal/repository"
)

type TransactionService interface {
	GetMonthlyTransactions(
		ctx context.Context,
		userId string,
		month string,
		year string,
	) (model.Transactions, error)
}

type transactionService struct {
	repo repository.Repository
}

func NewTransactionService(repo repository.Repository) TransactionService {
	return &transactionService{repo: repo}
}

func (t *transactionService) GetMonthlyTransactions(
	ctx context.Context,
	userId string,
	month string,
	year string,
) (model.Transactions, error) {
	slog.DebugContext(
		ctx,
		"get monthly transactions for",
		"userId", userId,
		"month", month,
		"year", year,
	)

	return t.repo.GetMonthlyTransactions(ctx, userId, month, year)
}
