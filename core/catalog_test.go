package core

import "testing"

func TestSRDSkills(t *testing.T) {
	if len(SRDSkills) != 18 {
		t.Fatalf("SRDSkills has %d entries, want 18", len(SRDSkills))
	}
	if SkillStealth.Ability != AbilityDexterity {
		t.Errorf("Stealth ability = %v, want Dexterity", SkillStealth.Ability)
	}
	if SkillArcana.Ability != AbilityIntelligence {
		t.Errorf("Arcana ability = %v, want Intelligence", SkillArcana.Ability)
	}
}

func TestDamageAny(t *testing.T) {
	if DamageAny.Id != "any" {
		t.Errorf("DamageAny.Id = %q, want \"any\"", DamageAny.Id)
	}
	if Slashing.Id != "slashing" {
		t.Errorf("Slashing.Id = %q, want \"slashing\"", Slashing.Id)
	}
}
