package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
	"github.com/kjj1998/kinji/bff/internal/transaction/service"
)

// Repository is the sqlite-backed implementation of service.TransactionRepository.
type Repository struct {
	client *sql.DB
}

// compile-time check that Repository satisfies the application port.
var _ service.TransactionRepository = (*Repository)(nil)

func NewRepository(client *sql.DB) *Repository {
	return &Repository{client: client}
}

func (d *Repository) getTransactionsWithinDateRange(
	ctx context.Context,
	userId, from, to string,
) ([]shared.Transaction, error) {
	rows, err := d.client.QueryContext(ctx, getAllTransactionsWithinDateRange, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying transactions userId %q from %q to %q: %w", userId, from, to, err)
	}
	defer rows.Close()

	txs := []shared.Transaction{}
	for rows.Next() {
		var t shared.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Date, &t.Merchant,
			&t.Category, &t.Amount, &t.Direction, &t.Notes, &t.Split); err != nil {
			return nil, fmt.Errorf("scanning transaction: %w", err)
		}
		txs = append(txs, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating transaction rows: %w", err)
	}

	return txs, nil
}

func (d *Repository) getPeriods(
	ctx context.Context,
	userId string,
) ([]domain.Period, error) {
	rows, err := d.client.QueryContext(ctx, getMonthAndYearWhichTransactionsOccur, userId)
	if err != nil {
		return nil, fmt.Errorf("getting periods: %w", err)
	}
	defer rows.Close()

	byYear := map[int]*domain.Period{}
	var order []int
	for rows.Next() {
		var y, m int
		if err := rows.Scan(&y, &m); err != nil {
			return nil, fmt.Errorf("scanning periods: %w", err)
		}
		if _, ok := byYear[y]; !ok {
			byYear[y] = &domain.Period{Year: y}
			order = append(order, y)
		}
		byYear[y].Months = append(byYear[y].Months, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating period rows: %w", err)
	}

	out := make([]domain.Period, 0, len(order))
	for _, y := range order {
		out = append(out, *byYear[y])
	}
	return out, nil
}

func (d *Repository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	from, to := shared.GetMonthRangeDateStrings(month, year)
	txs, err := d.getTransactionsWithinDateRange(ctx, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying monthly transactions: %w", err)
	}

	return txs, nil
}

func (d *Repository) SaveTransactions(
	ctx context.Context,
	userId string,
	transactions []shared.Transaction,
) error {
	tx, err := d.client.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db client error: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, saveTransactions)
	if err != nil {
		return fmt.Errorf("db context error: %w", err)
	}
	defer stmt.Close()

	for _, t := range transactions {
		if _, err := stmt.ExecContext(ctx, t.ID, t.UserID, t.Date, t.Merchant,
			t.Category, t.Amount, t.Direction, t.Notes, t.Split); err != nil {
			return fmt.Errorf("db execution context error: %w", err)
		}
	}
	return tx.Commit()
}

func (d *Repository) GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error) {
	periods, err := d.getPeriods(ctx, userId)
	if err != nil {
		return []domain.Period{}, fmt.Errorf("getting periods: %w", err)
	}

	return periods, nil
}
