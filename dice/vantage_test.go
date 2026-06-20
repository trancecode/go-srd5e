package dice

import (
	"reflect"
	"testing"
)

func TestRollD20(t *testing.T) {
	// none: single die.
	got := RollD20(&seqRoller{vals: []int{14}}, VantageNone) // face 15
	if !reflect.DeepEqual(got.Dice, []int{15}) || got.Total != 15 {
		t.Errorf("none = %+v, want Dice [15] Total 15", got)
	}
	// advantage: roll 15 and 10, keep 15; both dice recorded.
	got = RollD20(&seqRoller{vals: []int{14, 9}}, VantageAdvantage)
	if !reflect.DeepEqual(got.Dice, []int{15, 10}) || got.Total != 15 {
		t.Errorf("advantage = %+v, want Dice [15 10] Total 15", got)
	}
	// disadvantage: roll 15 and 10, keep 10.
	got = RollD20(&seqRoller{vals: []int{14, 9}}, VantageDisadvantage)
	if !reflect.DeepEqual(got.Dice, []int{15, 10}) || got.Total != 10 {
		t.Errorf("disadvantage = %+v, want Dice [15 10] Total 10", got)
	}
}

func TestCombineVantage(t *testing.T) {
	cases := []struct {
		adv, dis bool
		want     Vantage
	}{
		{false, false, VantageNone},
		{true, false, VantageAdvantage},
		{false, true, VantageDisadvantage},
		{true, true, VantageNone}, // any advantage + any disadvantage cancel
	}
	for _, c := range cases {
		if got := CombineVantage(c.adv, c.dis); got != c.want {
			t.Errorf("CombineVantage(%v,%v) = %v, want %v", c.adv, c.dis, got, c.want)
		}
	}
}
