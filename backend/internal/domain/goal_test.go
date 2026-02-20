package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewGoalCycleWithinLimit(t *testing.T) {
	entries := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalJava},
		{ID: GoalSleep},
	}

	cycle, err := NewGoalCycle(uuid.New(), entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cycle.ActiveGoals) != 3 {
		t.Errorf("expected 3 active goals, got %d", len(cycle.ActiveGoals))
	}

	intensive := cycle.IntensiveGoals()
	if len(intensive) != 2 {
		t.Errorf("expected 2 intensive goals, got %d", len(intensive))
	}
}

func TestNewGoalCycleExceedsLimit(t *testing.T) {
	entries := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalJava},
		{ID: GoalID("rust")}, // hypothetical third intensive
	}

	// Register "rust" as intensive for this test
	GoalClassifications[GoalID("rust")] = GoalClassIntensive
	defer delete(GoalClassifications, GoalID("rust"))

	_, err := NewGoalCycle(uuid.New(), entries)
	if err == nil {
		t.Fatal("expected error for exceeding intensive limit")
	}

	domErr, ok := err.(*DomainError)
	if !ok {
		t.Fatalf("expected DomainError, got %T", err)
	}
	if domErr.Code != ErrValidation {
		t.Errorf("expected VALIDATION_ERROR, got %s", domErr.Code)
	}
}

func TestGoalClassificationDefaults(t *testing.T) {
	tests := []struct {
		goal GoalID
		want GoalClassification
	}{
		{GoalEnglish, GoalClassIntensive},
		{GoalJava, GoalClassIntensive},
		{GoalSleep, GoalClassFoundation},
		{GoalHealth, GoalClassFoundation},
		{GoalSelfEsteem, GoalClassFoundation},
		{GoalSaaS, GoalClassWeeklyBet},
	}

	for _, tt := range tests {
		got := GoalClassifications[tt.goal]
		if got != tt.want {
			t.Errorf("goal %s: expected %s, got %s", tt.goal, tt.want, got)
		}
	}
}
