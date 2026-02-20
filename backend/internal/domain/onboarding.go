package domain

import (
	"time"

	"github.com/google/uuid"
)

type OnboardingStatus string

const (
	OnboardingNew              OnboardingStatus = "new"
	OnboardingInProgress       OnboardingStatus = "in_progress"
	OnboardingMinimumCompleted OnboardingStatus = "minimum_completed"
	OnboardingCompleted        OnboardingStatus = "completed"
)

type StepID string

const (
	StepWelcome        StepID = "welcome"
	StepSelectGoals    StepID = "select_goals"
	StepRestrictions   StepID = "restrictions"
	StepSleepBaseline  StepID = "sleep_baseline"
	StepEnglishBase    StepID = "english_baseline"
	StepJavaBaseline   StepID = "java_baseline"
	StepMVD            StepID = "mvd"
	StepPrivacy        StepID = "privacy"
	StepSummary        StepID = "summary"
)

// MinimumOnboardingSteps: steps required to reach minimum_completed.
var MinimumOnboardingSteps = []StepID{
	StepWelcome,
	StepSelectGoals,
	StepRestrictions,
	StepSleepBaseline,
	StepMVD,
	StepPrivacy,
	StepSummary,
}

type OnboardingAnswer struct {
	StepID    StepID         `json:"step_id"`
	Value     map[string]any `json:"value"`
	Timestamp time.Time      `json:"timestamp"`
}

type PendingItem struct {
	Domain      GoalID `json:"domain"`
	Description string `json:"description"`
	StepID      StepID `json:"step_id"`
}

type OnboardingSession struct {
	SessionID        uuid.UUID           `json:"session_id"`
	UserID           uuid.UUID           `json:"user_id"`
	Status           OnboardingStatus    `json:"status"`
	CurrentStepID    StepID              `json:"current_step_id"`
	Answers          []OnboardingAnswer  `json:"answers"`
	PendingItems     []PendingItem       `json:"pending_items"`
	StartedAt        time.Time           `json:"started_at"`
	LastInteraction  time.Time           `json:"last_interaction_at"`
	CompletedAt      *time.Time          `json:"completed_at,omitempty"`
}

func NewOnboardingSession(userID uuid.UUID) *OnboardingSession {
	now := time.Now()
	return &OnboardingSession{
		SessionID:       uuid.New(),
		UserID:          userID,
		Status:          OnboardingNew,
		CurrentStepID:   StepWelcome,
		Answers:         []OnboardingAnswer{},
		PendingItems:    []PendingItem{},
		StartedAt:       now,
		LastInteraction: now,
	}
}

func (s *OnboardingSession) HasAnswered(step StepID) bool {
	for _, a := range s.Answers {
		if a.StepID == step {
			return true
		}
	}
	return false
}

func (s *OnboardingSession) GetAnswer(step StepID) *OnboardingAnswer {
	for i := range s.Answers {
		if s.Answers[i].StepID == step {
			return &s.Answers[i]
		}
	}
	return nil
}

func (s *OnboardingSession) SubmitAnswer(stepID StepID, value map[string]any) error {
	if s.Status == OnboardingCompleted {
		return NewStateConflictError("Onboarding já concluído", "Use revisão para alterar campos")
	}

	if stepID != s.CurrentStepID {
		if s.HasAnswered(stepID) {
			return nil // idempotent: already answered this step
		}
		return NewStateConflictError("Step fora de ordem", string(stepID))
	}

	s.Answers = append(s.Answers, OnboardingAnswer{
		StepID:    stepID,
		Value:     value,
		Timestamp: time.Now(),
	})
	s.LastInteraction = time.Now()

	s.advanceStep()
	return nil
}

