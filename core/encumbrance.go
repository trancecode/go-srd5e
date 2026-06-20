package core

// Size is the closed set of creature sizes.
type Size int

const (
	SizeUnspecified Size = iota
	SizeTiny
	SizeSmall
	SizeMedium
	SizeLarge
	SizeHuge
	SizeGargantuan
)

// sizeFactor gives the carrying-capacity multiplier as a fraction (num/den):
// Tiny halves, Small and Medium x1, Large x2, Huge x4, Gargantuan x8. It panics
// on SizeUnspecified because Size is a must-set sentinel.
func sizeFactor(size Size) (num, den int) {
	switch size {
	case SizeTiny:
		return 1, 2
	case SizeSmall, SizeMedium:
		return 1, 1
	case SizeLarge:
		return 2, 1
	case SizeHuge:
		return 4, 1
	case SizeGargantuan:
		return 8, 1
	default: // SizeUnspecified or out of range
		panic("core: carrying capacity requires a Size")
	}
}

// CarryingCapacity is Strength x 15, scaled by size.
func CarryingCapacity(str AbilityScore, size Size) Weight {
	num, den := sizeFactor(size)
	return Weight(int(str) * 15 * num / den)
}

// PushDragLift is Strength x 30, scaled by size: the maximum to push, drag, or lift.
func PushDragLift(str AbilityScore, size Size) Weight {
	num, den := sizeFactor(size)
	return Weight(int(str) * 30 * num / den)
}

// Encumbrance is the optional SRD encumbrance variant tier.
type Encumbrance int

const (
	Unencumbered Encumbrance = iota
	Encumbered
	HeavilyEncumbered
)

// EncumbranceTier reports the optional variant: Encumbered above Strength x 5,
// HeavilyEncumbered above Strength x 10. Games that do not use the variant never
// call it.
func EncumbranceTier(str AbilityScore, carried Weight) Encumbrance {
	switch {
	case int(carried) > int(str)*10:
		return HeavilyEncumbered
	case int(carried) > int(str)*5:
		return Encumbered
	default:
		return Unencumbered
	}
}
