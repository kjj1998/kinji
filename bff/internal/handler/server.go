package handler

import (
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/service"
)

func New(repo service.TransactionRepository, parser service.StatementParser, corsOrigin string) http.Handler {
	mux := http.NewServeMux()

	// services
	summaryService := service.NewSummaryService(repo)
	txService := service.NewTransactionService(repo, parser)

	// handlers
	txHandler := NewTransactionHandler(txService)
	summaryHandler := NewSummaryHandler(summaryService)

	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("GET /api/v1/transactions/{id}", txHandler.GetMonthlyTransactions)
	mux.HandleFunc("POST /api/v1/transactions/{id}", txHandler.SaveTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}/periods", txHandler.GetPeriods)
	mux.HandleFunc("GET /api/v1/summary/{id}", summaryHandler.Summary)
	mux.HandleFunc("POST /api/v1/transactions/{id}/import", txHandler.ImportStatement)

	return Chain(mux,
		Recovery,
		Logging,
		CORS(corsOrigin),
	)
}
