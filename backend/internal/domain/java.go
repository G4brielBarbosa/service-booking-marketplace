package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventJavaPracticeSubmitted  EventType = "java_practice_submitted"
	EventJavaRetrievalCompleted EventType = "java_retrieval_completed"
	EventJavaLearningLogged     EventType = "java_learning_logged"
)

type JavaPracticeSession struct {
	SessionID           uuid.UUID `json:"session_id"`
	UserID              uuid.UUID `json:"user_id"`
	TaskID              uuid.UUID `json:"task_id"`
	LocalDate           string    `json:"local_date"`
	DurationEstMin      int       `json:"duration_est_min"`
	ObjectiveConstraint string    `json:"objective_constraint"`
	EvidenceShort       string    `json:"evidence_short"`
	Status              string    `json:"status"` // "complete" | "partial"
	CreatedAt           time.Time `json:"created_at"`
}

type JavaRetrieval struct {
	RetrievalID   uuid.UUID `json:"retrieval_id"`
	UserID        uuid.UUID `json:"user_id"`
	TaskID        uuid.UUID `json:"task_id"`
	LocalDate     string    `json:"local_date"`
	ItemsAnswered int       `json:"items_answered"`
	ItemsTotal    int       `json:"items_total"`
	Status        string    `json:"status"` // "ok" | "low"
	Targets       []string  `json:"targets"`
	CreatedAt     time.Time `json:"created_at"`
}

type JavaLearningLogEntry struct {
	EntryID            uuid.UUID `json:"entry_id"`
	UserID             uuid.UUID `json:"user_id"`
	TaskID             uuid.UUID `json:"task_id"`
	LocalDate          string    `json:"local_date"`
	ErrorOrLearning    string    `json:"error_or_learning"`
	FixOrNote          *string   `json:"fix_or_note,omitempty"`
	Category           *string   `json:"category,omitempty"`
	RecurringCount14d  int       `json:"recurring_count_14d"`
	IsRecurring        bool      `json:"is_recurring"`
	CreatedAt          time.Time `json:"created_at"`
}

type LearningLogResult struct {
	EntryID     uuid.UUID `json:"entry_id"`
	Label       string    `json:"label"`
	IsRecurring bool      `json:"is_recurring"`
	Count14d    int       `json:"count_14d"`
}

type JavaWeeklyTrend struct {
	WeekStart          string                 `json:"week_start"`
	PracticeSessions   int                    `json:"practice_sessions"`
	RetrievalCount     int                    `json:"retrieval_count"`
	RetrievalOkRate    float64                `json:"retrieval_ok_rate"`
	TopRecurringErrors []JavaLearningLogEntry `json:"top_recurring_errors,omitempty"`
}
