package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// errReader is a multipart.File-shaped reader whose Read always fails.
type errReader struct{}

func (errReader) Read([]byte) (int, error)          { return 0, errors.New("boom") }
func (errReader) ReadAt([]byte, int64) (int, error) { return 0, errors.New("boom") }
func (errReader) Seek(int64, int) (int64, error)    { return 0, errors.New("boom") }
func (errReader) Close() error                      { return nil }

// fakeFile adapts a *strings.Reader to the multipart.File interface.
type fakeFile struct{ *strings.Reader }

func (fakeFile) Close() error { return nil }

func newFakeFile(s string) fakeFile { return fakeFile{strings.NewReader(s)} }

func TestGetMonthlyTransactions(t *testing.T) {
	want := []shared.Transaction{{ID: "t1", Merchant: "Acme"}}
	repo := &MockRepository{
		GetMonthlyTransactionsFn: func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
			if userId != "u1" || month != "06" || year != "2026" {
				t.Errorf("unexpected args: %q %q %q", userId, month, year)
			}
			return want, nil
		},
	}
	svc := NewTransactionService(repo, &MockParser{})

	got, err := svc.GetMonthlyTransactions(context.Background(), "u1", "06", "2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "t1" {
		t.Errorf("result not delegated, got %+v", got)
	}
}

func TestGetMonthlyTransactions_Error(t *testing.T) {
	sentinel := errors.New("db down")
	repo := &MockRepository{
		GetMonthlyTransactionsFn: func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
			return nil, sentinel
		},
	}
	svc := NewTransactionService(repo, &MockParser{})

	_, err := svc.GetMonthlyTransactions(context.Background(), "u1", "06", "2026")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel propagated, got %v", err)
	}
}

func TestSaveTransactions(t *testing.T) {
	in := []shared.Transaction{{ID: "t1"}, {ID: "t2"}}
	repo := &MockRepository{
		SaveTransactionsFn: func(ctx context.Context, userId string, transactions []shared.Transaction) error {
			if userId != "u1" || len(transactions) != 2 {
				t.Errorf("unexpected args: %q %+v", userId, transactions)
			}
			return nil
		},
	}
	svc := NewTransactionService(repo, &MockParser{})

	got, err := svc.SaveTransactions(context.Background(), "u1", in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected saved slice returned, got %+v", got)
	}
}

func TestSaveTransactions_Error(t *testing.T) {
	repo := &MockRepository{
		SaveTransactionsFn: func(ctx context.Context, userId string, transactions []shared.Transaction) error {
			return errors.New("write failed")
		},
	}
	svc := NewTransactionService(repo, &MockParser{})

	got, err := svc.SaveTransactions(context.Background(), "u1", []shared.Transaction{{ID: "t1"}})
	if err == nil || !strings.Contains(err.Error(), "saving transactions") {
		t.Errorf("expected wrapped saving error, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil result on error, got %+v", got)
	}
}

func TestGetPeriods(t *testing.T) {
	want := []domain.Period{{Year: 2026, Months: []int{1, 6}}}
	sentinel := errors.New("period fail")

	repo := &MockRepository{
		GetTransactionPeriodsFn: func(ctx context.Context, userId string) ([]domain.Period, error) {
			return want, nil
		},
	}
	svc := NewTransactionService(repo, &MockParser{})
	got, err := svc.GetPeriods(context.Background(), "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Year != 2026 {
		t.Errorf("periods not delegated, got %+v", got)
	}

	repo.GetTransactionPeriodsFn = func(ctx context.Context, userId string) ([]domain.Period, error) {
		return nil, sentinel
	}
	if _, err := svc.GetPeriods(context.Background(), "u1"); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel propagated, got %v", err)
	}
}

func TestImportStatement_ReadError(t *testing.T) {
	svc := NewTransactionService(&MockRepository{}, &MockParser{})

	_, err := svc.ImportStatement(context.Background(), "u1", errReader{}, "", func(string) {})
	if err == nil || !strings.Contains(err.Error(), "read pdf") {
		t.Errorf("expected read pdf error, got %v", err)
	}
}

func TestImportStatement_ParserError(t *testing.T) {
	parser := &MockParser{
		ExtractFn: func(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error) {
			return nil, errors.New("llm down")
		},
	}
	svc := NewTransactionService(&MockRepository{}, parser)

	_, err := svc.ImportStatement(context.Background(), "u1", newFakeFile("pdf"), "", func(string) {})
	if err == nil || !strings.Contains(err.Error(), "extract statement") {
		t.Errorf("expected extract statement error, got %v", err)
	}
}

func TestImportStatement_BalanceMismatch(t *testing.T) {
	parser := &MockParser{
		ExtractFn: func(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error) {
			return []domain.StatementLine{
				{Txn: shared.Transaction{Merchant: "open", Direction: shared.Outflow, Amount: 0}, Balance: 1000},
				{Txn: shared.Transaction{Merchant: "coffee", Direction: shared.Outflow, Amount: 500}, Balance: 999}, // should be 500
			}, nil
		},
	}
	svc := NewTransactionService(&MockRepository{}, parser)

	_, err := svc.ImportStatement(context.Background(), "u1", newFakeFile("pdf"), "", func(string) {})
	if err == nil || !strings.Contains(err.Error(), "validate statement") {
		t.Errorf("expected validate statement error, got %v", err)
	}
	if !errors.Is(err, domain.ErrBalanceMismatch) {
		t.Errorf("expected ErrBalanceMismatch wrapped, got %v", err)
	}
}

func TestImportStatement_HappyPath(t *testing.T) {
	parser := &MockParser{
		ExtractFn: func(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error) {
			onProgress("validating")
			onProgress("parsing")
			return []domain.StatementLine{
				{Txn: shared.Transaction{Merchant: "open", Direction: shared.Outflow, Amount: 0}, Balance: 1000},
				{Txn: shared.Transaction{Merchant: "coffee", Direction: shared.Outflow, Amount: 500}, Balance: 500},
				{Txn: shared.Transaction{Merchant: "salary", Direction: shared.Inflow, Amount: 2000}, Balance: 2500},
			}, nil
		},
	}
	svc := NewTransactionService(&MockRepository{}, parser)

	var stages []string
	got, err := svc.ImportStatement(context.Background(), "u1", newFakeFile("pdf"), "", func(s string) {
		stages = append(stages, s)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 transactions, got %d", len(got))
	}
	for i, tx := range got {
		if tx.UserID != "u1" {
			t.Errorf("txn %d: UserID not stamped, got %q", i, tx.UserID)
		}
		if tx.ID == "" {
			t.Errorf("txn %d: ID not stamped", i)
		}
	}
	if got[0].ID == got[1].ID {
		t.Errorf("expected distinct ids, got %q twice", got[0].ID)
	}
	assertContains(t, stages, "uploaded")
	assertContains(t, stages, "checking_balances")
}

func assertContains(t *testing.T, stages []string, want string) {
	t.Helper()
	for _, s := range stages {
		if s == want {
			return
		}
	}
	t.Errorf("expected progress stage %q in %v", want, stages)
}
