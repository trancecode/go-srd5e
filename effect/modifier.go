package effect

import (
	"slices"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// ModifierBonus folds the active modifiers that apply to one roll, deduped by
// Source: each Source contributes at most once (same-source entries collapse;
// the highest Flat wins on conflict), distinct sources sum. Returns the total
// flat bonus and the per-source dice to roll (the game rolls them and adds the
// result, plus flat, to the roll's Bonus). Output order is first-seen source
// order, so the result is deterministic.
func ModifierBonus(active []ModifierSpec, target ModifierTarget, ability core.Ability) (flat int, dice []dice.Expr) {
	chosen := map[ModifierSource]ModifierSpec{}
	var order []ModifierSource
	for _, m := range active {
		if !appliesTo(m, target, ability) {
			continue
		}
		if prev, ok := chosen[m.Source]; ok {
			if m.Flat > prev.Flat {
				chosen[m.Source] = m
			}
			continue
		}
		chosen[m.Source] = m
		order = append(order, m.Source)
	}
	for _, src := range order {
		m := chosen[src]
		flat += m.Flat
		if m.Dice != (diceExprZero) {
			dice = append(dice, m.Dice)
		}
	}
	return flat, dice
}

var diceExprZero = dice.Expr{}

func appliesTo(m ModifierSpec, target ModifierTarget, ability core.Ability) bool {
	if !slices.Contains(m.Targets, target) {
		return false
	}
	return m.Ability == core.AbilityAny || m.Ability == ability
}
