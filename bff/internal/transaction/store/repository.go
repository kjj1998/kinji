package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
	sharedstore "github.com/kjj1998/kinji/bff/internal/shared/store"
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

func (d *Repository) getPeriods(
	ctx context.Context,
	userId string,
) ([]domain.Period, error) {
	type yearMonth struct{ year, month int }
	rows, err := database.QueryRows(ctx, d.client, getMonthAndYearWhichTransactionsOccur,
		func(r *sql.Rows) (yearMonth, error) {
			var ym yearMonth
			return ym, r.Scan(&ym.year, &ym.month)
		}, userId)
	if err != nil {
		return nil, fmt.Errorf("getting periods: %w", err)
	}

	byYear := map[int]*domain.Period{}
	var order []int
	for _, ym := range rows {
		if _, ok := byYear[ym.year]; !ok {
			byYear[ym.year] = &domain.Period{Year: ym.year}
			order = append(order, ym.year)
		}
		byYear[ym.year].Months = append(byYear[ym.year].Months, ym.month)
	}

	out := make([]domain.Period, 0, len(order))
	for _, y := range order {
		out = append(out, *byYear[y])
	}
	return out, nil
}

func (d *Repository) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	from, to := shared.GetMonthRangeDateStrings(month, year)
	return sharedstore.TransactionsInDateRange(ctx, d.client, userId, from, to)
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
