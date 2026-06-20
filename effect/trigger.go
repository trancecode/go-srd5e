package effect

// Outcome is a normalized resolution result: the game fills it from whichever
// resolution ran (attack, save, or contest). Saved and ContestWon are nil when
// that resolution did not occur.
type Outcome struct {
	Hit, Crit  bool
	Saved      *bool
	ContestWon *bool
}

// Triggered returns, in source order, the effects whose trigger fires for the
// outcome.
func Triggered(effs []ConditionalEffect, outcome Outcome) []Effect {
	var out []Effect
	for _, ce := range effs {
		if matches(ce.Trigger, outcome) {
			out = append(out, ce.Effect)
		}
	}
	return out
}

// matches reports whether a trigger fires for an outcome. OnHit also fires on a
// crit (a crit is a hit). OnMiss is meaningful only for an attack outcome: it
// fires when nothing hit and no save or contest was involved.
func matches(t Trigger, o Outcome) bool {
	switch t {
	case Always:
		return true
	case OnHit:
		return o.Hit || o.Crit
	case OnCrit:
		return o.Crit
	case OnMiss:
		return !o.Hit && !o.Crit && o.Saved == nil && o.ContestWon == nil
	case OnSaveFail:
		return o.Saved != nil && !*o.Saved
	case OnSaveSuccess:
		return o.Saved != nil && *o.Saved
	case OnSave:
		return o.Saved != nil
	case OnContestWin:
		return o.ContestWon != nil && *o.ContestWon
	case OnContestLose:
		return o.ContestWon != nil && !*o.ContestWon
	default: // TriggerUnspecified
		return false
	}
}
