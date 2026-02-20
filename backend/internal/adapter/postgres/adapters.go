package postgres

import (
	"context"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/google/uuid"
)

// UserRepo adapts Store to port.UserRepository.
type UserRepo struct{ s *Store }

func NewUserRepo(s *Store) *UserRepo { return &UserRepo{s: s} }

func (r *UserRepo) FindByTelegramID(ctx context.Context, id int64) (*domain.UserProfile, error) {
	return r.s.FindByTelegramID(ctx, id)
}
func (r *UserRepo) Save(ctx context.Context, u *domain.UserProfile) error {
	return r.s.SaveUser(ctx, u)
}
func (r *UserRepo) Update(ctx context.Context, u *domain.UserProfile) error {
	return r.s.UpdateUser(ctx, u)
}

// OnboardingRepo adapts Store to port.OnboardingRepository.
type OnboardingRepo struct{ s *Store }

func NewOnboardingRepo(s *Store) *OnboardingRepo { return &OnboardingRepo{s: s} }

func (r *OnboardingRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*domain.OnboardingSession, error) {
	return r.s.FindOnboardingByUserID(ctx, id)
}
func (r *OnboardingRepo) Save(ctx context.Context, sess *domain.OnboardingSession) error {
	return r.s.SaveOnboarding(ctx, sess)
}
func (r *OnboardingRepo) Update(ctx context.Context, sess *domain.OnboardingSession) error {
	return r.s.UpdateOnboarding(ctx, sess)
}

// GoalCycleRepo adapts Store to port.GoalCycleRepository.
type GoalCycleRepo struct{ s *Store }

func NewGoalCycleRepo(s *Store) *GoalCycleRepo { return &GoalCycleRepo{s: s} }

func (r *GoalCycleRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*domain.ActiveGoalCycle, error) {
	return r.s.FindGoalCycleByUserID(ctx, id)
}
func (r *GoalCycleRepo) Save(ctx context.Context, c *domain.ActiveGoalCycle) error {
	return r.s.SaveGoalCycle(ctx, c)
}
func (r *GoalCycleRepo) Update(ctx context.Context, c *domain.ActiveGoalCycle) error {
	return r.s.UpdateGoalCycle(ctx, c)
}

// PrivacyRepo adapts Store to port.PrivacyPolicyRepository.
type PrivacyRepo struct{ s *Store }

func NewPrivacyRepo(s *Store) *PrivacyRepo { return &PrivacyRepo{s: s} }

func (r *PrivacyRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*domain.PrivacyPolicy, error) {
	return r.s.FindPrivacyPolicy(ctx, id)
}
func (r *PrivacyRepo) Save(ctx context.Context, p *domain.PrivacyPolicy) error {
	return r.s.SavePrivacyPolicy(ctx, p)
}
func (r *PrivacyRepo) Update(ctx context.Context, p *domain.PrivacyPolicy) error {
	return r.s.UpdatePrivacyPolicy(ctx, p)
}

// BaselineRepo adapts Store to port.BaselineRepository.
type BaselineRepo struct{ s *Store }

func NewBaselineRepo(s *Store) *BaselineRepo { return &BaselineRepo{s: s} }

func (r *BaselineRepo) FindByUserAndDomain(ctx context.Context, uid uuid.UUID, d domain.GoalID) (*domain.BaselineSnapshot, error) {
	return r.s.FindBaselineByUserAndDomain(ctx, uid, d)
}
func (r *BaselineRepo) FindAllByUser(ctx context.Context, uid uuid.UUID) ([]domain.BaselineSnapshot, error) {
	return r.s.FindAllBaselinesByUser(ctx, uid)
}
func (r *BaselineRepo) Save(ctx context.Context, b *domain.BaselineSnapshot) error {
	return r.s.SaveBaseline(ctx, b)
}
func (r *BaselineRepo) Update(ctx context.Context, b *domain.BaselineSnapshot) error {
	return r.s.UpdateBaseline(ctx, b)
}

// MVDRepo adapts Store to port.MVDRepository.
type MVDRepo struct{ s *Store }

