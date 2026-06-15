// Package store holds transaction persistence reads shared by feature stores.
// Both the transaction and summary slices need "transactions within a date range";
// it lives here so the query and scan exist once without coupling the two features
// to each other.
package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
)

// TransactionsInDateRange returns a user's transactions with date in [from, to],
// most recent first.
func TransactionsInDateRange(ctx context.Context, db *sql.DB, userId, from, to string) ([]shared.Transaction, error) {
	txs, err := database.QueryRows(ctx, db, transactionsInDateRange, scanTransaction, userId, from, to)
	if err != nil {
		return nil, fmt.Errorf("querying transactions userId %q from %q to %q: %w", userId, from, to, err)
	}
	return txs, nil
}

func scanTransaction(r *sql.Rows) (shared.Transaction, error) {
	var t shared.Transaction
	err := r.Scan(&t.ID, &t.UserID, &t.Date, &t.Merchant,
		&t.Category, &t.Amount, &t.Direction, &t.Notes, &t.Split)
	return t, err
}
