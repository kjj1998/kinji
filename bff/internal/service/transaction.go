package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/kjj1998/kinji/bff/internal/claude"
	"github.com/kjj1998/kinji/bff/internal/models"
	"github.com/kjj1998/kinji/bff/internal/repository"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// ClientError is a service error caused by bad client input → maps to HTTP 4xx.
type ClientError struct {
	Reason string
}

func (e *ClientError) Error() string { return e.Reason }

type TransactionService interface {
	GetMonthlyTransactions(
		ctx context.Context,
		userId string,
		month string,
		year string,
	) (models.Transactions, error)
	ImportStatement(
		ctx context.Context,
		userId string,
		statement multipart.File,
		password string,
		onProgress func(stage string),
	) ([]models.Transaction, error)
	SaveTransactions(
		ctx context.Context,
		userId string,
		transactions []models.Transaction,
	) ([]models.Transaction, error)
}

type transactionService struct {
	repo   repository.Repository
	parser claude.Parser
}

func NewTransactionService(repo repository.Repository, parser claude.Parser) TransactionService {
	return &transactionService{repo: repo, parser: parser}
}

func (t *transactionService) GetMonthlyTransactions(
	ctx context.Context,
	userId string,
	month string,
	year string,
) (models.Transactions, error) {
	slog.DebugContext(
		ctx,
		"get monthly transactions for",
		"userId", userId,
		"month", month,
		"year", year,
	)

	return t.repo.GetMonthlyTransactions(ctx, userId, month, year)
}

func (t *transactionService) ImportStatement(
	ctx context.Context,
	userId string,
	statement multipart.File,
	password string,
	onProgress func(stage string),
) ([]models.Transaction, error) {
	pdfBytes, err := io.ReadAll(statement)
	if err != nil {
		return []models.Transaction{}, fmt.Errorf("read pdf: %w", err)
	}
	onProgress("uploaded")

	slog.Info("authenticatiing pdf")
	onProgress("validating")
	if password == "" {
		if err := api.Validate(bytes.NewReader(pdfBytes), model.NewDefaultConfiguration()); err != nil {
			if errors.Is(err, pdfcpu.ErrWrongPassword) {
				return []models.Transaction{}, &ClientError{Reason: "pdf password required"}
			}
			return []models.Transaction{}, &ClientError{Reason: "invalid/corrupt pdf file"}
		}
	} else {
		var out bytes.Buffer
		conf := model.NewDefaultConfiguration()
		conf.UserPW = password

		if err := api.Decrypt(bytes.NewReader(pdfBytes), &out, conf); err != nil {
			if errors.Is(err, pdfcpu.ErrWrongPassword) {
				return []models.Transaction{}, &ClientError{Reason: "wrong pdf password given"}
			}
			return []models.Transaction{}, &ClientError{Reason: "invalid/corrupt pdf file"}
		}

		pdfBytes = out.Bytes()
	}

	onProgress("parsing")

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	txns, err := t.parser.ParseStatement(ctx, pdfBytes, onProgress)
	for i := range txns {
		txns[i].UserID = userId
		txns[i].ID = uuid.Must(uuid.NewV7()).String()
	}
	if err != nil {
		slog.ErrorContext(ctx, "claude parse failed", "err", err)
		return nil, fmt.Errorf("parse statement: %w", err)
	}

	return txns, nil
}

func (t *transactionService) SaveTransactions(
	ctx context.Context,
	userId string,
	transactions []models.Transaction,
) ([]models.Transaction, error) {
	slog.Info(fmt.Sprintf("saving %d transactions into the database", len(transactions)))

	err := t.repo.SaveTransactions(ctx, userId, transactions)

	if err != nil {
		return []models.Transaction{}, fmt.Errorf("saving transactions: %w", err)
	}

	slog.Info(fmt.Sprintf("saved %d transactions into the database", len(transactions)))
	return transactions, nil
}
