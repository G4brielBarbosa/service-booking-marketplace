package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewEventEnvelope(t *testing.T) {
	userID := uuid.New()
	loc, _ := time.LoadLocation("America/Sao_Paulo")

	evt := NewEvent(userID, EventOnboardingStarted, SensitivityC5, map[string]any{"test": true}, loc)

	if evt.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, evt.UserID)
	}
	if evt.Type != EventOnboardingStarted {
		t.Errorf("expected type %s, got %s", EventOnboardingStarted, evt.Type)
	}
	if evt.Sensitivity != SensitivityC5 {
		t.Errorf("expected sensitivity C5, got %s", evt.Sensitivity)
	}
	if evt.LocalDate == "" {
		t.Error("expected local_date to be set")
	}
	if !strings.Contains(evt.WeekID, "-W") {
		t.Errorf("expected week_id format YYYY-WXX, got %s", evt.WeekID)
	}
	if evt.EventID == uuid.Nil {
		t.Error("expected event_id to be non-nil")
	}
}

func TestEventPayloadMinNoSensitive(t *testing.T) {
	evt := NewEvent(uuid.New(), EventOnboardingStepCompleted, SensitivityC5,
		map[string]any{"step_id": "welcome"}, time.UTC)

	if _, ok := evt.PayloadMin["step_id"]; !ok {
		t.Error("expected step_id in payload")
	}
}
