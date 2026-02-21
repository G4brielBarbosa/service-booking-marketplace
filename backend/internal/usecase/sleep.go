package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/abriesouza/super-assistente/internal/port"
	"github.com/google/uuid"
)

type SleepUseCase struct {
	users         port.UserRepository
	dailyState    port.DailyStateRepository
	diaries       port.SleepDiaryRepository
	routines      port.SleepRoutineRepository
	interventions port.SleepInterventionRepository
	gates         *GateUseCase
	events        port.EventRepository
	log           *slog.Logger
}

func NewSleepUseCase(
	users port.UserRepository,
	dailyState port.DailyStateRepository,
	diaries port.SleepDiaryRepository,
	routines port.SleepRoutineRepository,
	interventions port.SleepInterventionRepository,
	gates *GateUseCase,
	events port.EventRepository,
	log *slog.Logger,
) *SleepUseCase {
	return &SleepUseCase{
		users:         users,
		dailyState:    dailyState,
		diaries:       diaries,
		routines:      routines,
		interventions: interventions,
		gates:         gates,
		events:        events,
		log:           log,
	}
}

// SubmitSleepDiary records a sleep diary entry, computes duration/status, and evaluates the gate.
func (uc *SleepUseCase) SubmitSleepDiary(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	sleptAt *string,
	wokeAt *string,
	quality *int,
	energy *int,
	awakeningsNote *string,
) (*domain.SleepDiaryResult, *domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	status := domain.ClassifySleepDiaryStatus(sleptAt, wokeAt, quality, energy)

	var computedDuration *int
	if sleptAt != nil && wokeAt != nil {
		computedDuration = domain.ComputeSleepDuration(*sleptAt, *wokeAt)
	}

	entry := &domain.SleepDiaryEntry{
		EntryID:             uuid.New(),
		UserID:              user.UserID,
		TaskID:              taskID,
		LocalDate:           localDate,
		SleptAt:             sleptAt,
		WokeAt:              wokeAt,
		Quality0_10:         quality,
		MorningEnergy0_10:   energy,
		ComputedDurationMin: computedDuration,
		AwakeningsNote:      awakeningsNote,
		Status:              status,
		CreatedAt:           time.Now(),
	}
	if err := uc.diaries.Save(ctx, entry); err != nil {
		return nil, nil, fmt.Errorf("saving sleep diary entry: %w", err)
	}

	summary := buildDiarySummary(sleptAt, wokeAt, quality, energy)
	_, err = uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
		domain.EvidenceMetadata, domain.SensitivityC1, summary, nil)
	if err != nil {
		uc.log.Error("submitting sleep diary evidence", "error", err)
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		return nil, nil, fmt.Errorf("evaluating sleep diary gate: %w", err)
	}

	hasDuration := computedDuration != nil
	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventSleepDiarySubmitted, domain.SensitivityC1,
		map[string]any{
			"task_id":      taskID.String(),
			"status":       status,
			"has_duration": hasDuration,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	insight := "Sono registrado."
	if computedDuration != nil {
		hours := *computedDuration / 60
		mins := *computedDuration % 60
		insight = fmt.Sprintf("Duração estimada: %dh%02dmin.", hours, mins)
	}

	result := &domain.SleepDiaryResult{
		EntryID:             entry.EntryID,
		Status:              status,
		ComputedDurationMin: computedDuration,
		InsightShort:        insight,
	}
	return result, grView, nil
}

