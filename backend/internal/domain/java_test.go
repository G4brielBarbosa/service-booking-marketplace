package domain

import (
	"testing"

	"github.com/google/uuid"
)

// --- Java Template tests ---

func TestJavaTemplates_PlanA_Duration(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanA)
	if len(tasks) != 3 {
		t.Fatalf("PlanA should have 3 tasks, got %d", len(tasks))
	}
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total < 40 || total > 60 {
		t.Fatalf("PlanA total should be ~45 min (40-60 range), got %d", total)
	}
}

func TestJavaTemplates_PlanB_Duration(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanB)
	if len(tasks) != 3 {
		t.Fatalf("PlanB should have 3 tasks, got %d", len(tasks))
	}
	total := 0
	for _, task := range tasks {
		total += task.EstimatedMin
	}
	if total != 30 {
		t.Fatalf("PlanB total should be 30 min, got %d", total)
	}
}

func TestJavaTemplates_PlanC_FitsIn15Min(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanC)
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

func TestJavaTemplates_PlanA_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanA)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanA task %q should have a gate profile", task.Title)
		}
	}
}

func TestJavaTemplates_PlanB_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanB)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanB task %q should have a gate profile", task.Title)
		}
	}
}

func TestJavaTemplates_PlanC_AllBlocksHaveGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanC)
	for _, task := range tasks {
		if task.GateProfile == nil {
			t.Fatalf("PlanC task %q should have a gate profile", task.Title)
		}
	}
}

func TestJavaTemplates_PlanA_HasExpectedGateProfiles(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanA)

	gateRefs := map[string]bool{}
	for _, task := range tasks {
		if task.GateProfile != nil {
			gateRefs[*task.GateProfile] = true
		}
	}

	expected := []string{"java_practice", "java_retrieval", "java_learning_log"}
	for _, e := range expected {
		if !gateRefs[e] {
			t.Fatalf("PlanA should include gate profile %q", e)
		}
	}
}

func TestJavaTemplates_PlanC_HasMinimalGates(t *testing.T) {
	catalog := NewHardcodedCatalog()
	tasks := catalog.GetTasksForGoal(GoalJava, PlanC)

	gateRefs := map[string]bool{}
	for _, task := range tasks {
		if task.GateProfile != nil {
			gateRefs[*task.GateProfile] = true
		}
	}

	if !gateRefs["java_practice_min"] {
		t.Fatal("PlanC should include java_practice_min gate")
	}
	if !gateRefs["java_retrieval_min"] {
		t.Fatal("PlanC should include java_retrieval_min gate")
	}
}

// --- Gate profile tests for Java ---

func TestGateProfile_JavaPractice_1Answer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["java_practice"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "Resolvi exercício de streams e expliquei o uso de map e filter.")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on java_practice, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_JavaPracticeMin_1Answer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["java_practice_min"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "Li artigo sobre generics em Java.")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on java_practice_min, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_JavaRetrievalMin_1Answer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["java_retrieval_min"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "Streams permitem operações de map, filter e reduce sobre coleções.")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on java_retrieval_min, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_JavaLearningLog_1Answer_Satisfied(t *testing.T) {
	profile := GateProfileCatalog["java_learning_log"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "Confundi ArrayList com LinkedList na escolha de estrutura.")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for 1 text answer on java_learning_log, got %s (%s)", result.GateStatus, result.ReasonShort)
	}
}

func TestGateProfile_JavaPractice_NoEvidence_NotSatisfied(t *testing.T) {
	profile := GateProfileCatalog["java_practice"]
	result := EvaluateGate(profile, []Evidence{}, defaultPrivacy())
	if result.GateStatus != GateNotSatisfied {
		t.Fatal("expected not_satisfied with no evidence for java_practice")
	}
}

func TestGateProfile_JavaLearningLog_EmptyEvidence_NotSatisfied(t *testing.T) {
	profile := GateProfileCatalog["java_learning_log"]
	ev := Evidence{
		EvidenceID:  uuid.New(),
		Kind:        EvidenceText,
		Sensitivity: SensitivityC2,
		Summary:     "",
		ContentRef:  nil,
	}

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())
	if result.GateStatus != GateNotSatisfied {
		t.Fatal("expected not_satisfied for empty evidence on java_learning_log")
	}
}

// --- Reuse of ClassifyRetrievalStatus for Java context ---

func TestClassifyRetrievalStatus_JavaContext_3of5_Low(t *testing.T) {
	if got := ClassifyRetrievalStatus(3, 5); got != "low" {
		t.Fatalf("expected low for 3/5 (60%%), got %s", got)
	}
}

func TestClassifyRetrievalStatus_JavaContext_4of5_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(4, 5); got != "ok" {
		t.Fatalf("expected ok for 4/5 (80%%), got %s", got)
	}
}

func TestClassifyRetrievalStatus_JavaContext_1of1_Ok(t *testing.T) {
	if got := ClassifyRetrievalStatus(1, 1); got != "ok" {
		t.Fatalf("expected ok for 1/1, got %s", got)
	}
}

// --- Reuse of CheckRecurrence for Java context ---

func TestCheckRecurrence_JavaContext_Count2_False(t *testing.T) {
	if CheckRecurrence(2) {
		t.Fatal("expected false for count=2 in Java context")
	}
}

func TestCheckRecurrence_JavaContext_Count3_True(t *testing.T) {
	if !CheckRecurrence(3) {
		t.Fatal("expected true for count=3 in Java context")
	}
}

// --- All Java gate profiles exist ---

func TestGateProfile_AllJavaProfilesExist(t *testing.T) {
	expected := []string{
		"java_practice", "java_practice_min",
		"java_retrieval", "java_retrieval_min",
		"java_learning_log",
	}
	for _, id := range expected {
		if _, ok := GateProfileCatalog[id]; !ok {
			t.Fatalf("missing Java gate profile: %s", id)
		}
	}
}
