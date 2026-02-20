package domain

import (
	"testing"

	"github.com/google/uuid"
)

// --- SelectPlanType tests (PLAN-002 §8) ---

func TestSelectPlanType_PlanC_LowEnergy(t *testing.T) {
	pt, rationale := SelectPlanType(60, 3)
	if pt != PlanC {
		t.Errorf("expected PlanC for energy=3, got %s", pt)
	}
	if rationale == "" {
		t.Error("rationale should not be empty")
	}
}

func TestSelectPlanType_PlanC_LowTime(t *testing.T) {
	pt, _ := SelectPlanType(15, 7)
	if pt != PlanC {
		t.Errorf("expected PlanC for time=15, got %s", pt)
	}
}

func TestSelectPlanType_PlanC_BothLow(t *testing.T) {
	pt, _ := SelectPlanType(10, 2)
	if pt != PlanC {
		t.Errorf("expected PlanC for time=10 energy=2, got %s", pt)
	}
}

func TestSelectPlanType_PlanA(t *testing.T) {
	pt, _ := SelectPlanType(60, 7)
	if pt != PlanA {
		t.Errorf("expected PlanA for time=60 energy=7, got %s", pt)
	}
}

func TestSelectPlanType_PlanA_HighValues(t *testing.T) {
	pt, _ := SelectPlanType(120, 10)
	if pt != PlanA {
		t.Errorf("expected PlanA for time=120 energy=10, got %s", pt)
	}
}

func TestSelectPlanType_PlanB_MediumValues(t *testing.T) {
	pt, _ := SelectPlanType(30, 5)
	if pt != PlanB {
		t.Errorf("expected PlanB for time=30 energy=5, got %s", pt)
	}
}

func TestSelectPlanType_PlanB_HighTimeModerateEnergy(t *testing.T) {
	pt, _ := SelectPlanType(60, 6)
	if pt != PlanB {
		t.Errorf("expected PlanB for time=60 energy=6, got %s", pt)
	}
}

func TestSelectPlanType_PlanB_ModerateTimeLowishEnergy(t *testing.T) {
	pt, _ := SelectPlanType(45, 4)
	if pt != PlanB {
		t.Errorf("expected PlanB for time=45 energy=4, got %s", pt)
	}
}

func TestSelectPlanType_PlanB_NotQuitePlanA(t *testing.T) {
	pt, _ := SelectPlanType(59, 7)
	if pt != PlanB {
		t.Errorf("expected PlanB for time=59 energy=7, got %s", pt)
	}
}

// --- NormalizeEnergy tests ---

func TestNormalizeEnergy_Meh(t *testing.T) {
	if got := NormalizeEnergy("meh"); got != 4 {
		t.Errorf("expected 4 for 'meh', got %d", got)
	}
}

func TestNormalizeEnergy_TantoFaz(t *testing.T) {
	if got := NormalizeEnergy("tanto faz"); got != 4 {
		t.Errorf("expected 4 for 'tanto faz', got %d", got)
	}
}

func TestNormalizeEnergy_Ok(t *testing.T) {
	if got := NormalizeEnergy("ok"); got != 5 {
		t.Errorf("expected 5 for 'ok', got %d", got)
	}
}

func TestNormalizeEnergy_Bem(t *testing.T) {
	if got := NormalizeEnergy("bem"); got != 7 {
		t.Errorf("expected 7 for 'bem', got %d", got)
	}
}

func TestNormalizeEnergy_Pessimo(t *testing.T) {
	if got := NormalizeEnergy("péssimo"); got != 2 {
		t.Errorf("expected 2 for 'péssimo', got %d", got)
	}
}

func TestNormalizeEnergy_Numeric(t *testing.T) {
	if got := NormalizeEnergy("7"); got != 7 {
		t.Errorf("expected 7 for '7', got %d", got)
	}
}

func TestNormalizeEnergy_Empty(t *testing.T) {
	if got := NormalizeEnergy(""); got != 5 {
		t.Errorf("expected 5 for empty, got %d", got)
	}
}

