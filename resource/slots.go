package resource

// RestKind is the kind of rest taken. RestLong covers both short- and long-rest
// recharges; RestShort covers only short-rest recharges.
type RestKind int

const (
	RestShort RestKind = iota
	RestLong
)

// Restorable is the shared rest behavior of every pool.
type Restorable interface {
	Restore(rest RestKind)
}

// SpellSlotPool is a per-level slot pool. Index is the spell level (1..9); index
// 0 is unused. Cantrips cost no slot and never touch this pool.
type SpellSlotPool struct {
	Max     [10]int
	Current [10]int
}

// NewSpellSlotPool builds a full pool from per-level slot counts (indexed by
// spell level; index 0 is ignored), as pulled from a class's slot progression.
func NewSpellSlotPool(perLevel []int) SpellSlotPool {
	var p SpellSlotPool
	for i := 0; i < len(perLevel) && i < len(p.Max); i++ {
		p.Max[i] = perLevel[i]
		p.Current[i] = perLevel[i]
	}
	return p
}

// Available is the current slots of a spell level, or 0 for an out-of-range level.
func (p *SpellSlotPool) Available(level int) int {
	if level < 1 || level > 9 {
		return 0
	}
	return p.Current[level]
}

// Expend spends one slot of the given level; the caller upcasts by passing a
// higher level. Returns false if none are available.
func (p *SpellSlotPool) Expend(level int) bool {
	if level < 1 || level > 9 || p.Current[level] <= 0 {
		return false
	}
	p.Current[level]--
	return true
}

// RestoreAll refills every level to its max (e.g. Arcane Recovery, or a long rest).
func (p *SpellSlotPool) RestoreAll() {
	for i := range p.Current {
		p.Current[i] = p.Max[i]
	}
}

// Restore refills all slots on a long rest; a short rest is a no-op.
func (p *SpellSlotPool) Restore(rest RestKind) {
	if rest == RestLong {
		p.RestoreAll()
	}
}
