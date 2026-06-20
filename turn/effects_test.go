package turn

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

func roundsEffect(id string, rounds int) Active {
	return Active{Id: id, Target: "t", Duration: core.EffectDuration{Kind: core.DurationRounds, Rounds: rounds}, Remaining: rounds}
}

func TestScheduleBeforeAndAfter(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	// schedule before a turn-cursor starts so it is in round 1.
	tr.ScheduleEffect(roundsEffect("burn", 3), "a", After)
	tr.ScheduleEffect(roundsEffect("aura", 3), "a", Before)

	// round 1: aura (before), a's turn, burn (after).
	if e := tr.Current(); e.Kind != EventEffect || e.Effect.Id != "aura" {
		t.Fatalf("first = %+v, want aura effect", e)
	}
	if e := tr.Next(); e.Kind != EventTurn || e.Actor != "a" {
		t.Fatalf("second = %+v, want a turn", e)
	}
	if e := tr.Next(); e.Kind != EventEffect || e.Effect.Id != "burn" {
		t.Fatalf("third = %+v, want burn effect", e)
	}
}

func TestEffectExpiresAfterRemaining(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	tr.ScheduleEffect(roundsEffect("burn", 2), "a", After)

	count := 0
	for i := 0; i < 12; i++ {
		e := tr.Next()
		if e.Kind == EventEffect && e.Effect.Id == "burn" {
			count++
		}
	}
	if count != 2 { // fires in rounds 1 and 2 only
		t.Errorf("burn fired %d times, want 2", count)
	}
}

func TestCancelImmediate(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	tr.ScheduleEffect(roundsEffect("hex", 5), "a", After)
	tr.Current() // aura/turn order: at a's turn (no Before effect), so Current is a's turn
	tr.Cancel("hex")
	// hex was After a; cancelled, so Next loops to round 2 a's turn, never firing hex.
	for i := 0; i < 8; i++ {
		if e := tr.Next(); e.Kind == EventEffect && e.Effect.Id == "hex" {
			t.Fatal("cancelled effect still fired")
		}
	}
}

func TestUntilRemovedRecursForever(t *testing.T) {
	tr := NewTracker()
	tr.AddCombatant("a", 20, 10)
	tr.ScheduleEffect(Active{Id: "wall", Duration: core.EffectDuration{Kind: core.DurationUntilRemoved}}, "a", After)
	count := 0
	for i := 0; i < 10; i++ {
		if e := tr.Next(); e.Kind == EventEffect && e.Effect.Id == "wall" {
			count++
		}
	}
	if count < 4 { // fires every round; over ~5 rounds it should fire several times
		t.Errorf("until-removed fired %d times, want it to recur every round", count)
	}
}
