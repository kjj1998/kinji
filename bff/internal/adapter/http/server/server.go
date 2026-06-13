package server

import (
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/adapter/http/handler"
	"github.com/kjj1998/kinji/bff/internal/adapter/http/middleware"
	"github.com/kjj1998/kinji/bff/internal/app"
)

func New(repo app.TransactionRepository, parser app.StatementParser, corsOrigin string) http.Handler {
	mux := http.NewServeMux()

	// services
	summaryService := app.NewSummaryService(repo)
	txService := app.NewTransactionService(repo, parser)

	// handlers
	txHandler := handler.NewTransactionHandler(txService)
	summaryHandler := handler.NewSummaryHandler(summaryService)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /api/v1/transactions/{id}", txHandler.GetMonthlyTransactions)
	mux.HandleFunc("POST /api/v1/transactions/{id}", txHandler.SaveTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}/periods", txHandler.GetPeriods)
	mux.HandleFunc("GET /api/v1/summary/{id}", summaryHandler.Summary)
	mux.HandleFunc("POST /api/v1/transactions/{id}/import", txHandler.ImportStatement)

	return middleware.Chain(mux,
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS(corsOrigin),
	)
}
