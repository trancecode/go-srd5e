package combat

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

func TestResolveAttack(t *testing.T) {
	// natural 20 -> critical regardless of AC.
	if r := ResolveAttack(20, 0, 99); r.Outcome != AttackCritical {
		t.Errorf("nat 20 = %v, want Critical", r.Outcome)
	}
	// natural 1 -> miss regardless of bonus.
	if r := ResolveAttack(1, 100, 5); r.Outcome != AttackMiss {
		t.Errorf("nat 1 = %v, want Miss", r.Outcome)
	}
	// 15 + 4 = 19 vs AC 18 -> hit; Total recorded.
	r := ResolveAttack(15, 4, 18)
	if r.Outcome != AttackHit || r.Total != 19 || r.NaturalRoll != 15 {
		t.Errorf("hit = %+v, want Hit Total 19 NaturalRoll 15", r)
	}
	// 10 + 2 = 12 vs AC 18 -> miss.
	if r := ResolveAttack(10, 2, 18); r.Outcome != AttackMiss {
		t.Errorf("low total = %v, want Miss", r.Outcome)
	}
}

func TestConcentrationDc(t *testing.T) {
	cases := []struct {
		dmg  int
		want core.Dc
	}{{8, 10}, {19, 10}, {22, 11}, {30, 15}}
	for _, c := range cases {
		if got := ConcentrationDc(c.dmg); got != c.want {
			t.Errorf("ConcentrationDc(%d) = %d, want %d", c.dmg, got, c.want)
		}
	}
}
