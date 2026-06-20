package turn

import "testing"

func TestEconomy(t *testing.T) {
	e := Economy{MovementSpeed: 30}
	if e.MovementRemaining() != 30 || !e.CanReact() {
		t.Error("fresh economy should have full movement and a reaction")
	}
	e.MovementUsed = 20
	if e.MovementRemaining() != 10 || !e.CanMove(10) || e.CanMove(11) {
		t.Errorf("remaining = %d; CanMove wrong", e.MovementRemaining())
	}
	e.ReactionUsed = true
	if e.CanReact() {
		t.Error("used reaction: CanReact should be false")
	}
	// over-spent movement floors at 0.
	e.MovementUsed = 40
	if e.MovementRemaining() != 0 {
		t.Errorf("overspent remaining = %d, want 0", e.MovementRemaining())
	}
	// ResetTurn clears usage but keeps speed.
	e.ActionUsed, e.BonusUsed = true, true
	e.ResetTurn()
	if e.ActionUsed || e.BonusUsed || e.ReactionUsed || e.MovementUsed != 0 || e.MovementSpeed != 30 {
		t.Errorf("after reset = %+v, want all usage cleared, speed 30", e)
	}
}
