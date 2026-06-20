package content

import "github.com/trancecode/go-srd5e/core"

// SpellSlots is the slots available per spell level; index is the spell level
// (1..9), index 0 unused.
type SpellSlots []int

// SpellSlotProgression maps a class level (1..20) to its slots. An absent level
// means no casting yet, keeping late and partial casters free of filler zeros.
type SpellSlotProgression map[int]SpellSlots

// Class is a character class shape. SRD classes have two proficient saves; the
// slice is left open for custom classes. A zero SpellcastingAbility (AbilityNone)
// marks a non-caster.
type Class struct {
	Id, Name            string
	HitDie              int
	ProficientSaves     []core.Ability
	SkillChoiceCount    int
	AvailableSkills     []core.Skill
	SpellcastingAbility core.Ability
	Slots               SpellSlotProgression
}

// Race is a character race shape. AbilityBonuses is optional: empty models SRD
// 5.2 and reskinned-race settings where bonuses come from elsewhere.
type Race struct {
	Id, Name       string
	AbilityBonuses map[core.Ability]int
	MovementSpeed  core.Distance
	Traits         []string
}
