package core

import "testing"

func TestAbilityString(t *testing.T) {
	cases := []struct {
		a    Ability
		want string
	}{
		{AbilityStrength, "strength"},
		{AbilityDexterity, "dexterity"},
		{AbilityConstitution, "constitution"},
		{AbilityIntelligence, "intelligence"},
		{AbilityWisdom, "wisdom"},
		{AbilityCharisma, "charisma"},
		{AbilityNone, "none"},
		{AbilityAny, "any"},
	}
	for _, c := range cases {
		if got := c.a.String(); got != c.want {
			t.Errorf("Ability(%d).String() = %q, want %q", c.a, got, c.want)
		}
	}
}

func TestSkillString(t *testing.T) {
	if got := SkillPerception.String(); got != "perception" {
		t.Errorf("SkillPerception.String() = %q, want %q", got, "perception")
	}
	if got := SkillNone.String(); got != "none" {
		t.Errorf("SkillNone.String() = %q, want %q", got, "none")
	}
}

func TestSkillNoneIsZero(t *testing.T) {
	if SkillNone != (Skill{}) {
		t.Errorf("SkillNone = %+v, want the zero Skill", SkillNone)
	}
}

func TestSkillById(t *testing.T) {
	for _, want := range SRDSkills {
		got, ok := SkillById(want.Id)
		if !ok {
			t.Errorf("SkillById(%q) not found", want.Id)
			continue
		}
		if got != want {
			t.Errorf("SkillById(%q) = %+v, want %+v", want.Id, got, want)
		}
	}
	if got, ok := SkillById("nonesuch"); ok || got != SkillNone {
		t.Errorf("SkillById(%q) = %+v, %v; want SkillNone, false", "nonesuch", got, ok)
	}
}

// Every SRD skill round-trips through its own String.
func TestSkillStringRoundTrip(t *testing.T) {
	for _, s := range SRDSkills {
		got, ok := SkillById(s.String())
		if !ok || got != s {
			t.Errorf("round trip of %q failed: got %+v, %v", s.Id, got, ok)
		}
	}
}
