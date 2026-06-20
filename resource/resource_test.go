package resource

import "testing"

func TestResourceUseAndRecharge(t *testing.T) {
	ki := Resource{Id: "ki", Max: 3, Current: 3, Recharge: RechargeShortRest}
	if !ki.Use(2) || ki.Current != 1 {
		t.Fatalf("use 2 -> current %d, want 1", ki.Current)
	}
	if ki.Use(2) {
		t.Error("use 2 of 1 should be false")
	}
	ki.Restore(RestShort) // short covers short-rest recharge
	if ki.Current != 3 {
		t.Errorf("short rest -> %d, want 3", ki.Current)
	}

	feat := Resource{Id: "feat", Max: 1, Current: 0, Recharge: RechargeLongRest}
	feat.Restore(RestShort)
	if feat.Current != 0 {
		t.Error("long-rest resource must not refill on a short rest")
	}
	feat.Restore(RestLong)
	if feat.Current != 1 {
		t.Error("long-rest resource refills on a long rest")
	}

	none := Resource{Id: "x", Max: 1, Current: 0, Recharge: RechargeNone}
	none.Restore(RestLong)
	if none.Current != 0 {
		t.Error("RechargeNone never refills")
	}
}

func TestResourceSet(t *testing.T) {
	var s ResourceSet
	s.Add(Resource{Id: "ki", Max: 3, Current: 3, Recharge: RechargeShortRest})
	s.Add(Resource{Id: "channel", Max: 1, Current: 0, Recharge: RechargeShortRest})

	if !s.Use("ki", 1) {
		t.Error("use ki should succeed")
	}
	if r, _ := s.Get("ki"); r.Current != 2 {
		t.Errorf("ki current = %d, want 2", r.Current)
	}
	if s.Use("missing", 1) {
		t.Error("using a missing resource should be false")
	}
	s.Restore(RestShort) // both recharge on short rest
	if r, _ := s.Get("ki"); r.Current != 3 {
		t.Error("ki should be restored")
	}
	if r, _ := s.Get("channel"); r.Current != 1 {
		t.Error("channel should be restored")
	}
	// All returns deterministic, Id-sorted order.
	all := s.All()
	if len(all) != 2 || all[0].Id != "channel" || all[1].Id != "ki" {
		t.Errorf("All = %+v, want [channel, ki]", all)
	}
}
