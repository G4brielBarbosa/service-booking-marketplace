package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Event types for quality gates (SPEC-016, no C3 content) ---

const (
	EventGateRequired       EventType = "gate_required"
	EventEvidenceReceived   EventType = "evidence_received"
	EventGateEvaluated      EventType = "gate_evaluated"
	EventTaskBlockedByPolicy EventType = "task_blocked_by_policy"
)

// --- Evidence kind ---

type EvidenceKind string

const (
	EvidenceText     EvidenceKind = "text_answer"
	EvidenceRubric   EvidenceKind = "rubric"
	EvidenceAudio    EvidenceKind = "audio"
	EvidenceMetadata EvidenceKind = "metadata"
)

// --- Storage policy ---

type StoragePolicy string

const (
	StorageKept7d              StoragePolicy = "kept_7d"
	StorageKeptCustom          StoragePolicy = "kept_custom"
	StorageDiscardedAfterProc  StoragePolicy = "discarded_after_processing"
)

// --- Gate status ---

type GateStatus string

const (
	GateSatisfied    GateStatus = "satisfied"
	GateNotSatisfied GateStatus = "not_satisfied"
)

// --- Failure reason codes ---

type FailureReasonCode string

const (
	FailureMissing       FailureReasonCode = "missing"
	FailureInvalid       FailureReasonCode = "invalid"
	FailurePartial       FailureReasonCode = "partial"
	FailureNotEquivalent FailureReasonCode = "not_equivalent"
	FailureNone          FailureReasonCode = ""
)

// --- Task class (learning vs habit) ---

type TaskClass string

const (
	TaskClassLearning TaskClass = "learning"
	TaskClassHabit    TaskClass = "habit"
)

// --- Gate requirement ---

type GateRequirement struct {
	Kind        string `json:"kind"`
	Description string `json:"description"`
	MinCount    int    `json:"min_count,omitempty"`
}

// --- Equivalence policy ---

type EquivalencePolicy string

const (
	EquivTextEquivalent            EquivalencePolicy = "text_equivalent"
	EquivAudioRequiredNoEquivalent EquivalencePolicy = "audio_required_no_equivalent"
)

// --- QualityGateProfile ---

type QualityGateProfile struct {
	ProfileID         string            `json:"profile_id"`
	Domain            GoalID            `json:"domain"`
	TaskClass         TaskClass         `json:"task_class"`
	Requirements      []GateRequirement `json:"requirements"`
	ValidityRules     []string          `json:"validity_rules"`
	EquivalencePolicy EquivalencePolicy `json:"equivalence_policy"`
}

// --- Evidence ---

type Evidence struct {
	EvidenceID    uuid.UUID        `json:"evidence_id"`
	UserID        uuid.UUID        `json:"user_id"`
	TaskID        uuid.UUID        `json:"task_id"`
	Kind          EvidenceKind     `json:"kind"`
	Sensitivity   SensitivityLevel `json:"sensitivity"`
	StoragePolicy StoragePolicy    `json:"storage_policy"`
	ContentRef    *string          `json:"content_ref,omitempty"`
	Summary       string           `json:"summary"`
	Timestamp     time.Time        `json:"timestamp"`
}

// --- GateResult ---

type GateResult struct {
	GateResultID     uuid.UUID         `json:"gate_result_id"`
	UserID           uuid.UUID         `json:"user_id"`
	TaskID           uuid.UUID         `json:"task_id"`
	GateStatus       GateStatus        `json:"gate_status"`
	FailureReason    FailureReasonCode `json:"failure_reason_code"`
	ReasonShort      string            `json:"reason_short"`
	NextMinStep      string            `json:"next_min_step"`
	EvidenceIDs      []uuid.UUID       `json:"evidence_ids"`
	DerivedMetrics   map[string]any    `json:"derived_metrics"`
	Timestamp        time.Time         `json:"timestamp"`
}

// --- RubricScore ---

type RubricScore struct {
	RubricID   uuid.UUID      `json:"rubric_id"`
	UserID     uuid.UUID      `json:"user_id"`
	TaskID     uuid.UUID      `json:"task_id"`
	Domain     GoalID         `json:"domain"`
	Dimensions map[string]int `json:"dimensions"`
	Total      int            `json:"total"`
	Status     string         `json:"status"`
}

// --- Completion request result (use case output) ---

type CompletionRequestStatus string

const (
	CompletionCompleted        CompletionRequestStatus = "completed"
	CompletionEvidenceRequired CompletionRequestStatus = "evidence_required"
	CompletionAlreadyCompleted CompletionRequestStatus = "already_completed"
)

type EvidenceRequest struct {
	ProfileID          string            `json:"profile_id"`
	Requirements       []GateRequirement `json:"requirements"`
	ValidityRules      []string          `json:"validity_rules"`
	PrivacyDisclosure  string            `json:"privacy_disclosure_short"`
}

type CompletionRequestResult struct {
	Status          CompletionRequestStatus `json:"status"`
	EvidenceRequest *EvidenceRequest        `json:"evidence_request,omitempty"`
}

// --- Evidence receipt (use case output) ---

type EvidenceReceipt struct {
	EvidenceID uuid.UUID     `json:"evidence_id"`
	Valid      bool          `json:"valid"`
	Reason     string        `json:"reason,omitempty"`
	Stored     StoragePolicy `json:"stored"`
}

// --- Gate result view (use case output) ---

type GateResultView struct {
	GateStatus  GateStatus        `json:"gate_status"`
	ReasonShort string            `json:"reason_short"`
	NextMinStep string            `json:"next_min_step"`
	Metrics     map[string]any    `json:"result_artifacts,omitempty"`
}

// --- Gate summary item (for GetTodayGateSummary) ---

