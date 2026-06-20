package resource

import (
	"sort"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// HitDicePool is hit dice keyed by die size (mixed sizes for a multiclass
// character): Dice is the total per size, Spent is how many of each are used.
type HitDicePool struct {
	Dice  map[int]int
	Spent map[int]int
}

// NewHitDicePool builds a single-class pool: level dice of size hitDie.
func NewHitDicePool(level, hitDie int) HitDicePool {
	return HitDicePool{
		Dice:  map[int]int{hitDie: level},
		Spent: map[int]int{},
	}
}

// Available is the unspent dice of a given size.
func (p *HitDicePool) Available(die int) int { return p.Dice[die] - p.Spent[die] }

// Spend marks one die of the given size used (the player heals on a short rest);
// false if none are available.
func (p *HitDicePool) Spend(die int) bool {
	if p.Available(die) <= 0 {
		return false
	}
	if p.Spent == nil {
		p.Spent = map[int]int{}
	}
	p.Spent[die]++
	return true
}

// Recover regains up to n spent dice, largest size first (deterministic).
func (p *HitDicePool) Recover(n int) {
	sizes := make([]int, 0, len(p.Spent))
	for d := range p.Spent {
		sizes = append(sizes, d)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	for _, d := range sizes {
		for n > 0 && p.Spent[d] > 0 {
			p.Spent[d]--
			n--
		}
		if n == 0 {
			return
		}
	}
}

// Total is the pool size summed across die sizes.
func (p HitDicePool) Total() int {
	t := 0
	for _, c := range p.Dice {
		t += c
	}
	return t
}

// Restore recovers up to half the pool (round down, minimum 1) on a long rest; a
// short rest is a no-op.
func (p *HitDicePool) Restore(rest RestKind) {
	if rest != RestLong {
		return
	}
	half := p.Total() / 2
	if half < 1 {
		half = 1
	}
	p.Recover(half)
}

// HitDieHeal rolls one hit die plus the Constitution modifier, read from scores
// so the wrong modifier cannot be passed. Floors at zero.
func HitDieHeal(die int, scores core.AbilityScores, r dice.Roller) core.HitPoints {
	rolled := dice.Expr{Count: 1, Sides: die}.Roll(r).Total
	healed := rolled + int(core.AbilityModifier(scores.Constitution))
	if healed < 0 {
		healed = 0
	}
	return core.HitPoints(healed)
}
