package turn

import "testing"

func TestInitiativeOrderAndLoop(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 14)
	tr.AddCombatant("b", 15, 10)
	tr.AddCombatant("c", 18, 12)

	// order: a(20), c(18), b(15). Current is the first.
	if e := tr.Current(); e.Actor != "a" || e.Round != 1 || e.Kind != EventTurn {
		t.Fatalf("current = %+v, want a round 1 turn", e)
	}
	if tr.Next().Actor != "c" {
		t.Error("second should be c")
	}
	if tr.Next().Actor != "b" {
		t.Error("third should be b")
	}
	// loop to round 2, back to a.
	e := tr.Next()
	if e.Actor != "a" || e.Round != 2 {
		t.Errorf("loop = %+v, want a round 2", e)
	}
	if tr.Round() != 2 {
		t.Errorf("Round() = %d, want 2", tr.Round())
	}
}

func TestDexTieBreak(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("lowdex", 15, 8)
	tr.AddCombatant("highdex", 15, 16)
	if tr.Current().Actor != "highdex" {
		t.Error("equal initiative: higher Dex acts first")
	}
}

func TestRemoveCombatant(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	tr.AddCombatant("b", 15, 10)
	tr.Current() // a, round 1
	tr.RemoveCombatant("b")
	// b is gone; next is a in round 2.
	e := tr.Next()
	if e.Actor != "a" || e.Round != 2 {
		t.Errorf("after remove = %+v, want a round 2", e)
	}
}

func TestUpcoming(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	tr.AddCombatant("b", 15, 10)
	tr.Current() // at a
	up := tr.Upcoming()
	if len(up) != 1 || up[0].Actor != "b" {
		t.Errorf("upcoming = %+v, want [b]", up)
	}
}
