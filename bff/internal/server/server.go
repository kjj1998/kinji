package server

import (
	"net/http"

	"github.com/kjj1998/kinji/bff/internal/handler"
	"github.com/kjj1998/kinji/bff/internal/middleware"
	"github.com/kjj1998/kinji/bff/internal/repository"
	"github.com/kjj1998/kinji/bff/internal/service"
)

func New(repo repository.Repository, corsOrigin string) http.Handler {
	mux := http.NewServeMux()

	summaryService := service.NewSummaryService(repo)
	txService := service.NewTransactionService(repo)

	txHandler := handler.NewTransactionHandler(txService)
	summaryHandler := handler.NewSummaryHandler(summaryService)

	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /api/v1/transactions/{id}", txHandler.List)
	mux.HandleFunc("GET /api/v1/summary/{id}", summaryHandler.Summary)
	mux.HandleFunc("POST /api/v1/transactions", txHandler.Create)
	mux.HandleFunc("PATCH /api/v1/transactions/{id}", txHandler.Update)
	mux.HandleFunc("DELETE /api/v1/transactions/{id}", txHandler.Delete)

	return middleware.Chain(mux,
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS(corsOrigin),
	)
}
