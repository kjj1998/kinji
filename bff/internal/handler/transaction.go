package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kohjunjie/kinji/bff/internal/service"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	id, ok := requireUserID(w, r)
	if !ok {
		return
	}

	transactions, err := h.service.GetMonthlyTransactions(r.Context(), id, "", "")

	if err != nil {
		slog.Error(err.Error())
		writeError(w, http.StatusInternalServerError, "failed to fetch transactions")
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsAvailabilities(w http.ResponseWriter, r *http.Request) {
	id, ok := requireUserID(w, r)
	if !ok {
		return
	}

	availabilities, err := h.service.GetTransactionsAvailabilities(r.Context(), id)

	if err != nil {
		writeError(w, http.StatusInternalServerError,
			fmt.Sprintf("failed to retrieve transactions availabilites for id %s", id))
		return
	}
	writeJSON(w, http.StatusOK, availabilities)
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
