package dice

import "math/rand/v2"

// Roller is the source of randomness. *math/rand/v2.Rand satisfies it directly.
// It is always passed per call, never stored, so value types stay serializable.
type Roller interface {
	IntN(n int) int
}

// Result is the outcome of a roll: the individual die faces in roll order, and
// the total including any modifier.
type Result struct {
	Dice  []int
	Total int
}

// constRoller always reports the same face, clamped to each die's range.
type constRoller int

// IntN returns min(face, n) - 1, so a die of n sides shows min(face, n).
func (c constRoller) IntN(n int) int {
	face := int(c)
	if face > n {
		face = n
	}
	if face < 1 {
		face = 1
	}
	return face - 1
}

// Constant returns a Roller on which every die shows min(face, sides). Useful
// for assist modes (take 10, take 20), best-case analysis, and tests.
func Constant(face int) Roller { return constRoller(face) }

// Take10 and Take20 are fixed rollers: a d20 shows 10 or 20, smaller dice clamp
// to their maximum. Not SRD rules (5e uses passive checks); for house rules,
// difficulty modes, and tests.
var (
	Take10 Roller = Constant(10)
	Take20 Roller = Constant(20)
)

// NewRoller returns a fair, deterministically seeded roller backed by
// math/rand/v2 (a PCG source).
func NewRoller(seed uint64) Roller {
	return rand.New(rand.NewPCG(seed, seed))
}
