package dice

import (
	"reflect"
	"testing"
)

// seqRoller returns predetermined IntN values in order (a die face is value+1).
type seqRoller struct {
	vals []int
	i    int
}

func (s *seqRoller) IntN(int) int {
	v := s.vals[s.i]
	s.i++
	return v
}

func TestRoll(t *testing.T) {
	// 2d6+1 with a roller yielding IntN 2 then 4 -> faces 3 and 5 -> total 9.
	r := &seqRoller{vals: []int{2, 4}}
	got := Expr{2, 6, 1}.Roll(r)
	if !reflect.DeepEqual(got.Dice, []int{3, 5}) {
		t.Errorf("Dice = %v, want [3 5]", got.Dice)
	}
	if got.Total != 9 {
		t.Errorf("Total = %d, want 9", got.Total)
	}
}

func TestRollCritical(t *testing.T) {
	// 1d8+5 crit: dice doubled to 2d8, modifier added once. Constant(8) -> both 8.
	got := Expr{1, 8, 5}.RollCritical(Constant(8))
	if !reflect.DeepEqual(got.Dice, []int{8, 8}) {
		t.Errorf("crit Dice = %v, want [8 8]", got.Dice)
	}
	if got.Total != 21 { // 8 + 8 + 5
		t.Errorf("crit Total = %d, want 21", got.Total)
	}
}
