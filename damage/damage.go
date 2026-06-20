package damage

import (
	"github.com/trancecode/go-srd5e/combat"
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// DamagePart is a rolled amount of one damage type. Magical marks damage that
// overcomes "resistance to nonmagical" defenses.
type DamagePart struct {
	Amount  int
	Type    core.DamageType
	Magical bool
}

// Damage is the rolled, typed result of an attack or effect.
type Damage struct {
	Parts []DamagePart
}

// PartSpec describes one damage component before rolling.
type PartSpec struct {
	Dice    dice.Expr
	Type    core.DamageType
	Magical bool
}

// Spec is a weapon or effect's damage, possibly of several types.
type Spec struct {
	Parts []PartSpec
}

// Roll keys off the attack outcome: a miss yields zero damage, a critical
// doubles the dice per part. bonus is the attacker's flat modifier; it is added
// once to the primary part (so it carries that part's type for resistance and is
// not doubled on a crit), flooring that part at zero.
func Roll(spec Spec, bonus core.Modifier, outcome combat.AttackOutcome, r dice.Roller) Damage {
	if outcome == combat.AttackMiss {
		return Damage{}
	}
	crit := outcome == combat.AttackCritical
	parts := make([]DamagePart, len(spec.Parts))
	for i, ps := range spec.Parts {
		var res dice.Result
		if crit {
			res = ps.Dice.RollCritical(r)
		} else {
			res = ps.Dice.Roll(r)
		}
		amount := res.Total
		if i == 0 {
			amount += int(bonus)
			if amount < 0 {
				amount = 0
			}
		}
		parts[i] = DamagePart{Amount: amount, Type: ps.Type, Magical: ps.Magical}
	}
	return Damage{Parts: parts}
}

// Single builds a single-part damage spec (the common one-type weapon).
func Single(expr dice.Expr, t core.DamageType) Spec {
	return Spec{Parts: []PartSpec{{Dice: expr, Type: t}}}
}

// RollSingle rolls a single-part spec; convenience over Single + Roll.
func RollSingle(expr dice.Expr, t core.DamageType, bonus core.Modifier, outcome combat.AttackOutcome, r dice.Roller) Damage {
	return Roll(Single(expr, t), bonus, outcome, r)
}
