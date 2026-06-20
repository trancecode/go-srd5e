package resource

import "testing"

func TestSpellSlotPool(t *testing.T) {
	// perLevel indexed by spell level; index 0 is unused.
	p := NewSpellSlotPool([]int{0, 4, 3, 2})
	if p.Available(1) != 4 || p.Available(2) != 3 || p.Available(3) != 2 {
		t.Fatalf("available wrong: %+v", p.Current)
	}
	if p.Available(0) != 0 || p.Available(9) != 0 || p.Available(10) != 0 {
		t.Error("out-of-range levels should be 0")
	}
	if !p.Expend(1) || !p.Expend(1) || p.Available(1) != 2 {
		t.Errorf("expend wrong, available 1 = %d", p.Available(1))
	}
	if p.Expend(9) { // none at level 9
		t.Error("expend with no slot should be false")
	}
	p.RestoreAll()
	if p.Available(1) != 4 {
		t.Error("RestoreAll should refill")
	}
	// short rest is a no-op; long rest refills.
	p.Expend(1)
	p.Restore(RestShort)
	if p.Available(1) != 3 {
		t.Error("RestShort must not refill slots")
	}
	p.Restore(RestLong)
	if p.Available(1) != 4 {
		t.Error("RestLong must refill slots")
	}
}

// SpellSlotPool satisfies Restorable.
var _ Restorable = (*SpellSlotPool)(nil)
