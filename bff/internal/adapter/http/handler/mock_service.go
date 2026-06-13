package handler

import (
	"context"
	"mime/multipart"

	"github.com/kjj1998/kinji/bff/internal/app"
	"github.com/kjj1998/kinji/bff/internal/model"
)

// MockSummaryService is a function-backed test double for app.SummaryService.
type MockSummaryService struct {
	GenerateMonthlySummaryFn func(ctx context.Context, userId, month, year string) (*model.MonthlySummary, error)
}

var _ app.SummaryService = (*MockSummaryService)(nil)

func (m *MockSummaryService) GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*model.MonthlySummary, error) {
	return m.GenerateMonthlySummaryFn(ctx, userId, month, year)
}

// MockTransactionService is a function-backed test double for app.TransactionService.
type MockTransactionService struct {
	GetMonthlyTransactionsFn func(ctx context.Context, userId, month, year string) ([]model.Transaction, error)
	ImportStatementFn        func(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]model.Transaction, error)
	SaveTransactionsFn       func(ctx context.Context, userId string, transactions []model.Transaction) ([]model.Transaction, error)
	GetPeriodsFn             func(ctx context.Context, userId string) ([]model.Period, error)
}

var _ app.TransactionService = (*MockTransactionService)(nil)

func (m *MockTransactionService) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]model.Transaction, error) {
	return m.GetMonthlyTransactionsFn(ctx, userId, month, year)
}

func (m *MockTransactionService) ImportStatement(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]model.Transaction, error) {
	return m.ImportStatementFn(ctx, userId, statement, password, onProgress)
}

func (m *MockTransactionService) SaveTransactions(ctx context.Context, userId string, transactions []model.Transaction) ([]model.Transaction, error) {
	return m.SaveTransactionsFn(ctx, userId, transactions)
}

func (m *MockTransactionService) GetPeriods(ctx context.Context, userId string) ([]model.Period, error) {
	return m.GetPeriodsFn(ctx, userId)
}
