package effect

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

// Trigger is the resolution outcome that fires an effect arm.
type Trigger int

const (
	TriggerUnspecified Trigger = iota
	OnHit
	OnCrit
	OnMiss
	OnSaveFail
	OnSaveSuccess
	OnSave // either save outcome
	OnContestWin
	OnContestLose
	Always
)

// ConditionSpec describes a condition to apply for a duration. The game applies
// it (it touches entity state).
type ConditionSpec struct {
	Condition core.Condition
	Duration  core.EffectDuration
}

// MovementKind is a kind of forced movement, relative to the source.
type MovementKind int

const (
	MovementUnspecified MovementKind = iota
	MovementPush
	MovementPull
)

// MovementSpec describes forced movement; the game executes it (geometry).
type MovementSpec struct {
	Distance core.Distance
	Kind     MovementKind
}

// ModifierTarget is what a ModifierSpec adjusts.
type ModifierTarget int

const (
	ModTargetUnspecified ModifierTarget = iota
	ModTargetAc
	ModTargetAttack
	ModTargetSave
	ModTargetCheck
	ModTargetDamage
	ModTargetSpeed
)

// ModifierSource is a game-declared source name, used to dedupe (same source
// does not stack) and to key removal, e.g. const SourceBless ModifierSource = "bless".
type ModifierSource string

// ModifierSpec is an ongoing numeric buff or debuff (Bless +1d4 to attacks and
// saves, Shield of Faith +2 AC, Hunter's Mark +1d6 damage). It is declarative:
// the game stores active modifiers and folds them via ModifierBonus. Vantage is
// intentionally not here; advantage flows through conditions.
type ModifierSpec struct {
	Source   ModifierSource
	Targets  []ModifierTarget
	Ability  core.Ability // core.AbilityAny = all; a specific ability narrows save/check targets
	Flat     int
	Dice     dice.Expr // rolled per use; zero value = no dice
	Duration core.EffectDuration
}

// Effect is a bundle of consequences. Each arm is optional (nil pointer = absent).
type Effect struct {
	Damage     *damage.Spec
	HalfOnSave bool // full damage on a failed save, half on a success
	Healing    *dice.Expr
	Condition  *ConditionSpec
	Movement   *MovementSpec
	Modifier   *ModifierSpec
}

// ConditionalEffect pairs an Effect with the trigger that fires it.
type ConditionalEffect struct {
	Trigger Trigger
	Effect  Effect
}