func NewMVDRepo(s *Store) *MVDRepo { return &MVDRepo{s: s} }

func (r *MVDRepo) FindByUserID(ctx context.Context, id uuid.UUID) (*domain.MinimumViableDaily, error) {
	return r.s.FindMVDByUserID(ctx, id)
}
func (r *MVDRepo) Save(ctx context.Context, m *domain.MinimumViableDaily) error {
	return r.s.SaveMVD(ctx, m)
}
func (r *MVDRepo) Update(ctx context.Context, m *domain.MinimumViableDaily) error {
	return r.s.UpdateMVD(ctx, m)
}

// EventRepo adapts Store to port.EventRepository.
type EventRepo struct{ s *Store }

func NewEventRepo(s *Store) *EventRepo { return &EventRepo{s: s} }

func (r *EventRepo) Append(ctx context.Context, e domain.DomainEvent) error {
	return r.s.AppendEvent(ctx, e)
}

// IdempotencyRepo adapts Store to port.IdempotencyStore.
type IdempotencyRepo struct{ s *Store }

func NewIdempotencyRepo(s *Store) *IdempotencyRepo { return &IdempotencyRepo{s: s} }

func (r *IdempotencyRepo) Check(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	return r.s.CheckIdempotency(ctx, key)
}
func (r *IdempotencyRepo) Store(ctx context.Context, rec domain.IdempotencyRecord) error {
	return r.s.StoreIdempotency(ctx, rec)
}

// DailyStateRepo adapts Store to port.DailyStateRepository.
type DailyStateRepo struct{ s *Store }

func NewDailyStateRepo(s *Store) *DailyStateRepo { return &DailyStateRepo{s: s} }

func (r *DailyStateRepo) FindByUserAndDate(ctx context.Context, userID uuid.UUID, localDate string) (*domain.DailyState, error) {
	return r.s.FindDailyStateByUserAndDate(ctx, userID, localDate)
}
func (r *DailyStateRepo) Save(ctx context.Context, state *domain.DailyState) error {
	return r.s.SaveDailyState(ctx, state)
}
func (r *DailyStateRepo) Update(ctx context.Context, state *domain.DailyState) error {
	return r.s.UpdateDailyState(ctx, state)
}

// EvidenceRepo adapts Store to port.EvidenceRepository.
type EvidenceRepo struct{ s *Store }

func NewEvidenceRepo(s *Store) *EvidenceRepo { return &EvidenceRepo{s: s} }

func (r *EvidenceRepo) Save(ctx context.Context, ev *domain.Evidence) error {
	return r.s.SaveEvidence(ctx, ev)
}
func (r *EvidenceRepo) FindByTaskID(ctx context.Context, taskID uuid.UUID) ([]domain.Evidence, error) {
	return r.s.FindEvidenceByTaskID(ctx, taskID)
}

// GateResultRepo adapts Store to port.GateResultRepository.
type GateResultRepo struct{ s *Store }

func NewGateResultRepo(s *Store) *GateResultRepo { return &GateResultRepo{s: s} }

func (r *GateResultRepo) Save(ctx context.Context, gr *domain.GateResult) error {
	return r.s.SaveGateResult(ctx, gr)
}
func (r *GateResultRepo) FindByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.GateResult, error) {
	return r.s.FindGateResultByTaskID(ctx, taskID)
}
func (r *GateResultRepo) FindByUserAndDate(ctx context.Context, userID uuid.UUID, localDate string) ([]domain.GateResult, error) {
	return r.s.FindGateResultsByUserAndDate(ctx, userID, localDate)
}

// RubricRepo adapts Store to port.RubricRepository.
type RubricRepo struct{ s *Store }

func NewRubricRepo(s *Store) *RubricRepo { return &RubricRepo{s: s} }

func (r *RubricRepo) Save(ctx context.Context, rubric *domain.RubricScore) error {
	return r.s.SaveRubric(ctx, rubric)
}
func (r *RubricRepo) FindByTaskID(ctx context.Context, taskID uuid.UUID) (*domain.RubricScore, error) {
	return r.s.FindRubricByTaskID(ctx, taskID)
}
