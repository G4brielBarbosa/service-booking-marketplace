package domain

// HardcodedCatalog is the MVP TaskCatalog with generic templates per goal/plan.
// Domain-specific SPECs (SPEC-004/005/006) will refine these later.
type HardcodedCatalog struct{}

func NewHardcodedCatalog() *HardcodedCatalog {
	return &HardcodedCatalog{}
}

func (c *HardcodedCatalog) GetTasksForGoal(goal GoalID, pt PlanType) []TaskTemplate {
	templates, ok := catalog[goal]
	if !ok {
		return nil
	}
	result, ok := templates[pt]
	if !ok {
		return nil
	}
	return result
}

var (
	speakingGateRef        = "speaking_output"
	comprehensionGateRef   = "english_comprehension"
	retrievalGateRef       = "english_retrieval"
	comprehensionMinRef    = "english_comprehension_min"
	retrievalMinRef        = "english_retrieval_min"
)

var catalog = map[GoalID]map[PlanType][]TaskTemplate{
	GoalEnglish: {
		PlanA: {
			{Title: "Input em inglês", GoalDomain: GoalEnglish, EstimatedMin: 30, Instructions: "Assista ou leia conteúdo em inglês por 30 min. Depois, responda 3 perguntas de compreensão.", DoneCriteria: "3 respostas de compreensão validadas", GateProfile: &comprehensionGateRef},
			{Title: "Prática de speaking", GoalDomain: GoalEnglish, EstimatedMin: 15, Instructions: "Grave 1-2 min falando sobre o conteúdo consumido. Preencha a rubrica (clarity, fluency, accuracy, vocabulary: 0-2 cada).", DoneCriteria: "Áudio gravado + rubrica preenchida", GateProfile: &speakingGateRef},
			{Title: "Retrieval rápido", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Faça recall de 10 itens de vocabulário sem consultar.", DoneCriteria: "10 itens de recall registrados", GateProfile: &retrievalGateRef},
		},
		PlanB: {
			{Title: "Input em inglês", GoalDomain: GoalEnglish, EstimatedMin: 20, Instructions: "Leia ou ouça conteúdo em inglês por 20 min. Responda 3 perguntas de compreensão.", DoneCriteria: "3 respostas de compreensão validadas", GateProfile: &comprehensionGateRef},
			{Title: "Retrieval rápido", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Faça recall de 5 itens de vocabulário.", DoneCriteria: "5 itens de recall registrados", GateProfile: &retrievalGateRef},
		},
		PlanC: {
			{Title: "Listening rápido", GoalDomain: GoalEnglish, EstimatedMin: 10, Instructions: "Ouça 10 min de podcast/vídeo em inglês. Responda 1 pergunta de compreensão.", DoneCriteria: "1 resposta de compreensão validada", GateProfile: &comprehensionMinRef},
			{Title: "Retrieval mínimo", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Faça recall de 3 itens.", DoneCriteria: "3 itens de recall registrados", GateProfile: &retrievalMinRef},
		},
	},
	GoalJava: {
		PlanA: {
			{Title: "Estudo de Java", GoalDomain: GoalJava, EstimatedMin: 30, Instructions: "Estude um tópico do roteiro por 30 min", DoneCriteria: "30 min de estudo registrados"},
			{Title: "Prática de código", GoalDomain: GoalJava, EstimatedMin: 20, Instructions: "Resolva 1 exercício prático", DoneCriteria: "Exercício resolvido e testado"},
			{Title: "Revisão de conceitos", GoalDomain: GoalJava, EstimatedMin: 5, Instructions: "Revise flashcards de Java", DoneCriteria: "Revisão concluída"},
		},
		PlanB: {
			{Title: "Estudo de Java", GoalDomain: GoalJava, EstimatedMin: 20, Instructions: "Estude um tópico por 20 min", DoneCriteria: "20 min de estudo registrados"},
			{Title: "Revisão de conceitos", GoalDomain: GoalJava, EstimatedMin: 5, Instructions: "Revise flashcards de Java", DoneCriteria: "Revisão concluída"},
		},
		PlanC: {
			{Title: "Leitura rápida de Java", GoalDomain: GoalJava, EstimatedMin: 10, Instructions: "Leia 1 artigo curto sobre Java", DoneCriteria: "Artigo lido"},
		},
	},
	GoalSleep: {
		PlanA: {
			{Title: "Registro de sono", GoalDomain: GoalSleep, EstimatedMin: 2, Instructions: "Registre horário e qualidade do sono", DoneCriteria: "Sono registrado"},
		},
		PlanB: {
			{Title: "Registro de sono", GoalDomain: GoalSleep, EstimatedMin: 2, Instructions: "Registre horário e qualidade do sono", DoneCriteria: "Sono registrado"},
		},
		PlanC: {
			{Title: "Registro de sono", GoalDomain: GoalSleep, EstimatedMin: 1, Instructions: "Registre se dormiu bem ou mal", DoneCriteria: "Registro feito"},
		},
	},
	GoalHealth: {
		PlanA: {
			{Title: "Atividade física", GoalDomain: GoalHealth, EstimatedMin: 15, Instructions: "Faça 15 min de exercício ou caminhada", DoneCriteria: "15 min de atividade registrados"},
		},
		PlanB: {
			{Title: "Atividade leve", GoalDomain: GoalHealth, EstimatedMin: 10, Instructions: "Faça 10 min de alongamento ou caminhada", DoneCriteria: "10 min de atividade registrados"},
		},
		PlanC: {
			{Title: "Movimento mínimo", GoalDomain: GoalHealth, EstimatedMin: 5, Instructions: "Faça 5 min de alongamento", DoneCriteria: "Alongamento feito"},
		},
	},
	GoalSelfEsteem: {
		PlanA: {
			{Title: "Diário de gratidão", GoalDomain: GoalSelfEsteem, EstimatedMin: 5, Instructions: "Escreva 3 coisas boas do dia", DoneCriteria: "3 itens registrados"},
		},
		PlanB: {
			{Title: "Registro de gratidão", GoalDomain: GoalSelfEsteem, EstimatedMin: 3, Instructions: "Escreva 1-2 coisas boas", DoneCriteria: "Registro feito"},
		},
		PlanC: {
			{Title: "Gratidão rápida", GoalDomain: GoalSelfEsteem, EstimatedMin: 1, Instructions: "Pense em 1 coisa boa do dia", DoneCriteria: "Reflexão feita"},
		},
	},
	GoalSaaS: {
		PlanA: {
			{Title: "Avanço no SaaS", GoalDomain: GoalSaaS, EstimatedMin: 30, Instructions: "Trabalhe no próximo item do backlog", DoneCriteria: "Item do backlog avançado"},
		},
		PlanB: {
			{Title: "Micro-avanço no SaaS", GoalDomain: GoalSaaS, EstimatedMin: 15, Instructions: "Faça 1 tarefa pequena do SaaS", DoneCriteria: "Tarefa concluída"},
		},
		PlanC: {
			{Title: "Anotar próximo passo SaaS", GoalDomain: GoalSaaS, EstimatedMin: 5, Instructions: "Anote o próximo passo do SaaS", DoneCriteria: "Próximo passo anotado"},
		},
	},
}
