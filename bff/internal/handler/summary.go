package handler

import (
	"net/http"

	"github.com/kohjunjie/kinji/bff/internal/service"
)

type SummaryHandler struct {
	service service.SummaryService
}

func NewSummaryHandler(svc service.SummaryService) *SummaryHandler {
	return &SummaryHandler{service: svc}
}

func (h *SummaryHandler) Summary(w http.ResponseWriter, r *http.Request) {
	id, ok := requireUserID(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	from, ok := parseDate(w, q.Get("from"), "from")
	if !ok {
		return
	}
	to, ok := parseDate(w, q.Get("to"), "to")
	if !ok {
		return
	}

	summary, err := h.service.GenerateMonthlySummary(r.Context(), id, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to calculate monthly summary")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}
