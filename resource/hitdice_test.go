package resource

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestHitDicePool(t *testing.T) {
	p := NewHitDicePool(5, 8) // 5d8
	if p.Total() != 5 || p.Available(8) != 5 {
		t.Fatalf("fresh pool wrong: total %d avail %d", p.Total(), p.Available(8))
	}
	if !p.Spend(8) || !p.Spend(8) || p.Available(8) != 3 {
		t.Errorf("after 2 spends available = %d, want 3", p.Available(8))
	}
	p.Recover(1)
	if p.Available(8) != 4 {
		t.Errorf("after recover 1 available = %d, want 4", p.Available(8))
	}
	// spend down to 0 then over-spend.
	p.Spend(8)
	p.Spend(8)
	p.Spend(8)
	p.Spend(8)
	if p.Spend(8) {
		t.Error("over-spend should be false")
	}
}

func TestHitDicePoolRestore(t *testing.T) {
	p := NewHitDicePool(6, 10) // 6d10
	for i := 0; i < 6; i++ {
		p.Spend(10)
	}
	p.Restore(RestLong) // recovers up to half (3)
	if p.Available(10) != 3 {
		t.Errorf("long rest available = %d, want 3", p.Available(10))
	}
	// short rest does not recover dice.
	p2 := NewHitDicePool(6, 10)
	p2.Spend(10)
	p2.Restore(RestShort)
	if p2.Available(10) != 5 {
		t.Errorf("short rest available = %d, want 5 (unchanged)", p2.Available(10))
	}
}

func TestHitDieHeal(t *testing.T) {
	// d8 shows 8, Con 14 (+2): heal 10.
	sc := core.AbilityScores{Constitution: 14}
	if got := HitDieHeal(8, sc, dice.Constant(8)); got != 10 {
		t.Errorf("HitDieHeal = %d, want 10", got)
	}
}

var _ Restorable = (*HitDicePool)(nil)
