package content

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

// weapon builds one SRD weapon entry. The dice are parsed once at package
// initialisation; a bad expression is a programming error. dmg is empty for
// a weapon with no damage die (the net).
func weapon(id, name string, cat WeaponCategory, cost core.Coins, weight core.Weight, dmg string, dt core.DamageType, props ...WeaponProperty) Item {
	it := Item{Id: id, Name: name, Kind: ItemWeapon, WeaponCategory: cat, Cost: cost, Weight: weight, Properties: props, Melee: true}
	if dmg != "" {
		it.Damage = &damage.Spec{Parts: []damage.PartSpec{{Dice: dice.MustParse(dmg), Type: dt}}}
	}
	return it
}

// ranged sets a weapon's normal and long range in feet. A thrown melee
// weapon keeps Melee true; a weapon with the Ammunition property is not a
// melee weapon.
func (it Item) ranged(normal, long core.Distance) Item {
	it.Range, it.LongRange = normal, long
	if it.HasProperty(PropertyAmmunition) {
		it.Melee = false
	}
	return it
}

// versatile sets the two-handed damage of a versatile weapon.
func (it Item) versatile(dmg string) Item {
	it.VersatileDamage = &damage.Spec{Parts: []damage.PartSpec{{Dice: dice.MustParse(dmg), Type: it.Damage.Parts[0].Type}}}
	return it
}

