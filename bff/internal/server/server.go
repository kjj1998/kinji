package server

import (
	"database/sql"
	"net/http"

	platformhttp "github.com/kjj1998/kinji/bff/internal/platform/http"
	summaryhandler "github.com/kjj1998/kinji/bff/internal/summary/handler"
	summaryservice "github.com/kjj1998/kinji/bff/internal/summary/service"
	summarystore "github.com/kjj1998/kinji/bff/internal/summary/store"
	txnhandler "github.com/kjj1998/kinji/bff/internal/transaction/handler"
	txnservice "github.com/kjj1998/kinji/bff/internal/transaction/service"
	txnstore "github.com/kjj1998/kinji/bff/internal/transaction/store"
)

func New(db *sql.DB, parser txnservice.StatementParser, corsOrigin string) http.Handler {
	mux := http.NewServeMux()

	txnRepo := txnstore.NewRepository(db)
	summaryRepo := summarystore.NewRepository(db)

	summaryService := summaryservice.NewSummaryService(summaryRepo)
	txService := txnservice.NewTransactionService(txnRepo, parser)

	txHandler := txnhandler.NewTransactionHandler(txService)
	summaryHandler := summaryhandler.NewSummaryHandler(summaryService)

	mux.HandleFunc("GET /health", platformhttp.Health)
	mux.HandleFunc("GET /api/v1/transactions/{id}", txHandler.GetMonthlyTransactions)
	mux.HandleFunc("POST /api/v1/transactions/{id}", txHandler.SaveTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}/periods", txHandler.GetPeriods)
	mux.HandleFunc("GET /api/v1/summary/{id}", summaryHandler.Summary)
	mux.HandleFunc("POST /api/v1/transactions/{id}/import", txHandler.ImportStatement)

	return platformhttp.Chain(mux,
		platformhttp.Recovery,
		platformhttp.Logging,
		platformhttp.CORS(corsOrigin),
	)
}
