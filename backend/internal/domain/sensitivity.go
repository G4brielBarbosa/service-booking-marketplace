package domain

type SensitivityLevel string

const (
	SensitivityC1 SensitivityLevel = "C1" // Check-ins and plans
	SensitivityC2 SensitivityLevel = "C2" // Learning evidence (non-sensitive)
	SensitivityC3 SensitivityLevel = "C3" // Sensitive content (audio, emotional text)
	SensitivityC4 SensitivityLevel = "C4" // Aggregated metrics
	SensitivityC5 SensitivityLevel = "C5" // Governance and preferences
)