func (s *OnboardingSession) advanceStep() {
	steps := MinimumOnboardingSteps
	for i, step := range steps {
		if step == s.CurrentStepID && i+1 < len(steps) {
			s.CurrentStepID = steps[i+1]
			s.Status = OnboardingInProgress
			return
		}
	}

	// All minimum steps done
	if s.CurrentStepID == StepSummary && s.HasAnswered(StepSummary) {
		s.Status = OnboardingMinimumCompleted
		now := time.Now()
		s.CompletedAt = &now
		s.computePendingItems()
	}
}

func (s *OnboardingSession) computePendingItems() {
	goalsAnswer := s.GetAnswer(StepSelectGoals)
	if goalsAnswer == nil {
		return
	}

	goalIDs, ok := goalsAnswer.Value["goals"].([]any)
	if !ok {
		return
	}

	for _, gRaw := range goalIDs {
		gID, ok := gRaw.(string)
		if !ok {
			continue
		}
		goalID := GoalID(gID)
		class := GoalClassifications[goalID]

		if class == GoalClassIntensive && !s.hasBaselineFor(goalID) {
			s.PendingItems = append(s.PendingItems, PendingItem{
				Domain:      goalID,
				Description: "Baseline de " + gID,
				StepID:      baselineStepFor(goalID),
			})
		}
	}
}

func (s *OnboardingSession) hasBaselineFor(goal GoalID) bool {
	step := baselineStepFor(goal)
	return s.HasAnswered(step)
}

func baselineStepFor(goal GoalID) StepID {
	switch goal {
	case GoalEnglish:
		return StepEnglishBase
	case GoalJava:
		return StepJavaBaseline
	default:
		return StepID("baseline_" + string(goal))
	}
}

func (s *OnboardingSession) CompletePendingItem(item PendingItem, value map[string]any) {
	s.Answers = append(s.Answers, OnboardingAnswer{
		StepID:    item.StepID,
		Value:     value,
		Timestamp: time.Now(),
	})

	remaining := make([]PendingItem, 0, len(s.PendingItems))
	for _, p := range s.PendingItems {
		if p.StepID != item.StepID {
			remaining = append(remaining, p)
		}
	}
	s.PendingItems = remaining
	s.LastInteraction = time.Now()

	if len(s.PendingItems) == 0 {
		s.Status = OnboardingCompleted
		now := time.Now()
		s.CompletedAt = &now
	}
}

// MVD represents the Minimum Viable Daily plan.
type MVDItem struct {
	Domain   GoalID `json:"domain"`
	Action   string `json:"action"`
	Duration string `json:"duration"`
	Criteria string `json:"criteria"`
}

type MinimumViableDaily struct {
	MVDID     uuid.UUID `json:"mvd_id"`
	UserID    uuid.UUID `json:"user_id"`
	Items     []MVDItem `json:"items"`
	WhenToUse string    `json:"when_to_use"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BaselineSnapshot stores baseline data per domain.
type BaselineCompleteness string

const (
	BaselineMinimum  BaselineCompleteness = "minimum"
	BaselinePartial  BaselineCompleteness = "partial"
	BaselineComplete BaselineCompleteness = "complete"
)

type BaselineSnapshot struct {
	BaselineID   uuid.UUID            `json:"baseline_id"`
	UserID       uuid.UUID            `json:"user_id"`
	Domain       GoalID               `json:"domain"`
	Data         map[string]any       `json:"data"`
	Completeness BaselineCompleteness `json:"completeness"`
	CapturedAt   time.Time            `json:"captured_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

// OnboardingSummary is the consultable summary of onboarding state.
type OnboardingSummary struct {
	ActiveGoals   []GoalEntry        `json:"active_goals"`
	PausedGoals   []GoalEntry        `json:"paused_goals"`
	Restrictions  map[string]any     `json:"restrictions"`
	Baselines     []BaselineSnapshot `json:"baselines"`
	MVD           *MinimumViableDaily `json:"mvd"`
	PendingItems  []PendingItem      `json:"pending_items"`
	PrivacyPolicy *PrivacyPolicy     `json:"privacy_policy"`
}
