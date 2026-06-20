package damage

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

func TestApplyToHp(t *testing.T) {
	// damage that doesn't drop to zero.
	if o := ApplyToHp(20, 20, -5); o.Hp != 15 || o.DroppedToZero || o.InstantDeath {
		t.Errorf("partial = %+v, want Hp 15, no flags", o)
	}
	// damage to exactly 0: dropped, not dead.
	if o := ApplyToHp(5, 20, -5); o.Hp != 0 || !o.DroppedToZero || o.InstantDeath {
		t.Errorf("to zero = %+v, want Hp 0 DroppedToZero true InstantDeath false", o)
	}
	// overkill below max: dropped, not dead.
	if o := ApplyToHp(5, 20, -10); o.Hp != 0 || !o.DroppedToZero || o.InstantDeath {
		t.Errorf("overkill < max = %+v, want Hp 0 dropped, not dead", o)
	}
	// massive damage: remaining past 0 >= max -> instant death.
	if o := ApplyToHp(5, 20, -30); o.Hp != 0 || !o.DroppedToZero || !o.InstantDeath {
		t.Errorf("massive = %+v, want Hp 0 dropped dead", o)
	}
	// healing caps at max.
	if o := ApplyToHp(18, 20, 10); o.Hp != 20 {
		t.Errorf("heal cap = %+v, want Hp 20", o)
	}
	// partial heal.
	if o := ApplyToHp(5, 20, 7); o.Hp != 12 {
		t.Errorf("heal = %+v, want Hp 12", o)
	}
}

func TestApply(t *testing.T) {
	d := Damage{Parts: []DamagePart{part(8, core.Fire, false)}}
	out, res := Apply(d, Mitigation{Resist: map[string]bool{core.Fire.Id: true}}, 20, 20)
	if res.Final != 4 || out.Hp != 16 {
		t.Errorf("apply = out %+v res %+v, want Final 4 Hp 16", out, res)
	}
}
