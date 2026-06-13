package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kjj1998/kinji/bff/internal/adapter/http/dto"
	"github.com/kjj1998/kinji/bff/internal/app"
	"github.com/kjj1998/kinji/bff/internal/model"
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	service app.TransactionService
}

// NewTransactionHandler returns a TransactionHandler backed by svc.
func NewTransactionHandler(svc app.TransactionService) *TransactionHandler {
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
	writeJSON(w, http.StatusOK, dto.ToTransactions(transactions))
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

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	rc := http.NewResponseController(w)
	rc.SetWriteDeadline(time.Now().Add(60 * time.Second)) // extend deadline to give request longer time to complete

	send := func(event, data string) {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
		rc.Flush()
	}
	sendError := func(msg string) { send("error", fmt.Sprintf(`{"message":%q}`, msg)) }

	transactions, err := h.service.ImportStatement(r.Context(), userId, file, password,
		func(stage string) { send("progress", fmt.Sprintf(`{"stage":%q}`, stage)) })
	if err != nil {
		slog.ErrorContext(r.Context(), "import statement", "error", err)
		sendError(importErrorMessage(err))
		return
	}

	data, err := json.Marshal(dto.ToTransactions(transactions))
	if err != nil {
		sendError("failed to encode result")
		return
	}
	send("done", string(data))
}

// importErrorMessage maps an import error to a client-facing message, surfacing
// PDF and bad-input problems while hiding internal failures.
func importErrorMessage(err error) string {
	switch {
	case errors.Is(err, model.ErrPDFPasswordRequired):
		return "pdf password required"
	case errors.Is(err, model.ErrPDFWrongPassword):
		return "wrong pdf password given"
	case errors.Is(err, model.ErrPDFCorrupt):
		return "invalid/corrupt pdf file"
	}
	var ce *app.ClientError
	if errors.As(err, &ce) {
		return ce.Error()
	}
	return "failed to import statement"
}

// SaveTransactions saves transactions that has been reviewed and
// approved by the user into the database.
func (h *TransactionHandler) SaveTransactions(w http.ResponseWriter, r *http.Request) {
	userId, ok := requireUserId(w, r)
	if !ok {
		return
	}

	var transactions []dto.Transaction

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&transactions); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	saved, err := h.service.SaveTransactions(r.Context(), userId, dto.DomainTransactions(transactions))
	if err != nil {
		slog.ErrorContext(r.Context(), "save transactions", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save transactions")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToTransactions(saved))
}

// GetPeriods retrieves the years and months where transaction data is available
func (h *TransactionHandler) GetPeriods(w http.ResponseWriter, r *http.Request) {
	userId, ok := requireUserId(w, r)
	if !ok {
		return
	}

	periods, err := h.service.GetPeriods(r.Context(), userId)
	if err != nil {
		slog.ErrorContext(r.Context(), "get periods", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get periods")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToPeriods(periods))
}
