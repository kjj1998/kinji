package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/service"
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	service service.TransactionService
}

// NewTransactionHandler returns a TransactionHandler backed by svc.
func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: svc}
}

// GetMonthlyTransactions writes the user's monthly transactions as JSON, selected by the
// "month" and "year" query parameters.
func (h *TransactionHandler) GetMonthlyTransactions(w http.ResponseWriter, r *http.Request) {
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
		slog.ErrorContext(r.Context(), "get monthly transactions", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get monthly transactions")
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}

// ImportStatement parses an uploaded PDF bank statement and
// extracts and categorizes its transactions for the user.
func (h *TransactionHandler) ImportStatement(w http.ResponseWriter, r *http.Request) {
	userId, ok := requireUserId(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, _, err := r.FormFile("statement")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing statement file")
		return
	}
	defer file.Close()

	password := r.FormValue("password")

	transactions, err := h.service.ImportStatement(r.Context(), userId, file, password)
	if err != nil {
		var ce *service.ClientError
		if errors.As(err, &ce) {
			writeError(w, http.StatusBadRequest, ce.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to import statement")
		return
	}
	writeJSON(w, http.StatusOK, transactions)
}
