package core

import "testing"

func TestAbilityModifier(t *testing.T) {
	cases := []struct {
		score AbilityScore
		want  Modifier
	}{
		{1, -5}, {6, -2}, {7, -2}, {8, -1}, {9, -1}, {10, 0},
		{11, 0}, {12, 1}, {15, 2}, {20, 5}, {30, 10},
	}
	for _, c := range cases {
		if got := AbilityModifier(c.score); got != c.want {
			t.Errorf("AbilityModifier(%d) = %d, want %d", c.score, got, c.want)
		}
		if got := c.score.Modifier(); got != c.want {
			t.Errorf("AbilityScore(%d).Modifier() = %d, want %d", c.score, got, c.want)
		}
	}
}
