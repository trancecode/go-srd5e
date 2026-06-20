package damage

import "github.com/trancecode/go-srd5e/core"

// Mitigation is everything that reduces or modifies incoming damage. Resistance,
// vulnerability, and immunity are SRD (per-type sets keyed by DamageType.Id);
// FlatReduction is a generic extension keyed by DamageType.Id plus the wildcard
// core.DamageAny.Id that applies to every type. A zero-value Mitigation (nil
// maps) is a safe no-op.
type Mitigation struct {
	Immune, Resist, Vulnerable map[string]bool
	ResistNonmagicalPhysical   bool
	FlatReduction              map[string]int
}

// Result is the outcome of mitigation: the raw and final totals and the final
// amount per damage-type Id.
type Result struct {
	Raw    int
	Final  int
	ByType map[string]int
}

// ApplyMitigation reduces each typed part by immunity, vulnerability, resistance,
// and flat reduction (in that order), flooring each part at zero.
func ApplyMitigation(d Damage, m Mitigation) Result {
	res := Result{ByType: map[string]int{}}
	for _, p := range d.Parts {
		res.Raw += p.Amount
		amount := p.Amount
		switch {
		case m.Immune[p.Type.Id]:
			amount = 0
		default:
			if m.Vulnerable[p.Type.Id] {
				amount *= 2
			}
			if m.Resist[p.Type.Id] || (m.ResistNonmagicalPhysical && !p.Magical && isPhysical(p.Type)) {
				amount /= 2
			}
			amount -= m.FlatReduction[p.Type.Id] + m.FlatReduction[core.DamageAny.Id]
			if amount < 0 {
				amount = 0
			}
		}
		res.Final += amount
		res.ByType[p.Type.Id] += amount
	}
	return res
}

func isPhysical(t core.DamageType) bool {
	return t.Id == core.Bludgeoning.Id || t.Id == core.Piercing.Id || t.Id == core.Slashing.Id
}
