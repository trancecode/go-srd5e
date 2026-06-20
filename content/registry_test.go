package content

import "testing"

func TestRegistry(t *testing.T) {
	r := NewRegistry[int]()
	r.Register("a", 1)
	r.Register("b", 2)
	r.Register("c", 3)

	if v, ok := r.Get("b"); !ok || v != 2 {
		t.Errorf("Get(b) = %d,%v, want 2,true", v, ok)
	}
	if _, ok := r.Get("missing"); ok {
		t.Error("Get(missing) should be false")
	}
	// All in registration order.
	all := r.All()
	if len(all) != 3 || all[0] != 1 || all[1] != 2 || all[2] != 3 {
		t.Errorf("All = %v, want [1 2 3]", all)
	}
	// re-register updates in place, no duplicate entry.
	r.Register("b", 20)
	all = r.All()
	if len(all) != 3 || all[1] != 20 {
		t.Errorf("after update All = %v, want [1 20 3]", all)
	}
}

func TestRegistryZeroValueUsable(t *testing.T) {
	var r Registry[string]
	r.Register("x", "hello")
	if v, ok := r.Get("x"); !ok || v != "hello" {
		t.Error("zero-value Registry should be usable")
	}
}
