package domain

import (
	"testing"

	"github.com/google/uuid"
)

// --- Sleep Template tests ---

func TestSleepTemplates_PlanA_Has2Tasks(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanA)
	if len(tasks) != 2 {
		t.Fatalf("PlanA should have 2 sleep tasks, got %d", len(tasks))
	}
}

func TestSleepTemplates_PlanB_Has2Tasks(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanB)
	if len(tasks) != 2 {
		t.Fatalf("PlanB should have 2 sleep tasks, got %d", len(tasks))
	}
}

func TestSleepTemplates_PlanC_Has1Task(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanC)
	if len(tasks) != 1 {
		t.Fatalf("PlanC should have 1 sleep task, got %d", len(tasks))
	}
}

func TestSleepTemplates_PlanA_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanA)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanA sleep task %q should have a gate profile", task.Title)
		}
	}
}

func TestSleepTemplates_PlanB_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanB)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanB sleep task %q should have a gate profile", task.Title)
		}
	}
}

func TestSleepTemplates_PlanC_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanC)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanC sleep task %q should have a gate profile", task.Title)
		}
	}
}

func TestSleepTemplates_PlanC_FitsIn1Min(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanC)
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total > 1 {
		t.Fatalf("PlanC total should be <= 1 min, got %d", total)
	}
}

func TestSleepTemplates_PlanA_HasExpectedGateProfiles(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanA)
	gateRefs := map[string]bool{}
	for _, task := range tasks {
		if task.GateProfile != nil {
			gateRefs[*task.GateProfile] = true
		}
	}
	if !gateRefs["sleep_diary"] {
		t.Fatal("PlanA should include sleep_diary gate")
	}
	if !gateRefs["sleep_routine"] {
		t.Fatal("PlanA should include sleep_routine gate")
	}
}

func TestSleepTemplates_PlanC_HasDiaryMinGate(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalSleep, PlanC)
	if tasks[0].GateProfile == nil || *tasks[0].GateProfile != "sleep_diary_min" {
		t.Fatal("PlanC should use sleep_diary_min gate profile")
	}
}

// --- Gate profile tests for Sleep ---

func TestGateProfile_SleepDiary_1Metadata_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["sleep_diary"]
	ev := Evidence{
		EvidenceID:  uuid.New(),
		Kind:        EvidenceMetadata,
		Sensitivity: SensitivityC1,
		Summary:     "dormiu=23:30, acordou=07:00, qualidade=8",
	}
	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 metadata on sleep_diary, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_SleepDiaryMin_1TextAnswer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["sleep_diary_min"]
	ev := newTestEvidence(EvidenceText, SensitivityC1, "Dormi bem")
	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on sleep_diary_min, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_SleepRoutine_1TextAnswer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["sleep_routine"]
	ev := newTestEvidence(EvidenceText, SensitivityC1, "Desliguei telas, li por 10min")
	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on sleep_routine, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_SleepDiary_NoEvidence_NotSatisfied(t *testing.T) {
	profile := GateProfileCatalog["sleep_diary"]
	result := EvaluateGate(profile, []Evidence{}, defaultPrivacy())
	if result.GateStatus != GateNotSatisfied {
		t.Fatal("expected not_satisfied with no evidence for sleep_diary")
	}
}

func TestGateProfile_SleepRoutine_EmptyEvidence_NotSatisfied(t *testing.T) {
	profile := GateProfileCatalog["sleep_routine"]
	ev := Evidence{
		EvidenceID:  uuid.New(),
		Kind:        EvidenceText,
		Sensitivity: SensitivityC1,
		Summary:     "",
		ContentRef:  nil,
	}
	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateNotSatisfied {
		t.Fatal("expected not_satisfied for empty evidence on sleep_routine")
	}
}

func TestGateProfile_AllSleepProfilesExist(t *testing.T) {
	expected := []string{"sleep_diary", "sleep_diary_min", "sleep_routine"}
	for _, id := range expected {
		if _, ok := GateProfileCatalog[id]; !ok {
			t.Fatalf("missing Sleep gate profile: %s", id)
		}
	}
}

// --- ComputeSleepDuration tests ---

