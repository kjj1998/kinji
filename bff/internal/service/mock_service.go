package service

import (
	"context"
	"mime/multipart"

	"github.com/kjj1998/kinji/bff/internal/dto"
	"github.com/kjj1998/kinji/bff/internal/models"
)

type MockSummaryService struct {
	GenerateMonthlySummaryFn func(ctx context.Context, userId, month, year string) (*dto.TransactionSummary, error)
}

func (m *MockSummaryService) GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*dto.TransactionSummary, error) {
	return m.GenerateMonthlySummaryFn(ctx, userId, month, year)
}

var _ SummaryService = (*MockSummaryService)(nil)

type MockTransactionService struct {
	GetMonthlyTransactionsFn func(ctx context.Context, userId, month, year string) ([]models.Transaction, error)
	ImportStatementFn        func(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]models.Transaction, error)
	SaveTransactionsFn       func(ctx context.Context, userId string, transactions []models.Transaction) ([]models.Transaction, error)
	GetPeriodsFn             func(ctx context.Context, userId string) ([]models.Period, error)
}

func (m *MockTransactionService) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]models.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockTransactionService) ImportStatement(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]models.Transaction, error) {
	return m.ImportStatementFn(ctx, userId, statement, password, onProgress)
}

func (m *MockTransactionService) SaveTransactions(ctx context.Context, userId string, transactions []models.Transaction) ([]models.Transaction, error) {
	return m.SaveTransactionsFn(ctx, userId, transactions)
}
func (m *MockTransactionService) GetPeriods(ctx context.Context, userId string) ([]models.Period, error) {
	return m.GetPeriodsFn(ctx, userId)
}