func TestNormalizeEnergy_Unknown(t *testing.T) {
	if got := NormalizeEnergy("xyz"); got != 5 {
		t.Errorf("expected 5 for unknown, got %d", got)
	}
}

func TestNormalizeEnergy_ClampHigh(t *testing.T) {
	if got := NormalizeEnergy("15"); got != 10 {
		t.Errorf("expected 10 for 15, got %d", got)
	}
}

func TestNormalizeEnergy_ClampLow(t *testing.T) {
	if got := NormalizeEnergy("-3"); got != 0 {
		t.Errorf("expected 0 for -3, got %d", got)
	}
}

// --- ComposePlan tests ---

func TestComposePlan_PlanA_FullStructure(t *testing.T) {
	goals := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalJava},
		{ID: GoalSleep},
	}

	catalog := NewHardcodedCatalog()
	plan, tasks := ComposePlan(uuid.New(), "2026-02-20", PlanA, goals,
		func(g GoalID, pt PlanType) []TaskTemplate { return catalog.GetTasksForGoal(g, pt) }, 1)

	if plan.PlanType != PlanA {
		t.Errorf("expected PlanA, got %s", plan.PlanType)
	}

	if plan.PriorityTaskID == uuid.Nil {
		t.Error("priority task should be set")
	}

	if len(tasks) < 3 {
		t.Errorf("expected at least 3 tasks for PlanA, got %d", len(tasks))
	}

	if plan.FoundationTaskID == nil {
		t.Error("foundation task should be set when foundation goals exist")
	}
}

func TestComposePlan_PlanC_MaxTwoItems(t *testing.T) {
	goals := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalJava},
		{ID: GoalSleep},
		{ID: GoalHealth},
	}

	catalog := NewHardcodedCatalog()
	_, tasks := ComposePlan(uuid.New(), "2026-02-20", PlanC, goals,
		func(g GoalID, pt PlanType) []TaskTemplate { return catalog.GetTasksForGoal(g, pt) }, 1)

	if len(tasks) > 2 {
		t.Errorf("PlanC should have max 2 tasks, got %d", len(tasks))
	}
}

func TestComposePlan_MaxTwoIntensiveGoals(t *testing.T) {
	goals := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalJava},
		{ID: GoalSaaS},
	}

	catalog := NewHardcodedCatalog()
	_, tasks := ComposePlan(uuid.New(), "2026-02-20", PlanA, goals,
		func(g GoalID, pt PlanType) []TaskTemplate { return catalog.GetTasksForGoal(g, pt) }, 1)

	intensiveCount := 0
	seen := map[GoalID]bool{}
	for _, t := range tasks {
		class := GoalClassifications[t.GoalDomain]
		if class == GoalClassIntensive && !seen[t.GoalDomain] {
			intensiveCount++
			seen[t.GoalDomain] = true
		}
	}

	if intensiveCount > MaxIntensiveGoals {
		t.Errorf("should not have more than %d intensive goal domains, got %d", MaxIntensiveGoals, intensiveCount)
	}
}

func TestComposePlan_PlanB_HasComplementary(t *testing.T) {
	goals := []GoalEntry{
		{ID: GoalEnglish},
		{ID: GoalSleep},
	}

	catalog := NewHardcodedCatalog()
	plan, _ := ComposePlan(uuid.New(), "2026-02-20", PlanB, goals,
		func(g GoalID, pt PlanType) []TaskTemplate { return catalog.GetTasksForGoal(g, pt) }, 1)

	if len(plan.ComplementaryIDs) == 0 {
		t.Error("PlanB should have complementary tasks")
	}
}

func TestComposePlan_NoGoals_DefaultTask(t *testing.T) {
	plan, tasks := ComposePlan(uuid.New(), "2026-02-20", PlanC, nil,
		func(g GoalID, pt PlanType) []TaskTemplate { return nil }, 1)

	if len(tasks) == 0 {
		t.Error("should generate at least a default task")
	}

	if plan.PriorityTaskID == uuid.Nil {
		t.Error("priority task should be set even with no goals")
	}
}

// --- Task status transitions ---

