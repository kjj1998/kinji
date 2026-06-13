package handler

import (
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/adapter/http/dto"
	"github.com/kjj1998/kinji/bff/internal/app"
)

type SummaryHandler struct {
	service app.SummaryService
}

func NewSummaryHandler(svc app.SummaryService) *SummaryHandler {
	return &SummaryHandler{service: svc}
}

func (h *SummaryHandler) Summary(w http.ResponseWriter, r *http.Request) {
	id, ok := requireUserId(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	month, year, ok := parseMonthYear(w, q.Get("month"), q.Get("year"))
	if !ok {
		return
	}

	summary, err := h.service.GenerateMonthlySummary(r.Context(), id, month, year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to calculate monthly summary")
		return
	}
	writeJSON(w, http.StatusOK, dto.ToTransactionSummary(summary))
}
