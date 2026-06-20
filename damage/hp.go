package damage

import "github.com/trancecode/go-srd5e/core"

// HpOutcome is the result of applying a signed HP change. DroppedToZero and
// InstantDeath are meaningful only for damage (a negative delta).
type HpOutcome struct {
	Hp            core.HitPoints
	DroppedToZero bool
	InstantDeath  bool
}

// ApplyToHp applies a signed HP change: a negative delta is damage (floors at 0,
// sets DroppedToZero, and InstantDeath when the damage remaining past 0 is at
// least the hit-point maximum), a positive delta is healing (caps at max).
func ApplyToHp(current, max core.HitPoints, delta int) HpOutcome {
	next := int(current) + delta
	out := HpOutcome{}
	if delta < 0 {
		if next <= 0 {
			out.DroppedToZero = true
			if -next >= int(max) {
				out.InstantDeath = true
			}
			next = 0
		}
	} else if next > int(max) {
		next = int(max)
	}
	out.Hp = core.HitPoints(next)
	return out
}

// Apply bundles ApplyMitigation and ApplyToHp for the damage path, applying the
// mitigated total as a negative delta.
func Apply(d Damage, m Mitigation, cur, max core.HitPoints) (HpOutcome, Result) {
	res := ApplyMitigation(d, m)
	return ApplyToHp(cur, max, -res.Final), res
}
