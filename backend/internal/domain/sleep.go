package domain

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	EventSleepDiarySubmitted       EventType = "sleep_diary_submitted"
	EventSleepRoutineRecorded      EventType = "sleep_routine_recorded"
	EventSleepInterventionProposed EventType = "sleep_intervention_proposed"
	EventSleepInterventionAccepted EventType = "sleep_intervention_accepted"
	EventSleepInterventionRejected EventType = "sleep_intervention_rejected"
	EventSleepInterventionClosed   EventType = "sleep_intervention_closed"
)

type SleepDiaryEntry struct {
	EntryID            uuid.UUID `json:"entry_id"`
	UserID             uuid.UUID `json:"user_id"`
	TaskID             uuid.UUID `json:"task_id"`
	LocalDate          string    `json:"local_date"`
	SleptAt            *string   `json:"slept_at,omitempty"`
	WokeAt             *string   `json:"woke_at,omitempty"`
	Quality0_10        *int      `json:"quality_0_10,omitempty"`
	MorningEnergy0_10  *int      `json:"morning_energy_0_10,omitempty"`
	ComputedDurationMin *int     `json:"computed_duration_min,omitempty"`
	AwakeningsNote     *string   `json:"awakenings_note,omitempty"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}

type SleepRoutineRecord struct {
	RecordID  uuid.UUID `json:"record_id"`
	UserID    uuid.UUID `json:"user_id"`
	TaskID    uuid.UUID `json:"task_id"`
	LocalDate string    `json:"local_date"`
	Version   string    `json:"version"`
	StepsDone []string  `json:"steps_done"`
	Result    string    `json:"result"`
	NoteShort *string   `json:"note_short,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type WeeklySleepIntervention struct {
	InterventionID     uuid.UUID `json:"intervention_id"`
	UserID             uuid.UUID `json:"user_id"`
	WeekID             string    `json:"week_id"`
	Description        string    `json:"description"`
	WhyShort           string    `json:"why_short"`
	AdherenceRule      string    `json:"adherence_rule"`
	Status             string    `json:"status"`
	AdherenceCountDone int       `json:"adherence_count_done"`
	ClosingOutcome     *string   `json:"closing_outcome,omitempty"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
}

type SleepDiaryResult struct {
	EntryID            uuid.UUID `json:"entry_id"`
	Status             string    `json:"status"`
	ComputedDurationMin *int     `json:"computed_duration_min,omitempty"`
	InsightShort       string    `json:"insight_short"`
}

type SleepWeeklySummary struct {
	WeekStart          string                   `json:"week_start"`
	DaysWithDiary      int                      `json:"days_with_diary"`
	AvgQuality         *float64                 `json:"avg_quality,omitempty"`
	AvgEnergy          *float64                 `json:"avg_energy,omitempty"`
	AvgRegularityDelta *int                     `json:"avg_regularity_delta,omitempty"`
	Intervention       *WeeklySleepIntervention `json:"intervention,omitempty"`
	NextSuggestion     string                   `json:"next_suggestion"`
}

// --- Deterministic rules ---

// ComputeSleepDuration calculates sleep duration in minutes between two HH:MM times.
// Handles overnight sleep (e.g., 23:30 -> 07:00). Returns nil on parse failure.
func ComputeSleepDuration(sleptAt, wokeAt string) *int {
	sleptMin := parseTimeToMinutes(sleptAt)
	wokeMin := parseTimeToMinutes(wokeAt)
	if sleptMin < 0 || wokeMin < 0 {
		return nil
	}

	dur := wokeMin - sleptMin
	if dur <= 0 {
		dur += 24 * 60
	}
	return &dur
}

func parseTimeToMinutes(t string) int {
	t = strings.TrimSpace(t)
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return -1
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 || h > 23 {
		return -1
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return -1
	}
	return h*60 + m
}

// ClassifySleepDiaryStatus returns "complete" if >=2 fields are filled, "partial" otherwise.
func ClassifySleepDiaryStatus(sleptAt, wokeAt *string, quality, energy *int) string {
	count := 0
	if sleptAt != nil && *sleptAt != "" {
		count++
	}
	if wokeAt != nil && *wokeAt != "" {
		count++
	}
	if quality != nil {
		count++
	}
	if energy != nil {
		count++
	}
	if count >= 2 {
		return "complete"
	}
	return "partial"
}

// ComputeRegularityDelta calculates average deviation of sleep times (in minutes) vs median.
// Returns nil if fewer than 3 entries with slept_at data.
func ComputeRegularityDelta(entries []SleepDiaryEntry) *int {
	var minutes []int
	for _, e := range entries {
		if e.SleptAt != nil {
			m := parseTimeToMinutes(*e.SleptAt)
			if m >= 0 {
				if m < 12*60 {
					m += 24 * 60
				}
				minutes = append(minutes, m)
			}
		}
	}
	if len(minutes) < 3 {
		return nil
	}

	sum := 0
	for _, m := range minutes {
		sum += m
	}
	avg := sum / len(minutes)

	totalDev := 0
	for _, m := range minutes {
		diff := m - avg
		if diff < 0 {
			diff = -diff
		}
		totalDev += diff
	}
	result := int(math.Round(float64(totalDev) / float64(len(minutes))))
	return &result
}

// DefaultSleepRoutineSteps returns default pre-sleep routine steps for a given version.
func DefaultSleepRoutineSteps(version string) []string {
	if version == "minimal" {
		return []string{
			"Desligue telas 15min antes de dormir",
			"Faça 3 respirações profundas",
		}
	}
	return []string{
		"Desligue telas 30min antes de dormir",
		"Faça 1 atividade relaxante (leitura, música calma)",
		"Prepare o ambiente (escurecer, temperatura)",
		"Faça 5 respirações profundas",
	}
}

// SleepInterventionPool is the hardcoded pool of simple weekly interventions.
var SleepInterventionPool = []struct {
	Description   string
	WhyShort      string
	AdherenceRule string
}{
	{
		Description:   "Desligar telas 30 min antes de dormir por 3 dias esta semana",
		WhyShort:      "Luz azul suprime melatonina e atrasa o sono.",
		AdherenceRule: "3 dias",
	},
	{
		Description:   "Manter horário de dormir fixo (±30 min) por 4 dias esta semana",
		WhyShort:      "Regularidade de horário melhora a qualidade do sono.",
		AdherenceRule: "4 dias",
	},
	{
		Description:   "Evitar cafeína após 14h por 5 dias esta semana",
		WhyShort:      "Cafeína tem meia-vida de ~5h e pode atrapalhar o sono.",
		AdherenceRule: "5 dias",
	},
	{
		Description:   "Fazer 5 min de respiração antes de dormir por 3 dias esta semana",
		WhyShort:      "Respiração lenta ativa o parassimpático e facilita o sono.",
		AdherenceRule: "3 dias",
	},
	{
		Description:   "Reduzir luz do quarto 1h antes de dormir por 3 dias esta semana",
		WhyShort:      "Baixa luminosidade sinaliza ao corpo que é hora de dormir.",
		AdherenceRule: "3 dias",
	},
}

// ParseSleepDiaryInput parses user input for sleep diary in various formats.
// Supported: "23:30 07:00 8 7", "23:30, 07:00, 8, 7", "bem", "mal"
func ParseSleepDiaryInput(text string) (sleptAt, wokeAt *string, quality, energy *int) {
	text = strings.TrimSpace(text)
	lower := strings.ToLower(text)

	if lower == "bem" || lower == "mal" {
		q := 7
		e := 7
		if lower == "mal" {
			q = 3
			e = 3
		}
		return nil, nil, &q, &e
	}

	normalized := strings.ReplaceAll(text, ",", " ")
	parts := strings.Fields(normalized)

	if len(parts) >= 2 {
		s := parts[0]
		w := parts[1]
		sleptAt = &s
		wokeAt = &w
	}
	if len(parts) >= 3 {
		if v, err := strconv.Atoi(parts[2]); err == nil && v >= 0 && v <= 10 {
			quality = &v
		}
	}
	if len(parts) >= 4 {
		if v, err := strconv.Atoi(parts[3]); err == nil && v >= 0 && v <= 10 {
			energy = &v
		}
	}

	return
}
