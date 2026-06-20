package core

import "testing"

func TestRoundsConversion(t *testing.T) {
	if got := RoundsInMinutes(1); got != 10 {
		t.Errorf("RoundsInMinutes(1) = %d, want 10", got)
	}
	if got := RoundsInHours(1); got != 600 {
		t.Errorf("RoundsInHours(1) = %d, want 600", got)
	}
}

func TestEffectDurationZero(t *testing.T) {
	var d EffectDuration
	if d.Kind != DurationUnspecified {
		t.Errorf("zero EffectDuration.Kind = %v, want DurationUnspecified", d.Kind)
	}
}
