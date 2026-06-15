package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	shareddto "github.com/kjj1998/kinji/bff/internal/shared/dto"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// mockService is a function-backed test double for transactionsvc.TransactionService.
type mockService struct {
	GetMonthlyFn func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	ImportFn     func(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]shared.Transaction, error)
	SaveFn       func(ctx context.Context, userId string, transactions []shared.Transaction) ([]shared.Transaction, error)
	PeriodsFn    func(ctx context.Context, userId string) ([]domain.Period, error)
}

func (m *mockService) GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
	return m.GetMonthlyFn(ctx, userId, month, year)
}

func (m *mockService) ImportStatement(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]shared.Transaction, error) {
	return m.ImportFn(ctx, userId, statement, password, onProgress)
}

func (m *mockService) SaveTransactions(ctx context.Context, userId string, transactions []shared.Transaction) ([]shared.Transaction, error) {
	return m.SaveFn(ctx, userId, transactions)
}

func (m *mockService) GetPeriods(ctx context.Context, userId string) ([]domain.Period, error) {
	return m.PeriodsFn(ctx, userId)
}

// newRequest builds a request with the "id" path value set (empty id omits it).
func newRequest(method, target, id string, body []byte) *http.Request {
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, bytes.NewReader(body))
	}
	if id != "" {
		r.SetPathValue("id", id)
	}
	return r
}

func TestGetMonthlyTransactions_MissingID(t *testing.T) {
	h := NewTransactionHandler(&mockService{})
	w := httptest.NewRecorder()

	h.GetMonthlyTransactions(w, newRequest(http.MethodGet, "/api/v1/transactions/", "", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetMonthlyTransactions_InvalidMonth(t *testing.T) {
	h := NewTransactionHandler(&mockService{})
	w := httptest.NewRecorder()

	h.GetMonthlyTransactions(w, newRequest(http.MethodGet, "/api/v1/transactions/u1?month=13&year=2026", "u1", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetMonthlyTransactions_ServiceError(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		GetMonthlyFn: func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
			return nil, errors.New("boom")
		},
	})
	w := httptest.NewRecorder()

	h.GetMonthlyTransactions(w, newRequest(http.MethodGet, "/api/v1/transactions/u1?month=6&year=2026", "u1", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetMonthlyTransactions_Success(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		GetMonthlyFn: func(ctx context.Context, userId, month, year string) ([]shared.Transaction, error) {
			return []shared.Transaction{{ID: "t1", Merchant: "Acme"}}, nil
		},
	})
	w := httptest.NewRecorder()

	h.GetMonthlyTransactions(w, newRequest(http.MethodGet, "/api/v1/transactions/u1?month=6&year=2026", "u1", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var out []shareddto.Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(out) != 1 || out[0].ID != "t1" {
		t.Errorf("unexpected body: %+v", out)
	}
}

func TestSaveTransactions_BadJSON(t *testing.T) {
	h := NewTransactionHandler(&mockService{})
	w := httptest.NewRecorder()

	h.SaveTransactions(w, newRequest(http.MethodPost, "/api/v1/transactions/u1", "u1", []byte("not json")))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSaveTransactions_UnknownField(t *testing.T) {
	h := NewTransactionHandler(&mockService{})
	w := httptest.NewRecorder()

	body := []byte(`[{"id":"t1","bogus":true}]`)
	h.SaveTransactions(w, newRequest(http.MethodPost, "/api/v1/transactions/u1", "u1", body))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown field, got %d", w.Code)
	}
}

func TestSaveTransactions_ServiceError(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		SaveFn: func(ctx context.Context, userId string, transactions []shared.Transaction) ([]shared.Transaction, error) {
			return nil, errors.New("boom")
		},
	})
	w := httptest.NewRecorder()

	h.SaveTransactions(w, newRequest(http.MethodPost, "/api/v1/transactions/u1", "u1", []byte(`[{"id":"t1"}]`)))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestSaveTransactions_Success(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		SaveFn: func(ctx context.Context, userId string, transactions []shared.Transaction) ([]shared.Transaction, error) {
			return transactions, nil
		},
	})
	w := httptest.NewRecorder()

	h.SaveTransactions(w, newRequest(http.MethodPost, "/api/v1/transactions/u1", "u1", []byte(`[{"id":"t1"}]`)))

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetPeriods_ServiceError(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		PeriodsFn: func(ctx context.Context, userId string) ([]domain.Period, error) {
			return nil, errors.New("boom")
		},
	})
	w := httptest.NewRecorder()

	h.GetPeriods(w, newRequest(http.MethodGet, "/api/v1/transactions/u1/periods", "u1", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetPeriods_EmptyMarshalsAsArray(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		PeriodsFn: func(ctx context.Context, userId string) ([]domain.Period, error) {
			return nil, nil
		},
	})
	w := httptest.NewRecorder()

	h.GetPeriods(w, newRequest(http.MethodGet, "/api/v1/transactions/u1/periods", "u1", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if body := strings.TrimSpace(w.Body.String()); body != "[]" {
		t.Errorf("expected [] body, got %q", body)
	}
}

func TestImportStatement_MissingFile(t *testing.T) {
	h := NewTransactionHandler(&mockService{})
	w := httptest.NewRecorder()

	// multipart form with no "statement" file part.
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.WriteField("password", "")
	mw.Close()

	r := newRequest(http.MethodPost, "/api/v1/transactions/u1/import", "u1", buf.Bytes())
	r.Header.Set("Content-Type", mw.FormDataContentType())

	h.ImportStatement(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing file, got %d", w.Code)
	}
}

func TestImportStatement_SuccessEmitsDone(t *testing.T) {
	h := NewTransactionHandler(&mockService{
		ImportFn: func(ctx context.Context, userId string, statement multipart.File, password string, onProgress func(stage string)) ([]shared.Transaction, error) {
			onProgress("uploaded")
			return []shared.Transaction{{ID: "t1", Merchant: "Acme"}}, nil
		},
	})
	w := httptest.NewRecorder()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, _ := mw.CreateFormFile("statement", "stmt.pdf")
	fw.Write([]byte("%PDF-1.4 fake"))
	mw.Close()

	r := newRequest(http.MethodPost, "/api/v1/transactions/u1/import", "u1", buf.Bytes())
	r.Header.Set("Content-Type", mw.FormDataContentType())

	h.ImportStatement(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "event: done") {
		t.Errorf("expected a done event, got body:\n%s", body)
	}
	if !strings.Contains(body, "event: progress") {
		t.Errorf("expected a progress event, got body:\n%s", body)
	}
	if !strings.Contains(body, `"id":"t1"`) {
		t.Errorf("expected transaction in done payload, got body:\n%s", body)
	}
}

func TestImportErrorMessage(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"password required", domain.ErrPDFPasswordRequired, "pdf password required"},
		{"wrong password", domain.ErrPDFWrongPassword, "wrong pdf password given"},
		{"corrupt", domain.ErrPDFCorrupt, "invalid/corrupt pdf file"},
		{"client error", &shared.ClientError{Reason: "bad rows"}, "bad rows"},
		{"generic", errors.New("nope"), "failed to import statement"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := importErrorMessage(tc.err); got != tc.want {
				t.Errorf("importErrorMessage(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}
