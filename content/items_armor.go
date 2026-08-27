package content

import "github.com/trancecode/go-srd5e/core"

// armor builds one SRD armour entry. maxDex is -1 for light armour (the
// whole Dex modifier applies), the SRD's printed cap for medium armour, or
// 0 for heavy armour (no Dex modifier applies).
func armor(id, name string, cat ArmorCategory, cost core.Coins, weight core.Weight, ac, maxDex, str int, stealth bool) Item {
	return Item{
		Id: id, Name: name, Kind: ItemArmor, ArmorCategory: cat, Cost: cost, Weight: weight,
		ArmorClass: core.ArmorClass(ac), MaxDexBonus: maxDex, StrengthRequired: core.AbilityScore(str), StealthDisadvantage: stealth,
	}
}

// Armor returns the SRD 5.1 armour table and the shield, in the table's own
// order (light, medium, heavy, then the shield). Every field but weight is
// pinned to content/testdata/srd-2014-armor.json's data from Open5e's v1 API
// (api.open5e.com/v1/armor/?document__slug=wotc-srd), which carries AC, Dex
// cap, Strength requirement, stealth and cost but no weight; weight comes
// from Open5e's v2 API (api.open5e.com/v2/items/, filtered to document key
// "srd-2014"), hand-checked against the SRD 5.1 armour table.
func Armor() []Item {
	return []Item{
		// Light armor.
		armor("padded", "Padded", ArmorCategoryLight, core.Gp(5), 8, 11, -1, 0, true),
		armor("leather", "Leather", ArmorCategoryLight, core.Gp(10), 10, 11, -1, 0, false),
		armor("studded-leather", "Studded Leather", ArmorCategoryLight, core.Gp(45), 13, 12, -1, 0, false),

		// Medium armor.
		armor("hide", "Hide", ArmorCategoryMedium, core.Gp(10), 12, 12, 2, 0, false),
		armor("chain-shirt", "Chain Shirt", ArmorCategoryMedium, core.Gp(50), 20, 13, 2, 0, false),
		armor("scale-mail", "Scale mail", ArmorCategoryMedium, core.Gp(50), 45, 14, 2, 0, true),
		armor("breastplate", "Breastplate", ArmorCategoryMedium, core.Gp(400), 20, 14, 2, 0, false),
		armor("half-plate", "Half plate", ArmorCategoryMedium, core.Gp(750), 40, 15, 2, 0, true),

		// Heavy armor.
		armor("ring-mail", "Ring mail", ArmorCategoryHeavy, core.Gp(30), 40, 14, 0, 0, true),
		armor("chain-mail", "Chain mail", ArmorCategoryHeavy, core.Gp(75), 55, 16, 0, 13, true),
		armor("splint", "Splint", ArmorCategoryHeavy, core.Gp(200), 60, 17, 0, 15, true),
		armor("plate", "Plate", ArmorCategoryHeavy, core.Gp(1500), 65, 18, 0, 15, true),

		// Shield. The SRD gives it a flat +2 AC bonus, not a base AC
		// (Open5e's ac_string is "0 +2"); it is not "worn" for Dex-cap or
		// Strength-requirement purposes, so those fields stay zero.
		{Id: "shield", Name: "Shield", Kind: ItemShield, ArmorCategory: ArmorCategoryShield, Cost: core.Gp(10), Weight: 6, ArmorClass: 2},
	}
}