// RecordSleepRoutine records a pre-sleep routine execution and evaluates the gate.
func (uc *SleepUseCase) RecordSleepRoutine(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	stepsDone []string,
	result string,
	noteShort *string,
) (*domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	version := "normal"
	if len(stepsDone) <= 2 {
		version = "minimal"
	}

	record := &domain.SleepRoutineRecord{
		RecordID:  uuid.New(),
		UserID:    user.UserID,
		TaskID:    taskID,
		LocalDate: localDate,
		Version:   version,
		StepsDone: stepsDone,
		Result:    result,
		NoteShort: noteShort,
		CreatedAt: time.Now(),
	}
	if err := uc.routines.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("saving sleep routine record: %w", err)
	}

	summary := fmt.Sprintf("Rotina pré-sono: %d passos, resultado=%s", len(stepsDone), result)
	stepsText := ""
	for _, s := range stepsDone {
		if stepsText != "" {
			stepsText += "; "
		}
		stepsText += s
	}
	_, err = uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
		domain.EvidenceText, domain.SensitivityC1, summary, &stepsText)
	if err != nil {
		uc.log.Error("submitting sleep routine evidence", "error", err)
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		return nil, fmt.Errorf("evaluating sleep routine gate: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventSleepRoutineRecorded, domain.SensitivityC1,
		map[string]any{
			"task_id": taskID.String(),
			"result":  result,
			"version": version,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return grView, nil
}

// ProposeWeeklyIntervention creates a new weekly sleep intervention from the hardcoded pool.
func (uc *SleepUseCase) ProposeWeeklyIntervention(
	ctx context.Context,
	telegramUserID int64,
	weekID string,
) (*domain.WeeklySleepIntervention, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	existing, _ := uc.interventions.FindByUserAndWeek(ctx, user.UserID, weekID)
	if existing != nil {
		return existing, nil
	}

	pool := domain.SleepInterventionPool
	idx := rand.Intn(len(pool))
	chosen := pool[idx]

	intervention := &domain.WeeklySleepIntervention{
		InterventionID:     uuid.New(),
		UserID:             user.UserID,
		WeekID:             weekID,
		Description:        chosen.Description,
		WhyShort:           chosen.WhyShort,
		AdherenceRule:      chosen.AdherenceRule,
		Status:             "proposed",
		AdherenceCountDone: 0,
		CreatedAt:          time.Now(),
	}
	if err := uc.interventions.Save(ctx, intervention); err != nil {
		return nil, fmt.Errorf("saving weekly intervention: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventSleepInterventionProposed, domain.SensitivityC1,
		map[string]any{"week_id": weekID}, loc)
	_ = uc.events.Append(ctx, evt)

	return intervention, nil
}

// AcceptOrRejectIntervention updates the status of a weekly intervention.
func (uc *SleepUseCase) AcceptOrRejectIntervention(
	ctx context.Context,
	telegramUserID int64,
	weekID string,
	accepted bool,
) error {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return domain.NewNotFoundError("Usuário não encontrado", "")
	}

	intervention, err := uc.interventions.FindByUserAndWeek(ctx, user.UserID, weekID)
	if err != nil || intervention == nil {
		return domain.NewNotFoundError("Intervenção não encontrada para esta semana", "")
	}

	if accepted {
		intervention.Status = "accepted"
	} else {
		intervention.Status = "rejected"
	}
	if err := uc.interventions.Update(ctx, intervention); err != nil {
		return fmt.Errorf("updating intervention: %w", err)
	}

	evtType := domain.EventSleepInterventionAccepted
	if !accepted {
		evtType = domain.EventSleepInterventionRejected
	}
	loc := user.Location()
	evt := domain.NewEvent(user.UserID, evtType, domain.SensitivityC1,
		map[string]any{"week_id": weekID}, loc)
	_ = uc.events.Append(ctx, evt)

	return nil
}

// RecordInterventionAdherence increments the adherence counter for the current week's intervention.
func (uc *SleepUseCase) RecordInterventionAdherence(
	ctx context.Context,
	telegramUserID int64,
	weekID string,
) error {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return domain.NewNotFoundError("Usuário não encontrado", "")
	}

	intervention, err := uc.interventions.FindByUserAndWeek(ctx, user.UserID, weekID)
	if err != nil || intervention == nil {
		return domain.NewNotFoundError("Intervenção não encontrada para esta semana", "")
	}

	intervention.AdherenceCountDone++
	return uc.interventions.Update(ctx, intervention)
}

// CloseWeeklyIntervention closes the weekly intervention with an outcome.
func (uc *SleepUseCase) CloseWeeklyIntervention(
	ctx context.Context,
	telegramUserID int64,
	weekID string,
	outcome string,
) error {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return domain.NewNotFoundError("Usuário não encontrado", "")
	}

	intervention, err := uc.interventions.FindByUserAndWeek(ctx, user.UserID, weekID)
	if err != nil || intervention == nil {
		return domain.NewNotFoundError("Intervenção não encontrada para esta semana", "")
	}

	now := time.Now()
	intervention.ClosingOutcome = &outcome
	intervention.ClosedAt = &now
	if err := uc.interventions.Update(ctx, intervention); err != nil {
		return fmt.Errorf("closing intervention: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventSleepInterventionClosed, domain.SensitivityC1,
		map[string]any{"week_id": weekID, "outcome": outcome}, loc)
	_ = uc.events.Append(ctx, evt)

	return nil
}

// GetSleepWeekSummary returns aggregated sleep metrics for a given week.
func (uc *SleepUseCase) GetSleepWeekSummary(
	ctx context.Context,
	telegramUserID int64,
	weekStartDate string,
) (*domain.SleepWeeklySummary, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	endDate := addDays(weekStartDate, 6)

	diaries, _ := uc.diaries.FindByUserAndDateRange(ctx, user.UserID, weekStartDate, endDate)

	uniqueDates := map[string]bool{}
	var qualitySum, energySum float64
	var qualityCount, energyCount int

	for _, d := range diaries {
		uniqueDates[d.LocalDate] = true
		if d.Quality0_10 != nil {
			qualitySum += float64(*d.Quality0_10)
			qualityCount++
		}
		if d.MorningEnergy0_10 != nil {
			energySum += float64(*d.MorningEnergy0_10)
			energyCount++
		}
	}

	summary := &domain.SleepWeeklySummary{
		WeekStart:     weekStartDate,
		DaysWithDiary: len(uniqueDates),
	}

	if qualityCount >= 3 {
		avg := qualitySum / float64(qualityCount)
		summary.AvgQuality = &avg
	}
	if energyCount >= 3 {
		avg := energySum / float64(energyCount)
		summary.AvgEnergy = &avg
	}

	summary.AvgRegularityDelta = domain.ComputeRegularityDelta(diaries)

	intervention, _ := uc.interventions.FindByUserAndWeek(ctx, user.UserID, weekIDFromDate(weekStartDate))
	summary.Intervention = intervention

	summary.NextSuggestion = buildNextSuggestion(summary)

	return summary, nil
}

func buildDiarySummary(sleptAt, wokeAt *string, quality, energy *int) string {
	parts := []string{}
	if sleptAt != nil {
		parts = append(parts, fmt.Sprintf("dormiu=%s", *sleptAt))
	}
	if wokeAt != nil {
		parts = append(parts, fmt.Sprintf("acordou=%s", *wokeAt))
	}
	if quality != nil {
		parts = append(parts, fmt.Sprintf("qualidade=%d", *quality))
	}
	if energy != nil {
		parts = append(parts, fmt.Sprintf("energia=%d", *energy))
	}
	if len(parts) == 0 {
		return "registro de sono"
	}
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

func buildNextSuggestion(s *domain.SleepWeeklySummary) string {
	if s.DaysWithDiary < 3 {
		return "Registre sono em pelo menos 3 dias para ver tendências."
	}
	if s.AvgQuality != nil && *s.AvgQuality < 5.0 {
		return "Qualidade média baixa. Considere ajustar a rotina pré-sono."
	}
	if s.AvgRegularityDelta != nil && *s.AvgRegularityDelta > 60 {
		return "Horário de dormir irregular. Tente manter ±30min do mesmo horário."
	}
	return "Continue registrando para acompanhar tendências."
}

func weekIDFromDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}
