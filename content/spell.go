package content

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/effect"
)

// SpellResolution is how a spell is resolved against its targets.
type SpellResolution int

const (
	ResolveUnspecified SpellResolution = iota
	ResolveSpellAttack
	ResolveSave
	ResolveAuto
)

// RangeKind classifies a spell's range; the game interprets it spatially.
type RangeKind int

const (
	RangeUnspecified RangeKind = iota
	RangeSelf
	RangeTouch
	RangeRanged
)

// TargetKind classifies what a spell targets.
type TargetKind int

const (
	TargetUnspecified TargetKind = iota
	TargetSelf
	TargetSingle
	TargetMultiple
	TargetArea
)

// AreaShape is the shape of an area spell; AreaNone for non-area spells.
type AreaShape int

const (
	AreaNone AreaShape = iota
	AreaSphere
	AreaCone
	AreaLine
	AreaCube
)

// Targeting is the spatial descriptor the game interprets (geometry stays in the
// game).
type Targeting struct {
	Range      RangeKind
	Distance   core.Distance // when RangeRanged
	Target     TargetKind
	MaxTargets int       // when TargetMultiple
	Shape      AreaShape // when TargetArea
	AreaSize   core.Distance
}

// Spell is a self-describing spell shape: the game reads Resolution to pick an
// attack roll or a save, fills an effect.Outcome, then runs effect.Triggered.
type Spell struct {
	Id, Name    string
	Level       int
	Resolution  SpellResolution
	SaveAbility core.Ability // used when Resolution == ResolveSave
	Targeting   Targeting
	Effects     []effect.ConditionalEffect
}
