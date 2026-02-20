package domain

import (
	"time"

	"github.com/google/uuid"
)

// Default retention in days per category (PLAN-000 §8, SPEC-015).
const (
	RetentionC1Default = 90   // Check-ins/plans: 90 days
	RetentionC2Default = 90   // Evidence (non-sensitive): 90 days
	RetentionC3Default = 7    // Sensitive raw content: 7 days
	RetentionC4Default = 365  // Weekly aggregates: 12 months
	RetentionC5Default = 0    // Governance/prefs: while active (0 = no auto-expire)
)

type PrivacyPolicy struct {
	UserID            uuid.UUID        `json:"user_id"`
	OptOutCategories  []SensitivityLevel `json:"opt_out_categories"`
	RetentionDays     map[SensitivityLevel]int `json:"retention_days"`
	MinimalMode       bool             `json:"minimal_mode_enabled"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

func NewDefaultPrivacyPolicy(userID uuid.UUID) PrivacyPolicy {
	return PrivacyPolicy{
		UserID:           userID,
		OptOutCategories: []SensitivityLevel{},
		RetentionDays: map[SensitivityLevel]int{
			SensitivityC1: RetentionC1Default,
			SensitivityC2: RetentionC2Default,
			SensitivityC3: RetentionC3Default,
			SensitivityC4: RetentionC4Default,
			SensitivityC5: RetentionC5Default,
		},
		MinimalMode: false,
		UpdatedAt:   time.Now(),
	}
}

func (p *PrivacyPolicy) IsOptedOut(cat SensitivityLevel) bool {
	for _, c := range p.OptOutCategories {
		if c == cat {
			return true
		}
	}
	return false
}

func (p *PrivacyPolicy) SetOptOut(cat SensitivityLevel, optOut bool) {
	if optOut && !p.IsOptedOut(cat) {
		p.OptOutCategories = append(p.OptOutCategories, cat)
	} else if !optOut {
		filtered := make([]SensitivityLevel, 0, len(p.OptOutCategories))
		for _, c := range p.OptOutCategories {
			if c != cat {
				filtered = append(filtered, c)
			}
		}
		p.OptOutCategories = filtered
	}
	p.UpdatedAt = time.Now()
}
