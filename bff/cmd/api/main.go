package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kjj1998/kinji/bff/internal/platform/config"
	"github.com/kjj1998/kinji/bff/internal/platform/database"
	"github.com/kjj1998/kinji/bff/internal/server"
	"github.com/kjj1998/kinji/bff/internal/transaction/parser"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := database.NewClient(cfg.SQLitePath)
	if err != nil {
		slog.Error("failed to create sqlite client",
			"error", err,
			"path", cfg.SQLitePath)
		os.Exit(1)
	}
	parser := parser.NewParser(cfg.AnthropicModel)

	h := server.New(db, parser, cfg.CORSOrigin)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Port, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
