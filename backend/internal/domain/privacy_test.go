package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestDefaultPrivacyPolicy(t *testing.T) {
	pp := NewDefaultPrivacyPolicy(uuid.New())

	if pp.MinimalMode {
		t.Error("expected minimal_mode to be false by default")
	}

	if len(pp.OptOutCategories) != 0 {
		t.Errorf("expected 0 opt-out categories, got %d", len(pp.OptOutCategories))
	}

	expectedRetention := map[SensitivityLevel]int{
		SensitivityC1: 90,
		SensitivityC2: 90,
		SensitivityC3: 7,
		SensitivityC4: 365,
		SensitivityC5: 0,
	}

	for cat, want := range expectedRetention {
		got := pp.RetentionDays[cat]
		if got != want {
			t.Errorf("retention for %s: expected %d, got %d", cat, want, got)
		}
	}
}

func TestOptOutToggle(t *testing.T) {
	pp := NewDefaultPrivacyPolicy(uuid.New())

	// Opt out of C3
	pp.SetOptOut(SensitivityC3, true)
	if !pp.IsOptedOut(SensitivityC3) {
		t.Error("expected C3 to be opted out")
	}

	// Double opt-out should not duplicate
	pp.SetOptOut(SensitivityC3, true)
	count := 0
	for _, c := range pp.OptOutCategories {
		if c == SensitivityC3 {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 C3 entry, got %d", count)
	}

	// Opt back in
	pp.SetOptOut(SensitivityC3, false)
	if pp.IsOptedOut(SensitivityC3) {
		t.Error("expected C3 to not be opted out after removing")
	}
}

func TestOptOutDoesNotAffectOthers(t *testing.T) {
	pp := NewDefaultPrivacyPolicy(uuid.New())

	pp.SetOptOut(SensitivityC3, true)

	if pp.IsOptedOut(SensitivityC1) {
		t.Error("C1 should not be affected by C3 opt-out")
	}
	if pp.IsOptedOut(SensitivityC2) {
		t.Error("C2 should not be affected by C3 opt-out")
	}
}
