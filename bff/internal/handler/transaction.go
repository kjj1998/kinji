package handler

import (
	"net/http"

	"github.com/kohjunjie/kinji/bff/internal/repository"
	"github.com/kohjunjie/kinji/bff/internal/service"
)

type TransactionHandler struct {
	repo    repository.TransactionRepository
	service service.TransactionService
}

func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repo:    repo,
		service: service.NewTransactionService(repo),
	}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID not provided")
		return
	}
	transactions, err := h.repo.List(r.Context(), id, "", "")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch transactions")
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Summary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "User ID not provided")
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

	summary, err := h.service.Summary(r.Context(), id, from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to calculate monthly summary")
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
