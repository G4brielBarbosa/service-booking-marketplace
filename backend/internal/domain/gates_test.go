package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestEvidence(kind EvidenceKind, sensitivity SensitivityLevel, summary string) Evidence {
	ref := "test-ref"
	return Evidence{
		EvidenceID:  uuid.New(),
		UserID:      uuid.New(),
		TaskID:      uuid.New(),
		Kind:        kind,
		Sensitivity: sensitivity,
		Summary:     summary,
		ContentRef:  &ref,
		Timestamp:   time.Now(),
	}
}

func defaultPrivacy() PrivacyPolicy {
	return NewDefaultPrivacyPolicy(uuid.New())
}

func optedOutC3Privacy() PrivacyPolicy {
	p := NewDefaultPrivacyPolicy(uuid.New())
	p.SetOptOut(SensitivityC3, true)
	return p
}

// --- ValidateEvidence tests ---

func TestValidateEvidence_Empty(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	ev := Evidence{
		EvidenceID: uuid.New(),
		Kind:       EvidenceText,
	}

	valid, reason := ValidateEvidence(ev, profile)
	if valid {
		t.Fatal("expected empty evidence to be invalid")
	}
	if reason == "" {
		t.Fatal("expected a reason for invalid evidence")
	}
}

func TestValidateEvidence_WrongType_AudioRequired(t *testing.T) {
	profile := GateProfileCatalog["speaking_output"]
	ev := newTestEvidence(EvidenceText, SensitivityC3, "some text")

	valid, reason := ValidateEvidence(ev, profile)
	if valid {
		t.Fatal("expected text evidence to be invalid for speaking gate (audio required)")
	}
	if reason == "" {
		t.Fatal("expected reason for not_equivalent")
	}
}

func TestValidateEvidence_TextForTextProfile(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "answer 1")

	valid, reason := ValidateEvidence(ev, profile)
	if !valid {
		t.Fatalf("expected text evidence to be valid for comprehension, got reason: %s", reason)
	}
}

func TestValidateEvidence_AudioForSpeaking(t *testing.T) {
	profile := GateProfileCatalog["speaking_output"]
	ev := newTestEvidence(EvidenceAudio, SensitivityC3, "áudio 45s")

	valid, reason := ValidateEvidence(ev, profile)
	if !valid {
		t.Fatalf("expected audio evidence to be valid for speaking, got reason: %s", reason)
	}
}

// --- EvaluateGate tests ---

func TestEvaluateGate_Missing(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	result := EvaluateGate(profile, nil, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatalf("expected not_satisfied, got %s", result.GateStatus)
	}
	if result.FailureReason != FailureMissing {
		t.Fatalf("expected missing, got %s", result.FailureReason)
	}
	if result.ReasonShort == "" {
		t.Fatal("expected non-empty reason_short")
	}
	if result.NextMinStep == "" {
		t.Fatal("expected non-empty next_min_step")
	}
}

func TestEvaluateGate_Invalid(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	ev := Evidence{EvidenceID: uuid.New(), Kind: EvidenceText}

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatalf("expected not_satisfied, got %s", result.GateStatus)
	}
	if result.FailureReason != FailureInvalid {
		t.Fatalf("expected invalid, got %s", result.FailureReason)
	}
}

func TestEvaluateGate_Partial(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "answer 1")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatalf("expected not_satisfied, got %s", result.GateStatus)
	}
	if result.FailureReason != FailurePartial {
		t.Fatalf("expected partial, got %s", result.FailureReason)
	}
	if result.ReasonShort == "" {
		t.Fatal("expected non-empty reason_short for partial")
	}
	if result.NextMinStep == "" {
		t.Fatal("expected non-empty next_min_step for partial")
	}
}

func TestEvaluateGate_NotEquivalent_SpeakingText(t *testing.T) {
	profile := GateProfileCatalog["speaking_output"]
	ev := newTestEvidence(EvidenceText, SensitivityC3, "eu tentei falar")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatalf("expected not_satisfied, got %s", result.GateStatus)
	}
	if result.FailureReason != FailureNotEquivalent {
		t.Fatalf("expected not_equivalent, got %s", result.FailureReason)
	}
}

func TestEvaluateGate_Satisfied_Comprehension(t *testing.T) {
	profile := GateProfileCatalog["english_comprehension"]
	evs := []Evidence{
		newTestEvidence(EvidenceText, SensitivityC2, "answer 1"),
		newTestEvidence(EvidenceText, SensitivityC2, "answer 2"),
		newTestEvidence(EvidenceText, SensitivityC2, "answer 3"),
	}

	result := EvaluateGate(profile, evs, defaultPrivacy())

	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied, got %s (reason: %s)", result.GateStatus, result.ReasonShort)
	}
	if result.FailureReason != FailureNone {
		t.Fatalf("expected no failure reason, got %s", result.FailureReason)
	}
}

func TestEvaluateGate_Satisfied_JavaPractice(t *testing.T) {
	profile := GateProfileCatalog["java_practice"]
	ev := newTestEvidence(EvidenceText, SensitivityC2, "code + explanation")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied, got %s (reason: %s)", result.GateStatus, result.ReasonShort)
	}
}

