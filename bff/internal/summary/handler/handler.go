package handler

import (
	"net/http"

	platformhttp "github.com/kjj1998/kinji/bff/internal/platform/http"
	"github.com/kjj1998/kinji/bff/internal/summary/dto"
	"github.com/kjj1998/kinji/bff/internal/summary/service"
)

type SummaryHandler struct {
	svc service.SummaryService
}

func NewSummaryHandler(svc service.SummaryService) *SummaryHandler {
	return &SummaryHandler{svc: svc}
}

func (h *SummaryHandler) Summary(w http.ResponseWriter, r *http.Request) {
	id, ok := platformhttp.RequireUserId(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	month, year, ok := platformhttp.ParseMonthYear(w, q.Get("month"), q.Get("year"))
	if !ok {
		return
	}

	summary, err := h.svc.GenerateMonthlySummary(r.Context(), id, month, year)
	if err != nil {
		platformhttp.WriteError(w, http.StatusInternalServerError, "failed to calculate monthly summary")
		return
	}
	platformhttp.WriteJSON(w, http.StatusOK, dto.ToTransactionSummary(summary))
}
