package content

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

// ItemKind classifies an item. ItemUnspecified is the zero value and must
// be set.
type ItemKind int

const (
	ItemUnspecified ItemKind = iota
	ItemWeapon
	ItemArmor
	ItemShield
	ItemAmmunition
	ItemGear
	ItemTool
	ItemPotion
	ItemScroll
	ItemWand
	ItemRing
	ItemWondrous
	// ItemWorn is a mundane worn accessory the SRD's equipment chapter
	// does not price as armour: a helmet, boots, bracers, a necklace, a
	// cloak. Games define the entries they need; the SRD ships none.
	ItemWorn
)

// WeaponProperty is one of the SRD's weapon properties. PropertyNone is the
// zero value, for weapons with none.
type WeaponProperty int

const (
	PropertyNone WeaponProperty = iota
	PropertyAmmunition
	PropertyFinesse
	PropertyHeavy
	PropertyLight
	PropertyLoading
	PropertyReach
	PropertySpecial
	PropertyThrown
	PropertyTwoHanded
	PropertyVersatile
)

// WeaponCategory is the SRD's proficiency grouping of a weapon.
type WeaponCategory int

const (
	WeaponCategoryNone WeaponCategory = iota
	WeaponCategorySimple
	WeaponCategoryMartial
)

// ArmorCategory is the SRD's grouping of armor, which decides the Dex cap
// and the proficiency needed.
type ArmorCategory int

const (
	ArmorCategoryNone ArmorCategory = iota
	ArmorCategoryLight
	ArmorCategoryMedium
	ArmorCategoryHeavy
	ArmorCategoryShield
)

// Item is an SRD equipment entry: what a kind of thing is, never where one
// is. Weapons carry a damage.Spec, armor an armor class; the other fields
// apply by Kind and are zero otherwise. Games add their own flavor and
// rules data beside it.
type Item struct {
	// Id is the item's stable identifier.
	Id string
	// Name is the item's display name.
	Name string
	// Kind classifies the item; see ItemKind.
	Kind ItemKind
	// Weight is the item's weight in pounds.
	Weight core.Weight
	// Cost is the item's price.
	Cost core.Coins
	// Stackable reports whether multiple units of the item occupy a single
	// inventory slot.
	Stackable bool
	// Quantity is the number of units the SRD entry prices and weighs as
	// one purchase: arrows are sold twenty at a time. Zero means a single
	// thing. Cost and Weight describe the whole bundle, not one unit,
	// while Name is the singular unit ("Arrow").
	Quantity int

	// Damage is the weapon's base damage. Set when Kind == ItemWeapon.
	Damage *damage.Spec
	// VersatileDamage is the two-handed damage for a weapon carrying
	// PropertyVersatile; nil otherwise.
	VersatileDamage *damage.Spec
	// Properties lists the weapon's SRD properties.
	Properties []WeaponProperty
	// WeaponCategory is the weapon's proficiency grouping.
	WeaponCategory WeaponCategory
	// Melee reports whether the weapon is wielded in melee.
	Melee bool
	// Range is the weapon's normal range, for ranged and thrown weapons.
	Range core.Distance
	// LongRange is the weapon's long range, for ranged and thrown weapons.
	LongRange core.Distance

	// ArmorCategory is the armor's SRD grouping.
	ArmorCategory ArmorCategory
	// ArmorClass is the armor class the armor or shield grants.
	ArmorClass core.ArmorClass
	// MaxDexBonus caps the Dexterity modifier added to armor class;
	// interpreted per ArmorCategory (light leaves it unrestricted, medium
	// caps it, heavy allows none).
	MaxDexBonus int
	// StrengthRequired is the minimum Strength score needed to wear the
	// armor without a speed penalty; zero means none.
	StrengthRequired core.AbilityScore
	// StealthDisadvantage reports whether wearing the armor imposes
	// disadvantage on Stealth checks.
	StealthDisadvantage bool

	// Charges is the number of uses remaining, for items with charges such
	// as wands.
	Charges int
}

// HasProperty reports whether the weapon has a property.
func (it Item) HasProperty(p WeaponProperty) bool {
	for _, q := range it.Properties {
		if q == p {
			return true
		}
	}
	return false
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
