package damage

import (
	"testing"

	"github.com/trancecode/go-srd5e/combat"
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestRollHit(t *testing.T) {
	// 1d6 slashing + bonus 3, hit, Constant(6) -> one part amount 9.
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackHit, dice.Constant(6))
	if len(d.Parts) != 1 || d.Parts[0].Amount != 9 || d.Parts[0].Type != core.Slashing {
		t.Errorf("hit = %+v, want one slashing part amount 9", d)
	}
}

func TestRollCriticalDoublesDiceNotBonus(t *testing.T) {
	// 1d6 slashing + bonus 3, crit, Constant(6) -> 2d6 = 12, + 3 = 15.
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackCritical, dice.Constant(6))
	if d.Parts[0].Amount != 15 {
		t.Errorf("crit amount = %d, want 15 (doubled dice, bonus once)", d.Parts[0].Amount)
	}
}

func TestRollMiss(t *testing.T) {
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackMiss, dice.Constant(6))
	if len(d.Parts) != 0 {
		t.Errorf("miss = %+v, want no parts", d)
	}
}

func TestRollBonusOnPrimaryPartOnly(t *testing.T) {
	// Two parts: 1d8 slashing + 1d6 fire, bonus 2, hit, Constant(8)/Constant(6) won't
	// vary per part with one Constant, so use Constant(10): slashing die shows 8, fire die shows 6.
	spec := Spec{Parts: []PartSpec{
		{Dice: dice.Expr{Count: 1, Sides: 8}, Type: core.Slashing},
		{Dice: dice.Expr{Count: 1, Sides: 6}, Type: core.Fire},
	}}
	d := Roll(spec, 2, combat.AttackHit, dice.Constant(10))
	if d.Parts[0].Amount != 10 { // 8 + bonus 2
		t.Errorf("primary part = %d, want 10 (8 + bonus 2)", d.Parts[0].Amount)
	}
	if d.Parts[1].Amount != 6 { // fire, no bonus
		t.Errorf("secondary part = %d, want 6 (no bonus)", d.Parts[1].Amount)
	}
}
