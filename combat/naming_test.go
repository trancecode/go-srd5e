package combat

import "testing"

func TestAttackOutcomeString(t *testing.T) {
	cases := []struct {
		o    AttackOutcome
		want string
	}{
		{AttackMiss, "miss"},
		{AttackHit, "hit"},
		{AttackCritical, "critical"},
		{AttackOutcomeNone, "none"},
	}
	for _, c := range cases {
		if got := c.o.String(); got != c.want {
			t.Errorf("AttackOutcome(%d).String() = %q, want %q", c.o, got, c.want)
		}
	}
}

// AttackOutcomeNone is the zero value, so an unset outcome never reads as a miss.
func TestAttackOutcomeNoneIsZero(t *testing.T) {
	var zero AttackOutcome
	if zero != AttackOutcomeNone {
		t.Errorf("zero AttackOutcome = %d, want AttackOutcomeNone (%d)", zero, AttackOutcomeNone)
	}
}
