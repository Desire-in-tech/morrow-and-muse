package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/catalog"
	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/config"
	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/dbpool"
	"github.com/Desire-in-tech/morrow-and-muse/backend/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	os.Exit(run(logger))
}

func run(logger *slog.Logger) int {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config", "err", err.Error())
		return 1
	}

	pool, err := dbpool.New(cfg.DatabaseURL)
	if err != nil {
		logger.Error("database", "err", err.Error())
		return 1
	}
	defer pool.Close()

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: httpapi.NewRouter(httpapi.Options{
			Catalog: catalog.NewPostgres(pool),
			Ready:   pool,
			Origins: cfg.CORSOrigins,
			Logger:  logger,
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", server.Addr, "cors_origins", len(cfg.CORSOrigins))
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "err", err)
			return 1
		}
	case <-ctx.Done():
		logger.Info("shutdown requested")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown", "err", err)
			return 1
		}
	}
	return 0
}
