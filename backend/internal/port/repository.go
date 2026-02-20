package port

import (
	"context"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByTelegramID(ctx context.Context, telegramUserID int64) (*domain.UserProfile, error)
	Save(ctx context.Context, user *domain.UserProfile) error
	Update(ctx context.Context, user *domain.UserProfile) error
}

type OnboardingRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.OnboardingSession, error)
	Save(ctx context.Context, session *domain.OnboardingSession) error
	Update(ctx context.Context, session *domain.OnboardingSession) error
}

type GoalCycleRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.ActiveGoalCycle, error)
	Save(ctx context.Context, cycle *domain.ActiveGoalCycle) error
	Update(ctx context.Context, cycle *domain.ActiveGoalCycle) error
}

type PrivacyPolicyRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.PrivacyPolicy, error)
	Save(ctx context.Context, policy *domain.PrivacyPolicy) error
	Update(ctx context.Context, policy *domain.PrivacyPolicy) error
}

type BaselineRepository interface {
	FindByUserAndDomain(ctx context.Context, userID uuid.UUID, d domain.GoalID) (*domain.BaselineSnapshot, error)
	FindAllByUser(ctx context.Context, userID uuid.UUID) ([]domain.BaselineSnapshot, error)
	Save(ctx context.Context, baseline *domain.BaselineSnapshot) error
	Update(ctx context.Context, baseline *domain.BaselineSnapshot) error
}

type MVDRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.MinimumViableDaily, error)
	Save(ctx context.Context, mvd *domain.MinimumViableDaily) error
	Update(ctx context.Context, mvd *domain.MinimumViableDaily) error
}

type EventRepository interface {
	Append(ctx context.Context, event domain.DomainEvent) error
}

type IdempotencyStore interface {
	Check(ctx context.Context, key string) (*domain.IdempotencyRecord, error)
	Store(ctx context.Context, record domain.IdempotencyRecord) error
}

type DailyStateRepository interface {
	FindByUserAndDate(ctx context.Context, userID uuid.UUID, localDate string) (*domain.DailyState, error)
	Save(ctx context.Context, state *domain.DailyState) error
	Update(ctx context.Context, state *domain.DailyState) error
}

type TaskCatalog interface {
	GetTasksForGoal(goalDomain domain.GoalID, planType domain.PlanType) []domain.TaskTemplate
}
