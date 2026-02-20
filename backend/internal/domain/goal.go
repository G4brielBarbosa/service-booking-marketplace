package domain

import (
	"time"

	"github.com/google/uuid"
)

type GoalID string

const (
	GoalEnglish    GoalID = "english"
	GoalJava       GoalID = "java"
	GoalSleep      GoalID = "sleep"
	GoalHealth     GoalID = "health"
	GoalSelfEsteem GoalID = "self_esteem"
	GoalSaaS       GoalID = "saas"
)

type GoalClassification string

const (
	GoalClassIntensive GoalClassification = "intensive"
	GoalClassFoundation GoalClassification = "foundation"
	GoalClassWeeklyBet GoalClassification = "weekly_bet"
)

const MaxIntensiveGoals = 2

var GoalClassifications = map[GoalID]GoalClassification{
	GoalEnglish:    GoalClassIntensive,
	GoalJava:       GoalClassIntensive,
	GoalSleep:      GoalClassFoundation,
	GoalHealth:     GoalClassFoundation,
	GoalSelfEsteem: GoalClassFoundation,
	GoalSaaS:       GoalClassWeeklyBet,
}

var AllGoals = []GoalID{GoalEnglish, GoalJava, GoalSleep, GoalHealth, GoalSelfEsteem, GoalSaaS}

type GoalEntry struct {
	ID       GoalID `json:"id"`
	Reason   string `json:"reason,omitempty"`
}

type ActiveGoalCycle struct {
	CycleID        uuid.UUID   `json:"cycle_id"`
	UserID         uuid.UUID   `json:"user_id"`
	ActiveGoals    []GoalEntry `json:"active_goals"`
	PausedGoals    []GoalEntry `json:"paused_goals"`
	StartedAt      time.Time   `json:"started_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func NewGoalCycle(userID uuid.UUID, activeGoals []GoalEntry) (*ActiveGoalCycle, error) {
	intensive := countIntensive(activeGoals)
	if intensive > MaxIntensiveGoals {
		return nil, NewValidationError(
			"Limite de metas intensivas excedido",
			"No máximo 2 metas intensivas por ciclo. Escolha quais manter e quais pausar.",
		)
	}

	now := time.Now()
	return &ActiveGoalCycle{
		CycleID:     uuid.New(),
		UserID:      userID,
		ActiveGoals: activeGoals,
		PausedGoals: []GoalEntry{},
		StartedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (c *ActiveGoalCycle) IntensiveGoals() []GoalEntry {
	var result []GoalEntry
	for _, g := range c.ActiveGoals {
		if GoalClassifications[g.ID] == GoalClassIntensive {
			result = append(result, g)
		}
	}
	return result
}

func countIntensive(goals []GoalEntry) int {
	count := 0
	for _, g := range goals {
		if GoalClassifications[g.ID] == GoalClassIntensive {
			count++
		}
	}
	return count
}
