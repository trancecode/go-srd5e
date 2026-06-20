package combat

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// Cover is the degree of cover a target has.
type Cover int

const (
	CoverNone Cover = iota
	CoverHalf
	CoverThreeQuarters
	CoverTotal
)

// AcBonus is the bonus cover grants to armor class (+0/+2/+5).
func (c Cover) AcBonus() core.Modifier {
	switch c {
	case CoverHalf:
		return 2
	case CoverThreeQuarters:
		return 5
	default:
		return 0
	}
}

// DexSaveBonus is the bonus cover grants to Dexterity saving throws (same as AC).
func (c Cover) DexSaveBonus() core.Modifier { return c.AcBonus() }

// BlocksTargeting reports whether the target cannot be targeted directly.
func (c Cover) BlocksTargeting() bool { return c == CoverTotal }

// Range is a weapon's normal and long range, in feet.
type Range struct {
	Normal, Long core.Distance
}

// Band classifies a distance relative to a range.
type Band int

const (
	BandNormal Band = iota
	BandLong
	BandOutOf
)

// MeleeRange builds a Range for a melee weapon: reach with no long band, so any
// distance beyond reach is out of range.
func MeleeRange(reach core.Distance) Range { return Range{Normal: reach, Long: reach} }

// Band reports the band a distance falls in: within normal, in the long-range
// (disadvantage) window, or out of range.
func (r Range) Band(distance core.Distance) Band {
	switch {
	case distance <= r.Normal:
		return BandNormal
	case distance <= r.Long:
		return BandLong
	default:
		return BandOutOf
	}
}

// Attack bundles the pre-roll facts that determine whether an attack is possible
// and at what vantage. The game supplies Distance, cover, and any pre-aggregated
// advantage/disadvantage; long-range disadvantage is added by Setup.
type Attack struct {
	Range        Range
	Distance     core.Distance
	TargetCover  Cover
	Advantage    bool
	Disadvantage bool
}

// AttackSetup is the pre-roll result: whether the attack is possible, the vantage
// to roll the d20 with, and the cover-adjusted armor class.
type AttackSetup struct {
	Possible    bool
	Vantage     dice.Vantage
	EffectiveAc core.ArmorClass
}

// Setup computes the pre-roll setup against the target's full armor class. Total
// cover or a distance beyond long range makes the attack impossible.
func (a Attack) Setup(baseAc core.ArmorClass) AttackSetup {
	if a.TargetCover.BlocksTargeting() {
		return AttackSetup{Possible: false}
	}
	band := a.Range.Band(a.Distance)
	if band == BandOutOf {
		return AttackSetup{Possible: false}
	}
	disadvantage := a.Disadvantage || band == BandLong
	return AttackSetup{
		Possible:    true,
		Vantage:     dice.CombineVantage(a.Advantage, disadvantage),
		EffectiveAc: baseAc + core.ArmorClass(a.TargetCover.AcBonus()),
	}
}
