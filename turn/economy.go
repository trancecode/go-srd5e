package turn

import "github.com/trancecode/go-srd5e/core"

// Economy is one participant's action economy for a round. SRD 5e has no
// separate move action: movement is a budget (MovementSpeed) split around the
// action, so this models movement as MovementUsed against MovementSpeed rather
// than an action slot.
type Economy struct {
	ActionUsed, BonusUsed, ReactionUsed bool
	MovementUsed                        core.Distance
	MovementSpeed                       core.Distance
}

// ResetTurn clears per-turn usage at the start of the creature's turn (including
// the once-per-round reaction), keeping the movement speed.
func (e *Economy) ResetTurn() {
	e.ActionUsed = false
	e.BonusUsed = false
	e.ReactionUsed = false
	e.MovementUsed = 0
}

// CanReact reports whether the reaction is available this round.
func (e Economy) CanReact() bool { return !e.ReactionUsed }

// MovementRemaining is the movement budget left, floored at zero.
func (e Economy) MovementRemaining() core.Distance {
	if e.MovementUsed >= e.MovementSpeed {
		return 0
	}
	return e.MovementSpeed - e.MovementUsed
}

// CanMove reports whether a measured cost fits in the remaining budget.
func (e Economy) CanMove(cost core.Distance) bool { return cost <= e.MovementRemaining() }
