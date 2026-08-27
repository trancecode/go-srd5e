package core

import "testing"

func TestCoinsConstructorsAndString(t *testing.T) {
	if Gp(15) != 1500 || Sp(5) != 50 || Cp(7) != 7 {
		t.Errorf("Gp/Sp/Cp: %d %d %d", Gp(15), Sp(5), Cp(7))
	}
	cases := map[Coins]string{Gp(15): "15 gp", Sp(5): "5 sp", Cp(7): "7 cp", Gp(1) + Sp(5): "1 gp 5 sp", 0: "0 cp"}
	for c, want := range cases {
		if got := c.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", c, got, want)
		}
	}
}
