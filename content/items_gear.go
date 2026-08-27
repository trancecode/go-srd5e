package content

import "github.com/trancecode/go-srd5e/core"

// gear builds one non-ammunition gear entry: adventuring gear or a tool.
func gear(id, name string, kind ItemKind, cost core.Coins, weight core.Weight) Item {
	return Item{Id: id, Name: name, Kind: kind, Cost: cost, Weight: weight}
}

// ammunition builds one SRD ammunition entry, priced and weighed for the
// bundle the SRD sells it in (Arrows (20), Blowgun needles (50), ...), not
// per unit; ammunition stacks in inventory.
func ammunition(id, name string, cost core.Coins, weight core.Weight) Item {
	return Item{Id: id, Name: name, Kind: ItemAmmunition, Stackable: true, Cost: cost, Weight: weight}
}

// Gear returns the ammunition, tools and adventuring gear the engine ships
// as content: a chosen subset of the SRD 5.1 gear table, not the whole
// table. Every value is hand-checked against the SRD 5.1 text; Open5e's v2
// API (api.open5e.com/v2/items/, filtered to document key "srd-2014") prices
// and weighs ammunition per unit rather than per SRD bundle, and its scaled
// numbers depart from the SRD text for the crossbow bolt's weight and the
// sling bullet's cost. See content/items_srd_test.go's srdItemCorrections
// for both.
func Gear() []Item {
	return []Item{
		ammunition("arrow", "Arrows (20)", core.Gp(1), 1),
		ammunition("blowgun-needle", "Blowgun needles (50)", core.Gp(1), 1),
		ammunition("crossbow-bolt", "Crossbow bolts (20)", core.Gp(1), 1.5),
		ammunition("sling-bullet", "Sling bullets (20)", core.Cp(4), 1.5),

		gear("thieves-tools", "Thieves' Tools", ItemTool, core.Gp(25), 1),

		gear("rations", "Rations (1 day)", ItemGear, core.Sp(5), 2),
		gear("torch", "Torch", ItemGear, core.Cp(1), 1),
		gear("rope-hempen", "Rope, hempen (50 feet)", ItemGear, core.Gp(1), 10),
		gear("waterskin", "Waterskin", ItemGear, core.Sp(2), 5),
		gear("backpack", "Backpack", ItemGear, core.Gp(2), 5),
	}
}
