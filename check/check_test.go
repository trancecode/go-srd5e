package check

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestResolve(t *testing.T) {
	c := Check{Bonus: 5, Dc: 15}
	// rolled a 12 -> total 17 -> success, margin 2.
	got := c.Resolve(dice.Result{Dice: []int{12}, Total: 12})
	if !got.Success || got.Total != 17 || got.Margin != 2 {
		t.Errorf("success case = %+v, want Total 17 Success true Margin 2", got)
	}
	// rolled a 9 -> total 14 -> fail, margin -1.
	got = c.Resolve(dice.Result{Dice: []int{9}, Total: 9})
	if got.Success || got.Total != 14 || got.Margin != -1 {
		t.Errorf("fail case = %+v, want Total 14 Success false Margin -1", got)
	}
}

func TestConstructors(t *testing.T) {
	sc := core.AbilityScores{Strength: 16, Dexterity: 16, Intelligence: 16}
	// STR 16 (+3), level 5 (prof +3).
	if Ability(sc, core.AbilityStrength, 12).Bonus != 3 {
		t.Errorf("Ability bonus = %d, want 3", Ability(sc, core.AbilityStrength, 12).Bonus)
	}
	if Skill(sc, core.SkillAthletics, true, 5, 12).Bonus != 6 {
		t.Errorf("Skill proficient bonus = %d, want 6", Skill(sc, core.SkillAthletics, true, 5, 12).Bonus)
	}
	if Save(sc, core.AbilityDexterity, false, 5, 12).Bonus != 3 {
		t.Errorf("Save non-proficient bonus = %d, want 3", Save(sc, core.AbilityDexterity, false, 5, 12).Bonus)
	}
}

func TestContest(t *testing.T) {
	if !Contest(18, 12).InitiatorWins {
		t.Error("18 vs 12 should be initiator win")
	}
	if Contest(12, 18).InitiatorWins {
		t.Error("12 vs 18 should be responder win")
	}
	if Contest(15, 15).InitiatorWins {
		t.Error("tie should favor responder (InitiatorWins false)")
	}
}

func TestPassiveScore(t *testing.T) {
	if PassiveScore(3, dice.VantageNone) != 13 {
		t.Errorf("passive none = %d, want 13", PassiveScore(3, dice.VantageNone))
	}
	if PassiveScore(3, dice.VantageAdvantage) != 18 {
		t.Errorf("passive advantage = %d, want 18", PassiveScore(3, dice.VantageAdvantage))
	}
	if PassiveScore(3, dice.VantageDisadvantage) != 8 {
		t.Errorf("passive disadvantage = %d, want 8", PassiveScore(3, dice.VantageDisadvantage))
	}
}
