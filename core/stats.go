package core

// ProficiencyBonus is 2 + (level-1)/4.
func ProficiencyBonus(level Level) Modifier { return Modifier(2 + (int(level)-1)/4) }

// SkillBonus is the ability modifier plus the proficiency bonus when proficient.
func SkillBonus(score AbilityScore, proficient bool, level Level) Modifier {
	return abilityCheckBonus(score, proficient, level)
}

// SavingThrowBonus is the ability modifier plus the proficiency bonus when proficient.
func SavingThrowBonus(score AbilityScore, proficient bool, level Level) Modifier {
	return abilityCheckBonus(score, proficient, level)
}

func abilityCheckBonus(score AbilityScore, proficient bool, level Level) Modifier {
	b := AbilityModifier(score)
	if proficient {
		b += ProficiencyBonus(level)
	}
	return b
}

// SpellSaveDc is 8 + proficiency bonus + spellcasting ability modifier.
func SpellSaveDc(score AbilityScore, level Level) Dc {
	return Dc(8 + int(ProficiencyBonus(level)) + int(AbilityModifier(score)))
}

// SpellAttackBonus is proficiency bonus + spellcasting ability modifier.
func SpellAttackBonus(score AbilityScore, level Level) Modifier {
	return ProficiencyBonus(level) + AbilityModifier(score)
}

// MaxHp sums the per-level hit-point rolls and adds the Constitution modifier
// once per level.
func MaxHp(hpRolls []int, conMod Modifier) HitPoints {
	total := 0
	for _, r := range hpRolls {
		total += r
	}
	return HitPoints(total + int(conMod)*len(hpRolls))
}
