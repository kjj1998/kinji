package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/shared"
)

func newTestRepo(t *testing.T) (*Repository, *sql.DB) {
	t.Helper()
	db, err := database.NewClient(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewRepository(db), db
}

func insert(t *testing.T, db *sql.DB, id, date, merchant string, cat shared.Category, amount int, dir shared.Direction) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO transactions (id, user_id, date, merchant, category, amount, direction, notes, split)
		 VALUES (?, 'u1', ?, ?, ?, ?, ?, '', 0)`,
		id, date, merchant, cat, amount, dir)
	if err != nil {
		t.Fatalf("insert %s: %v", id, err)
	}
}

// seed loads a fixed June-2026 (current) / May-2026 (previous) dataset for u1.
func seed(t *testing.T, db *sql.DB) {
	t.Helper()
	// Current month: June 2026
	insert(t, db, "j-inc", "2026-06-01", "Employer", shared.CategoryIncome, 500000, shared.Inflow)
	insert(t, db, "j-f1", "2026-06-05", "Cafe", shared.CategoryFood, 30000, shared.Outflow)
	insert(t, db, "j-f2", "2026-06-06", "Diner", shared.CategoryFood, 20000, shared.Outflow)
	insert(t, db, "j-t1", "2026-06-07", "Metro", shared.CategoryTransport, 10000, shared.Outflow)
	// Previous month: May 2026
	insert(t, db, "m-inc", "2026-05-01", "Employer", shared.CategoryIncome, 450000, shared.Inflow)
	insert(t, db, "m-f1", "2026-05-10", "Cafe", shared.CategoryFood, 15000, shared.Outflow)
	insert(t, db, "m-t1", "2026-05-11", "Metro", shared.CategoryTransport, 25000, shared.Outflow)
}

func TestGetMonthlyTopMerchants(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	got, err := repo.GetMonthlyTopMerchants(context.Background(), "u1", "06", "2026", 5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Three distinct merchants, ranked by spend desc: Cafe 30000, Diner 20000, Metro 10000.
	if len(got) != 3 {
		t.Fatalf("expected 3 merchants, got %d: %+v", len(got), got)
	}
	if got[0].Name != "Cafe" || got[0].Amount != 30000 {
		t.Errorf("top merchant wrong: %+v", got[0])
	}
	if got[2].Name != "Metro" || got[2].Amount != 10000 {
		t.Errorf("last merchant wrong: %+v", got[2])
	}
}

func TestGetMonthlyTopMerchants_RespectsLimit(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	got, err := repo.GetMonthlyTopMerchants(context.Background(), "u1", "06", "2026", 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Cafe" {
		t.Errorf("expected only top merchant, got %+v", got)
	}
}

func TestGetMonthlyTopCategories(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	got, err := repo.GetMonthlyTopCategories(context.Background(), "u1", "06", "2026", 5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 categories, got %d: %+v", len(got), got)
	}
	// Food (30000+20000=50000) ranks above Transport (10000).
	if got[0].Category != shared.CategoryFood || got[0].Amount != 50000 {
		t.Errorf("top category wrong: %+v", got[0])
	}
}

func TestGetTotalIncomeTotalSpentAndNetSavings(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	income, spent, savings, lastMonthSpent, err := repo.GetTotalIncomeTotalSpentAndNetSavings(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// June: income 500000, spent 60000, net 440000.
	if income.Value != 500000 || spent.Value != 60000 || savings.Value != 440000 {
		t.Errorf("June totals wrong: income=%+v spent=%+v savings=%+v", income, spent, savings)
	}
	// Change vs May (income 450000, spent 40000, net 410000).
	if income.Change != 50000 || spent.Change != 20000 || savings.Change != 30000 {
		t.Errorf("changes wrong: income=%d spent=%d savings=%d", income.Change, spent.Change, savings.Change)
	}
	if lastMonthSpent != 40000 {
		t.Errorf("lastMonthSpent = %d, want 40000", lastMonthSpent)
	}
}

func TestGetTotalIncomeTotalSpentAndNetSavings_NoPriorMonth(t *testing.T) {
	repo, db := newTestRepo(t)
	insert(t, db, "j-inc", "2026-06-01", "Employer", shared.CategoryIncome, 500000, shared.Inflow)
	insert(t, db, "j-f1", "2026-06-05", "Cafe", shared.CategoryFood, 30000, shared.Outflow)

	income, spent, _, lastMonthSpent, err := repo.GetTotalIncomeTotalSpentAndNetSavings(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if income.Change != 0 || spent.Change != 0 {
		t.Errorf("expected zero change with no prior month, got income=%d spent=%d", income.Change, spent.Change)
	}
	if lastMonthSpent != 0 {
		t.Errorf("expected lastMonthSpent 0 with no prior month, got %d", lastMonthSpent)
	}
}

func TestGetCategorySpendingForLastTwoMonths(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	cur, prev, err := repo.GetCategorySpendingForLastTwoMonths(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if cur[shared.CategoryFood] != 50000 || cur[shared.CategoryTransport] != 10000 {
		t.Errorf("current spending wrong: %+v", cur)
	}
	if prev[shared.CategoryFood] != 15000 || prev[shared.CategoryTransport] != 25000 {
		t.Errorf("previous spending wrong: %+v", prev)
	}
	// Income must be excluded from spending.
	if _, ok := cur[shared.CategoryIncome]; ok {
		t.Errorf("income should not appear in category spending: %+v", cur)
	}
}

func TestGetLastSixMonthsExpenses(t *testing.T) {
	repo, db := newTestRepo(t)
	seed(t, db)

	got, err := repo.GetLastSixMonthsExpenses(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got["2026-06"] != 60000 {
		t.Errorf("June expenses = %d, want 60000", got["2026-06"])
	}
	if got["2026-05"] != 40000 {
		t.Errorf("May expenses = %d, want 40000", got["2026-05"])
	}
	// Income is an inflow and must not be counted as an expense.
	if got["2026-06"] == 560000 {
		t.Error("income leaked into expenses")
	}
}
