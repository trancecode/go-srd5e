package core

// Named units. Distinct SRD quantities are named types so the compiler rejects
// argument swaps; untyped constant literals still convert freely.
type (
	AbilityScore int // 1..30
	Modifier     int // a value added to a d20
	Level        int // 1..20
	ArmorClass   int
	Dc           int // difficulty class
	Distance     int // feet
	HitPoints    int
	Xp           int
	Weight       float64 // pounds; the SRD lists fractional weights (e.g. a dart at 1/4 lb.)
)

// Ability is the closed set of the six SRD abilities. AbilityAny is a wildcard
// for qualifiers only (e.g. a buff to all saves); never use it where a concrete
// ability is required.
type Ability int

const (
	AbilityNone Ability = iota
	AbilityStrength
	AbilityDexterity
	AbilityConstitution
	AbilityIntelligence
	AbilityWisdom
	AbilityCharisma
)

const AbilityAny Ability = -1

// AbilityScores is the six ability scores.
type AbilityScores struct {
	Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma AbilityScore
}

// AbilityModifier is (score-10)/2 rounded down.
func AbilityModifier(score AbilityScore) Modifier {
	d := int(score) - 10
	if d >= 0 {
		return Modifier(d / 2)
	}
	return Modifier(-((-d + 1) / 2))
}

// Modifier returns the score's ability modifier.
func (s AbilityScore) Modifier() Modifier { return AbilityModifier(s) }

// Score returns the score for a concrete ability. It panics on AbilityNone or
// AbilityAny, which are not concrete abilities.
func (s AbilityScores) Score(a Ability) AbilityScore {
	switch a {
	case AbilityStrength:
		return s.Strength
	case AbilityDexterity:
		return s.Dexterity
	case AbilityConstitution:
		return s.Constitution
	case AbilityIntelligence:
		return s.Intelligence
	case AbilityWisdom:
		return s.Wisdom
	case AbilityCharisma:
		return s.Charisma
	default:
		panic("core: Score requires a concrete ability")
	}
}

// SetScore sets the score for a concrete ability. It panics on AbilityNone or
// AbilityAny, which are not concrete abilities.
func (s *AbilityScores) SetScore(a Ability, v AbilityScore) {
	switch a {
	case AbilityStrength:
		s.Strength = v
	case AbilityDexterity:
		s.Dexterity = v
	case AbilityConstitution:
		s.Constitution = v
	case AbilityIntelligence:
		s.Intelligence = v
	case AbilityWisdom:
		s.Wisdom = v
	case AbilityCharisma:
		s.Charisma = v
	default:
		panic("core: SetScore requires a concrete ability")
	}
}

// String returns the lowercase identifier of the ability, or "none" for
// AbilityNone and "any" for AbilityAny.
func (a Ability) String() string {
	switch a {
	case AbilityStrength:
		return "strength"
	case AbilityDexterity:
		return "dexterity"
	case AbilityConstitution:
		return "constitution"
	case AbilityIntelligence:
		return "intelligence"
	case AbilityWisdom:
		return "wisdom"
	case AbilityCharisma:
		return "charisma"
	case AbilityAny:
		return "any"
	default:
		return "none"
	}
}
