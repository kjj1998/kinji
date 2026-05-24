package handler

import (
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/service"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	id, ok := requireUserId(w, r)
	if !ok {
		return
	}

	q := r.URL.Query()
	month, year, ok := parseMonthYear(w, q.Get("month"), q.Get("year"))
	if !ok {
		return
	}

	transactions, err := h.service.GetMonthlyTransactions(r.Context(), id, month, year)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get monthly transactions")
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, nil)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, nil)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, nil)
}
