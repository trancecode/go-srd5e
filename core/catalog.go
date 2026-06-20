package core

// Skill is a named proficiency governed by an ability. It is an open value:
// the eighteen SRD skills ship as predefined values and a game declares more.
type Skill struct {
	Id, Name string
	Ability  Ability
}

// Condition is an open named state (prone, poisoned, plus game-defined ones).
type Condition struct{ Id, Name string }

// DamageType is an open damage category. Games add their own (e.g. thermal).
type DamageType struct{ Id, Name string }

// The eighteen SRD skills.
var (
	SkillAcrobatics     = Skill{"acrobatics", "Acrobatics", AbilityDexterity}
	SkillAnimalHandling = Skill{"animal_handling", "Animal Handling", AbilityWisdom}
	SkillArcana         = Skill{"arcana", "Arcana", AbilityIntelligence}
	SkillAthletics      = Skill{"athletics", "Athletics", AbilityStrength}
	SkillDeception      = Skill{"deception", "Deception", AbilityCharisma}
	SkillHistory        = Skill{"history", "History", AbilityIntelligence}
	SkillInsight        = Skill{"insight", "Insight", AbilityWisdom}
	SkillIntimidation   = Skill{"intimidation", "Intimidation", AbilityCharisma}
	SkillInvestigation  = Skill{"investigation", "Investigation", AbilityIntelligence}
	SkillMedicine       = Skill{"medicine", "Medicine", AbilityWisdom}
	SkillNature         = Skill{"nature", "Nature", AbilityIntelligence}
	SkillPerception     = Skill{"perception", "Perception", AbilityWisdom}
	SkillPerformance    = Skill{"performance", "Performance", AbilityCharisma}
	SkillPersuasion     = Skill{"persuasion", "Persuasion", AbilityCharisma}
	SkillReligion       = Skill{"religion", "Religion", AbilityIntelligence}
	SkillSleightOfHand  = Skill{"sleight_of_hand", "Sleight of Hand", AbilityDexterity}
	SkillStealth        = Skill{"stealth", "Stealth", AbilityDexterity}
	SkillSurvival       = Skill{"survival", "Survival", AbilityWisdom}
)

// SRDSkills is the eighteen SRD skills.
var SRDSkills = []Skill{
	SkillAcrobatics, SkillAnimalHandling, SkillArcana, SkillAthletics,
	SkillDeception, SkillHistory, SkillInsight, SkillIntimidation,
	SkillInvestigation, SkillMedicine, SkillNature, SkillPerception,
	SkillPerformance, SkillPersuasion, SkillReligion, SkillSleightOfHand,
	SkillStealth, SkillSurvival,
}

// SRD physical and common elemental damage types.
var (
	Bludgeoning = DamageType{"bludgeoning", "Bludgeoning"}
	Piercing    = DamageType{"piercing", "Piercing"}
	Slashing    = DamageType{"slashing", "Slashing"}
	Fire        = DamageType{"fire", "Fire"}
	Cold        = DamageType{"cold", "Cold"}
	Lightning   = DamageType{"lightning", "Lightning"}
	Acid        = DamageType{"acid", "Acid"}
	Poison      = DamageType{"poison", "Poison"}
	Necrotic    = DamageType{"necrotic", "Necrotic"}
	Radiant     = DamageType{"radiant", "Radiant"}
	Psychic     = DamageType{"psychic", "Psychic"}
	Thunder     = DamageType{"thunder", "Thunder"}
	Force       = DamageType{"force", "Force"}
)

// DamageAny is a wildcard used only as a key in damage.Mitigation maps. It is
// not a real damage type and must never be put on a damage part.
var DamageAny = DamageType{"any", "Any"}

// A few common SRD conditions; games declare more as open values.
var (
	Prone      = Condition{"prone", "Prone"}
	Poisoned   = Condition{"poisoned", "Poisoned"}
	Restrained = Condition{"restrained", "Restrained"}
	Stunned    = Condition{"stunned", "Stunned"}
	Blinded    = Condition{"blinded", "Blinded"}
	Invisible  = Condition{"invisible", "Invisible"}
)