func TestTaskTransition_PlannedToInProgress(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskPlanned}
	view, err := task.ApplyAction(ActionStart, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskInProgress {
		t.Errorf("expected in_progress, got %s", view.Status)
	}
}

func TestTaskTransition_InProgressToEvidencePending(t *testing.T) {
	gate := "speaking_output"
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskInProgress, GateRef: &gate}
	view, err := task.ApplyAction(ActionMarkDoneReq, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskEvidencePending {
		t.Errorf("expected evidence_pending, got %s", view.Status)
	}
}

func TestTaskTransition_InProgressToCompleted_NoGate(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskInProgress}
	view, err := task.ApplyAction(ActionMarkDoneReq, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskCompleted {
		t.Errorf("expected completed, got %s", view.Status)
	}
}

func TestTaskTransition_PlannedToBlocked(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskPlanned}
	view, err := task.ApplyAction(ActionBlock, "sem privacidade")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskBlocked {
		t.Errorf("expected blocked, got %s", view.Status)
	}
	if task.BlockReason == nil || *task.BlockReason != "sem privacidade" {
		t.Error("block reason should be set")
	}
}

func TestTaskTransition_PlannedToDeferred(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskPlanned}
	view, err := task.ApplyAction(ActionDefer, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskDeferred {
		t.Errorf("expected deferred, got %s", view.Status)
	}
}

func TestTaskTransition_InvalidTransition(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskCompleted}
	_, err := task.ApplyAction(ActionStart, "")
	if err == nil {
		t.Error("expected error for completed → in_progress")
	}
}

func TestTaskTransition_AddNote(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskInProgress}
	view, err := task.ApplyAction(ActionAddNote, "minha nota")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Status != TaskInProgress {
		t.Errorf("add_note should not change status, got %s", view.Status)
	}
	if task.Note == nil || *task.Note != "minha nota" {
		t.Error("note should be set")
	}
}

func TestTaskTransition_EvidencePendingToCompleted(t *testing.T) {
	task := &PlannedTask{TaskID: uuid.New(), Status: TaskEvidencePending}
	if !task.CanTransitionTo(TaskCompleted) {
		t.Error("evidence_pending should be able to transition to completed")
	}
}

// --- HardcodedCatalog tests ---

func TestHardcodedCatalog_EnglishPlanA(t *testing.T) {
	c := NewHardcodedCatalog()
	templates := c.GetTasksForGoal(GoalEnglish, PlanA)
	if len(templates) == 0 {
		t.Error("expected templates for English/PlanA")
	}
	if templates[0].EstimatedMin == 0 {
		t.Error("template should have estimated time")
	}
}

func TestHardcodedCatalog_UnknownGoal(t *testing.T) {
	c := NewHardcodedCatalog()
	templates := c.GetTasksForGoal(GoalID("unknown"), PlanA)
	if templates != nil {
		t.Error("expected nil for unknown goal")
	}
}

// --- DailyState helpers ---

func TestDailyState_BuildStepsSummary(t *testing.T) {
	ds := NewDailyState(uuid.New(), "2026-02-20")
	ds.Tasks = []PlannedTask{
		{TaskID: uuid.New(), Title: "Task 1", Status: TaskCompleted},
		{TaskID: uuid.New(), Title: "Task 2", Status: TaskInProgress},
		{TaskID: uuid.New(), Title: "Task 3", Status: TaskPlanned},
		{TaskID: uuid.New(), Title: "Task 4", Status: TaskBlocked},
	}

	summary := ds.BuildStepsSummary()
	if len(summary.Done) != 1 {
		t.Errorf("expected 1 done, got %d", len(summary.Done))
	}
	if len(summary.InProgress) != 1 {
		t.Errorf("expected 1 in_progress, got %d", len(summary.InProgress))
	}
	if len(summary.Pending) != 1 {
		t.Errorf("expected 1 pending, got %d", len(summary.Pending))
	}
	if len(summary.Blocked) != 1 {
		t.Errorf("expected 1 blocked, got %d", len(summary.Blocked))
	}
}
