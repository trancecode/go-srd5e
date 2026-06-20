package combat

import "github.com/trancecode/go-srd5e/core"

// AttackOutcome is the result kind of an attack roll.
type AttackOutcome int

const (
	AttackMiss AttackOutcome = iota
	AttackHit
	AttackCritical
)

// AttackResult is the outcome of resolving an attack roll against armor class.
type AttackResult struct {
	Outcome     AttackOutcome
	NaturalRoll int
	Total       int
}

// ResolveAttack resolves a d20 attack: natural 20 is a critical hit and natural
// 1 an automatic miss, regardless of modifiers; otherwise the total (natural +
// modifier) is compared to the target's armor class.
func ResolveAttack(naturalRoll int, mod core.Modifier, ac core.ArmorClass) AttackResult {
	total := naturalRoll + int(mod)
	switch {
	case naturalRoll == 20:
		return AttackResult{Outcome: AttackCritical, NaturalRoll: naturalRoll, Total: total}
	case naturalRoll == 1:
		return AttackResult{Outcome: AttackMiss, NaturalRoll: naturalRoll, Total: total}
	case total >= int(ac):
		return AttackResult{Outcome: AttackHit, NaturalRoll: naturalRoll, Total: total}
	default:
		return AttackResult{Outcome: AttackMiss, NaturalRoll: naturalRoll, Total: total}
	}
}

// ConcentrationDc is the save DC to maintain concentration after taking damage:
// the greater of 10 and half the damage taken.
func ConcentrationDc(damage int) core.Dc {
	dc := damage / 2
	if dc < 10 {
		dc = 10
	}
	return core.Dc(dc)
}
