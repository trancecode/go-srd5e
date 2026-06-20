package effect

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

func TestEnumZeroValues(t *testing.T) {
	if TriggerUnspecified != 0 || MovementUnspecified != 0 || ModTargetUnspecified != 0 {
		t.Error("Unspecified enum values must be zero")
	}
	// ordering sanity: Always is the last trigger.
	if !(OnHit < OnCrit && OnCrit < OnMiss && Always > OnContestLose) {
		t.Error("trigger ordering unexpected")
	}
}

func TestEffectArms(t *testing.T) {
	heal := dice.Expr{Count: 1, Sides: 8}
	dmg := damage.Single(dice.Expr{Count: 8, Sides: 6}, core.Fire)
	e := Effect{
		Damage:     &dmg,
		HalfOnSave: true,
		Healing:    &heal,
		Condition:  &ConditionSpec{Condition: core.Prone, Duration: core.EffectDuration{}},
		Movement:   &MovementSpec{Distance: 10, Kind: MovementPush},
		Modifier:   &ModifierSpec{Source: "bless", Targets: []ModifierTarget{ModTargetAttack, ModTargetSave}, Ability: core.AbilityAny, Dice: dice.Expr{Count: 1, Sides: 4}},
	}
	if e.Damage == nil || !e.HalfOnSave || e.Healing == nil || e.Condition.Condition != core.Prone || e.Movement.Kind != MovementPush || e.Modifier.Source != "bless" {
		t.Errorf("effect arms not wired: %+v", e)
	}
	ce := ConditionalEffect{Trigger: OnSave, Effect: e}
	if ce.Trigger != OnSave {
		t.Error("conditional effect trigger not set")
	}
}
