package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/abriesouza/super-assistente/internal/adapter/postgres"
	"github.com/abriesouza/super-assistente/internal/adapter/telegram"
	"github.com/abriesouza/super-assistente/internal/config"
	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/abriesouza/super-assistente/internal/usecase"
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

	if err := pool.Ping(ctx); err != nil {
		log.Error("pinging database", "error", err)
		os.Exit(1)
	}
	log.Info("database connected")

	store := postgres.NewStore(pool)

	userRepo := postgres.NewUserRepo(store)
	onbRepo := postgres.NewOnboardingRepo(store)
	goalRepo := postgres.NewGoalCycleRepo(store)
	privRepo := postgres.NewPrivacyRepo(store)
	baseRepo := postgres.NewBaselineRepo(store)
	mvdRepo := postgres.NewMVDRepo(store)
	eventRepo := postgres.NewEventRepo(store)
	idempRepo := postgres.NewIdempotencyRepo(store)
	dailyStateRepo := postgres.NewDailyStateRepo(store)

	onboardingUC := usecase.NewOnboardingUseCase(
		userRepo, onbRepo, goalRepo, privRepo, baseRepo, mvdRepo,
		eventRepo, idempRepo, log,
	)

	catalog := domain.NewHardcodedCatalog()
	dailyRoutineUC := usecase.NewDailyRoutineUseCase(
		userRepo, goalRepo, dailyStateRepo, catalog, eventRepo, log,
	)

	evidenceRepo := postgres.NewEvidenceRepo(store)
	gateResultRepo := postgres.NewGateResultRepo(store)
	rubricRepo := postgres.NewRubricRepo(store)

	gateUC := usecase.NewGateUseCase(
		userRepo, dailyStateRepo, privRepo, evidenceRepo, gateResultRepo, rubricRepo, eventRepo, log,
	)

	englishInputRepo := postgres.NewEnglishInputRepo(store)
	englishRetrievalRepo := postgres.NewEnglishRetrievalRepo(store)
	englishErrorLogRepo := postgres.NewEnglishErrorLogRepo(store)

	englishUC := usecase.NewEnglishUseCase(
		userRepo, dailyStateRepo, englishInputRepo, englishRetrievalRepo, englishErrorLogRepo,
		rubricRepo, gateUC, eventRepo, log,
	)

	// HTTP health endpoint for admin/monitoring
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	srv := &http.Server{Addr: httpAddr, Handler: mux}

	go func() {
		log.Info("http server starting", "addr", httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "error", err)
		}
	}()

	// Start Telegram bot (long polling)
	bot, err := telegram.NewBot(cfg.TelegramBotToken, onboardingUC, dailyRoutineUC, gateUC, englishUC, idempRepo, cfg.FeatureOnboardingV1, cfg.FeatureDailyRoutineV1, cfg.FeatureQualityGatesV1, log)
	if err != nil {
		log.Error("creating telegram bot", "error", err)
		os.Exit(1)
	}

	go bot.StartPolling(ctx)

	<-ctx.Done()
	log.Info("shutting down")
	_ = srv.Shutdown(context.Background())
}