// Weapons returns the SRD 5.1 weapon table, in the order the SRD lists it:
// simple melee, simple ranged, martial melee, martial ranged. Every value is
// pinned to content/testdata/srd-2014-weapons.json, Open5e's v1 API data for
// the SRD 5.1 weapon table (api.open5e.com/v1/weapons/?document__slug=wotc-srd).
func Weapons() []Item {
	return []Item{
		// Simple melee weapons.
		weapon("club", "Club", WeaponCategorySimple, core.Sp(1), 2, "1d4", core.Bludgeoning, PropertyLight),
		weapon("dagger", "Dagger", WeaponCategorySimple, core.Gp(2), 1, "1d4", core.Piercing, PropertyFinesse, PropertyLight, PropertyThrown).ranged(20, 60),
		weapon("greatclub", "Greatclub", WeaponCategorySimple, core.Sp(2), 10, "1d8", core.Bludgeoning, PropertyTwoHanded),
		weapon("handaxe", "Handaxe", WeaponCategorySimple, core.Gp(5), 2, "1d6", core.Slashing, PropertyLight, PropertyThrown).ranged(20, 60),
		weapon("javelin", "Javelin", WeaponCategorySimple, core.Sp(5), 2, "1d6", core.Piercing, PropertyThrown).ranged(30, 120),
		weapon("light-hammer", "Light hammer", WeaponCategorySimple, core.Gp(2), 2, "1d4", core.Bludgeoning, PropertyLight, PropertyThrown).ranged(20, 60),
		weapon("mace", "Mace", WeaponCategorySimple, core.Gp(5), 4, "1d6", core.Bludgeoning),
		weapon("quarterstaff", "Quarterstaff", WeaponCategorySimple, core.Sp(2), 4, "1d6", core.Bludgeoning, PropertyVersatile).versatile("1d8"),
		weapon("sickle", "Sickle", WeaponCategorySimple, core.Gp(1), 2, "1d4", core.Slashing, PropertyLight),
		weapon("spear", "Spear", WeaponCategorySimple, core.Gp(1), 3, "1d6", core.Piercing, PropertyThrown, PropertyVersatile).ranged(20, 60).versatile("1d8"),

		// Simple ranged weapons.
		weapon("crossbow-light", "Crossbow, light", WeaponCategorySimple, core.Gp(25), 5, "1d8", core.Piercing, PropertyAmmunition, PropertyLoading, PropertyTwoHanded).ranged(80, 320),
		weapon("dart", "Dart", WeaponCategorySimple, core.Cp(5), 0.25, "1d4", core.Piercing, PropertyFinesse, PropertyThrown).ranged(20, 60),
		weapon("shortbow", "Shortbow", WeaponCategorySimple, core.Gp(25), 2, "1d6", core.Piercing, PropertyAmmunition, PropertyTwoHanded).ranged(80, 320),
		weapon("sling", "Sling", WeaponCategorySimple, core.Sp(1), 0, "1d4", core.Bludgeoning, PropertyAmmunition).ranged(30, 120),

		// Martial melee weapons.
		weapon("battleaxe", "Battleaxe", WeaponCategoryMartial, core.Gp(10), 4, "1d8", core.Slashing, PropertyVersatile).versatile("1d10"),
		weapon("flail", "Flail", WeaponCategoryMartial, core.Gp(10), 2, "1d8", core.Bludgeoning),
		weapon("glaive", "Glaive", WeaponCategoryMartial, core.Gp(20), 6, "1d10", core.Slashing, PropertyHeavy, PropertyReach, PropertyTwoHanded),
		weapon("greataxe", "Greataxe", WeaponCategoryMartial, core.Gp(30), 7, "1d12", core.Slashing, PropertyHeavy, PropertyTwoHanded),
		weapon("greatsword", "Greatsword", WeaponCategoryMartial, core.Gp(50), 6, "2d6", core.Slashing, PropertyHeavy, PropertyTwoHanded),
		weapon("halberd", "Halberd", WeaponCategoryMartial, core.Gp(20), 6, "1d10", core.Slashing, PropertyHeavy, PropertyReach, PropertyTwoHanded),
		weapon("lance", "Lance", WeaponCategoryMartial, core.Gp(10), 6, "1d12", core.Piercing, PropertyReach, PropertySpecial),
		weapon("longsword", "Longsword", WeaponCategoryMartial, core.Gp(15), 3, "1d8", core.Slashing, PropertyVersatile).versatile("1d10"),
		weapon("maul", "Maul", WeaponCategoryMartial, core.Gp(10), 10, "2d6", core.Bludgeoning, PropertyHeavy, PropertyTwoHanded),
		weapon("morningstar", "Morningstar", WeaponCategoryMartial, core.Gp(15), 4, "1d8", core.Piercing),
		weapon("pike", "Pike", WeaponCategoryMartial, core.Gp(5), 18, "1d10", core.Piercing, PropertyHeavy, PropertyReach, PropertyTwoHanded),
		weapon("rapier", "Rapier", WeaponCategoryMartial, core.Gp(25), 2, "1d8", core.Piercing, PropertyFinesse),
		weapon("scimitar", "Scimitar", WeaponCategoryMartial, core.Gp(25), 3, "1d6", core.Slashing, PropertyFinesse, PropertyLight),
		weapon("shortsword", "Shortsword", WeaponCategoryMartial, core.Gp(10), 2, "1d6", core.Piercing, PropertyFinesse, PropertyLight),
		weapon("trident", "Trident", WeaponCategoryMartial, core.Gp(5), 4, "1d6", core.Piercing, PropertyThrown, PropertyVersatile).ranged(20, 60).versatile("1d8"),
		weapon("war-pick", "War pick", WeaponCategoryMartial, core.Gp(5), 2, "1d8", core.Piercing),
		weapon("warhammer", "Warhammer", WeaponCategoryMartial, core.Gp(15), 2, "1d8", core.Bludgeoning, PropertyVersatile).versatile("1d10"),
		weapon("whip", "Whip", WeaponCategoryMartial, core.Gp(2), 3, "1d4", core.Slashing, PropertyFinesse, PropertyReach),

		// Martial ranged weapons.
		// Blowgun deals a flat 1 damage; "1d1" is the dice expression for a
		// fixed value (a one-sided die always shows 1), the same idiom the
		// creature blocks use for a flat amount.
		weapon("blowgun", "Blowgun", WeaponCategoryMartial, core.Gp(10), 1, "1d1", core.Piercing, PropertyAmmunition, PropertyLoading).ranged(25, 100),
		weapon("crossbow-hand", "Crossbow, hand", WeaponCategoryMartial, core.Gp(75), 3, "1d6", core.Piercing, PropertyAmmunition, PropertyLight, PropertyLoading).ranged(30, 120),
		weapon("crossbow-heavy", "Crossbow, heavy", WeaponCategoryMartial, core.Gp(50), 18, "1d10", core.Piercing, PropertyAmmunition, PropertyHeavy, PropertyLoading, PropertyTwoHanded).ranged(100, 400),
		weapon("longbow", "Longbow", WeaponCategoryMartial, core.Gp(50), 2, "1d8", core.Piercing, PropertyAmmunition, PropertyHeavy, PropertyTwoHanded).ranged(150, 600),
		// The net has no damage die: it restrains rather than hurts.
		weapon("net", "Net", WeaponCategoryMartial, core.Gp(1), 3, "", core.Bludgeoning, PropertySpecial, PropertyThrown).ranged(5, 15),
	}
}
