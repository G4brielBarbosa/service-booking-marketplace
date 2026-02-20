package domain

import (
	"testing"
)

// --- ClassifyRetrievalStatus tests ---

func TestClassifyRetrievalStatus_7of10_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(7, 10); got != "ok" {
		t.Fatalf("expected ok for 7/10, got %s", got)
	}
}

func TestClassifyRetrievalStatus_6of10_Low(t *testing.T) {
	if got := ClassifyRetrievalStatus(6, 10); got != "low" {
		t.Fatalf("expected low for 6/10, got %s", got)
	}
}

func TestClassifyRetrievalStatus_3of3_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(3, 3); got != "ok" {
		t.Fatalf("expected ok for 3/3, got %s", got)
	}
}

func TestClassifyRetrievalStatus_0of5_Low(t *testing.T) {
	if got := ClassifyRetrievalStatus(0, 5); got != "low" {
		t.Fatalf("expected low for 0/5, got %s", got)
	}
}

func TestClassifyRetrievalStatus_0of0_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(0, 0); got != "ok" {
		t.Fatalf("expected ok for 0/0 (edge case), got %s", got)
	}
}

func TestClassifyRetrievalStatus_10of10_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(10, 10); got != "ok" {
		t.Fatalf("expected ok for 10/10, got %s", got)
	}
}

func TestClassifyRetrievalStatus_1of1_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(1, 1); got != "ok" {
		t.Fatalf("expected ok for 1/1, got %s", got)
	}
}

func TestClassifyRetrievalStatus_Boundary70Percent(t *testing.T) {
	// 7/10 = 70% -> ok
	if got := ClassifyRetrievalStatus(7, 10); got != "ok" {
		t.Fatalf("expected ok for exactly 70%%, got %s", got)
	}
	// 14/20 = 70% -> ok
	if got := ClassifyRetrievalStatus(14, 20); got != "ok" {
		t.Fatalf("expected ok for 14/20, got %s", got)
	}
	// 13/20 = 65% -> low
	if got := ClassifyRetrievalStatus(13, 20); got != "low" {
		t.Fatalf("expected low for 13/20 (65%%), got %s", got)
	}
}

// --- CheckRecurrence tests ---

func TestCheckRecurrence_Count2_False(t *testing.T) {
	if CheckRecurrence(2) {
		t.Fatal("expected false for count=2")
	}
}

func TestCheckRecurrence_Count3_True(t *testing.T) {
	if !CheckRecurrence(3) {
		t.Fatal("expected true for count=3")
	}
}

func TestCheckRecurrence_Count5_True(t *testing.T) {
	if !CheckRecurrence(5) {
		t.Fatal("expected true for count=5")
	}
}

func TestCheckRecurrence_Count0_False(t *testing.T) {
	if CheckRecurrence(0) {
		t.Fatal("expected false for count=0")
	}
}

func TestCheckRecurrence_Count1_False(t *testing.T) {
	if CheckRecurrence(1) {
		t.Fatal("expected false for count=1")
	}
}

// --- Template tests ---

func TestEnglishTemplates_PlanA_Duration(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanA)
	if len(tasks) != 3 {
		t.Fatalf("PlanA should have 3 tasks, got %d", len(tasks))
	}
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total != 50 {
		t.Fatalf("PlanA total should be 50 min, got %d", total)
	}
}

func TestEnglishTemplates_PlanB_Duration(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanB)
	if len(tasks) != 2 {
		t.Fatalf("PlanB should have 2 tasks, got %d", len(tasks))
	}
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total != 25 {
		t.Fatalf("PlanB total should be 25 min, got %d", total)
	}
}

func TestEnglishTemplates_PlanC_FitsIn15Min(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanC)
	if len(tasks) != 2 {
		t.Fatalf("PlanC should have 2 tasks, got %d", len(tasks))
	}
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total > 15 {
		t.Fatalf("PlanC total should be <= 15 min, got %d", total)
	}
}

func TestEnglishTemplates_PlanA_SpeakingHasGateProfile(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanA)

	found := false
	for _, task := range tasks {
		if task.Title == "Prática de speaking" {
			found = true
			if task.GateProfile == nil {
				t.Fatal("speaking task should have a gate profile")
			}
			if *task.GateProfile != "speaking_output" {
				t.Fatalf("speaking gate profile should be speaking_output, got %s", *task.GateProfile)
			}
		}
	}
	if !found {
		t.Fatal("PlanA should include Prática de speaking")
	}
}

func TestEnglishTemplates_PlanA_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanA)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanA task %q should have a gate profile", task.Title)
		}
	}
}

func TestEnglishTemplates_PlanB_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanB)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanB task %q should have a gate profile", task.Title)
		}
	}
}

func TestEnglishTemplates_PlanC_HasMinimalGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalEnglish, PlanC)

	gateRefs := map[string]bool{}
	for _, task := range tasks {
		if task.GateProfile != nil {
			gateRefs[*task.GateProfile] = true
		}
	}

	if !gateRefs["english_comprehension_min"] {
		t.Fatal("PlanC should include english_comprehension_min gate")
	}
	if !gateRefs["english_retrieval_min"] {
		t.Fatal("PlanC should include english_retrieval_min gate")
	}
}

// --- New gate profiles tests ---

func TestGateProfile_ComprehensionMin_1Answer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension_min"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "resposta 1")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 answer on comprehension_min, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_RetrievalMin_3Items_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["english_retrieval_min"]
	evs := []Evidence{
		newTestEvidence(EvidenceText, SensitivityC2, "item 1"),
		newTestEvidence(EvidenceText, SensitivityC2, "item 2"),
		newTestEvidence(EvidenceText, SensitivityC2, "item 3"),
	}

	result := EvaluateGate(profile, evs, defaultPrivacy())

	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 3 items on retrieval_min, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_RetrievalMin_2Items_Partial(t *testing.T) {
	profile := GateProfileCatalog["english_retrieval_min"]
	evs := []Evidence{
		newTestEvidence(EvidenceText, SensitivityC2, "item 1"),
		newTestEvidence(EvidenceText, SensitivityC2, "item 2"),
	}

	result := EvaluateGate(profile, evs, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatalf("expected not_satisfied for 2 items on retrieval_min, got %s", result.GateStatus)
	}
	if result.FailureReason != FailurePartial {
		t.Fatalf("expected partial failure, got %s", result.FailureReason)
	}
}

func TestGateProfile_AllProfilesExist_WithMinProfiles(t *testing.T) {
	expected := []string{
		"speaking_output", "english_comprehension", "english_retrieval",
		"english_comprehension_min", "english_retrieval_min",
		"java_practice", "java_retrieval", "sleep_diary",
	}
	for _, id := range expected {
		if _, ok := GateProfileCatalog[id]; !ok {
			t.Fatalf("missing profile: %s", id)
		}
	}
}

func TestSpeakingGate_NeverPassesWithoutAudio(t *testing.T) {
	profile := GateProfileCatalog["speaking_output"]

	rubricEvs := []Evidence{
		newTestEvidence(EvidenceRubric, SensitivityC2, "rubric total=6"),
	}

	result := EvaluateGate(profile, rubricEvs, defaultPrivacy())
	if result.GateStatus == GateSatisfied {
		t.Fatal("speaking gate should not pass with rubric only (no audio)")
	}

	textEvs := []Evidence{
		newTestEvidence(EvidenceText, SensitivityC2, "tentativa falar"),
	}
	result2 := EvaluateGate(profile, textEvs, defaultPrivacy())
	if result2.GateStatus == GateSatisfied {
		t.Fatal("speaking gate should not pass with text only (no audio)")
	}
}
