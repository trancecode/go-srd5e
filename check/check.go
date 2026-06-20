package check

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// Check is a configured d20 check: a total bonus against a difficulty class.
type Check struct {
	Bonus core.Modifier
	Dc    core.Dc
}

// Result is the outcome of resolving a Check against a d20 roll.
type Result struct {
	Roll    dice.Result
	Total   int
	Success bool
	Margin  int // Total minus Dc
}

// Resolve interprets an already-made d20 roll against the check.
func (c Check) Resolve(roll dice.Result) Result {
	total := roll.Total + int(c.Bonus)
	margin := total - int(c.Dc)
	return Result{Roll: roll, Total: total, Success: margin >= 0, Margin: margin}
}

// Ability builds an ability check: the ability modifier against the DC.
func Ability(scores core.AbilityScores, a core.Ability, dc core.Dc) Check {
	return Check{Bonus: core.AbilityModifier(scores.Score(a)), Dc: dc}
}

// Skill builds a skill check, reading the governing ability from the skill.
func Skill(scores core.AbilityScores, s core.Skill, proficient bool, level core.Level, dc core.Dc) Check {
	return Check{Bonus: core.SkillBonus(scores.Score(s.Ability), proficient, level), Dc: dc}
}

// Save builds a saving throw.
func Save(scores core.AbilityScores, a core.Ability, proficient bool, level core.Level, dc core.Dc) Check {
	return Check{Bonus: core.SavingThrowBonus(scores.Score(a), proficient, level), Dc: dc}
}

// ContestResult is the outcome of an opposed check. Ties favor the responder.
type ContestResult struct {
	InitiatorTotal, ResponderTotal int
	InitiatorWins                  bool
}

// Contest resolves an opposed check from the two already-computed totals.
func Contest(initiatorTotal, responderTotal int) ContestResult {
	return ContestResult{
		InitiatorTotal: initiatorTotal,
		ResponderTotal: responderTotal,
		InitiatorWins:  initiatorTotal > responderTotal,
	}
}

// PassiveScore is 10 + modifier, +5 for advantage, -5 for disadvantage. Generic
// across passive Perception, Investigation, and Insight.
func PassiveScore(modifier core.Modifier, v dice.Vantage) int {
	p := 10 + int(modifier)
	switch v {
	case dice.VantageAdvantage:
		p += 5
	case dice.VantageDisadvantage:
		p -= 5
	}
	return p
}
