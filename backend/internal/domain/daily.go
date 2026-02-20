package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// --- Plan types ---

type PlanType string

const (
	PlanA PlanType = "A"
	PlanB PlanType = "B"
	PlanC PlanType = "C"
)

// --- Task status (PLAN-002 D-004) ---

type TaskStatus string

const (
	TaskPlanned          TaskStatus = "planned"
	TaskInProgress       TaskStatus = "in_progress"
	TaskCompleted        TaskStatus = "completed"
	TaskBlocked          TaskStatus = "blocked"
	TaskDeferred         TaskStatus = "deferred"
	TaskEvidencePending  TaskStatus = "evidence_pending"
	TaskAttempt          TaskStatus = "attempt"
)

// --- Task action commands ---

type TaskAction string

const (
	ActionStart          TaskAction = "start"
	ActionBlock          TaskAction = "block"
	ActionDefer          TaskAction = "defer"
	ActionMarkDoneReq    TaskAction = "mark_done_request"
	ActionAddNote        TaskAction = "add_note"
)

// --- Event types for daily routine (SPEC-016) ---

const (
	EventDailyCheckInSubmitted EventType = "daily_check_in_submitted"
	EventDailyPlanGenerated    EventType = "daily_plan_generated"
	EventDayReplanned          EventType = "day_replanned"
	EventTaskStatusChanged     EventType = "task_status_changed"
)

// --- Domain entities ---