func TestComputeSleepDuration_OvernightNormal(t *testing.T) {
	result := ComputeSleepDuration("23:30", "07:00")
	if result == nil {
		t.Fatal("expected non-nil result for 23:30 -> 07:00")
	}
	if *result != 450 {
		t.Fatalf("expected 450 min, got %d", *result)
	}
}

func TestComputeSleepDuration_AfterMidnight(t *testing.T) {
	result := ComputeSleepDuration("01:00", "08:00")
	if result == nil {
		t.Fatal("expected non-nil result for 01:00 -> 08:00")
	}
	if *result != 420 {
		t.Fatalf("expected 420 min, got %d", *result)
	}
}

func TestComputeSleepDuration_InvalidInput(t *testing.T) {
	result := ComputeSleepDuration("abc", "07:00")
	if result != nil {
		t.Fatalf("expected nil for invalid input, got %d", *result)
	}
}

func TestComputeSleepDuration_EmptyInput(t *testing.T) {
	result := ComputeSleepDuration("", "")
	if result != nil {
		t.Fatalf("expected nil for empty input, got %d", *result)
	}
}

func TestComputeSleepDuration_SameTime(t *testing.T) {
	result := ComputeSleepDuration("08:00", "08:00")
	if result == nil {
		t.Fatal("expected non-nil result for same time")
	}
	if *result != 24*60 {
		t.Fatalf("expected 1440 min (24h), got %d", *result)
	}
}

// --- ClassifySleepDiaryStatus tests ---

func TestClassifySleepDiaryStatus_Complete_2Fields(t *testing.T) {
	slept := "23:30"
	woke := "07:00"
	status := ClassifySleepDiaryStatus(&slept, &woke, nil, nil)
	if status != "complete" {
		t.Fatalf("expected complete with 2 fields, got %s", status)
	}
}

func TestClassifySleepDiaryStatus_Complete_AllFields(t *testing.T) {
	slept := "23:30"
	woke := "07:00"
	q := 8
	e := 7
	status := ClassifySleepDiaryStatus(&slept, &woke, &q, &e)
	if status != "complete" {
		t.Fatalf("expected complete with 4 fields, got %s", status)
	}
}

func TestClassifySleepDiaryStatus_Partial_1Field(t *testing.T) {
	q := 7
	status := ClassifySleepDiaryStatus(nil, nil, &q, nil)
	if status != "partial" {
		t.Fatalf("expected partial with 1 field, got %s", status)
	}
}

func TestClassifySleepDiaryStatus_Partial_0Fields(t *testing.T) {
	status := ClassifySleepDiaryStatus(nil, nil, nil, nil)
	if status != "partial" {
		t.Fatalf("expected partial with 0 fields, got %s", status)
	}
}

// --- DefaultSleepRoutineSteps tests ---

func TestDefaultSleepRoutineSteps_Normal(t *testing.T) {
	steps := DefaultSleepRoutineSteps("normal")
	if len(steps) < 3 || len(steps) > 4 {
		t.Fatalf("normal version should return 3-4 steps, got %d", len(steps))
	}
}

func TestDefaultSleepRoutineSteps_Minimal(t *testing.T) {
	steps := DefaultSleepRoutineSteps("minimal")
	if len(steps) < 1 || len(steps) > 2 {
		t.Fatalf("minimal version should return 1-2 steps, got %d", len(steps))
	}
}

// --- ComputeRegularityDelta tests ---

func TestComputeRegularityDelta_LessThan3Entries_Nil(t *testing.T) {
	entries := []SleepDiaryEntry{
		{SleptAt: strPtr("23:00")},
		{SleptAt: strPtr("23:30")},
	}
	result := ComputeRegularityDelta(entries)
	if result != nil {
		t.Fatalf("expected nil with <3 entries, got %d", *result)
	}
}

func TestComputeRegularityDelta_3Entries_ReturnsValue(t *testing.T) {
	entries := []SleepDiaryEntry{
		{SleptAt: strPtr("23:00")},
		{SleptAt: strPtr("23:30")},
		{SleptAt: strPtr("23:00")},
	}
	result := ComputeRegularityDelta(entries)
	if result == nil {
		t.Fatal("expected non-nil with 3 entries")
	}
	if *result < 0 {
		t.Fatalf("regularity delta should be >= 0, got %d", *result)
	}
}

