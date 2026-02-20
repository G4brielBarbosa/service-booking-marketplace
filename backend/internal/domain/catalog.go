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

var speakingGateRef = "speaking_output"

var catalog = map[GoalID]map[PlanType][]TaskTemplate{
	GoalEnglish: {
		PlanA: {
			{Title: "Input em inglês", GoalDomain: GoalEnglish, EstimatedMin: 30, Instructions: "Assista ou leia conteúdo em inglês por 30 min", DoneCriteria: "30 min de input registrados"},
			{Title: "Prática de speaking", GoalDomain: GoalEnglish, EstimatedMin: 15, Instructions: "Grave 1 min falando sobre o conteúdo consumido", DoneCriteria: "Áudio gravado", GateProfile: &speakingGateRef},
			{Title: "Retrieval rápido", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Revise 10 cards de vocabulário", DoneCriteria: "10 cards revisados"},
		},
		PlanB: {
			{Title: "Input em inglês", GoalDomain: GoalEnglish, EstimatedMin: 20, Instructions: "Leia ou ouça conteúdo em inglês por 20 min", DoneCriteria: "20 min de input registrados"},
			{Title: "Retrieval rápido", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Revise 5 cards de vocabulário", DoneCriteria: "5 cards revisados"},
		},
		PlanC: {
			{Title: "Listening rápido", GoalDomain: GoalEnglish, EstimatedMin: 5, Instructions: "Ouça 5 min de podcast/vídeo em inglês", DoneCriteria: "5 min de listening concluídos"},
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