type GateSummaryItem struct {
	TaskID      uuid.UUID  `json:"task_id"`
	TaskTitle   string     `json:"task_title"`
	GateStatus  GateStatus `json:"gate_status"`
	ReasonShort string     `json:"reason_short,omitempty"`
	NextMinStep string     `json:"next_min_step,omitempty"`
}

// --- Deterministic rules (PLAN-003 §8, pure functions) ---

// ValidateEvidence checks a single evidence against its gate profile requirements.
func ValidateEvidence(evidence Evidence, profile QualityGateProfile) (valid bool, reason string) {
	if evidence.Summary == "" && evidence.ContentRef == nil {
		return false, "Evidência vazia. Envie o conteúdo solicitado."
	}

	if profile.EquivalencePolicy == EquivAudioRequiredNoEquivalent {
		hasAudioReq := false
		for _, r := range profile.Requirements {
			if r.Kind == "audio" {
				hasAudioReq = true
				break
			}
		}
		if hasAudioReq && evidence.Kind != EvidenceAudio {
			return false, "Este gate exige áudio. Texto não é equivalente para speaking."
		}
	}

	for _, rule := range profile.ValidityRules {
		switch rule {
		case "not_empty":
			if evidence.Summary == "" && evidence.ContentRef == nil {
				return false, "Evidência não pode ser vazia."
			}
		case "audio_required":
			if evidence.Kind != EvidenceAudio {
				return false, "Áudio é obrigatório para esta tarefa."
			}
		}
	}

	return true, ""
}

// EvaluateGate runs deterministic evaluation of a gate given profile, collected evidences, and privacy policy.
func EvaluateGate(profile QualityGateProfile, evidences []Evidence, privacyPolicy PrivacyPolicy) GateResult {
	now := time.Now()
	result := GateResult{
		GateResultID:   uuid.New(),
		DerivedMetrics: map[string]any{},
		Timestamp:      now,
	}

	if len(evidences) == 0 {
		result.GateStatus = GateNotSatisfied
		result.FailureReason = FailureMissing
		result.ReasonShort = "Nenhuma evidência enviada."
		result.NextMinStep = nextMinStepForProfile(profile)
		return result
	}

	var validEvidences []Evidence
	var lastInvalidReason string
	hasEquivalenceIssue := false

	for _, ev := range evidences {
		valid, reason := ValidateEvidence(ev, profile)
		if valid {
			validEvidences = append(validEvidences, ev)
		} else {
			lastInvalidReason = reason
			if profile.EquivalencePolicy == EquivAudioRequiredNoEquivalent && ev.Kind != EvidenceAudio {
				hasEquivalenceIssue = true
			}
		}
	}

	if hasEquivalenceIssue && len(validEvidences) == 0 {
		result.GateStatus = GateNotSatisfied
		result.FailureReason = FailureNotEquivalent
		result.ReasonShort = "Texto não substitui áudio para speaking."
		result.NextMinStep = "Grave um áudio curto falando sobre o conteúdo."
		result.EvidenceIDs = evidenceIDs(evidences)
		return result
	}

	if len(validEvidences) == 0 {
		result.GateStatus = GateNotSatisfied
		result.FailureReason = FailureInvalid
		result.ReasonShort = lastInvalidReason
		if result.ReasonShort == "" {
			result.ReasonShort = "Evidência inválida."
		}
		result.NextMinStep = nextMinStepForProfile(profile)
		result.EvidenceIDs = evidenceIDs(evidences)
		return result
	}

	totalRequired := totalMinCount(profile)
	if totalRequired > 0 && len(validEvidences) < totalRequired {
		result.GateStatus = GateNotSatisfied
		result.FailureReason = FailurePartial
		result.ReasonShort = fmt.Sprintf("Evidência parcial: %d/%d itens válidos.", len(validEvidences), totalRequired)
		result.NextMinStep = fmt.Sprintf("Envie mais %d item(ns) para completar.", totalRequired-len(validEvidences))
		result.EvidenceIDs = evidenceIDs(evidences)
		result.DerivedMetrics["valid_count"] = len(validEvidences)
		result.DerivedMetrics["required_count"] = totalRequired
		return result
	}

	result.GateStatus = GateSatisfied
	result.FailureReason = FailureNone
	result.ReasonShort = "Gate satisfeito."
	result.NextMinStep = ""
	result.EvidenceIDs = evidenceIDs(validEvidences)
	result.DerivedMetrics["valid_count"] = len(validEvidences)

	return result
}

// ApplyStoragePolicy determines and applies the storage policy based on evidence sensitivity and user privacy.
func ApplyStoragePolicy(evidence Evidence, privacyPolicy PrivacyPolicy) Evidence {
	switch evidence.Sensitivity {
	case SensitivityC3:
		if privacyPolicy.IsOptedOut(SensitivityC3) {
			evidence.StoragePolicy = StorageDiscardedAfterProc
			evidence.ContentRef = nil
		} else {
			evidence.StoragePolicy = StorageKept7d
		}
	case SensitivityC2:
		evidence.StoragePolicy = StorageKeptCustom
	default:
		evidence.StoragePolicy = StorageKeptCustom
	}
	return evidence
}

// --- helpers ---

func nextMinStepForProfile(profile QualityGateProfile) string {
	if len(profile.Requirements) > 0 {
		return profile.Requirements[0].Description
	}
	return "Envie a evidência mínima solicitada."
}

func totalMinCount(profile QualityGateProfile) int {
	total := 0
	for _, r := range profile.Requirements {
		if r.MinCount > 0 {
			total += r.MinCount
		}
	}
	return total
}

func evidenceIDs(evs []Evidence) []uuid.UUID {
	ids := make([]uuid.UUID, len(evs))
	for i, e := range evs {
		ids[i] = e.EvidenceID
	}
	return ids
}
