package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventOnboardingStarted          EventType = "onboarding_started"
	EventOnboardingMinCompleted     EventType = "onboarding_minimum_completed"
	EventOnboardingCompleted        EventType = "onboarding_completed"
	EventOnboardingStepCompleted    EventType = "onboarding_step_completed"
	EventOnboardingResumed          EventType = "onboarding_resumed"
	EventOnboardingFieldRevised     EventType = "onboarding_field_revised"
	EventGoalCycleSet               EventType = "goal_cycle_set"
	EventMVDDefined                 EventType = "mvd_defined"
	EventPrivacyPolicySet           EventType = "privacy_policy_set"
)

// DomainEvent is the minimal event envelope (PLAN-000 §6).
// PayloadMin MUST NOT contain C3 sensitive content.
type DomainEvent struct {
	EventID     uuid.UUID        `json:"event_id"`
	UserID      uuid.UUID        `json:"user_id"`
	Timestamp   time.Time        `json:"timestamp"`
	LocalDate   string           `json:"local_date"`
	WeekID      string           `json:"week_id"`
	Type        EventType        `json:"type"`
	PayloadMin  map[string]any   `json:"payload_min"`
	Sensitivity SensitivityLevel `json:"sensitivity"`
}

func NewEvent(userID uuid.UUID, eventType EventType, sensitivity SensitivityLevel, payload map[string]any, loc *time.Location) DomainEvent {
	now := time.Now().In(loc)
	year, week := now.ISOWeek()

	return DomainEvent{
		EventID:     uuid.New(),
		UserID:      userID,
		Timestamp:   now,
		LocalDate:   now.Format("2006-01-02"),
		WeekID:      weekID(year, week),
		Type:        eventType,
		PayloadMin:  payload,
		Sensitivity: sensitivity,
	}
}

func weekID(year, week int) string {
	return fmt.Sprintf("%d-W%02d", year, week)
}
