package store

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
)

func newTestRepo(t *testing.T) *Repository {
	t.Helper()
	db, err := database.NewClient(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewRepository(db)
}

func txn(id, date string) shared.Transaction {
	return shared.Transaction{
		ID: id, UserID: "u1", Date: date, Merchant: "Acme",
		Category: shared.CategoryFood, Amount: 1000, Direction: shared.Outflow,
	}
}

func TestSaveAndGetMonthlyTransactions(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.SaveTransactions(ctx, "u1", []shared.Transaction{txn("a", "2026-06-10"), txn("b", "2026-06-20")}); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := repo.GetMonthlyTransactions(ctx, "u1", "06", "2026")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	// round-trips all fields
	if got[0].ID != "b" || got[0].Merchant != "Acme" || got[0].Category != shared.CategoryFood ||
		got[0].Amount != 1000 || got[0].Direction != shared.Outflow {
		t.Errorf("fields not round-tripped: %+v", got[0])
	}
}

func TestGetMonthlyTransactions_FiltersByMonthAndUser(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.SaveTransactions(ctx, "u1", []shared.Transaction{
		txn("jun", "2026-06-15"),
		txn("jul", "2026-07-01"),
	}); err != nil {
		t.Fatalf("save: %v", err)
	}
	other := txn("other", "2026-06-15")
	other.UserID = "u2"
	if err := repo.SaveTransactions(ctx, "u2", []shared.Transaction{other}); err != nil {
		t.Fatalf("save other: %v", err)
	}

	got, err := repo.GetMonthlyTransactions(ctx, "u1", "06", "2026")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got) != 1 || got[0].ID != "jun" {
		t.Errorf("expected only u1's June txn, got %+v", got)
	}
}

func TestSaveTransactions_RollsBackOnFailure(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	// Two rows sharing a primary key -> the second Exec fails, the whole tx rolls back.
	err := repo.SaveTransactions(ctx, "u1", []shared.Transaction{txn("dup", "2026-06-01"), txn("dup", "2026-06-02")})
	if err == nil {
		t.Fatal("expected error on duplicate primary key")
	}

	got, err := repo.GetMonthlyTransactions(ctx, "u1", "06", "2026")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected rollback (0 rows), got %d", len(got))
	}
}

func TestGetTransactionPeriods(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	if err := repo.SaveTransactions(ctx, "u1", []shared.Transaction{
		txn("a", "2025-12-05"),
		txn("b", "2026-01-10"),
		txn("c", "2026-06-20"),
		txn("d", "2026-06-25"), // same period as c -> deduped into one month entry
	}); err != nil {
		t.Fatalf("save: %v", err)
	}

	periods, err := repo.GetTransactionPeriods(ctx, "u1")
	if err != nil {
		t.Fatalf("periods: %v", err)
	}

	if len(periods) != 2 {
		t.Fatalf("expected 2 years, got %d: %+v", len(periods), periods)
	}
	if periods[0].Year != 2025 || len(periods[0].Months) != 1 || periods[0].Months[0] != 12 {
		t.Errorf("2025 period wrong: %+v", periods[0])
	}
	if periods[1].Year != 2026 || len(periods[1].Months) != 2 ||
		periods[1].Months[0] != 1 || periods[1].Months[1] != 6 {
		t.Errorf("2026 period wrong: %+v", periods[1])
	}
}
