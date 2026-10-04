package turn

import (
	"fmt"

	"github.com/trancecode/go-srd5e/core"
)

// TrackerState is a complete snapshot of a Tracker, for saving one mid-combat
// and restoring it with RestoreTracker. It is a plain value type that
// round-trips through encoding/json (or any other encoder), so a game can store
// it in its own save struct or ECS component.
type TrackerState struct {
	// Combatants lists every combatant in the order it was first added,
	// including removed ones (Active false), whose anchored effects keep their
	// slot.
	Combatants []CombatantState

	// Effects lists the scheduled effects in scheduling order, with their
	// current Remaining rounds.
	Effects []ScheduledEffect

	// NextOrder is the insertion tie-break the next new combatant receives.
	NextOrder int

	// Round is the current round number; 0 until the tracker has started.
	Round int

	// Sequence is the current round's timeline as built at the start of the
	// round. It is restored as is rather than rebuilt, because combatants and
	// effects added mid-round only join from the next round. Each effect
	// event carries its own copy of the Active, taken when the round was built.
	Sequence []Event

	// Position is the cursor into Sequence: the index of the current event.
	Position int

	// Started reports whether the first round has begun. A tracker starts on
	// the first Current, Next, Round, or Upcoming call.
	Started bool
}

// CombatantState is one combatant of a TrackerState.
type CombatantState struct {
	// Id is the game's opaque identifier for the combatant.
	Id string

	// Initiative is the combatant's initiative value.
	Initiative int

	// Dex is the Dexterity score that breaks initiative ties.
	Dex core.AbilityScore

	// Order is the insertion tie-break used when initiative and Dexterity are
	// both equal; lower goes first. Unique within a TrackerState and below
	// NextOrder.
	Order int

	// Active is false for a removed combatant, which takes no turns.
	Active bool
}

// ScheduledEffect is one effect scheduled on a TrackerState, as passed to
// Tracker.ScheduleEffect.
type ScheduledEffect struct {
	// Effect is the scheduled effect, with its current Remaining rounds.
	Effect Active

	// Anchor is the Id of the combatant whose turn the effect rides.
	Anchor string

	// When places the effect before or after the anchor's turn.
	When Timing
}

// State returns a snapshot of the tracker. It does not start the tracker: a
// snapshot taken before the first Current has Started false, and the restored
// tracker starts round 1 exactly as the original would. The snapshot shares no
// memory with the tracker.
func (t *Tracker) State() TrackerState {
	s := TrackerState{
		NextOrder: t.nextOrder,
		Round:     t.round,
		Sequence:  copyEvents(t.seq),
		Position:  t.pos,
		Started:   t.started,
	}
	for _, c := range t.combatants {
		s.Combatants = append(s.Combatants, CombatantState{Id: c.id, Initiative: c.initiative, Dex: c.dex, Order: c.order, Active: c.active})
	}
	for _, se := range t.effects {
		s.Effects = append(s.Effects, ScheduledEffect{Effect: se.a, Anchor: se.anchor, When: se.when})
	}
	return s
}

// RestoreTracker rebuilds a tracker from a snapshot taken by Tracker.State.
// The restored tracker produces the same Current, Next, Round, and Upcoming
// results as the original from the point of the snapshot on, including for
// combatants added afterwards. It shares no memory with s. A snapshot is
// external input, so an inconsistent one is reported as an error rather than a
// panic.
func RestoreTracker(s TrackerState) (*Tracker, error) {
	if err := validateState(s); err != nil {
		return nil, fmt.Errorf("restoring tracker: %w", err)
	}
	t := &Tracker{
		nextOrder: s.NextOrder,
		round:     s.Round,
		seq:       copyEvents(s.Sequence),
		pos:       s.Position,
		started:   s.Started,
	}
	for _, c := range s.Combatants {
		t.combatants = append(t.combatants, &combatant{id: c.Id, initiative: c.Initiative, dex: c.Dex, order: c.Order, active: c.Active})
	}
	for _, se := range s.Effects {
		t.effects = append(t.effects, scheduledEffect{a: se.Effect, anchor: se.Anchor, when: se.When})
	}
	return t, nil
}

func validateState(s TrackerState) error {
	ids := map[string]bool{}
	orders := map[int]bool{}
	for _, c := range s.Combatants {
		if ids[c.Id] {
			return fmt.Errorf("duplicate combatant Id %q", c.Id)
		}
		ids[c.Id] = true
		if c.Order < 0 || c.Order >= s.NextOrder {
			return fmt.Errorf("combatant %q has order %d outside [0, %d)", c.Id, c.Order, s.NextOrder)
		}
		if orders[c.Order] {
			return fmt.Errorf("combatant %q has duplicate order %d", c.Id, c.Order)
		}
		orders[c.Order] = true
	}
	for _, se := range s.Effects {
		if se.When != Before && se.When != After {
			return fmt.Errorf("effect %q has unknown timing %d", se.Effect.Id, se.When)
		}
	}
	if !s.Started {
		if s.Round != 0 || len(s.Sequence) != 0 || s.Position != 0 {
			return fmt.Errorf("not started but has round %d, %d events, position %d", s.Round, len(s.Sequence), s.Position)
		}
		return nil
	}
	if s.Round < 1 {
		return fmt.Errorf("started but has round %d", s.Round)
	}
	// An empty sequence while started is legal: every combatant was removed.
	if len(s.Sequence) == 0 && s.Position != 0 || len(s.Sequence) > 0 && (s.Position < 0 || s.Position >= len(s.Sequence)) {
		return fmt.Errorf("position %d outside a sequence of %d events", s.Position, len(s.Sequence))
	}
	// Effect events are not checked against Effects: an effect that fired
	// earlier this round may since have been cancelled.
	for i, e := range s.Sequence {
		if e.Round != s.Round {
			return fmt.Errorf("event %d has round %d, want %d", i, e.Round, s.Round)
		}
		switch e.Kind {
		case EventTurn:
			if !ids[e.Actor] {
				return fmt.Errorf("turn event %d names unknown combatant %q", i, e.Actor)
			}
			if e.Effect != nil {
				return fmt.Errorf("turn event %d carries effect %q", i, e.Effect.Id)
			}
		case EventEffect:
			if e.Effect == nil {
				return fmt.Errorf("effect event %d carries no effect", i)
			}
		default:
			return fmt.Errorf("event %d has unknown kind %d", i, e.Kind)
		}
	}
	return nil
}

// copyEvents deep-copies events, so no *Active is shared between a tracker and
// a snapshot.
func copyEvents(events []Event) []Event {
	var out []Event
	for _, e := range events {
		if e.Effect != nil {
			a := *e.Effect
			e.Effect = &a
		}
		out = append(out, e)
	}
	return out
}
