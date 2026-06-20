package core

// xpThresholds is the cumulative XP required to reach each level (index 1..20).
var xpThresholds = [21]Xp{
	0, 0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000,
	85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000,
}

// XpForLevel is the cumulative XP needed to reach the given level (1..20).
func XpForLevel(level Level) Xp { return xpThresholds[level] }

// XpForNextLevel is the cumulative XP needed to reach the level after the
// current one. At level 20 it returns the level-20 threshold.
func XpForNextLevel(current Level) Xp {
	if current >= 20 {
		return xpThresholds[20]
	}
	return xpThresholds[current+1]
}
