package repository

import (
	"context"

	"github.com/kohjunjie/kinji/bff/internal/model"
)

type Repository interface {
	List(ctx context.Context, userID string, month string, year string) ([]model.Transaction, error)
	ListRange(ctx context.Context, userID, from, to string) ([]model.Transaction, error)
	Create(ctx context.Context, tx model.Transaction) error
	Update(ctx context.Context, tx model.Transaction) error
	Delete(ctx context.Context, userID, id string) error
}
