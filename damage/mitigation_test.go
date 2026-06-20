package damage

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

func part(amount int, t core.DamageType, magical bool) DamagePart {
	return DamagePart{Amount: amount, Type: t, Magical: magical}
}

func TestApplyMitigation(t *testing.T) {
	d := Damage{Parts: []DamagePart{part(10, core.Fire, false)}}

	// no mitigation: full.
	if r := ApplyMitigation(d, Mitigation{}); r.Final != 10 || r.Raw != 10 {
		t.Errorf("none = %+v, want Raw 10 Final 10", r)
	}
	// resistance halves (round down).
	if r := ApplyMitigation(Damage{Parts: []DamagePart{part(9, core.Fire, false)}},
		Mitigation{Resist: map[string]bool{core.Fire.Id: true}}); r.Final != 4 {
		t.Errorf("resist 9 = %d, want 4", r.Final)
	}
	// vulnerability doubles.
	if r := ApplyMitigation(d, Mitigation{Vulnerable: map[string]bool{core.Fire.Id: true}}); r.Final != 20 {
		t.Errorf("vulnerable = %d, want 20", r.Final)
	}
	// immunity zeroes.
	if r := ApplyMitigation(d, Mitigation{Immune: map[string]bool{core.Fire.Id: true}}); r.Final != 0 {
		t.Errorf("immune = %d, want 0", r.Final)
	}
	// flat reduction per type plus the DamageAny wildcard, floored at 0.
	if r := ApplyMitigation(d, Mitigation{FlatReduction: map[string]int{core.Fire.Id: 3, core.DamageAny.Id: 2}}); r.Final != 5 {
		t.Errorf("flat reduction = %d, want 5 (10-3-2)", r.Final)
	}
	if r := ApplyMitigation(d, Mitigation{FlatReduction: map[string]int{core.DamageAny.Id: 100}}); r.Final != 0 {
		t.Errorf("flat reduction floor = %d, want 0", r.Final)
	}
}

func TestResistNonmagicalPhysical(t *testing.T) {
	m := Mitigation{ResistNonmagicalPhysical: true}
	// nonmagical slashing: resisted (halved).
	if r := ApplyMitigation(Damage{Parts: []DamagePart{part(10, core.Slashing, false)}}, m); r.Final != 5 {
		t.Errorf("nonmagical slashing = %d, want 5", r.Final)
	}
	// magical slashing: not resisted.
	if r := ApplyMitigation(Damage{Parts: []DamagePart{part(10, core.Slashing, true)}}, m); r.Final != 10 {
		t.Errorf("magical slashing = %d, want 10", r.Final)
	}
	// nonmagical fire (not physical): not resisted.
	if r := ApplyMitigation(Damage{Parts: []DamagePart{part(10, core.Fire, false)}}, m); r.Final != 10 {
		t.Errorf("nonmagical fire = %d, want 10", r.Final)
	}
}

func TestByType(t *testing.T) {
	d := Damage{Parts: []DamagePart{part(8, core.Slashing, false), part(6, core.Fire, false)}}
	r := ApplyMitigation(d, Mitigation{Resist: map[string]bool{core.Fire.Id: true}})
	if r.ByType[core.Slashing.Id] != 8 || r.ByType[core.Fire.Id] != 3 || r.Final != 11 {
		t.Errorf("by type = %+v, want slashing 8 fire 3 final 11", r)
	}
}
