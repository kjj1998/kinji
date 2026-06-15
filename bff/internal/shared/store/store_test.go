package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.NewClient(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func insertTxn(t *testing.T, db *sql.DB, id, userId, date string, dir shared.Direction) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO transactions (id, user_id, date, merchant, category, amount, direction, notes, split)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, userId, date, "Merchant", shared.CategoryFood, 1000, dir, "", 0)
	if err != nil {
		t.Fatalf("insert txn: %v", err)
	}
}

func TestTransactionsInDateRange(t *testing.T) {
	db := newTestDB(t)
	insertTxn(t, db, "a", "u1", "2026-06-05", shared.Outflow)
	insertTxn(t, db, "b", "u1", "2026-06-20", shared.Outflow)
	insertTxn(t, db, "c", "u1", "2026-05-31", shared.Outflow) // before range
	insertTxn(t, db, "d", "u1", "2026-07-01", shared.Outflow) // after range
	insertTxn(t, db, "e", "u2", "2026-06-10", shared.Outflow) // other user

	got, err := TransactionsInDateRange(context.Background(), db, "u1", "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only u1's in-range rows, most recent first.
	if len(got) != 2 {
		t.Fatalf("expected 2 rows, got %d: %+v", len(got), got)
	}
	if got[0].ID != "b" || got[1].ID != "a" {
		t.Errorf("expected [b, a] ordered date DESC, got [%s, %s]", got[0].ID, got[1].ID)
	}
}

func TestTransactionsInDateRange_EmptyIsNonNil(t *testing.T) {
	db := newTestDB(t)

	got, err := TransactionsInDateRange(context.Background(), db, "nobody", "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %d", len(got))
	}
}

func TestTransactionsInDateRange_BoundsInclusive(t *testing.T) {
	db := newTestDB(t)
	insertTxn(t, db, "start", "u1", "2026-06-01", shared.Outflow)
	insertTxn(t, db, "end", "u1", "2026-06-30", shared.Outflow)

	got, err := TransactionsInDateRange(context.Background(), db, "u1", "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected both boundary rows included, got %d", len(got))
	}
}
