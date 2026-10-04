package turn

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

// saveAndRestore takes a tracker's state through a JSON round trip, as a game
// saving between actions would, and restores it.
func saveAndRestore(t *testing.T, tr *Tracker) *Tracker {
	t.Helper()
	data, err := json.Marshal(tr.State())
	if err != nil {
		t.Fatalf("marshalling state: %v", err)
	}
	var s TrackerState
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("unmarshalling state: %v", err)
	}
	if !reflect.DeepEqual(s, tr.State()) {
		t.Fatalf("state changed through JSON:\n got %+v\nwant %+v", s, tr.State())
	}
	restored, err := RestoreTracker(s)
	if err != nil {
		t.Fatalf("restoring: %v", err)
	}
	if !reflect.DeepEqual(restored.State(), tr.State()) {
		t.Fatalf("restored state differs:\n got %+v\nwant %+v", restored.State(), tr.State())
	}
	return restored
}

// driveSideBySide advances both trackers in lockstep, comparing every Current,
// Round, Upcoming, and Next result by value (reflect.DeepEqual follows the
// *Active in each event).
func driveSideBySide(t *testing.T, orig, restored *Tracker, steps int) {
	t.Helper()
	for i := 0; i < steps; i++ {
		if a, b := orig.Current(), restored.Current(); !reflect.DeepEqual(a, b) {
			t.Fatalf("step %d: Current = %s, want %s", i, fmtEvent(b), fmtEvent(a))
		}
		if a, b := orig.Round(), restored.Round(); a != b {
			t.Fatalf("step %d: Round = %d, want %d", i, b, a)
		}
		if a, b := orig.Upcoming(), restored.Upcoming(); !reflect.DeepEqual(a, b) {
			t.Fatalf("step %d: Upcoming = %s, want %s", i, fmtEvents(b), fmtEvents(a))
		}
		if a, b := orig.Next(), restored.Next(); !reflect.DeepEqual(a, b) {
			t.Fatalf("step %d: Next = %s, want %s", i, fmtEvent(b), fmtEvent(a))
		}
	}
}

func fmtEvent(e Event) string {
	if e.Effect == nil {
		return fmt.Sprintf("{turn %s r%d}", e.Actor, e.Round)
	}
	return fmt.Sprintf("{effect %s r%d rem %d}", e.Effect.Id, e.Round, e.Effect.Remaining)
}

