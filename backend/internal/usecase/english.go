package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/abriesouza/super-assistente/internal/domain"
	"github.com/abriesouza/super-assistente/internal/port"
	"github.com/google/uuid"
)

type EnglishUseCase struct {
	users      port.UserRepository
	dailyState port.DailyStateRepository
	inputs     port.EnglishInputRepository
	retrievals port.EnglishRetrievalRepository
	errorLog   port.EnglishErrorLogRepository
	rubrics    port.RubricRepository
	gates      *GateUseCase
	events     port.EventRepository
	log        *slog.Logger
}

func NewEnglishUseCase(
	users port.UserRepository,
	dailyState port.DailyStateRepository,
	inputs port.EnglishInputRepository,
	retrievals port.EnglishRetrievalRepository,
	errorLog port.EnglishErrorLogRepository,
	rubrics port.RubricRepository,
	gates *GateUseCase,
	events port.EventRepository,
	log *slog.Logger,
) *EnglishUseCase {
	return &EnglishUseCase{
		users:      users,
		dailyState: dailyState,
		inputs:     inputs,
		retrievals: retrievals,
		errorLog:   errorLog,
		rubrics:    rubrics,
		gates:      gates,
		events:     events,
		log:        log,
	}
}

// SubmitInputCheck records an English input session and submits comprehension answers
// as evidence to the gate system. Each answer becomes one text_answer Evidence.
func (uc *EnglishUseCase) SubmitInputCheck(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	answers []string,
) (*domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	status := "complete"
	if len(answers) == 0 {
		status = "partial"
	}

	sess := &domain.EnglishInputSession{
		SessionID:            uuid.New(),
		UserID:               user.UserID,
		TaskID:               taskID,
		LocalDate:            localDate,
		ComprehensionAnswers: answers,
		Status:               status,
		CreatedAt:            time.Now(),
	}
	if err := uc.inputs.Save(ctx, sess); err != nil {
		return nil, fmt.Errorf("saving english input session: %w", err)
	}

	for _, answer := range answers {
		summary := answer
		if len(summary) > 100 {
			summary = summary[:100] + "..."
		}
		_, err := uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
			domain.EvidenceText, domain.SensitivityC2, summary, &answer)
		if err != nil {
			uc.log.Error("submitting input evidence", "error", err)
		}
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		return nil, fmt.Errorf("evaluating input gate: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventEnglishInputCompleted, domain.SensitivityC1,
		map[string]any{
			"task_id":       taskID.String(),
			"gate_status":   string(grView.GateStatus),
			"answers_count": len(answers),
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return grView, nil
}

// SubmitSpeakingRubric records a speaking rubric for a task.
// If audio was already submitted, it evaluates the gate. Otherwise returns a prompt.
func (uc *EnglishUseCase) SubmitSpeakingRubric(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	dimensions map[string]int,
) (*domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	total := 0
	for _, v := range dimensions {
		total += v
	}

	rubric := &domain.RubricScore{
		RubricID:   uuid.New(),
		UserID:     user.UserID,
		TaskID:     taskID,
		Domain:     domain.GoalEnglish,
		Dimensions: dimensions,
		Total:      total,
		Status:     "complete",
	}
	if err := uc.rubrics.Save(ctx, rubric); err != nil {
		return nil, fmt.Errorf("saving rubric: %w", err)
	}

	dimSummary := fmt.Sprintf("rubric total=%d", total)
	_, err = uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
		domain.EvidenceRubric, domain.SensitivityC2, dimSummary, nil)
	if err != nil {
		return nil, fmt.Errorf("submitting rubric evidence: %w", err)
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		return nil, fmt.Errorf("evaluating speaking gate: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventEnglishSpeakingCompleted, domain.SensitivityC1,
		map[string]any{
			"task_id":      taskID.String(),
			"gate_status":  string(grView.GateStatus),
			"rubric_total": total,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return grView, nil
}

// SubmitRetrieval records a retrieval session, classifies performance, and submits evidence.
func (uc *EnglishUseCase) SubmitRetrieval(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	itemsAnswered int,
	itemsTotal int,
	targets []string,
) (*domain.RetrievalResult, *domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	status := domain.ClassifyRetrievalStatus(itemsAnswered, itemsTotal)

	retrieval := &domain.EnglishRetrieval{
		RetrievalID:   uuid.New(),
		UserID:        user.UserID,
		TaskID:        taskID,
		LocalDate:     localDate,
		ItemsAnswered: itemsAnswered,
		ItemsTotal:    itemsTotal,
		Status:        status,
		Targets:       targets,
		CreatedAt:     time.Now(),
	}
	if err := uc.retrievals.Save(ctx, retrieval); err != nil {
		return nil, nil, fmt.Errorf("saving english retrieval: %w", err)
	}

	for i := 0; i < itemsAnswered; i++ {
		summary := fmt.Sprintf("recall item %d/%d", i+1, itemsTotal)
		if i < len(targets) {
			summary = targets[i]
		}
		_, err := uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
			domain.EvidenceText, domain.SensitivityC2, summary, nil)
		if err != nil {
			uc.log.Error("submitting retrieval evidence", "error", err)
		}
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		uc.log.Error("evaluating retrieval gate", "error", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventEnglishRetrievalCompleted, domain.SensitivityC1,
		map[string]any{
			"task_id":        taskID.String(),
			"items_answered": itemsAnswered,
			"items_total":    itemsTotal,
			"status":         status,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	result := &domain.RetrievalResult{
		Status:  status,
		Targets: targets,
	}
	return result, grView, nil
}

// LogErrorOfDay records an English error of the day and checks recurrence (>= 3 in 14 days).
func (uc *EnglishUseCase) LogErrorOfDay(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	label string,
	noteShort *string,
) (*domain.ErrorLogResult, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	sinceDate := subtractDays(localDate, 14)

	count, err := uc.errorLog.CountByUserLabelSince(ctx, user.UserID, label, sinceDate)
	if err != nil {
		count = 0
	}
	newCount := count + 1
	isRecurring := domain.CheckRecurrence(newCount)

	entry := &domain.EnglishErrorLogEntry{
		ErrorID:          uuid.New(),
		UserID:           user.UserID,
		LocalDate:        localDate,
		Label:            label,
		NoteShort:        noteShort,
		RecurringCount14d: newCount,
		IsRecurring:      isRecurring,
		CreatedAt:        time.Now(),
	}
	if err := uc.errorLog.Save(ctx, entry); err != nil {
		return nil, fmt.Errorf("saving english error log: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventEnglishErrorLogged, domain.SensitivityC1,
		map[string]any{
			"label":        label,
			"is_recurring": isRecurring,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return &domain.ErrorLogResult{
		ErrorID:     entry.ErrorID,
		Label:       label,
		IsRecurring: isRecurring,
		Count14d:    newCount,
	}, nil
}

// GetWeekTrend returns aggregated English metrics for a week.
func (uc *EnglishUseCase) GetWeekTrend(
	ctx context.Context,
	telegramUserID int64,
	weekStartDate string,
) (*domain.EnglishWeeklyTrend, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	endDate := addDays(weekStartDate, 6)

	retrievals, _ := uc.retrievals.FindByUserAndDateRange(ctx, user.UserID, weekStartDate, endDate)
	errors, _ := uc.errorLog.FindByUserAndDateRange(ctx, user.UserID, weekStartDate, endDate)

	okCount := 0
	for _, r := range retrievals {
		if r.Status == "ok" {
			okCount++
		}
	}

	okRate := 0.0
	if len(retrievals) > 0 {
		okRate = float64(okCount) / float64(len(retrievals))
	}

	var recurring []domain.EnglishErrorLogEntry
	for _, e := range errors {
		if e.IsRecurring {
			recurring = append(recurring, e)
		}
	}

	trend := &domain.EnglishWeeklyTrend{
		WeekStart:          weekStartDate,
		RetrievalCount:     len(retrievals),
		RetrievalOkRate:    okRate,
		TopRecurringErrors: recurring,
	}

	return trend, nil
}

func subtractDays(dateStr string, days int) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.AddDate(0, 0, -days).Format("2006-01-02")
}

func addDays(dateStr string, days int) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.AddDate(0, 0, days).Format("2006-01-02")
}
