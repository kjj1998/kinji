package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// mockService is a function-backed test double for service.SummaryService.
type mockService struct {
	GenerateFn func(ctx context.Context, userId, month, year string) (*domain.MonthlySummary, error)
}

func (m *mockService) GenerateMonthlySummary(ctx context.Context, userId, month, year string) (*domain.MonthlySummary, error) {
	return m.GenerateFn(ctx, userId, month, year)
}

func newRequest(target, id string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, target, nil)
	if id != "" {
		r.SetPathValue("id", id)
	}
	return r
}

func TestSummary_MissingID(t *testing.T) {
	h := NewSummaryHandler(&mockService{})
	w := httptest.NewRecorder()

	h.Summary(w, newRequest("/api/v1/summary/", ""))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSummary_InvalidMonthYear(t *testing.T) {
	h := NewSummaryHandler(&mockService{})
	w := httptest.NewRecorder()

	h.Summary(w, newRequest("/api/v1/summary/u1?month=0&year=2026", "u1"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSummary_ServiceError(t *testing.T) {
	h := NewSummaryHandler(&mockService{
		GenerateFn: func(ctx context.Context, userId, month, year string) (*domain.MonthlySummary, error) {
			return nil, errors.New("boom")
		},
	})
	w := httptest.NewRecorder()

	h.Summary(w, newRequest("/api/v1/summary/u1?month=6&year=2026", "u1"))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestSummary_Success(t *testing.T) {
	h := NewSummaryHandler(&mockService{
		GenerateFn: func(ctx context.Context, userId, month, year string) (*domain.MonthlySummary, error) {
			return &domain.MonthlySummary{SummaryStatement: "you did fine", LastMonthSpent: 350}, nil
		},
	})
	w := httptest.NewRecorder()

	h.Summary(w, newRequest("/api/v1/summary/u1?month=6&year=2026", "u1"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var out struct {
		MonthlySummary string `json:"monthlySummary"`
		LastMonthSpent int    `json:"lastMonthSpent"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if out.MonthlySummary != "you did fine" || out.LastMonthSpent != 350 {
		t.Errorf("unexpected body: %+v", out)
	}
}