func TestEvaluateGate_SleepDiary_MinimalRegistration(t *testing.T) {
	profile := GateProfileCatalog["sleep_diary"]
	ev := newTestEvidence(EvidenceMetadata, SensitivityC2, "dormiu 23h, acordou 7h, qualidade ok")

	result := EvaluateGate(profile, []Evidence{ev}, defaultPrivacy())

	if result.GateStatus != GateSatisfied {
		t.Fatalf("expected satisfied for sleep diary minimal, got %s (reason: %s)", result.GateStatus, result.ReasonShort)
	}
}

func TestEvaluateGate_SpeakingNeverPassesWithoutAudio(t *testing.T) {
	profile := GateProfileCatalog["speaking_output"]

	textEvs := []Evidence{
		newTestEvidence(EvidenceText, SensitivityC3, "tentativa 1"),
		newTestEvidence(EvidenceText, SensitivityC3, "tentativa 2"),
		newTestEvidence(EvidenceText, SensitivityC3, "tentativa 3"),
	}

	result := EvaluateGate(profile, textEvs, defaultPrivacy())

	if result.GateStatus != GateNotSatisfied {
		t.Fatal("speaking gate should never pass with only text evidence")
	}
	if result.FailureReason != FailureNotEquivalent {
		t.Fatalf("expected not_equivalent for text on speaking, got %s", result.FailureReason)
	}
}

func TestEvaluateGate_ReasonShort_NonEmpty_WhenNotSatisfied(t *testing.T) {
	profiles := []string{"speaking_output", "english_comprehension", "english_retrieval", "java_practice", "java_retrieval", "sleep_diary"}

	for _, pid := range profiles {
		profile := GateProfileCatalog[pid]
		result := EvaluateGate(profile, nil, defaultPrivacy())

		if result.GateStatus == GateNotSatisfied {
			if result.ReasonShort == "" {
				t.Fatalf("profile %s: expected non-empty reason_short", pid)
			}
			if result.NextMinStep == "" {
				t.Fatalf("profile %s: expected non-empty next_min_step", pid)
			}
		}
	}
}

// --- ApplyStoragePolicy tests ---

func TestApplyStoragePolicy_C3_OptOut_Discarded(t *testing.T) {
	ev := newTestEvidence(EvidenceAudio, SensitivityC3, "áudio 30s")
	privacy := optedOutC3Privacy()

	result := ApplyStoragePolicy(ev, privacy)

	if result.StoragePolicy != StorageDiscardedAfterProc {
		t.Fatalf("expected discarded_after_processing, got %s", result.StoragePolicy)
	}
	if result.ContentRef != nil {
		t.Fatal("expected ContentRef to be nil after discard")
	}
}

func TestApplyStoragePolicy_C3_NoOptOut_Kept7d(t *testing.T) {
	ev := newTestEvidence(EvidenceAudio, SensitivityC3, "áudio 30s")
	privacy := defaultPrivacy()

	result := ApplyStoragePolicy(ev, privacy)

	if result.StoragePolicy != StorageKept7d {
		t.Fatalf("expected kept_7d, got %s", result.StoragePolicy)
	}
}

func TestApplyStoragePolicy_C2_KeptCustom(t *testing.T) {
	ev := newTestEvidence(EvidenceText, SensitivityC2, "resposta")
	privacy := defaultPrivacy()

	result := ApplyStoragePolicy(ev, privacy)

	if result.StoragePolicy != StorageKeptCustom {
		t.Fatalf("expected kept_custom, got %s", result.StoragePolicy)
	}
}

// --- Gate profiles catalog tests ---

func TestGateProfileCatalog_SpeakingRequiresAudio(t *testing.T) {
	p := GateProfileCatalog["speaking_output"]
	if p.EquivalencePolicy != EquivAudioRequiredNoEquivalent {
		t.Fatal("speaking_output should have audio_required_no_equivalent")
	}
	if p.TaskClass != TaskClassLearning {
		t.Fatalf("expected learning, got %s", p.TaskClass)
	}
}

func TestGateProfileCatalog_ComprehensionAcceptsText(t *testing.T) {
	p := GateProfileCatalog["english_comprehension"]
	if p.EquivalencePolicy != EquivTextEquivalent {
		t.Fatal("english_comprehension should accept text_equivalent")
	}
}

func TestGateProfileCatalog_SleepDiaryIsHabit(t *testing.T) {
	p := GateProfileCatalog["sleep_diary"]
	if p.TaskClass != TaskClassHabit {
		t.Fatalf("expected habit, got %s", p.TaskClass)
	}
}

func TestGateProfileCatalog_AllProfilesExist(t *testing.T) {
	expected := []string{"speaking_output", "english_comprehension", "english_retrieval", "java_practice", "java_retrieval", "sleep_diary"}
	for _, id := range expected {
		if _, ok := GateProfileCatalog[id]; !ok {
			t.Fatalf("missing profile: %s", id)
		}
	}
}

func TestLookupGateProfile_Found(t *testing.T) {
	p := LookupGateProfile("speaking_output")
	if p == nil {
		t.Fatal("expected non-nil profile")
	}
	if p.ProfileID != "speaking_output" {
		t.Fatalf("expected speaking_output, got %s", p.ProfileID)
	}
}

func TestLookupGateProfile_NotFound(t *testing.T) {
	p := LookupGateProfile("nonexistent")
	if p != nil {
		t.Fatal("expected nil for unknown profile")
	}
}
