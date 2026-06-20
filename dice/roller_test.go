package dice

import "testing"

func TestConstant(t *testing.T) {
	// Constant(face).IntN(n) yields min(face,n)-1, so a die of `sides` shows min(face,sides).
	cases := []struct {
		face, sides, wantIntN int
	}{
		{10, 20, 9},  // d20 shows 10
		{20, 20, 19}, // d20 shows 20
		{20, 6, 5},   // d6 clamps to 6
		{10, 8, 7},   // d8 clamps to 8
		{3, 6, 2},    // d6 shows 3
	}
	for _, c := range cases {
		if got := Constant(c.face).IntN(c.sides); got != c.wantIntN {
			t.Errorf("Constant(%d).IntN(%d) = %d, want %d", c.face, c.sides, got, c.wantIntN)
		}
	}
}

func TestTakeRollers(t *testing.T) {
	if got := Take10.IntN(20); got != 9 {
		t.Errorf("Take10.IntN(20) = %d, want 9 (face 10)", got)
	}
	if got := Take20.IntN(20); got != 19 {
		t.Errorf("Take20.IntN(20) = %d, want 19 (face 20)", got)
	}
}

func TestNewRollerInRange(t *testing.T) {
	r := NewRoller(42)
	for i := 0; i < 1000; i++ {
		v := r.IntN(6)
		if v < 0 || v >= 6 {
			t.Fatalf("NewRoller IntN(6) = %d, out of [0,6)", v)
		}
	}
	// Deterministic for a given seed.
	a, b := NewRoller(7), NewRoller(7)
	for i := 0; i < 10; i++ {
		if a.IntN(100) != b.IntN(100) {
			t.Fatal("NewRoller not deterministic for equal seeds")
		}
	}
}
