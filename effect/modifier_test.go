package effect

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestModifierBonus(t *testing.T) {
	bless := ModifierSpec{Source: "bless", Targets: []ModifierTarget{ModTargetAttack, ModTargetSave}, Ability: core.AbilityAny, Dice: dice.Expr{Count: 1, Sides: 4}}
	shield := ModifierSpec{Source: "shield_of_faith", Targets: []ModifierTarget{ModTargetAc}, Ability: core.AbilityAny, Flat: 2}

	// bless applies to attack rolls: no flat, one die.
	flat, ds := ModifierBonus([]ModifierSpec{bless, shield}, ModTargetAttack, core.AbilityAny)
	if flat != 0 || len(ds) != 1 || ds[0] != (dice.Expr{Count: 1, Sides: 4}) {
		t.Errorf("bless on attack = flat %d dice %v, want 0 and [1d4]", flat, ds)
	}
	// neither applies to damage.
	if flat, ds := ModifierBonus([]ModifierSpec{bless, shield}, ModTargetDamage, core.AbilityAny); flat != 0 || len(ds) != 0 {
		t.Errorf("damage = flat %d dice %v, want 0 and none", flat, ds)
	}
}

func TestModifierBonusDedupAndStack(t *testing.T) {
	a := ModifierSpec{Source: "src_a", Targets: []ModifierTarget{ModTargetAc}, Ability: core.AbilityAny, Flat: 2}
	b := ModifierSpec{Source: "src_b", Targets: []ModifierTarget{ModTargetAc}, Ability: core.AbilityAny, Flat: 2}
	// distinct sources stack.
	if flat, _ := ModifierBonus([]ModifierSpec{a, b}, ModTargetAc, core.AbilityAny); flat != 4 {
		t.Errorf("distinct sources = %d, want 4", flat)
	}
	// same source does not stack: collapses to one.
	if flat, _ := ModifierBonus([]ModifierSpec{a, a}, ModTargetAc, core.AbilityAny); flat != 2 {
		t.Errorf("same source = %d, want 2", flat)
	}
	// same source conflict: highest flat wins.
	hi := ModifierSpec{Source: "src_a", Targets: []ModifierTarget{ModTargetAc}, Ability: core.AbilityAny, Flat: 3}
	if flat, _ := ModifierBonus([]ModifierSpec{a, hi}, ModTargetAc, core.AbilityAny); flat != 3 {
		t.Errorf("same source conflict = %d, want 3 (highest)", flat)
	}
}

func TestModifierBonusAbilityNarrowing(t *testing.T) {
	strOnly := ModifierSpec{Source: "guidance_str", Targets: []ModifierTarget{ModTargetSave}, Ability: core.AbilityStrength, Flat: 1}
	// applies to a Strength save.
	if flat, _ := ModifierBonus([]ModifierSpec{strOnly}, ModTargetSave, core.AbilityStrength); flat != 1 {
		t.Errorf("str save = %d, want 1", flat)
	}
	// does not apply to a Dexterity save.
	if flat, _ := ModifierBonus([]ModifierSpec{strOnly}, ModTargetSave, core.AbilityDexterity); flat != 0 {
		t.Errorf("dex save = %d, want 0", flat)
	}
}
