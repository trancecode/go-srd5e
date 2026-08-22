package core

// xpThresholds is the cumulative XP required to reach each level (index 1..20).
var xpThresholds = [21]Xp{
	0, 0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000,
	85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000,
}

// MaxLevel is the highest character level the SRD progression defines.
const MaxLevel Level = 20

// XpForLevel is the cumulative XP needed to reach the given level. Levels below
// 1 report 0 and levels above MaxLevel report the MaxLevel threshold, so an
// out-of-range level clamps rather than panicking.
func XpForLevel(level Level) Xp {
	if level < 1 {
		return 0
	}
	if level > MaxLevel {
		return xpThresholds[MaxLevel]
	}
	return xpThresholds[level]
}

// XpForNextLevel is the cumulative XP needed to reach the level after the
// current one. At MaxLevel and above there is no next level, so it reports
// ok=false; callers must not treat the returned Xp as a threshold in that case.
func XpForNextLevel(current Level) (Xp, bool) {
	if current >= MaxLevel {
		return 0, false
	}
	return XpForLevel(current + 1), true
}
