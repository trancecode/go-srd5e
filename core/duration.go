package core

// DurationKind classifies how long an applied effect or condition lasts.
type DurationKind int

const (
	DurationUnspecified DurationKind = iota
	DurationInstant
	DurationRounds
	DurationEndOfNextTurn
	DurationConcentration
	DurationUntilRemoved
)

// EffectDuration describes how long an applied effect lasts. Named EffectDuration
// (not Duration) to stay clear of time.Duration. The game ticks the countdown;
// SaveEnds means the bearer repeats the save at the end of each of its turns.
type EffectDuration struct {
	Kind        DurationKind
	Rounds      int // when DurationRounds
	SaveEnds    bool
	SaveAbility Ability // which save, when SaveEnds
}

// RoundsInMinutes converts minutes to rounds (1 round = 6 seconds).
func RoundsInMinutes(m int) int { return m * 10 }

// RoundsInHours converts hours to rounds.
func RoundsInHours(h int) int { return h * 600 }
