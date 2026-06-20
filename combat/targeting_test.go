package combat

import (
	"testing"

	"github.com/trancecode/go-srd5e/dice"
)

func TestCover(t *testing.T) {
	if CoverNone.AcBonus() != 0 || CoverHalf.AcBonus() != 2 || CoverThreeQuarters.AcBonus() != 5 {
		t.Error("cover AC bonuses wrong")
	}
	if CoverHalf.DexSaveBonus() != 2 {
		t.Error("cover dex save bonus should match AC bonus")
	}
	if !CoverTotal.BlocksTargeting() || CoverHalf.BlocksTargeting() {
		t.Error("only total cover blocks targeting")
	}
}

func TestRangeBand(t *testing.T) {
	r := Range{Normal: 80, Long: 320}
	if r.Band(40) != BandNormal || r.Band(200) != BandLong || r.Band(400) != BandOutOf {
		t.Error("ranged bands wrong")
	}
	// melee: reach with no long band.
	m := MeleeRange(5)
	if m.Band(5) != BandNormal || m.Band(10) != BandOutOf {
		t.Error("melee bands wrong")
	}
}

func TestAttackSetup(t *testing.T) {
	// total cover -> not possible.
	if (Attack{TargetCover: CoverTotal, Range: MeleeRange(5)}).Setup(15).Possible {
		t.Error("total cover should block")
	}
	// out of range -> not possible.
	if (Attack{Range: Range{Normal: 80, Long: 320}, Distance: 400}).Setup(15).Possible {
		t.Error("beyond long range should block")
	}
	// long range -> disadvantage; half cover -> +2 AC.
	s := Attack{Range: Range{80, 320}, Distance: 200, TargetCover: CoverHalf}.Setup(15)
	if !s.Possible || s.Vantage != dice.VantageDisadvantage || s.EffectiveAc != 17 {
		t.Errorf("long+halfcover = %+v, want Possible Vantage Disadvantage EffectiveAc 17", s)
	}
	// long range disadvantage cancels with a supplied advantage.
	s = Attack{Range: Range{80, 320}, Distance: 200, Advantage: true}.Setup(15)
	if s.Vantage != dice.VantageNone {
		t.Errorf("advantage + long-range disadvantage = %v, want None", s.Vantage)
	}
}
