package core

import "testing"

func TestProficiencyBonus(t *testing.T) {
	cases := []struct {
		level Level
		want  Modifier
	}{{1, 2}, {4, 2}, {5, 3}, {8, 3}, {9, 4}, {13, 5}, {17, 6}, {20, 6}}
	for _, c := range cases {
		if got := ProficiencyBonus(c.level); got != c.want {
			t.Errorf("ProficiencyBonus(%d) = %d, want %d", c.level, got, c.want)
		}
	}
}

func TestBonuses(t *testing.T) {
	// STR 16 (+3) at level 5 (prof +3): skill proficient = 6, non-proficient = 3.
	if got := SkillBonus(16, true, 5); got != 6 {
		t.Errorf("SkillBonus proficient = %d, want 6", got)
	}
	if got := SkillBonus(16, false, 5); got != 3 {
		t.Errorf("SkillBonus non-proficient = %d, want 3", got)
	}
	if got := SavingThrowBonus(16, true, 5); got != 6 {
		t.Errorf("SavingThrowBonus = %d, want 6", got)
	}
	// INT 16 (+3) at level 5 (prof +3): spell save DC = 8+3+3 = 14, attack = 6.
	if got := SpellSaveDc(16, 5); got != 14 {
		t.Errorf("SpellSaveDc = %d, want 14", got)
	}
	if got := SpellAttackBonus(16, 5); got != 6 {
		t.Errorf("SpellAttackBonus = %d, want 6", got)
	}
}

func TestMaxHp(t *testing.T) {
	// three levels of rolls 6,5,5 with CON +2 = 16 + 6 = 22.
	if got := MaxHp([]int{6, 5, 5}, 2); got != 22 {
		t.Errorf("MaxHp = %d, want 22", got)
	}
}

func TestXp(t *testing.T) {
	if got := XpForLevel(1); got != 0 {
		t.Errorf("XpForLevel(1) = %d, want 0", got)
	}
	if got := XpForLevel(5); got != 6500 {
		t.Errorf("XpForLevel(5) = %d, want 6500", got)
	}
	if got := XpForLevel(20); got != 355000 {
		t.Errorf("XpForLevel(20) = %d, want 355000", got)
	}
	if got := XpForNextLevel(4); got != 6500 {
		t.Errorf("XpForNextLevel(4) = %d, want 6500", got)
	}
}
