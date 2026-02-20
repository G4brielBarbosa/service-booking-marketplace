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

type JavaUseCase struct {
	users       port.UserRepository
	dailyState  port.DailyStateRepository
	practices   port.JavaPracticeRepository
	retrievals  port.JavaRetrievalRepository
	learningLog port.JavaLearningLogRepository
	gates       *GateUseCase
	events      port.EventRepository
	log         *slog.Logger
}

func NewJavaUseCase(
	users port.UserRepository,
	dailyState port.DailyStateRepository,
	practices port.JavaPracticeRepository,
	retrievals port.JavaRetrievalRepository,
	learningLog port.JavaLearningLogRepository,
	gates *GateUseCase,
	events port.EventRepository,
	log *slog.Logger,
) *JavaUseCase {
	return &JavaUseCase{
		users:       users,
		dailyState:  dailyState,
		practices:   practices,
		retrievals:  retrievals,
		learningLog: learningLog,
		gates:       gates,
		events:      events,
		log:         log,
	}
}

// SubmitPracticeEvidence records a Java practice session and submits evidence to the gate system.
func (uc *JavaUseCase) SubmitPracticeEvidence(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	evidenceShort string,
) (*domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	status := "complete"
	if evidenceShort == "" {
		status = "partial"
	}

	sess := &domain.JavaPracticeSession{
		SessionID:      uuid.New(),
		UserID:         user.UserID,
		TaskID:         taskID,
		LocalDate:      localDate,
		EvidenceShort:  evidenceShort,
		Status:         status,
		CreatedAt:      time.Now(),
	}
	if err := uc.practices.Save(ctx, sess); err != nil {
		return nil, fmt.Errorf("saving java practice session: %w", err)
	}

	summary := evidenceShort
	if len(summary) > 100 {
		summary = summary[:100] + "..."
	}
	_, err = uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
		domain.EvidenceText, domain.SensitivityC2, summary, &evidenceShort)
	if err != nil {
		uc.log.Error("submitting java practice evidence", "error", err)
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		return nil, fmt.Errorf("evaluating java practice gate: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventJavaPracticeSubmitted, domain.SensitivityC1,
		map[string]any{
			"task_id":     taskID.String(),
			"gate_status": string(grView.GateStatus),
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return grView, nil
}

// SubmitRetrieval records a Java retrieval session, classifies performance, and submits evidence.
func (uc *JavaUseCase) SubmitRetrieval(
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

	retrieval := &domain.JavaRetrieval{
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
		return nil, nil, fmt.Errorf("saving java retrieval: %w", err)
	}

	for i := 0; i < itemsAnswered; i++ {
		summary := fmt.Sprintf("recall item %d/%d", i+1, itemsTotal)
		if i < len(targets) {
			summary = targets[i]
		}
		_, err := uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
			domain.EvidenceText, domain.SensitivityC2, summary, nil)
		if err != nil {
			uc.log.Error("submitting java retrieval evidence", "error", err)
		}
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		uc.log.Error("evaluating java retrieval gate", "error", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventJavaRetrievalCompleted, domain.SensitivityC1,
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

// LogLearning records a Java learning log entry and checks recurrence (>= 3 in 14 days).
func (uc *JavaUseCase) LogLearning(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	localDate string,
	errorOrLearning string,
	fixOrNote *string,
	category *string,
) (*domain.LearningLogResult, *domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	sinceDate := subtractDays(localDate, 14)

	count, err := uc.learningLog.CountByUserLabelSince(ctx, user.UserID, errorOrLearning, sinceDate)
	if err != nil {
		count = 0
	}
	newCount := count + 1
	isRecurring := domain.CheckRecurrence(newCount)

	entry := &domain.JavaLearningLogEntry{
		EntryID:           uuid.New(),
		UserID:            user.UserID,
		TaskID:            taskID,
		LocalDate:         localDate,
		ErrorOrLearning:   errorOrLearning,
		FixOrNote:         fixOrNote,
		Category:          category,
		RecurringCount14d: newCount,
		IsRecurring:       isRecurring,
		CreatedAt:         time.Now(),
	}
	if err := uc.learningLog.Save(ctx, entry); err != nil {
		return nil, nil, fmt.Errorf("saving java learning log: %w", err)
	}

	summary := errorOrLearning
	if len(summary) > 100 {
		summary = summary[:100] + "..."
	}
	_, err = uc.gates.SubmitEvidence(ctx, telegramUserID, taskID,
		domain.EvidenceText, domain.SensitivityC2, summary, &errorOrLearning)
	if err != nil {
		uc.log.Error("submitting java learning evidence", "error", err)
	}

	grView, err := uc.gates.EvaluateGate(ctx, telegramUserID, localDate, taskID)
	if err != nil {
		uc.log.Error("evaluating java learning gate", "error", err)
	}

	categoryStr := ""
	if category != nil {
		categoryStr = *category
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventJavaLearningLogged, domain.SensitivityC1,
		map[string]any{
			"task_id":      taskID.String(),
			"category":     categoryStr,
			"is_recurring": isRecurring,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	logResult := &domain.LearningLogResult{
		EntryID:     entry.EntryID,
		Label:       errorOrLearning,
		IsRecurring: isRecurring,
		Count14d:    newCount,
	}
	return logResult, grView, nil
}

// GetWeekTrend returns aggregated Java metrics for a week.
func (uc *JavaUseCase) GetWeekTrend(
	ctx context.Context,
	telegramUserID int64,
	weekStartDate string,
) (*domain.JavaWeeklyTrend, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	endDate := addDays(weekStartDate, 6)

	retrievals, _ := uc.retrievals.FindByUserAndDateRange(ctx, user.UserID, weekStartDate, endDate)
	learningEntries, _ := uc.learningLog.FindByUserAndDateRange(ctx, user.UserID, weekStartDate, endDate)

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

	var recurring []domain.JavaLearningLogEntry
	for _, e := range learningEntries {
		if e.IsRecurring {
			recurring = append(recurring, e)
		}
	}

	trend := &domain.JavaWeeklyTrend{
		WeekStart:          weekStartDate,
		RetrievalCount:     len(retrievals),
		RetrievalOkRate:    okRate,
		TopRecurringErrors: recurring,
	}

	return trend, nil
}
