package dice

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want Expr
	}{
		{"2d6+3", Expr{2, 6, 3}},
		{"1d20", Expr{1, 20, 0}},
		{"d8", Expr{1, 8, 0}},
		{"2d6-1", Expr{2, 6, -1}},
		{" 3d4 ", Expr{3, 4, 0}},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
	for _, bad := range []string{"", "abc", "2x6", "d", "2d", "d0", "2d6+"} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", bad)
		}
	}
}

func TestBounds(t *testing.T) {
	e := Expr{2, 6, 3} // 2d6+3
	if e.Min() != 5 {
		t.Errorf("Min = %d, want 5", e.Min())
	}
	if e.Max() != 15 {
		t.Errorf("Max = %d, want 15", e.Max())
	}
	if e.Average() != 10.0 {
		t.Errorf("Average = %v, want 10.0", e.Average())
	}
	d20 := Expr{1, 20, 0}
	if d20.Average() != 10.5 {
		t.Errorf("d20 Average = %v, want 10.5", d20.Average())
	}
}
