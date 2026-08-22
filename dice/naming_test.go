package dice

import "testing"

func TestVantageString(t *testing.T) {
	cases := []struct {
		v    Vantage
		want string
	}{
		{VantageAdvantage, "advantage"},
		{VantageDisadvantage, "disadvantage"},
		{VantageNone, "none"},
	}
	for _, c := range cases {
		if got := c.v.String(); got != c.want {
			t.Errorf("Vantage(%d).String() = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestExprString(t *testing.T) {
	cases := []struct {
		e    Expr
		want string
	}{
		{Expr{Count: 2, Sides: 6, Modifier: 3}, "2d6+3"},
		{Expr{Count: 1, Sides: 20}, "1d20"},
		{Expr{Count: 1, Sides: 8, Modifier: -1}, "1d8-1"},
	}
	for _, c := range cases {
		if got := c.e.String(); got != c.want {
			t.Errorf("%+v.String() = %q, want %q", c.e, got, c.want)
		}
	}
}

// A parsed expression prints back to an equivalent string.
func TestExprStringRoundTrip(t *testing.T) {
	for _, s := range []string{"2d6+3", "1d20", "1d8-1", "3d4"} {
		e, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", s, err)
			continue
		}
		if got := e.String(); got != s {
			t.Errorf("Parse(%q).String() = %q, want %q", s, got, s)
		}
	}
}
