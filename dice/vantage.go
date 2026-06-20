package dice

// Vantage is the advantage state of a d20 roll. Advantage rolls two d20 and
// keeps the higher; disadvantage keeps the lower.
type Vantage int

const (
	VantageNone Vantage = iota
	VantageAdvantage
	VantageDisadvantage
)

// RollD20 rolls a d20 under the given vantage. With advantage or disadvantage it
// rolls twice; Result.Dice holds both faces and Total is the kept one. There is
// no modifier on a raw d20, so Total is the kept die.
func RollD20(r Roller, v Vantage) Result {
	a := r.IntN(20) + 1
	switch v {
	case VantageAdvantage, VantageDisadvantage:
		b := r.IntN(20) + 1
		keep := a
		if (v == VantageAdvantage && b > keep) || (v == VantageDisadvantage && b < keep) {
			keep = b
		}
		return Result{Dice: []int{a, b}, Total: keep}
	default:
		return Result{Dice: []int{a}, Total: a}
	}
}

// CombineVantage applies the SRD rule: any advantage source plus any
// disadvantage source cancel to a straight roll, regardless of count.
func CombineVantage(advantage, disadvantage bool) Vantage {
	switch {
	case advantage == disadvantage:
		return VantageNone
	case advantage:
		return VantageAdvantage
	default:
		return VantageDisadvantage
	}
}
