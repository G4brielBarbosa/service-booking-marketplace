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

type DailyRoutineUseCase struct {
	users      port.UserRepository
	goals      port.GoalCycleRepository
	dailyState port.DailyStateRepository
	catalog    port.TaskCatalog
	events     port.EventRepository
	log        *slog.Logger
}

func NewDailyRoutineUseCase(
	users port.UserRepository,
	goals port.GoalCycleRepository,
	dailyState port.DailyStateRepository,
	catalog port.TaskCatalog,
	events port.EventRepository,
	log *slog.Logger,
) *DailyRoutineUseCase {
	return &DailyRoutineUseCase{
		users:      users,
		goals:      goals,
		dailyState: dailyState,
		catalog:    catalog,
		events:     events,
		log:        log,
	}
}

func (uc *DailyRoutineUseCase) SubmitDailyCheckIn(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	timeAvailMin int,
	energy int,
	moodStress *int,
	constraintText *string,
) (*domain.DailyPlanView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "Inicie com /start")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	isNew := false
	if err != nil || state == nil {
		state = domain.NewDailyState(user.UserID, localDate)
		isNew = true
	}

	checkIn := &domain.DailyCheckIn{
		CheckInID:      uuid.New(),
		UserID:         user.UserID,
		LocalDate:      localDate,
		TimeAvailMin:   timeAvailMin,
		Energy:         energy,
		MoodStress:     moodStress,
		ConstraintText: constraintText,
		CreatedAt:      time.Now(),
	}
	state.CheckIn = checkIn

	cycle, _ := uc.goals.FindByUserID(ctx, user.UserID)
	var activeGoals []domain.GoalEntry
	onboardingMissing := false
	if cycle != nil {
		activeGoals = cycle.ActiveGoals
	} else {
		onboardingMissing = true
		activeGoals = []domain.GoalEntry{{ID: domain.GoalHealth}}
	}

	planType, rationale := domain.SelectPlanType(timeAvailMin, energy)

	version := 1
	if state.Plan != nil {
		version = state.Plan.Version + 1
	}

	plan, tasks := domain.ComposePlan(
		user.UserID, localDate, planType, activeGoals,
		func(goal domain.GoalID, pt domain.PlanType) []domain.TaskTemplate {
			return uc.catalog.GetTasksForGoal(goal, pt)
		},
		version,
	)
	plan.Rationale = rationale

	state.Plan = plan
	state.Tasks = tasks
	state.UpdatedAt = time.Now()

	if isNew {
		if err := uc.dailyState.Save(ctx, state); err != nil {
			return nil, fmt.Errorf("saving daily state: %w", err)
		}
	} else {
		if err := uc.dailyState.Update(ctx, state); err != nil {
			return nil, fmt.Errorf("updating daily state: %w", err)
		}
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventDailyCheckInSubmitted, domain.SensitivityC1,
		map[string]any{
			"time_available_min": timeAvailMin,
			"energy":            energy,
			"plan_type":         string(planType),
			"onboarding_missing": onboardingMissing,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	evt2 := domain.NewEvent(user.UserID, domain.EventDailyPlanGenerated, domain.SensitivityC1,
		map[string]any{
			"plan_type":   string(planType),
			"task_count":  len(tasks),
			"version":     version,
		}, loc)
	_ = uc.events.Append(ctx, evt2)

	view := state.BuildPlanView()
	if onboardingMissing && view != nil {
		view.Rationale += " (sem onboarding — próximo passo: fazer onboarding)"
	}

	return view, nil
}

func (uc *DailyRoutineUseCase) GetTodayPlan(ctx context.Context, telegramUserID int64, localDate string) (*domain.DailyPlanView, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil || state.Plan == nil {
		return nil, domain.NewNotFoundError("Nenhum plano para hoje", "Faça o check-in com /checkin")
	}

	return state.BuildPlanView(), nil
}

func (uc *DailyRoutineUseCase) GetTodayStepsSummary(ctx context.Context, telegramUserID int64, localDate string) (*domain.DailyStepsSummary, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil {
		return nil, domain.NewNotFoundError("Nenhum registro para hoje", "Faça o check-in com /checkin")
	}

	return state.BuildStepsSummary(), nil
}

func (uc *DailyRoutineUseCase) ReplanDay(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	newTimeAvail *int,
	newEnergy *int,
	timeRemaining *int,
) (*domain.DailyPlanView, string, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, "", domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil || state.Plan == nil {
		return nil, "", domain.NewNotFoundError("Nenhum plano para replanejar", "Faça o check-in primeiro")
	}

	timeAvail := state.CheckIn.TimeAvailMin
	energy := state.CheckIn.Energy

	if newTimeAvail != nil {
		timeAvail = *newTimeAvail
	}
	if timeRemaining != nil {
		timeAvail = *timeRemaining
	}
	if newEnergy != nil {
		energy = *newEnergy
	}

	oldType := state.Plan.PlanType
	planType, rationale := domain.SelectPlanType(timeAvail, energy)
	version := state.Plan.Version + 1

	var doneTasks []domain.PlannedTask
	for _, t := range state.Tasks {
		if t.Status == domain.TaskCompleted || t.Status == domain.TaskEvidencePending ||
			t.Status == domain.TaskInProgress || t.Status == domain.TaskAttempt {
			doneTasks = append(doneTasks, t)
		}
	}

	cycle, _ := uc.goals.FindByUserID(ctx, user.UserID)
	var activeGoals []domain.GoalEntry
	if cycle != nil {
		activeGoals = cycle.ActiveGoals
	} else {
		activeGoals = []domain.GoalEntry{{ID: domain.GoalHealth}}
	}

	plan, tasks := domain.ComposePlan(
		user.UserID, localDate, planType, activeGoals,
		func(goal domain.GoalID, pt domain.PlanType) []domain.TaskTemplate {
			return uc.catalog.GetTasksForGoal(goal, pt)
		},
		version,
	)
	plan.Rationale = rationale

	allTasks := append(doneTasks, tasks...)
	state.Plan = plan
	state.Tasks = allTasks
	state.UpdatedAt = time.Now()

	if err := uc.dailyState.Update(ctx, state); err != nil {
		return nil, "", fmt.Errorf("updating daily state: %w", err)
	}

	loc := user.Location()
	evt := domain.NewEvent(user.UserID, domain.EventDayReplanned, domain.SensitivityC1,
		map[string]any{
			"old_plan_type": string(oldType),
			"new_plan_type": string(planType),
			"version":       version,
		}, loc)
	_ = uc.events.Append(ctx, evt)

	explanation := fmt.Sprintf("Plano ajustado de %s para %s. %s", oldType, planType, rationale)
	view := state.BuildPlanView()
	return view, explanation, nil
}

func (uc *DailyRoutineUseCase) UpdateTaskStatus(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	taskID uuid.UUID,
	action domain.TaskAction,
	actionContext string,
) (*domain.TaskStatusView, error) {
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

	oldStatus := task.Status
	view, err := task.ApplyAction(action, actionContext)
	if err != nil {
		return nil, err
	}

	state.UpdatedAt = time.Now()
	if err := uc.dailyState.Update(ctx, state); err != nil {
		return nil, fmt.Errorf("updating daily state: %w", err)
	}

	if task.Status != oldStatus {
		loc := user.Location()
		evt := domain.NewEvent(user.UserID, domain.EventTaskStatusChanged, domain.SensitivityC1,
			map[string]any{
				"task_id":    taskID.String(),
				"old_status": string(oldStatus),
				"new_status": string(task.Status),
				"action":     string(action),
			}, loc)
		_ = uc.events.Append(ctx, evt)
	}

	return view, nil
}

// FindTaskByTitle searches today's tasks by a partial title match.
func (uc *DailyRoutineUseCase) FindTaskByTitle(
	ctx context.Context,
	telegramUserID int64,
	localDate string,
	titleQuery string,
) (*domain.PlannedTask, error) {
	user, err := uc.users.FindByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, domain.NewNotFoundError("Usuário não encontrado", "")
	}

	state, err := uc.dailyState.FindByUserAndDate(ctx, user.UserID, localDate)
	if err != nil || state == nil {
		return nil, domain.NewNotFoundError("Nenhum plano para hoje", "")
	}

	for i := range state.Tasks {
		if containsIgnoreCase(state.Tasks[i].Title, titleQuery) {
			return &state.Tasks[i], nil
		}
	}

	return nil, domain.NewNotFoundError("Tarefa não encontrada", titleQuery)
}

func containsIgnoreCase(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) &&
		contains(toLower(s), toLower(substr))
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
