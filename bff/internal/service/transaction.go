package service

import (
	"context"

	"github.com/kohjunjie/kinji/bff/internal/model"
	"github.com/kohjunjie/kinji/bff/internal/repository"
)

type TransactionService interface {
	GetMonthlyTransactions(
		ctx context.Context,
		userId string,
		month string,
		year string,
	) ([]model.Transaction, error)

	GetTransactionsAvailabilities(
		ctx context.Context,
		userId string,
	) ([]model.TransactionsAvailability, error)
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
) ([]model.Transaction, error) {
	return t.repo.List(ctx, userId, month, year)
}

func (t *transactionService) GetTransactionsAvailabilities(
	ctx context.Context,
	userId string,
) ([]model.TransactionsAvailability, error) {
	return nil, nil
}
