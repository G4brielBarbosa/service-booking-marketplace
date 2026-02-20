package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL,required"`
	RedisURL         string `env:"REDIS_URL,required"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	LogLevel         string `env:"LOG_LEVEL,default=info"`
	HTTPPort         int    `env:"HTTP_PORT,default=8080"`

	// Feature flags (PLAN-000 §11)
	FeatureOnboardingV1  bool `env:"FEATURE_ONBOARDING_V1,default=false"`
	FeatureDailyRoutineV1 bool `env:"FEATURE_DAILY_ROUTINE_V1,default=false"`
	FeatureQualityGatesV1 bool `env:"FEATURE_QUALITY_GATES_V1,default=false"`
	FeatureNudgesV1      bool `env:"FEATURE_NUDGES_V1,default=false"`
}

func Load(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &cfg, nil
}
