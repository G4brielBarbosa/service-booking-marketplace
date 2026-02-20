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

type GateUseCase struct {
	users       port.UserRepository
	dailyState  port.DailyStateRepository
	privacy     port.PrivacyPolicyRepository
	evidences   port.EvidenceRepository
	gateResults port.GateResultRepository
	rubrics     port.RubricRepository
	events      port.EventRepository
	log         *slog.Logger
}

func NewGateUseCase(
	users port.UserRepository,
	dailyState port.DailyStateRepository,
	privacy port.PrivacyPolicyRepository,
	evidences port.EvidenceRepository,
	gateResults port.GateResultRepository,
	rubrics port.RubricRepository,
	events port.EventRepository,
	log *slog.Logger,
) *GateUseCase {
	return &GateUseCase{
		users:       users,
		dailyState:  dailyState,
		privacy:     privacy,
		evidences:   evidences,
		gateResults: gateResults,
		rubrics:     rubrics,
		events:      events,
		log:         log,
	}
}

func (uc *GateUseCase) RequestTaskCompletion(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	taskID uuid.UUID,
) (*domain.CompletionRequestResult, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil {
		return nil, domain.NewNotFoundError("Nenhum plano para hoje", "")
	}

	task := state.FindTask(taskID)
	if task == nil {
		return nil, domain.NewNotFoundError("Tarefa não encontrada", taskID.String())
	}

	if task.Status == domain.TaskCompleted {
		return &domain.CompletionRequestResult{
			Status: domain.CompletionAlreadyCompleted,
		}, nil
	}

	if task.GateRef == nil {
		if task.Status != domain.TaskInProgress && task.Status != domain.TaskPlanned {
			return nil, domain.NewStateConflictError(
				"Tarefa não pode ser concluída neste estado",
				fmt.Sprintf("status atual: %s", task.Status),
			)
		}

		task.Status = domain.TaskCompleted
		task.UpdatedAt = time.Now()
		state.UpdatedAt = time.Now()
		if err := uc.dailyState.Update(ctx, state); err != nil {
			return nil, fmt.Errorf("updating daily state: %w", err)
		}

		return &domain.CompletionRequestResult{
			Status: domain.CompletionCompleted,
		}, nil
	}

	profile := domain.LookupGateProfile(*task.GateRef)
	if profile == nil {
		return nil, domain.NewNotFoundError("Gate profile não encontrado", *task.GateRef)
	}

	if task.Status != domain.TaskEvidencePending {
		if task.Status == domain.TaskInProgress || task.Status == domain.TaskPlanned {
			task.Status = domain.TaskEvidencePending
			task.UpdatedAt = time.Now()
			state.UpdatedAt = time.Now()
			if err := uc.dailyState.Update(ctx, state); err != nil {
				return nil, fmt.Errorf("updating daily state: %w", err)
			}
		}
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventGateRequired, domain.SensitivityC1,
		map[string]any{
			"task_id":    taskID.String(),
			"profile_id": profile.ProfileID,
			"domain":     string(profile.Domain),
		}, loc)
	_ = uc.events.Append(ctx, evt)

	privacyDisclosure := "Dados coletados para validar progresso. Conteúdo sensível retido por 7 dias (ajustável)."
	pp, _ := uc.privacy.FindByUserID(ctx, user.UserID)
	if pp != nil && pp.IsOptedOut(domain.SensitivityC3) {
		privacyDisclosure = "Modo mínimo ativo: conteúdo processado e descartado. Apenas o resultado do gate é mantido."
	}

	return &domain.CompletionRequestResult{
		Status: domain.CompletionEvidenceRequired,
		EvidenceRequest: &domain.EvidenceRequest{
			ProfileID:         profile.ProfileID,
			Requirements:      profile.Requirements,
			ValidityRules:     profile.ValidityRules,
			PrivacyDisclosure: privacyDisclosure,
		},
	}, nil
}

