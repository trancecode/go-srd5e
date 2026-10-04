package damage

import (
	"reflect"
	"testing"

	"github.com/trancecode/go-srd5e/combat"
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestRollHit(t *testing.T) {
	// 1d6 slashing + bonus 3, hit, Constant(6) -> one part amount 9.
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackHit, dice.Constant(6))
	if len(d.Parts) != 1 || d.Parts[0].Amount != 9 || d.Parts[0].Type != core.Slashing {
		t.Errorf("hit = %+v, want one slashing part amount 9", d)
	}
}

func TestRollCriticalDoublesDiceNotBonus(t *testing.T) {
	// 1d6 slashing + bonus 3, crit, Constant(6) -> 2d6 = 12, + 3 = 15.
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackCritical, dice.Constant(6))
	if d.Parts[0].Amount != 15 {
		t.Errorf("crit amount = %d, want 15 (doubled dice, bonus once)", d.Parts[0].Amount)
	}
}

func TestRollMiss(t *testing.T) {
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackMiss, dice.Constant(6))
	if len(d.Parts) != 0 {
		t.Errorf("miss = %+v, want no parts", d)
	}
}

func TestRollBonusOnPrimaryPartOnly(t *testing.T) {
	// Two parts: 1d8 slashing + 1d6 fire, bonus 2, hit, Constant(8)/Constant(6) won't
	// vary per part with one Constant, so use Constant(10): slashing die shows 8, fire die shows 6.
	spec := Spec{Parts: []PartSpec{
		{Dice: dice.Expr{Count: 1, Sides: 8}, Type: core.Slashing},
		{Dice: dice.Expr{Count: 1, Sides: 6}, Type: core.Fire},
	}}
	d := Roll(spec, 2, combat.AttackHit, dice.Constant(10))
	if d.Parts[0].Amount != 10 { // 8 + bonus 2
		t.Errorf("primary part = %d, want 10 (8 + bonus 2)", d.Parts[0].Amount)
	}
	if d.Parts[1].Amount != 6 { // fire, no bonus
		t.Errorf("secondary part = %d, want 6 (no bonus)", d.Parts[1].Amount)
	}
}

// seqRoller shows the given faces in order, one per die, so a test can tell
// which face landed in which part.
type seqRoller struct {
	faces []int
	next  int
}

func (r *seqRoller) IntN(n int) int {
	face := r.faces[r.next]
	r.next++
	if face < 1 || face > n {
		panic("seqRoller: face out of range for die")
	}
	return face - 1
}

func TestRollRecordsDiceOnHit(t *testing.T) {
	r := &seqRoller{faces: []int{5}}
	d := RollSingle(dice.Expr{Count: 1, Sides: 8}, core.Slashing, 3, combat.AttackHit, r)
	if p := d.Parts[0]; p.Amount != 8 || !reflect.DeepEqual(p.Dice, []int{5}) {
		t.Errorf("hit part = %+v, want amount 8 with dice [5]", p)
	}
}

func TestRollRecordsDoubledDiceOnCritical(t *testing.T) {
	// 1d8+1 slashing, bonus 3, crit: two d8s showing 8 and 7, expression
	// modifier and bonus each added once.
	r := &seqRoller{faces: []int{8, 7}}
	d := RollSingle(dice.Expr{Count: 1, Sides: 8, Modifier: 1}, core.Slashing, 3, combat.AttackCritical, r)
	if p := d.Parts[0]; p.Amount != 19 || !reflect.DeepEqual(p.Dice, []int{8, 7}) {
		t.Errorf("crit part = %+v, want amount 19 with dice [8 7]", p)
	}
}

func TestRollRecordsDicePerPart(t *testing.T) {
	spec := Spec{Parts: []PartSpec{
		{Dice: dice.Expr{Count: 2, Sides: 6}, Type: core.Slashing},
		{Dice: dice.Expr{Count: 1, Sides: 4}, Type: core.Fire},
	}}
	cases := []struct {
		name    string
		outcome combat.AttackOutcome
		faces   []int
		want    [][]int
		amounts []int
	}{
		{"hit", combat.AttackHit, []int{3, 6, 2}, [][]int{{3, 6}, {2}}, []int{11, 2}},
		{"critical", combat.AttackCritical, []int{1, 2, 3, 4, 4, 1}, [][]int{{1, 2, 3, 4}, {4, 1}}, []int{12, 5}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := Roll(spec, 2, tc.outcome, &seqRoller{faces: tc.faces})
			if len(d.Parts) != 2 {
				t.Fatalf("parts = %+v, want 2", d.Parts)
			}
			for i, p := range d.Parts {
				if !reflect.DeepEqual(p.Dice, tc.want[i]) || p.Amount != tc.amounts[i] {
					t.Errorf("part %d = %+v, want dice %v amount %d", i, p, tc.want[i], tc.amounts[i])
				}
			}
		})
	}
}

func TestRollMissRollsNoDice(t *testing.T) {
	r := &seqRoller{}
	d := RollSingle(dice.Expr{Count: 1, Sides: 6}, core.Slashing, 3, combat.AttackMiss, r)
	if len(d.Parts) != 0 || r.next != 0 {
		t.Errorf("miss = %+v after %d draws, want no parts and no draws", d, r.next)
	}
}