func fmtEvents(es []Event) string {
	var parts []string
	for _, e := range es {
		parts = append(parts, fmtEvent(e))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

// newBattle returns a tracker with three combatants, two of them tied on
// initiative and Dexterity, and effects of each duration kind on both sides of
// a turn.
func newBattle() *Tracker {
	tr := NewTracker()
	tr.AddCombatant("a", 18, 14)
	tr.AddCombatant("b", 12, 10)
	tr.AddCombatant("c", 12, 10)
	tr.ScheduleEffect(roundsEffect("burn", 3), "b", After)
	tr.ScheduleEffect(roundsEffect("aura", 2), "a", Before)
	tr.ScheduleEffect(Active{Id: "wall", TargetKind: TargetArea, Duration: core.EffectDuration{Kind: core.DurationUntilRemoved}, SaveDc: 13}, "c", Before)
	return tr
}

func TestRestoreBeforeStart(t *testing.T) {
	orig := newBattle()
	restored := saveAndRestore(t, orig)
	if s := restored.State(); s.Started || s.Round != 0 {
		t.Fatalf("restored tracker started early: %+v", s)
	}
	// A combatant added after the restore ties with b and c and must get the
	// same insertion order as on the original.
	for _, tr := range []*Tracker{orig, restored} {
		tr.AddCombatant("d", 12, 10)
	}
	driveSideBySide(t, orig, restored, 25)
}

func TestRestoreMidRound(t *testing.T) {
	orig := newBattle()
	orig.Current()
	orig.Next()
	orig.Next()
	// Mid-round additions join only from round 2; rebuilding the sequence on
	// restore would pull them into round 1.
	orig.AddCombatant("d", 20, 10)
	orig.ScheduleEffect(roundsEffect("hex", 2), "a", After)
	restored := saveAndRestore(t, orig)
	driveSideBySide(t, orig, restored, 25)
}

func TestRestoreAtRoundBoundary(t *testing.T) {
	orig := newBattle()
	orig.Current()
	for orig.Upcoming() != nil {
		orig.Next()
	}
	restored := saveAndRestore(t, orig)
	driveSideBySide(t, orig, restored, 25)
}

func TestRestoreAfterCancel(t *testing.T) {
	orig := newBattle()
	// Round 1 sequence: aura, a, b, burn, wall, c.
	orig.Current() // aura, before a's turn
	orig.Next()    // a's turn
	orig.Next()    // b's turn
	orig.Next()    // burn, after b's turn
	// aura has fired and burn is firing: both stay in the sequence at or before
	// the cursor. wall is still ahead and is dropped.
	orig.Cancel("aura")
	orig.Cancel("burn")
	orig.Cancel("wall")

	s := orig.State()
	for _, id := range []string{"aura", "burn"} {
		idx := -1
		for i, e := range s.Sequence {
			if e.Effect != nil && e.Effect.Id == id {
				idx = i
			}
		}
		if idx < 0 || idx > s.Position {
			t.Fatalf("setup: %s at sequence index %d, want at or before position %d", id, idx, s.Position)
		}
	}
	for _, e := range s.Sequence {
		if e.Effect != nil && e.Effect.Id == "wall" {
			t.Fatalf("setup: cancelled wall still in the sequence ahead of the cursor")
		}
	}
	for _, se := range s.Effects {
		if id := se.Effect.Id; id == "aura" || id == "burn" || id == "wall" {
			t.Fatalf("setup: cancelled effect %s still in Effects", id)
		}
	}

	restored := saveAndRestore(t, orig)
	driveSideBySide(t, orig, restored, 25)
}

func TestRestoreAfterRemove(t *testing.T) {
	orig := newBattle()
	orig.Current()
	orig.Next()
	orig.RemoveCombatant("b") // b's turn is ahead and is dropped; burn keeps its slot
	restored := saveAndRestore(t, orig)
	driveSideBySide(t, orig, restored, 10)
	// Reactivating b updates the inactive entry rather than adding a new one.
	for _, tr := range []*Tracker{orig, restored} {
		tr.AddCombatant("b", 25, 10)
	}
	driveSideBySide(t, orig, restored, 15)
}

func TestRestoreEmptySequence(t *testing.T) {
	orig := NewTracker()
	orig.AddCombatant("a", 10, 10)
	orig.Current()
	orig.RemoveCombatant("a")
	orig.Next() // rolls into round 2 with no one left
	if len(orig.State().Sequence) != 0 {
		t.Fatalf("setup: sequence = %+v, want empty", orig.State().Sequence)
	}
	restored := saveAndRestore(t, orig)
	driveSideBySide(t, orig, restored, 5)
}

func TestStateDoesNotStartTracker(t *testing.T) {
	tr := newBattle()
	if s := tr.State(); s.Started || s.Round != 0 || len(s.Sequence) != 0 {
		t.Fatalf("State on a fresh tracker = %+v, want not started", s)
	}
	if tr.started {
		t.Fatal("State started the tracker")
	}
}

func TestStateDoesNotAlias(t *testing.T) {
	tr := newBattle()
	tr.Current()
	want := tr.State()

	s := tr.State()
	s.Sequence[0].Effect.Remaining = 99
	s.Effects[0].Effect.Remaining = 99
	s.Combatants[0].Initiative = 99
	if !reflect.DeepEqual(tr.State(), want) {
		t.Fatal("mutating a snapshot changed the tracker")
	}

	s = tr.State()
	restored, err := RestoreTracker(s)
	if err != nil {
		t.Fatalf("restoring: %v", err)
	}
	s.Sequence[0].Effect.Remaining = 99
	s.Combatants[0].Initiative = 99
	if !reflect.DeepEqual(restored.State(), want) {
		t.Fatal("mutating the restored-from snapshot changed the tracker")
	}
}

func TestRestoreTrackerRejectsInvalidState(t *testing.T) {
	valid := func() TrackerState {
		tr := newBattle()
		tr.Current()
		return tr.State()
	}
	cases := []struct {
		name   string
		mutate func(*TrackerState)
		want   string
	}{
		{"duplicate id", func(s *TrackerState) { s.Combatants[1].Id = "a" }, `duplicate combatant Id "a"`},
		{"duplicate order", func(s *TrackerState) { s.Combatants[1].Order = 0 }, "duplicate order 0"},
		{"order at next order", func(s *TrackerState) { s.Combatants[0].Order = s.NextOrder }, "outside"},
		{"negative order", func(s *TrackerState) { s.Combatants[0].Order = -1 }, "outside"},
		{"unknown timing", func(s *TrackerState) { s.Effects[0].When = 7 }, "unknown timing 7"},
		{"not started with round", func(s *TrackerState) { *s = newBattle().State(); s.Round = 1 }, "not started"},
		{"not started with sequence", func(s *TrackerState) {
			seq := s.Sequence
			*s = newBattle().State()
			s.Sequence = seq
		}, "not started"},
		{"started at round 0", func(s *TrackerState) { s.Round = 0 }, "started but has round 0"},
		{"position past end", func(s *TrackerState) { s.Position = len(s.Sequence) }, "outside a sequence"},
		{"negative position", func(s *TrackerState) { s.Position = -1 }, "outside a sequence"},
		{"position in empty sequence", func(s *TrackerState) { s.Sequence = nil; s.Position = 1 }, "outside a sequence of 0"},
		{"event from another round", func(s *TrackerState) { s.Sequence[1].Round = 2 }, "event 1 has round 2"},
		{"turn for unknown combatant", func(s *TrackerState) { s.Sequence[1].Actor = "z" }, `unknown combatant "z"`},
		{"turn with effect", func(s *TrackerState) { s.Sequence[1].Effect = &Active{Id: "x"} }, "carries effect"},
		{"effect without effect", func(s *TrackerState) { s.Sequence[0].Effect = nil }, "carries no effect"},
		{"unknown kind", func(s *TrackerState) { s.Sequence[1].Kind = 5 }, "unknown kind 5"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := valid()
			tc.mutate(&s)
			tr, err := RestoreTracker(s)
			if err == nil {
				t.Fatalf("RestoreTracker returned %+v, want error", tr)
			}
			if !strings.HasPrefix(err.Error(), "restoring tracker: ") || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error = %q, want it to contain %q", err, tc.want)
			}
		})
	}
}