type DailyCheckIn struct {
	CheckInID      uuid.UUID `json:"check_in_id"`
	UserID         uuid.UUID `json:"user_id"`
	LocalDate      string    `json:"local_date"`
	TimeAvailMin   int       `json:"time_available_min"`
	Energy         int       `json:"energy_0_10"`
	MoodStress     *int      `json:"mood_stress_0_10,omitempty"`
	ConstraintText *string   `json:"constraints_text,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type DailyPlan struct {
	PlanID              uuid.UUID   `json:"plan_id"`
	UserID              uuid.UUID   `json:"user_id"`
	LocalDate           string      `json:"local_date"`
	PlanType            PlanType    `json:"plan_type"`
	Rationale           string      `json:"rationale"`
	PriorityTaskID      uuid.UUID   `json:"priority_task_id"`
	ComplementaryIDs    []uuid.UUID `json:"complementary_task_ids"`
	FoundationTaskID    *uuid.UUID  `json:"foundation_task_id,omitempty"`
	Version             int         `json:"version"`
	CreatedAt           time.Time   `json:"created_at"`
}

type PlannedTask struct {
	TaskID       uuid.UUID  `json:"task_id"`
	UserID       uuid.UUID  `json:"user_id"`
	LocalDate    string     `json:"local_date"`
	Title        string     `json:"title"`
	GoalDomain   GoalID     `json:"goal_domain"`
	EstimatedMin int        `json:"estimated_min"`
	Instructions string     `json:"instructions"`
	DoneCriteria string     `json:"done_criteria"`
	Status       TaskStatus `json:"status"`
	BlockReason  *string    `json:"block_reason,omitempty"`
	Note         *string    `json:"note,omitempty"`
	GateRef      *string    `json:"gate_ref,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type DailyState struct {
	UserID    uuid.UUID    `json:"user_id"`
	LocalDate string       `json:"local_date"`
	CheckIn   *DailyCheckIn `json:"check_in,omitempty"`
	Plan      *DailyPlan   `json:"plan,omitempty"`
	Tasks     []PlannedTask `json:"tasks"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type TaskTemplate struct {
	Title        string  `json:"title"`
	GoalDomain   GoalID  `json:"goal_domain"`
	EstimatedMin int     `json:"estimated_min"`
	Instructions string  `json:"instructions"`
	DoneCriteria string  `json:"done_criteria"`
	GateProfile  *string `json:"gate_profile,omitempty"`
}

// --- View types (use case outputs) ---

type DailyPlanView struct {
	PlanType          PlanType          `json:"plan_type"`
	Rationale         string            `json:"rationale"`
	PriorityTask      PlannedTask       `json:"priority_task"`
	ComplementaryTasks []PlannedTask    `json:"complementary_tasks"`
	FoundationTask    *PlannedTask      `json:"foundation_task,omitempty"`
	TotalEstimatedMin int               `json:"total_estimated_min"`
}

type DailyStepsSummary struct {
	Done       []PlannedTask `json:"done"`
	Pending    []PlannedTask `json:"pending"`
	Blocked    []PlannedTask `json:"blocked"`
	InProgress []PlannedTask `json:"in_progress"`
}

type TaskStatusView struct {
	TaskID   uuid.UUID  `json:"task_id"`
	Status   TaskStatus `json:"status"`
	NextStep string     `json:"next_step,omitempty"`
}

// --- Deterministic rules (PLAN-002 §8) ---

func SelectPlanType(timeAvailMin, energy int) (PlanType, string) {
	if timeAvailMin <= 15 || energy <= 3 {
		return PlanC, fmt.Sprintf("Plano C porque tempo=%d e energia=%d", timeAvailMin, energy)
	}
	if timeAvailMin >= 60 && energy >= 7 {
		return PlanA, fmt.Sprintf("Plano A porque tempo=%d e energia=%d", timeAvailMin, energy)
	}
	return PlanB, fmt.Sprintf("Plano B porque tempo=%d e energia=%d", timeAvailMin, energy)
}

// NormalizeEnergy parses informal energy descriptions into 0-10 scale.
func NormalizeEnergy(text string) int {
	text = strings.TrimSpace(strings.ToLower(text))
	if text == "" {
		return 5
	}

	if n, err := strconv.Atoi(text); err == nil {
		if n < 0 {
			return 0
		}
		if n > 10 {
			return 10
		}
		return n
	}

	switch {
	case text == "péssimo" || text == "pessimo" || text == "horrível" || text == "horrivel":
		return 2
	case text == "meh" || text == "tanto faz" || text == "tanto_faz":
		return 4
	case text == "ok" || text == "okay" || text == "normal":
		return 5
	case text == "bem" || text == "bom" || text == "boa":
		return 7
	case text == "ótimo" || text == "otimo" || text == "excelente":
		return 9
	default:
		return 5
	}
}

// --- Plan composition (PLAN-002 §8) ---

// ComposePlan generates a DailyPlan selecting tasks from candidates.
// Rules:
//   - 1 priority (absolute), 1-2 complementary, 1 foundation
//   - Max 2 intensive goals per day
//   - Plan C: max 2 items (1 priority min + 1 foundation min)
func ComposePlan(
	userID uuid.UUID,
	localDate string,
	planType PlanType,
	activeGoals []GoalEntry,
	getCandidates func(goal GoalID, pt PlanType) []TaskTemplate,
	version int,
) (*DailyPlan, []PlannedTask) {
	now := time.Now()

	var intensiveGoals, foundationGoals []GoalEntry
	for _, g := range activeGoals {
		class := GoalClassifications[g.ID]
		switch class {
		case GoalClassIntensive:
			intensiveGoals = append(intensiveGoals, g)
		case GoalClassFoundation:
			foundationGoals = append(foundationGoals, g)
		default:
			intensiveGoals = append(intensiveGoals, g)
		}
	}

	if len(intensiveGoals) > MaxIntensiveGoals {
		intensiveGoals = intensiveGoals[:MaxIntensiveGoals]
	}

	var allTasks []PlannedTask
	var priorityTask *PlannedTask
	var complementary []PlannedTask
	var foundationTask *PlannedTask

	if len(intensiveGoals) > 0 {
		candidates := getCandidates(intensiveGoals[0].ID, planType)
		if len(candidates) > 0 {
			t := templateToTask(candidates[0], userID, localDate, now)
			priorityTask = &t
			allTasks = append(allTasks, t)
		}
	}

	if planType != PlanC {
		for i := 0; i < len(intensiveGoals) && len(complementary) < 2; i++ {
			candidates := getCandidates(intensiveGoals[i].ID, planType)
			for j, c := range candidates {
				if i == 0 && j == 0 {
					continue // already used as priority
				}
				if len(complementary) >= 2 {
					break
				}
				t := templateToTask(c, userID, localDate, now)
				complementary = append(complementary, t)
				allTasks = append(allTasks, t)
			}
		}

		// SaaS/weekly bet as complementary if room
		for _, g := range activeGoals {
			if GoalClassifications[g.ID] == GoalClassWeeklyBet && len(complementary) < 2 {
				candidates := getCandidates(g.ID, planType)
				if len(candidates) > 0 {
					t := templateToTask(candidates[0], userID, localDate, now)
					complementary = append(complementary, t)
					allTasks = append(allTasks, t)
				}
			}
		}
	}

	if len(foundationGoals) > 0 {
		candidates := getCandidates(foundationGoals[0].ID, planType)
		if len(candidates) > 0 {
			t := templateToTask(candidates[0], userID, localDate, now)
			foundationTask = &t
			allTasks = append(allTasks, t)
		}
	}

	if priorityTask == nil && len(allTasks) > 0 {
		priorityTask = &allTasks[0]
	}

	if priorityTask == nil {
		t := PlannedTask{
			TaskID:       uuid.New(),
			UserID:       userID,
			LocalDate:    localDate,
			Title:        "Check-in do dia",
			GoalDomain:   GoalHealth,
			EstimatedMin: 5,
			Instructions: "Registre como você está hoje",
			DoneCriteria: "Check-in registrado",
			Status:       TaskPlanned,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		priorityTask = &t
		allTasks = append(allTasks, t)
	}

	if planType == PlanC && len(allTasks) > 2 {
		allTasks = allTasks[:2]
		complementary = nil
		if len(allTasks) > 1 {
			foundationTask = &allTasks[1]
		} else {
			foundationTask = nil
		}
	}

	plan := &DailyPlan{
		PlanID:    uuid.New(),
		UserID:    userID,
		LocalDate: localDate,
		PlanType:  planType,
		Version:   version,
		CreatedAt: now,
	}

	plan.PriorityTaskID = priorityTask.TaskID

	var compIDs []uuid.UUID
	for _, c := range complementary {
		compIDs = append(compIDs, c.TaskID)
	}
	plan.ComplementaryIDs = compIDs

	if foundationTask != nil {
		fid := foundationTask.TaskID
		plan.FoundationTaskID = &fid
	}

	return plan, allTasks
}

func templateToTask(tmpl TaskTemplate, userID uuid.UUID, localDate string, now time.Time) PlannedTask {
	return PlannedTask{
		TaskID:       uuid.New(),
		UserID:       userID,
		LocalDate:    localDate,
		Title:        tmpl.Title,
		GoalDomain:   tmpl.GoalDomain,
		EstimatedMin: tmpl.EstimatedMin,
		Instructions: tmpl.Instructions,
		DoneCriteria: tmpl.DoneCriteria,
		Status:       TaskPlanned,
		GateRef:      tmpl.GateProfile,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// --- Task status transitions ---

var validTransitions = map[TaskStatus][]TaskStatus{
	TaskPlanned:         {TaskInProgress, TaskBlocked, TaskDeferred},
	TaskInProgress:      {TaskCompleted, TaskEvidencePending, TaskBlocked, TaskDeferred, TaskAttempt},
	TaskEvidencePending: {TaskCompleted},
	TaskBlocked:         {TaskPlanned, TaskInProgress, TaskDeferred},
	TaskDeferred:        {TaskPlanned, TaskInProgress},
}

func (t *PlannedTask) CanTransitionTo(target TaskStatus) bool {
	allowed, ok := validTransitions[t.Status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}

func (t *PlannedTask) ApplyAction(action TaskAction, context string) (*TaskStatusView, error) {
	var targetStatus TaskStatus
	var nextStep string

	switch action {
	case ActionStart:
		targetStatus = TaskInProgress
		nextStep = t.Instructions
	case ActionBlock:
		targetStatus = TaskBlocked
		if context != "" {
			t.BlockReason = &context
		}
		nextStep = "Tarefa bloqueada. Tente uma alternativa ou adie."
	case ActionDefer:
		targetStatus = TaskDeferred
		nextStep = "Tarefa adiada para outro momento."
	case ActionMarkDoneReq:
		if t.GateRef != nil {
			targetStatus = TaskEvidencePending
			nextStep = "Pendente de evidência para validação."
		} else {
			targetStatus = TaskCompleted
			nextStep = "Tarefa concluída!"
		}
	case ActionAddNote:
		if context != "" {
			t.Note = &context
		}
		t.UpdatedAt = time.Now()
		return &TaskStatusView{
			TaskID: t.TaskID,
			Status: t.Status,
		}, nil
	default:
		return nil, NewValidationError("Ação inválida", string(action))
	}

	if !t.CanTransitionTo(targetStatus) {
		return nil, NewStateConflictError(
			"Transição de status inválida",
			fmt.Sprintf("%s → %s", t.Status, targetStatus),
		)
	}

	t.Status = targetStatus
	t.UpdatedAt = time.Now()

	return &TaskStatusView{
		TaskID:   t.TaskID,
		Status:   t.Status,
		NextStep: nextStep,
	}, nil
}

// --- DailyState helpers ---

func NewDailyState(userID uuid.UUID, localDate string) *DailyState {
	now := time.Now()
	return &DailyState{
		UserID:    userID,
		LocalDate: localDate,
		Tasks:     []PlannedTask{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (ds *DailyState) FindTask(taskID uuid.UUID) *PlannedTask {
	for i := range ds.Tasks {
		if ds.Tasks[i].TaskID == taskID {
			return &ds.Tasks[i]
		}
	}
	return nil
}

func (ds *DailyState) BuildPlanView() *DailyPlanView {
	if ds.Plan == nil {
		return nil
	}

	view := &DailyPlanView{
		PlanType:  ds.Plan.PlanType,
		Rationale: ds.Plan.Rationale,
	}

	totalMin := 0
	for _, t := range ds.Tasks {
		totalMin += t.EstimatedMin

		if t.TaskID == ds.Plan.PriorityTaskID {
			view.PriorityTask = t
		}

		for _, cid := range ds.Plan.ComplementaryIDs {
			if t.TaskID == cid {
				view.ComplementaryTasks = append(view.ComplementaryTasks, t)
			}
		}

		if ds.Plan.FoundationTaskID != nil && t.TaskID == *ds.Plan.FoundationTaskID {
			cp := t
			view.FoundationTask = &cp
		}
	}
	view.TotalEstimatedMin = totalMin

	return view
}

func (ds *DailyState) BuildStepsSummary() *DailyStepsSummary {
	summary := &DailyStepsSummary{}

	for _, t := range ds.Tasks {
		switch t.Status {
		case TaskCompleted, TaskAttempt:
			summary.Done = append(summary.Done, t)
		case TaskBlocked:
			summary.Blocked = append(summary.Blocked, t)
		case TaskInProgress, TaskEvidencePending:
			summary.InProgress = append(summary.InProgress, t)
		default:
			summary.Pending = append(summary.Pending, t)
		}
	}

	return summary
}