func TestComputeRegularityDelta_AllSameTime_Zero(t *testing.T) {
	entries := []SleepDiaryEntry{
		{SleptAt: strPtr("23:00")},
		{SleptAt: strPtr("23:00")},
		{SleptAt: strPtr("23:00")},
	}
	result := ComputeRegularityDelta(entries)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if *result != 0 {
		t.Fatalf("expected 0 for all same times, got %d", *result)
	}
}

// --- ParseSleepDiaryInput tests ---

func TestParseSleepDiaryInput_FullFormat(t *testing.T) {
	sleptAt, wokeAt, quality, energy := ParseSleepDiaryInput("23:30 07:00 8 7")
	if sleptAt == nil || *sleptAt != "23:30" {
		t.Fatal("expected sleptAt=23:30")
	}
	if wokeAt == nil || *wokeAt != "07:00" {
		t.Fatal("expected wokeAt=07:00")
	}
	if quality == nil || *quality != 8 {
		t.Fatal("expected quality=8")
	}
	if energy == nil || *energy != 7 {
		t.Fatal("expected energy=7")
	}
}

func TestParseSleepDiaryInput_CommaFormat(t *testing.T) {
	sleptAt, wokeAt, quality, energy := ParseSleepDiaryInput("23:30, 07:00, 8, 7")
	if sleptAt == nil || *sleptAt != "23:30" {
		t.Fatal("expected sleptAt=23:30")
	}
	if wokeAt == nil || *wokeAt != "07:00" {
		t.Fatal("expected wokeAt=07:00")
	}
	if quality == nil || *quality != 8 {
		t.Fatal("expected quality=8")
	}
	if energy == nil || *energy != 7 {
		t.Fatal("expected energy=7")
	}
}

func TestParseSleepDiaryInput_Bem(t *testing.T) {
	sleptAt, wokeAt, quality, energy := ParseSleepDiaryInput("bem")
	if sleptAt != nil {
		t.Fatal("expected nil sleptAt for 'bem'")
	}
	if wokeAt != nil {
		t.Fatal("expected nil wokeAt for 'bem'")
	}
	if quality == nil || *quality != 7 {
		t.Fatal("expected quality=7 for 'bem'")
	}
	if energy == nil || *energy != 7 {
		t.Fatal("expected energy=7 for 'bem'")
	}
}

func TestParseSleepDiaryInput_Mal(t *testing.T) {
	sleptAt, wokeAt, quality, energy := ParseSleepDiaryInput("mal")
	if sleptAt != nil {
		t.Fatal("expected nil sleptAt for 'mal'")
	}
	if wokeAt != nil {
		t.Fatal("expected nil wokeAt for 'mal'")
	}
	if quality == nil || *quality != 3 {
		t.Fatal("expected quality=3 for 'mal'")
	}
	if energy == nil || *energy != 3 {
		t.Fatal("expected energy=3 for 'mal'")
	}
}

func TestParseSleepDiaryInput_PartialTwoFields(t *testing.T) {
	sleptAt, wokeAt, quality, energy := ParseSleepDiaryInput("23:30 07:00")
	if sleptAt == nil || *sleptAt != "23:30" {
		t.Fatal("expected sleptAt=23:30")
	}
	if wokeAt == nil || *wokeAt != "07:00" {
		t.Fatal("expected wokeAt=07:00")
	}
	if quality != nil {
		t.Fatal("expected nil quality")
	}
	if energy != nil {
		t.Fatal("expected nil energy")
	}
}

// --- Sleep Intervention Pool test ---

func TestSleepInterventionPool_HasEntries(t *testing.T) {
	if len(SleepInterventionPool) < 5 {
		t.Fatalf("expected at least 5 interventions in pool, got %d", len(SleepInterventionPool))
	}
	for _, item := range SleepInterventionPool {
		if item.Description == "" {
			t.Fatal("intervention description should not be empty")
		}
		if item.AdherenceRule == "" {
			t.Fatal("intervention adherence_rule should not be empty")
		}
	}
}

func strPtr(s string) *string {
	return &s
}
