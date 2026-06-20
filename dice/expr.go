package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Expr is a dice expression NdS+M: Count dice of Sides faces, plus Modifier.
type Expr struct {
	Count, Sides, Modifier int
}

var diceRe = regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)

// Parse reads a dice expression such as "2d6+3", "1d20", "d8" (count defaults to
// 1), or "2d6-1".
func Parse(s string) (Expr, error) {
	m := diceRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Expr{}, fmt.Errorf("parsing dice %q: invalid format", s)
	}
	count := 1
	if m[1] != "" {
		count, _ = strconv.Atoi(m[1])
	}
	sides, _ := strconv.Atoi(m[2])
	if sides < 1 {
		return Expr{}, fmt.Errorf("parsing dice %q: sides must be at least 1", s)
	}
	mod := 0
	if m[3] != "" {
		mod, _ = strconv.Atoi(m[3])
	}
	return Expr{Count: count, Sides: sides, Modifier: mod}, nil
}

// Min is the lowest possible total (every die shows 1, plus the modifier).
func (e Expr) Min() int { return e.Count + e.Modifier }

// Max is the highest possible total (every die shows its max, plus the modifier).
func (e Expr) Max() int { return e.Count*e.Sides + e.Modifier }

// Average is the expected total.
func (e Expr) Average() float64 {
	return float64(e.Count)*(float64(e.Sides)+1)/2 + float64(e.Modifier)
}

// Roll rolls each die and sums, adding the modifier once. Result.Dice holds the
// individual faces in roll order.
func (e Expr) Roll(r Roller) Result {
	dice := make([]int, e.Count)
	sum := 0
	for i := range dice {
		face := r.IntN(e.Sides) + 1
		dice[i] = face
		sum += face
	}
	return Result{Dice: dice, Total: sum + e.Modifier}
}

// RollCritical rolls a critical hit: the number of dice is doubled, the modifier
// is added once (SRD rule: double the dice, not the modifier).
func (e Expr) RollCritical(r Roller) Result {
	res := Expr{Count: e.Count * 2, Sides: e.Sides}.Roll(r)
	res.Total += e.Modifier
	return res
}
