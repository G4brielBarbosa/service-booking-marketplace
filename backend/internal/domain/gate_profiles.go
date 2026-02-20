package domain

// GateProfileCatalog maps profile_id to QualityGateProfile.
// MVP profiles: speaking, comprehension, retrieval (english), java practice/retrieval, sleep diary.
var GateProfileCatalog = map[string]QualityGateProfile{
	"speaking_output": {
		ProfileID: "speaking_output",
		Domain:    GoalEnglish,
		TaskClass: TaskClassLearning,
		Requirements: []GateRequirement{
			{Kind: "rubric", Description: "Preencha a rubrica de speaking (4 dimensões: clarity, fluency, accuracy, vocabulary, cada 0-2).", MinCount: 1},
			{Kind: "audio", Description: "Grave um áudio curto falando sobre o conteúdo.", MinCount: 1},
		},
		ValidityRules:     []string{"not_empty", "audio_required"},
		EquivalencePolicy: EquivAudioRequiredNoEquivalent,
	},

	"english_comprehension": {
		ProfileID: "english_comprehension",
		Domain:    GoalEnglish,
		TaskClass: TaskClassLearning,
		Requirements: []GateRequirement{
			{Kind: "text_answer", Description: "Responda 3 perguntas de compreensão sobre o conteúdo.", MinCount: 3},
		},
		ValidityRules:     []string{"not_empty"},
		EquivalencePolicy: EquivTextEquivalent,
	},

	"english_retrieval": {
		ProfileID: "english_retrieval",
		Domain:    GoalEnglish,
		TaskClass: TaskClassLearning,
		Requirements: []GateRequirement{
			{Kind: "text_answer", Description: "Faça recall de pelo menos 5 itens de vocabulário.", MinCount: 5},
		},
		ValidityRules:     []string{"not_empty"},
		EquivalencePolicy: EquivTextEquivalent,
	},

	"java_practice": {
		ProfileID: "java_practice",
		Domain:    GoalJava,
		TaskClass: TaskClassLearning,
		Requirements: []GateRequirement{
			{Kind: "text_answer", Description: "Envie o código do exercício e uma breve explicação.", MinCount: 1},
		},
		ValidityRules:     []string{"not_empty"},
		EquivalencePolicy: EquivTextEquivalent,
	},

	"java_retrieval": {
		ProfileID: "java_retrieval",
		Domain:    GoalJava,
		TaskClass: TaskClassLearning,
		Requirements: []GateRequirement{
			{Kind: "text_answer", Description: "Faça recall dos conceitos estudados (texto curto).", MinCount: 1},
		},
		ValidityRules:     []string{"not_empty"},
		EquivalencePolicy: EquivTextEquivalent,
	},

	"sleep_diary": {
		ProfileID: "sleep_diary",
		Domain:    GoalSleep,
		TaskClass: TaskClassHabit,
		Requirements: []GateRequirement{
			{Kind: "metadata", Description: "Registre horário de dormir, acordar e qualidade (mínimo).", MinCount: 1},
		},
		ValidityRules:     []string{"not_empty"},
		EquivalencePolicy: EquivTextEquivalent,
	},
}

// LookupGateProfile returns the profile for a given gate_ref, or nil if not found.
func LookupGateProfile(gateRef string) *QualityGateProfile {
	p, ok := GateProfileCatalog[gateRef]
	if !ok {
		return nil
	}
	return &p
}
