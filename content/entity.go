package content

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

// ItemKind classifies an item.
type ItemKind int

const (
	ItemUnspecified ItemKind = iota
	ItemWeapon
	ItemArmor
	ItemGear
)

// Item is a weapon, armor, or gear shape. Weapons carry a damage.Spec; armor
// carries an armor class. Games add their own flavor and rules data.
type Item struct {
	Id, Name   string
	Kind       ItemKind
	Damage     *damage.Spec    // when ItemWeapon
	ArmorClass core.ArmorClass // when ItemArmor
}

// Creature carries the stat-block fields the kernel pipeline reads (notably
// Mitigation, which damage.ApplyMitigation consumes, and HitDice for HP). Games
// extend it with their own data or define their own type entirely.
type Creature struct {
	Id, Name      string
	Size          core.Size
	Abilities     core.AbilityScores
	ArmorClass    core.ArmorClass
	HitDice       dice.Expr
	Mitigation    damage.Mitigation
	MovementSpeed core.Distance
}
