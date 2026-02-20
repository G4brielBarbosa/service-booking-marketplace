package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abriesouza/super-assistente/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load(ctx)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := config.NewLogger(cfg.LogLevel)

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("connecting to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("worker started", "db", "connected")

	// Periodic cleanup of expired idempotency records
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("worker shutting down")
			return
		case <-ticker.C:
			log.Debug("running periodic cleanup")
			_, err := pool.Exec(ctx, "DELETE FROM idempotency_records WHERE expires_at <= now()")
			if err != nil {
				log.Error("cleanup failed", "error", err)
			}
		}
	}
}
