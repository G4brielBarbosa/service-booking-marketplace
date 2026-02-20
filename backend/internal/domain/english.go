package domain

import (
	"time"

	"github.com/google/uuid"
)

// --- Event types for English daily (SPEC-016, no C3 content) ---

const (
	EventEnglishInputCompleted     EventType = "english_input_completed"
	EventEnglishSpeakingCompleted  EventType = "english_speaking_completed"
	EventEnglishRetrievalCompleted EventType = "english_retrieval_completed"
	EventEnglishErrorLogged        EventType = "english_error_logged"
)

// --- English Input Session ---

type EnglishInputSession struct {
	SessionID            uuid.UUID `json:"session_id"`
	UserID               uuid.UUID `json:"user_id"`
	TaskID               uuid.UUID `json:"task_id"`
	LocalDate            string    `json:"local_date"`
	DurationEstMin       int       `json:"duration_est_min"`
	ContentDescriptor    string    `json:"content_descriptor"`
	ComprehensionAnswers []string  `json:"comprehension_answers"`
	Status               string    `json:"status"` // "complete" | "partial"
	CreatedAt            time.Time `json:"created_at"`
}

// --- English Retrieval ---

type EnglishRetrieval struct {
	RetrievalID  uuid.UUID `json:"retrieval_id"`
	UserID       uuid.UUID `json:"user_id"`
	TaskID       uuid.UUID `json:"task_id"`
	LocalDate    string    `json:"local_date"`
	ItemsAnswered int      `json:"items_answered"`
	ItemsTotal   int       `json:"items_total"`
	Status       string    `json:"status"` // "ok" | "low"
	Targets      []string  `json:"targets"`
	CreatedAt    time.Time `json:"created_at"`
}

// --- English Error Log Entry ---

type EnglishErrorLogEntry struct {
	ErrorID          uuid.UUID `json:"error_id"`
	UserID           uuid.UUID `json:"user_id"`
	LocalDate        string    `json:"local_date"`
	Label            string    `json:"label"`
	NoteShort        *string   `json:"note_short,omitempty"`
	RecurringCount14d int      `json:"recurring_count_14d"`
	IsRecurring      bool      `json:"is_recurring"`
	CreatedAt        time.Time `json:"created_at"`
}

// --- Use case result types ---

type RetrievalResult struct {
	Status  string   `json:"status"`
	Targets []string `json:"targets,omitempty"`
}

type ErrorLogResult struct {
	ErrorID     uuid.UUID `json:"error_id"`
	Label       string    `json:"label"`
	IsRecurring bool      `json:"is_recurring"`
	Count14d    int       `json:"count_14d"`
}

type EnglishWeeklyTrend struct {
	WeekStart          string             `json:"week_start"`
	InputSessions      int                `json:"input_sessions"`
	SpeakingCount      int                `json:"speaking_count"`
	RetrievalCount     int                `json:"retrieval_count"`
	RetrievalOkRate    float64            `json:"retrieval_ok_rate"`
	AvgRubricTotal     float64            `json:"avg_rubric_total"`
	TopRecurringErrors []EnglishErrorLogEntry `json:"top_recurring_errors,omitempty"`
}

// --- Deterministic rules (PLAN-004 §8) ---

// ClassifyRetrievalStatus returns "ok" if itemsAnswered >= 70% of itemsTotal, "low" otherwise.
// Edge case: 0/0 is treated as "ok" (nothing was required).
func ClassifyRetrievalStatus(itemsAnswered, itemsTotal int) string {
	if itemsTotal <= 0 {
		return "ok"
	}
	if itemsAnswered*10 >= itemsTotal*7 {
		return "ok"
	}
	return "low"
}

// CheckRecurrence returns true if the existing 14-day count (including current) is >= 3.
func CheckRecurrence(existingCount int) bool {
	return existingCount >= 3
}
