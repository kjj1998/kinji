package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// okRepo returns a MockRepository whose every method succeeds with benign data.
// Individual tests override one field to exercise its error path.
func okRepo() *MockRepository {
	return &MockRepository{
		GetMonthlyTransactionsFn: func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
			return []shared.Transaction{{ID: "t1", Direction: shared.Outflow, Amount: 100}}, nil
		},
		GetMonthlyTopMerchantsFn: func(ctx context.Context, userId, month, year string, limit int) ([]domain.MerchantSpending, error) {
			return []domain.MerchantSpending{{Name: "Acme", Amount: 100, Category: shared.CategoryShopping}}, nil
		},
		GetMonthlyTopCategoriesFn: func(ctx context.Context, userId, month, year string, limit int) ([]domain.CategorySpending, error) {
			return []domain.CategorySpending{{Category: shared.CategoryShopping, Amount: 100}}, nil
		},
		GetTotalsFn: func(ctx context.Context, userId, month, year string) (domain.ValueAndChange[int], domain.ValueAndChange[int], domain.ValueAndChange[int], int, error) {
			return domain.ValueAndChange[int]{Value: 1000}, domain.ValueAndChange[int]{Value: 400}, domain.ValueAndChange[int]{Value: 600}, 350, nil
		},
		GetLastSixMonthsExpensesFn: func(ctx context.Context, userId, month, year string) (map[string]int, error) {
			return map[string]int{"2026-06": 400}, nil
		},
		GetCategorySpendingFn: func(ctx context.Context, userId, month, year string) (map[shared.Category]int, map[shared.Category]int, error) {
			return map[shared.Category]int{shared.CategoryShopping: 100}, map[shared.Category]int{shared.CategoryShopping: 80}, nil
		},
	}
}

func TestGenerateMonthlySummary_HappyPath(t *testing.T) {
	svc := NewSummaryService(okRepo())

	got, err := svc.GenerateMonthlySummary(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil summary")
	}
	if got.LastMonthSpent != 350 {
		t.Errorf("LastMonthSpent not wired into SummaryInput, got %d", got.LastMonthSpent)
	}
}

func TestGenerateMonthlySummary_ErrorPaths(t *testing.T) {
	sentinel := errors.New("repo boom")

	cases := []struct {
		name     string
		breaks   func(r *MockRepository)
		wantWrap string
	}{
		{
			name:     "current month transactions",
			breaks:   func(r *MockRepository) { r.GetMonthlyTransactionsFn = func(context.Context, string, string, string) ([]shared.Transaction, error) { return nil, sentinel } },
			wantWrap: "get current month transactions",
		},
		{
			name: "totals",
			breaks: func(r *MockRepository) {
				r.GetTotalsFn = func(context.Context, string, string, string) (domain.ValueAndChange[int], domain.ValueAndChange[int], domain.ValueAndChange[int], int, error) {
					return domain.ValueAndChange[int]{}, domain.ValueAndChange[int]{}, domain.ValueAndChange[int]{}, 0, sentinel
				}
			},
			wantWrap: "get total income, total spent, net savings",
		},
		{
			name:     "top merchants",
			breaks:   func(r *MockRepository) { r.GetMonthlyTopMerchantsFn = func(context.Context, string, string, string, int) ([]domain.MerchantSpending, error) { return nil, sentinel } },
			wantWrap: "get top merchants",
		},
		{
			name:     "top categories",
			breaks:   func(r *MockRepository) { r.GetMonthlyTopCategoriesFn = func(context.Context, string, string, string, int) ([]domain.CategorySpending, error) { return nil, sentinel } },
			wantWrap: "get top categories",
		},
		{
			name:     "last six months expenses",
			breaks:   func(r *MockRepository) { r.GetLastSixMonthsExpensesFn = func(context.Context, string, string, string) (map[string]int, error) { return nil, sentinel } },
			wantWrap: "get last six month expenses",
		},
		{
			name: "category spending",
			breaks: func(r *MockRepository) {
				r.GetCategorySpendingFn = func(context.Context, string, string, string) (map[shared.Category]int, map[shared.Category]int, error) {
					return nil, nil, sentinel
				}
			},
			wantWrap: "get category spending",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := okRepo()
			tc.breaks(repo)
			svc := NewSummaryService(repo)

			got, err := svc.GenerateMonthlySummary(context.Background(), "u1", "06", "2026")
			if got != nil {
				t.Errorf("expected nil summary on error, got %+v", got)
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantWrap) {
				t.Errorf("expected error containing %q, got %v", tc.wantWrap, err)
			}
			if !errors.Is(err, sentinel) {
				t.Errorf("expected sentinel wrapped, got %v", err)
			}
		})
	}
}