func (uc *GateUseCase) SubmitEvidence(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
	kind domain.EvidenceKind,
	sensitivity domain.SensitivityLevel,
	summary string,
	contentRef *string,
) (*domain.EvidenceReceipt, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	pp, _ := uc.privacy.FindByUserID(ctx, user.UserID)
	if pp == nil {
		def := domain.NewDefaultPrivacyPolicy(user.UserID)
		pp = &def
	}

	ev := domain.Evidence{
		EvidenceID:  uuid.New(),
		UserID:      user.UserID,
		TaskID:      taskID,
		Kind:        kind,
		Sensitivity: sensitivity,
		Summary:     summary,
		ContentRef:  contentRef,
		Timestamp:   time.Now(),
	}

	ev = domain.ApplyStoragePolicy(ev, *pp)

	if err := uc.evidences.Save(ctx, &ev); err != nil {
		return nil, fmt.Errorf("saving evidence: %w", err)
	}

	storedOrDiscarded := "kept"
	if ev.StoragePolicy == domain.StorageDiscardedAfterProc {
		storedOrDiscarded = "discarded"
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventEvidenceReceived, domain.SensitivityC1,
		map[string]any{
			"task_id":             taskID.String(),
			"kind":                string(kind),
			"valid":               summary != "",
			"stored_or_discarded": storedOrDiscarded,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return &domain.EvidenceReceipt{
		EvidenceID: ev.EvidenceID,
		Valid:      summary != "",
		Stored:     ev.StoragePolicy,
	}, nil
}

func (uc *GateUseCase) EvaluateGate(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	taskID uuid.UUID,
) (*domain.GateResultView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil {
		return nil, domain.NewNotFoundError("Nenhum plano para hoje", "")
	}

	task := state.FindTask(taskID)
	if task == nil {
		return nil, domain.NewNotFoundError("Tarefa não encontrada", taskID.String())
	}

	if task.GateRef == nil {
		return nil, domain.NewValidationError("Tarefa não possui gate", task.Title)
	}

	profile := domain.LookupGateProfile(*task.GateRef)
	if profile == nil {
		return nil, domain.NewNotFoundError("Gate profile não encontrado", *task.GateRef)
	}

	evidences, err := uc.evidences.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("loading evidences: %w", err)
	}

	pp, _ := uc.privacy.FindByUserID(ctx, user.UserID)
	if pp == nil {
		def := domain.NewDefaultPrivacyPolicy(user.UserID)
		pp = &def
	}

	result := domain.EvaluateGate(*profile, evidences, *pp)
	result.UserID = user.UserID
	result.TaskID = taskID

	if err := uc.gateResults.Save(ctx, &result); err != nil {
		return nil, fmt.Errorf("saving gate result: %w", err)
	}

	if result.GateStatus == domain.GateSatisfied {
		task.Status = domain.TaskCompleted
		task.UpdatedAt = time.Now()
		state.UpdatedAt = time.Now()
		if err := uc.dailyState.Update(ctx, state); err != nil {
			return nil, fmt.Errorf("updating daily state: %w", err)
		}
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventGateEvaluated, domain.SensitivityC1,
		map[string]any{
			"task_id":            taskID.String(),
			"gate_status":        string(result.GateStatus),
			"failure_reason_code": string(result.FailureReason),
		}, loc)
	_ = uc.events.Append(ctx, evt)

	return &domain.GateResultView{
		GateStatus:  result.GateStatus,
		ReasonShort: result.ReasonShort,
		NextMinStep: result.NextMinStep,
		Metrics:     result.DerivedMetrics,
	}, nil
}

func (uc *GateUseCase) GetGateStatus(
	ctx context.Context,
	telegramUserID int64,
	taskID uuid.UUID,
) (*domain.GateResultView, error) {
	_, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	gr, err := uc.gateResults.FindByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("loading gate result: %w", err)
	}
	if gr == nil {
		return nil, domain.NewNotFoundError("Nenhum resultado de gate encontrado", taskID.String())
	}

	return &domain.GateResultView{
		GateStatus:  gr.GateStatus,
		ReasonShort: gr.ReasonShort,
		NextMinStep: gr.NextMinStep,
		Metrics:     gr.DerivedMetrics,
	}, nil
}

func (uc *GateUseCase) GetTodayGateSummary(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
) ([]domain.GateSummaryItem, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil {
		return nil, domain.NewNotFoundError("Nenhum plano para hoje", "")
	}

	gateResults, err := uc.gateResults.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil {
		return nil, fmt.Errorf("loading gate results: %w", err)
	}

	grMap := make(map[uuid.UUID]*domain.GateResult)
	for i := range gateResults {
		grMap[gateResults[i].TaskID] = &gateResults[i]
	}

	var items []domain.GateSummaryItem
	for _, task := range state.Tasks {
		if task.GateRef == nil {
			continue
		}

		item := domain.GateSummaryItem{
			TaskID:    task.TaskID,
			TaskTitle: task.Title,
		}

		if gr, ok := grMap[task.TaskID]; ok {
			item.GateStatus = gr.GateStatus
			item.ReasonShort = gr.ReasonShort
			item.NextMinStep = gr.NextMinStep
		} else {
			item.GateStatus = domain.GateNotSatisfied
			item.ReasonShort = "Aguardando evidência."
			item.NextMinStep = "Envie a evidência para validar."
		}

		items = append(items, item)
	}

	return items, nil
}
