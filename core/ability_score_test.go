package core

import "testing"

func TestAbilityScoresScore(t *testing.T) {
	s := AbilityScores{Strength: 15, Dexterity: 12, Constitution: 14, Intelligence: 8, Wisdom: 10, Charisma: 13}
	cases := []struct {
		a    Ability
		want AbilityScore
	}{
		{AbilityStrength, 15}, {AbilityDexterity, 12}, {AbilityConstitution, 14},
		{AbilityIntelligence, 8}, {AbilityWisdom, 10}, {AbilityCharisma, 13},
	}
	for _, c := range cases {
		if got := s.Score(c.a); got != c.want {
			t.Errorf("Score(%v) = %d, want %d", c.a, got, c.want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("Score(AbilityNone) should panic")
		}
	}()
	s.Score(AbilityNone)
}
