package turn

import (
	"sort"

	"github.com/trancecode/go-srd5e/core"
)

// EventKind distinguishes a creature's turn from a scheduled effect tick.
type EventKind int

const (
	EventTurn EventKind = iota
	EventEffect
)

// TargetKind is what a scheduled effect targets.
type TargetKind int

const (
	TargetCreature TargetKind = iota
	TargetArea
)

// Timing places a scheduled effect before or after its anchor creature's turn.
type Timing int

const (
	Before Timing = iota
	After
)

// Active is a scheduled, recurring effect on the timeline. The game maps Id back
// to what the effect does; the tracker carries only identity, target, and timing.
type Active struct {
	Id         string
	Target     string
	TargetKind TargetKind
	Duration   core.EffectDuration
	Remaining  int // rounds left, when Duration.Kind == core.DurationRounds
	SaveDc     core.Dc
}

// Event is one entry on the timeline: a creature's turn (Actor set) or a due
// effect (Effect set).
type Event struct {
	Kind   EventKind
	Round  int
	Actor  string
	Effect *Active
}

type combatant struct {
	id         string
	initiative int
	dex        core.AbilityScore
	order      int
	active     bool
}

type scheduledEffect struct {
	a      Active
	anchor string
	when   Timing
}

// Tracker is an initiative-ordered, looping timeline of creature turns and
// scheduled effect ticks.
type Tracker struct {
	combatants []*combatant
	effects    []scheduledEffect
	nextOrder  int
	round      int
	seq        []Event
	pos        int
	started    bool
}

// NewTracker returns an empty tracker.
func NewTracker() *Tracker { return &Tracker{} }

// AddCombatant adds a combatant (or reactivates and updates one with the same
// Id). Dexterity score breaks initiative ties. Takes effect from the next round
// if combat is already underway.
func (t *Tracker) AddCombatant(id string, initiative int, dex core.AbilityScore) {
	for _, c := range t.combatants {
		if c.id == id {
			c.initiative = initiative
			c.dex = dex
			c.active = true
			return
		}
	}
	t.combatants = append(t.combatants, &combatant{id: id, initiative: initiative, dex: dex, order: t.nextOrder, active: true})
	t.nextOrder++
}

// RemoveCombatant stops a combatant's turns. Its anchored effects keep their
// slot. Any not-yet-reached turn for it this round is dropped immediately.
func (t *Tracker) RemoveCombatant(id string) {
	for _, c := range t.combatants {
		if c.id == id {
			c.active = false
			break
		}
	}
	if t.started {
		var ns []Event
		for i, e := range t.seq {
			if i > t.pos && e.Kind == EventTurn && e.Actor == id {
				continue
			}
			ns = append(ns, e)
		}
		t.seq = ns
	}
}

func (t *Tracker) ordered() []*combatant {
	cs := make([]*combatant, len(t.combatants))
	copy(cs, t.combatants)
	sort.SliceStable(cs, func(i, j int) bool {
		if cs[i].initiative != cs[j].initiative {
			return cs[i].initiative > cs[j].initiative
		}
		if cs[i].dex != cs[j].dex {
			return cs[i].dex > cs[j].dex
		}
		return cs[i].order < cs[j].order
	})
	return cs
}

// buildSeq composes the current round's event sequence. (Turns only; Task 4
// interleaves scheduled effects.)
func (t *Tracker) buildSeq() {
	t.seq = nil
	for _, c := range t.ordered() {
		if c.active {
			t.seq = append(t.seq, Event{Kind: EventTurn, Round: t.round, Actor: c.id})
		}
	}
}

func (t *Tracker) start() {
	t.started = true
	t.round = 1
	t.buildSeq()
	t.pos = 0
}

func (t *Tracker) cur() Event {
	if len(t.seq) == 0 {
		return Event{}
	}
	return t.seq[t.pos]
}

func (t *Tracker) rollover() {
	t.round++
	t.buildSeq()
	t.pos = 0
}

// Current returns the event the cursor is on (starting the first round on first
// use).
func (t *Tracker) Current() Event {
	if !t.started {
		t.start()
	}
	return t.cur()
}

// Next advances to and returns the next event, looping into the next round at
// the end of the order.
func (t *Tracker) Next() Event {
	if !t.started {
		t.start()
		return t.cur()
	}
	if len(t.seq) == 0 {
		return Event{}
	}
	t.pos++
	if t.pos >= len(t.seq) {
		t.rollover()
	}
	return t.cur()
}

// Round is the current round number (1-based).
func (t *Tracker) Round() int {
	if !t.started {
		t.start()
	}
	return t.round
}

// Upcoming is a copy of the events after the current one through the end of this
// round, for logging or rendering.
func (t *Tracker) Upcoming() []Event {
	if !t.started {
		t.start()
	}
	if t.pos+1 >= len(t.seq) {
		return nil
	}
	return append([]Event{}, t.seq[t.pos+1:]...)
}
