package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOnboardingSession(t *testing.T) {
	userID := uuid.New()
	s := NewOnboardingSession(userID)

	if s.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, s.UserID)
	}
	if s.Status != OnboardingNew {
		t.Errorf("expected status 'new', got %s", s.Status)
	}
	if s.CurrentStepID != StepWelcome {
		t.Errorf("expected step 'welcome', got %s", s.CurrentStepID)
	}
	if len(s.Answers) != 0 {
		t.Errorf("expected 0 answers, got %d", len(s.Answers))
	}
}

func TestOnboardingStepProgression(t *testing.T) {
	s := NewOnboardingSession(uuid.New())

	steps := MinimumOnboardingSteps
	for i, step := range steps[:len(steps)-1] {
		err := s.SubmitAnswer(step, map[string]any{"test": true})
		if err != nil {
			t.Fatalf("step %d (%s): unexpected error: %v", i, step, err)
		}
		if s.CurrentStepID != steps[i+1] {
			t.Errorf("after step %s, expected next step %s, got %s", step, steps[i+1], s.CurrentStepID)
		}
	}

	if s.Status != OnboardingInProgress {
		t.Errorf("expected in_progress before final step, got %s", s.Status)
	}
}

func TestOnboardingMinimumCompleted(t *testing.T) {
	s := NewOnboardingSession(uuid.New())

	for _, step := range MinimumOnboardingSteps {
		val := map[string]any{"test": true}
		if step == StepSelectGoals {
			val = map[string]any{"goals": []any{"sleep", "health"}}
		}
		err := s.SubmitAnswer(step, val)
		if err != nil {
			t.Fatalf("step %s: unexpected error: %v", step, err)
		}
	}

	if s.Status != OnboardingMinimumCompleted {
		t.Errorf("expected minimum_completed, got %s", s.Status)
	}
	if s.CompletedAt == nil {
		t.Error("expected completed_at to be set")
	}
}

func TestOnboardingIdempotentAnswer(t *testing.T) {
	s := NewOnboardingSession(uuid.New())

	err := s.SubmitAnswer(StepWelcome, map[string]any{"ok": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Submit again for same step (idempotent)
	err = s.SubmitAnswer(StepWelcome, map[string]any{"ok": true})
	if err != nil {
		t.Fatalf("idempotent answer should not error: %v", err)
	}
}

func TestOnboardingOutOfOrderStep(t *testing.T) {
	s := NewOnboardingSession(uuid.New())

	// Try to submit for a step that's not current and not previously answered
	err := s.SubmitAnswer(StepRestrictions, map[string]any{})
	if err == nil {
		t.Error("expected error for out-of-order step")
	}

	domErr, ok := err.(*DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}
	if domErr.Code != ErrStateConflct {
		t.Errorf("expected STATE_CONFLICT, got %s", domErr.Code)
	}
}

func TestOnboardingCompletedRejectsAnswer(t *testing.T) {
	s := NewOnboardingSession(uuid.New())
	for _, step := range MinimumOnboardingSteps {
		val := map[string]any{"test": true}
		if step == StepSelectGoals {
			val = map[string]any{"goals": []any{"sleep"}}
		}
		_ = s.SubmitAnswer(step, val)
	}

	// Force completed state
	s.Status = OnboardingCompleted
	err := s.SubmitAnswer(StepWelcome, map[string]any{})
	if err == nil {
		t.Error("expected error when onboarding is completed")
	}
}

func TestOnboardingPendingItemsForIntensiveGoals(t *testing.T) {
	s := NewOnboardingSession(uuid.New())

	for _, step := range MinimumOnboardingSteps {
		val := map[string]any{"test": true}
		if step == StepSelectGoals {
			val = map[string]any{"goals": []any{"english", "java", "sleep"}}
		}
		_ = s.SubmitAnswer(step, val)
	}

	if len(s.PendingItems) != 2 {
		t.Errorf("expected 2 pending items (english + java baselines), got %d", len(s.PendingItems))
	}
}
