package content

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

// srdWeapon is one weapon as Open5e's v1 API serves it
// (api.open5e.com/v1/weapons/?document__slug=wotc-srd), a complete dataset
// for the SRD 5.1 weapon table. An earlier version of this fixture used the
// v2 API, whose weapon data has gaps (missing properties on several
// entries, one wrong damage type); v1 does not share those gaps, so no
// per-row correction table is needed here.
type srdWeapon struct {
	Slug       string   `json:"slug"`
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	Cost       string   `json:"cost"`
	Weight     string   `json:"weight"`
	DamageDice string   `json:"damage_dice"`
	DamageType string   `json:"damage_type"`
	Properties []string `json:"properties"`
}

func loadFixture[T any](t *testing.T, name string) []T {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var rows []T
	if err := json.Unmarshal(b, &rows); err != nil {
		t.Fatal(err)
	}
	return rows
}

func byId(items []Item) map[string]Item {
	m := make(map[string]Item, len(items))
	for _, it := range items {
		m[it.Id] = it
	}
	return m
}

// srdCost parses Open5e's cost text ("50 gp", "1 sp", "5 cp") and renders
// it the way core.Coins does.
func srdCost(cost string) string {
	amount, unit, ok := strings.Cut(cost, " ")
	if !ok {
		panic(fmt.Sprintf("parsing SRD cost %q: want \"N unit\"", cost))
	}
	n, err := strconv.Atoi(amount)
	if err != nil {
		panic(fmt.Sprintf("parsing SRD cost %q: %v", cost, err))
	}
	var c core.Coins
	switch unit {
	case "gp":
		c = core.Gp(n)
	case "sp":
		c = core.Sp(n)
	case "cp":
		c = core.Cp(n)
	default:
		panic(fmt.Sprintf("parsing SRD cost %q: unknown unit %q", cost, unit))
	}
	return c.String()
}

// srdWeight parses Open5e's weight text ("18 lb.", "0 lb.", "1/4 lb.")
// into core.Weight.
func srdWeight(w string) core.Weight {
	s := strings.TrimSuffix(strings.TrimSpace(w), "lb.")
	s = strings.TrimSpace(s)
	if before, after, ok := strings.Cut(s, "/"); ok {
		num, errN := strconv.ParseFloat(before, 64)
		den, errD := strconv.ParseFloat(after, 64)
		if errN != nil || errD != nil || den == 0 {
			panic(fmt.Sprintf("parsing SRD weight %q", w))
		}
		return core.Weight(num / den)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(fmt.Sprintf("parsing SRD weight %q: %v", w, err))
	}
	return core.Weight(f)
}

// srdDiceExpr normalizes Open5e's damage_dice text into the notation
// dice.Expr.String() produces. Every SRD 5.1 weapon but the blowgun already
// uses "NdS" notation; the blowgun's flat 1 damage is written as the bare
// integer "1", which content/items_weapons.go represents as the dice
// expression "1d1" (a one-sided die always shows 1), the same zero-variance
// idiom the creature blocks use for a flat value.
func srdDiceExpr(s string) string {
	if !strings.Contains(s, "d") {
		return s + "d" + s
	}
	return s
}

// srdHasDamage reports whether the SRD fixture's damage_dice describes real
// damage. Open5e writes "no damage" as "0" (the net has none): a net
// restrains rather than hurts.
func srdHasDamage(dice string) bool { return dice != "" && dice != "0" }

// srdCategory parses Open5e's category text ("Simple Melee Weapons",
// "Martial Ranged Weapons") into the SRD's proficiency grouping. The
// melee/ranged half of the heading is not a field on Item: it names a table
// in the book, not a rule Item.Melee enforces. A thrown weapon (the
// handaxe, the javelin, the trident, ...) lives in a melee table but keeps
// Item.Melee true when thrown; the net lives in the martial ranged table
// but is thrown, not fired, so it too keeps Item.Melee true. Item.Melee is
// instead derived from the Ammunition property; see the invariant check in
// TestWeaponsMatchTheSrd.
func srdCategory(cat string) WeaponCategory {
	switch {
	case strings.HasPrefix(cat, "Simple "):
		return WeaponCategorySimple
	case strings.HasPrefix(cat, "Martial "):
		return WeaponCategoryMartial
	default:
		panic(fmt.Sprintf("parsing SRD category %q", cat))
	}
}

// splitProperty splits one SRD property string ("versatile (1d10)",
// "light") into its bare name and parenthetical detail, if any: a dice
// expression for Versatile, "range N/M" for Ammunition or Thrown.
func splitProperty(s string) (name, detail string) {
	if i := strings.Index(s, " ("); i >= 0 {
		return s[:i], s[i+2 : len(s)-1]
	}
	return s, ""
}

// propertyByName maps an SRD weapon property's bare name (as splitProperty
// returns it) to its enum value. Every SRD 5.1 property name is known; the
// false case exists only as a defensive guard.
func propertyByName(name string) (WeaponProperty, bool) {
	switch name {
	case "ammunition":
		return PropertyAmmunition, true
	case "finesse":
		return PropertyFinesse, true
	case "heavy":
		return PropertyHeavy, true
	case "light":
		return PropertyLight, true
	case "loading":
		return PropertyLoading, true
	case "reach":
		return PropertyReach, true
	case "special":
		return PropertySpecial, true
	case "thrown":
		return PropertyThrown, true
	case "two-handed":
		return PropertyTwoHanded, true
	case "versatile":
		return PropertyVersatile, true
	default:
		return PropertyNone, false
	}
}

// countKnown counts the fixture properties propertyByName recognizes.
func countKnown(props []string) int {
	n := 0
	for _, p := range props {
		name, _ := splitProperty(p)
		if _, ok := propertyByName(name); ok {
			n++
		}
	}
	return n
}

// srdRange returns a weapon's normal and long range in feet, parsed from
// its Ammunition or Thrown property (the SRD never gives a weapon both);
// zero, zero for a weapon with neither.
func srdRange(props []string) (normal, long core.Distance) {
	for _, p := range props {
		name, detail := splitProperty(p)
		if name != "ammunition" && name != "thrown" {
			continue
		}
		detail = strings.TrimPrefix(detail, "range ")
		before, after, ok := strings.Cut(detail, "/")
		if !ok {
			panic(fmt.Sprintf("parsing SRD range %q", p))
		}
		n, errN := strconv.Atoi(before)
		l, errL := strconv.Atoi(after)
		if errN != nil || errL != nil {
			panic(fmt.Sprintf("parsing SRD range %q", p))
		}
		return core.Distance(n), core.Distance(l)
	}
	return 0, 0
}

// srdArmor is one armour entry as this fixture pins it, generated by
// content/testdata/fetch-srd-2014.py armor. The SRD's own fields (category,
// AC, Dex cap, Strength requirement, stealth, cost) come from Open5e's v1
// API (api.open5e.com/v1/armor/?document__slug=wotc-srd), which has no
// weight field; Weight comes from Open5e's v2 API (api.open5e.com/v2/items/,
// filtered to document key "srd-2014", category "armor" or "shield"),
// hand-checked against the SRD 5.1 armour table.
type srdArmor struct {
	Slug                string `json:"slug"`
	Name                string `json:"name"`
	Category            string `json:"category"`
	BaseAc              int    `json:"base_ac"`
	PlusDexMod          bool   `json:"plus_dex_mod"`
	PlusMax             int    `json:"plus_max"`
	PlusFlatMod         int    `json:"plus_flat_mod"`
	StrengthRequirement *int   `json:"strength_requirement"`
	StealthDisadvantage bool   `json:"stealth_disadvantage"`
	Cost                string `json:"cost"`
	Weight              string `json:"weight"`
}

// srdItem is one entry from Open5e's v2 items API (api.open5e.com/v2/items/),
// filtered to document key "srd-2014" and sorted by key: the full SRD 5.1
// item table, generated by content/testdata/fetch-srd-2014.py items. Cost
// and Weight are Open5e's raw decimal strings (gold pieces, pounds) for a
// single unit; ammunition is priced and weighed per unit here, not per the
// SRD's bundle, so the test scales them by Item.Quantity.
type srdItem struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Cost     string `json:"cost"`
	Weight   string `json:"weight"`
}

// byKey indexes a fixture's rows by Key.
func byKey(items []srdItem) map[string]srdItem {
	m := make(map[string]srdItem, len(items))
	for _, it := range items {
		m[it.Key] = it
	}
	return m
}

// itemCost parses Open5e's v2 decimal gold-piece string ("15.00", "0.05")
// into core.Coins, rounding to the nearest copper piece (the SRD's smallest
// denomination).
func itemCost(gp string) core.Coins {
	f, err := strconv.ParseFloat(gp, 64)
	if err != nil {
		panic(fmt.Sprintf("parsing item cost %q: %v", gp, err))
	}
	return core.Coins(math.Round(f * 100))
}

// itemWeight parses Open5e's v2 decimal pound-weight string ("1.000",
// "0.050") into core.Weight.
func itemWeight(w string) core.Weight {
	f, err := strconv.ParseFloat(w, 64)
	if err != nil {
		panic(fmt.Sprintf("parsing item weight %q: %v", w, err))
	}
	return core.Weight(f)
}

// gearSrdKey maps a Gear() item id to the srdItem.Key (in
// testdata/srd-2014-items.json) it corresponds to. Open5e's own product
// names don't follow the SRD's bundle-based ammunition naming (its key for
// "Arrows (20)" is "srd_arrow-bow"), so the mapping is explicit; it also
// doubles as the list of ids Gear() is required to ship, so
// TestGearMatchesTheSrd catches both a wrong value on an existing entry and
// an entry going missing.
var gearSrdKey = map[string]string{
	"arrow":          "srd_arrow-bow",
	"blowgun-needle": "srd_blowgun-needles",
	"crossbow-bolt":  "srd_crossbow-bolt",
	"sling-bullet":   "srd_sling-bullets",
	"thieves-tools":  "srd_thieves-tools",
	"rations":        "srd_rations-1-day",
	"torch":          "srd_torch",
	"rope-hempen":    "srd_rope-hempen-50-feet",
	"waterskin":      "srd_waterskin",
	"backpack":       "srd_backpack",
}

// srdGearName reduces a gear name to the form both sides of the pin can be
// compared in, because Open5e's product names and the SRD's table entries
// spell the same thing differently: case ("Thieves' Tools"), a trailing
// parenthetical that qualifies rather than names ("Arrow (bow)", the SRD's
// own bundle count in "Arrows (20)"), and the plural Open5e uses for some
// ammunition ("Sling bullets") where Item.Name is the singular unit and
// Item.Quantity carries the count.
func srdGearName(s string) string {
	if i := strings.Index(s, " ("); i >= 0 && strings.HasSuffix(s, ")") {
		s = s[:i]
	}
	return strings.TrimSuffix(strings.ToLower(strings.TrimSpace(s)), "s")
}

// srdItemCorrection documents one place Open5e v2's per-unit item data,
// scaled to the SRD's bundle, demonstrably departs from the SRD 5.1 text
// itself. The SRD text is authoritative: production (content/items_gear.go)
// carries the SRD's own value, and this is the single place documenting why
// it differs from the scaled fetch. wasCost/wasWeight (zero means "not
// corrected") record the scaled value this correction expects to still
// see; if a refreshed fetch no longer produces it, the correction is
// obsolete and TestGearMatchesTheSrd fails so it gets removed.
type srdItemCorrection struct {
	reason            string
	wasCost, cost     core.Coins
	wasWeight, weight core.Weight
}

var srdItemCorrections = map[string]srdItemCorrection{
	"crossbow-bolt": {
		reason:    `crossbow bolts: weight 1.5 lb for the bundle of 20 (SRD 5.1 "Crossbow bolts (20), 1 gp, 1 1/2 lb."); Open5e's per-unit weight (0.08 lb) scales to 1.6 lb for 20`,
		wasWeight: 1.6,
		weight:    1.5,
	},
	"sling-bullet": {
		reason:  `sling bullets: cost 4 cp for the bundle of 20 (SRD 5.1 "Sling bullets (20), 4 cp, 1 1/2 lb."); Open5e's per-unit cost (0.01 gp) scales to 20 cp for 20`,
		wasCost: core.Cp(20),
		cost:    core.Cp(4),
	},
}

// srdGearCostWeight returns the cost and weight the gear entry it must
// carry to match the SRD 5.1 text: raw.Cost/raw.Weight, which Open5e prices
// per unit, scaled by the entry's own Quantity (the SRD's bundle size, zero
// for a single thing), then corrected per srdItemCorrections. It fails t if
// a correction's recorded "was" value no longer matches the scaled figure,
// since that means the correction has become obsolete.
func srdGearCostWeight(t *testing.T, it Item, raw srdItem) (core.Coins, core.Weight) {
	t.Helper()
	id := it.Id
	bundle := it.Quantity
	if bundle == 0 {
		bundle = 1
	}
	cost := itemCost(raw.Cost) * core.Coins(bundle)
	weight := core.Weight(math.Round(float64(itemWeight(raw.Weight))*float64(bundle)*1000) / 1000)
	c, ok := srdItemCorrections[id]
	if !ok {
		return cost, weight
	}
	if c.wasCost != 0 {
		if cost != c.wasCost {
			t.Errorf("%s: correction %q is obsolete, scaled cost is now %s", id, c.reason, cost)
		}
		cost = c.cost
	}
	if c.wasWeight != 0 {
		if weight != c.wasWeight {
			t.Errorf("%s: correction %q is obsolete, scaled weight is now %v", id, c.reason, weight)
		}
		weight = c.weight
	}
	return cost, weight
}

func TestSrdCost(t *testing.T) {
	cases := []struct{ in, want string }{
		{"50 gp", "50 gp"},
		{"1 sp", "1 sp"},
		{"5 cp", "5 cp"},
	}
	for _, c := range cases {
		if got := srdCost(c.in); got != c.want {
			t.Errorf("srdCost(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Every SRD 5.1 weapon is present, once, and every pinned field matches the
// SRD 5.1 text as Open5e's v1 API serves it.
func TestWeaponsMatchTheSrd(t *testing.T) {
	fixture := loadFixture[srdWeapon](t, "srd-2014-weapons.json")
	got := byId(Weapons())
	if len(Weapons()) != len(fixture) || len(got) != len(fixture) {
		t.Errorf("Weapons() has %d entries (%d distinct ids), fixture has %d", len(Weapons()), len(got), len(fixture))
	}
	for _, w := range fixture {
		it, ok := got[w.Slug]
		if !ok {
			t.Errorf("%s: missing", w.Slug)
			continue
		}
		if it.Kind != ItemWeapon || it.Name != w.Name {
			t.Errorf("%s: kind/name %v %q", w.Slug, it.Kind, it.Name)
		}
		if it.Cost.String() != srdCost(w.Cost) {
			t.Errorf("%s: cost %s, SRD %s", w.Slug, it.Cost, srdCost(w.Cost))
		}
		if it.Weight != srdWeight(w.Weight) {
			t.Errorf("%s: weight %v, SRD %v", w.Slug, it.Weight, srdWeight(w.Weight))
		}
		if !srdHasDamage(w.DamageDice) {
			if it.Damage != nil {
				t.Errorf("%s: has damage, SRD has none", w.Slug)
			}
		} else if it.Damage == nil || len(it.Damage.Parts) != 1 || it.Damage.Parts[0].Dice.String() != srdDiceExpr(w.DamageDice) || !strings.EqualFold(it.Damage.Parts[0].Type.Id, w.DamageType) {
			t.Errorf("%s: damage %+v, SRD %s %s", w.Slug, it.Damage, w.DamageDice, w.DamageType)
		}
		if it.WeaponCategory != srdCategory(w.Category) {
			t.Errorf("%s: category %v, SRD %s", w.Slug, it.WeaponCategory, w.Category)
		}
		wantRange, wantLongRange := srdRange(w.Properties)
		if it.Range != wantRange || it.LongRange != wantLongRange {
			t.Errorf("%s: range %v/%v, SRD %v/%v", w.Slug, it.Range, it.LongRange, wantRange, wantLongRange)
		}
		for _, raw := range w.Properties {
			name, detail := splitProperty(raw)
			prop, ok := propertyByName(name)
			if !ok {
				continue // defensive: every SRD 5.1 property name is known
			}
			if !it.HasProperty(prop) {
				t.Errorf("%s: missing property %s", w.Slug, raw)
			}
			if prop == PropertyVersatile && (it.VersatileDamage == nil || it.VersatileDamage.Parts[0].Dice.String() != detail) {
				t.Errorf("%s: versatile damage %+v, SRD %s", w.Slug, it.VersatileDamage, detail)
			}
		}
		if len(it.Properties) != countKnown(w.Properties) {
			t.Errorf("%s: %d properties, SRD %d", w.Slug, len(it.Properties), countKnown(w.Properties))
		}
		// The SRD's rule: a weapon with the Ammunition property is not a
		// melee weapon; every other weapon (including a thrown one) is.
		if it.Melee == it.HasProperty(PropertyAmmunition) {
			t.Errorf("%s: melee=%v is inconsistent with ammunition=%v", w.Slug, it.Melee, it.HasProperty(PropertyAmmunition))
		}
	}
}

// armorDexCap derives Item.MaxDexBonus from the SRD's own Dex-mod fields:
// -1 for an uncapped Dex bonus (light armour), the printed cap for a capped
// one (medium armour), 0 for none (heavy armour).
func armorDexCap(a srdArmor) int {
	if !a.PlusDexMod {
		return 0
	}
	if a.PlusMax > 0 {
		return a.PlusMax
	}
	return -1
}

// Every SRD 5.1 armour entry and the shield are present, once, and every
// pinned field matches the SRD 5.1 text.
func TestArmorMatchesTheSrd(t *testing.T) {
	want := loadFixture[srdArmor](t, "srd-2014-armor.json")
	got := byId(Armor())
	if len(Armor()) != len(want) || len(got) != len(want) {
		t.Errorf("Armor() has %d entries (%d distinct ids), fixture has %d", len(Armor()), len(got), len(want))
	}
	categories := map[string]ArmorCategory{"Light Armor": ArmorCategoryLight, "Medium Armor": ArmorCategoryMedium, "Heavy Armor": ArmorCategoryHeavy}
	for _, a := range want {
		it, ok := got[a.Slug]
		if !ok {
			t.Errorf("%s: missing", a.Slug)
			continue
		}
		wantKind, wantCat, wantAc := ItemArmor, categories[a.Category], a.BaseAc
		if a.Category == "Shield" {
			// The SRD gives the shield a flat AC bonus, not a base AC
			// (Open5e's ac_string is "0 +2"); PlusFlatMod carries it.
			wantKind, wantCat, wantAc = ItemShield, ArmorCategoryShield, a.PlusFlatMod
		}
		if it.Kind != wantKind || it.ArmorCategory != wantCat || it.Name != a.Name {
			t.Errorf("%s: kind/category/name %v %v %q, SRD %v %v %q", a.Slug, it.Kind, it.ArmorCategory, it.Name, wantKind, wantCat, a.Name)
		}
		if int(it.ArmorClass) != wantAc {
			t.Errorf("%s: AC %d, SRD %d", a.Slug, it.ArmorClass, wantAc)
		}
		if wantCap := armorDexCap(a); it.MaxDexBonus != wantCap {
			t.Errorf("%s: Dex cap %d, SRD %d", a.Slug, it.MaxDexBonus, wantCap)
		}
		wantStr := 0
		if a.StrengthRequirement != nil {
			wantStr = *a.StrengthRequirement
		}
		if int(it.StrengthRequired) != wantStr || it.StealthDisadvantage != a.StealthDisadvantage {
			t.Errorf("%s: Str %d stealth %v, SRD %d %v", a.Slug, it.StrengthRequired, it.StealthDisadvantage, wantStr, a.StealthDisadvantage)
		}
		if it.Cost.String() != srdCost(a.Cost) || it.Weight != srdWeight(a.Weight) {
			t.Errorf("%s: cost/weight %s/%v, SRD %s/%v", a.Slug, it.Cost, it.Weight, srdCost(a.Cost), srdWeight(a.Weight))
		}
	}
}

// wantGearCategory maps Gear()'s ItemKind to the category string
// testdata/srd-2014-items.json uses for it.
var wantGearCategory = map[ItemKind]string{ItemAmmunition: "ammunition", ItemTool: "tools", ItemGear: "adventuring-gear"}

// Gear is a chosen subset of the full SRD 5.1 item table (gearSrdKey is the
// list of ids it must ship); every entry's cost, weight and category are
// pinned to testdata/srd-2014-items.json, scaled from Open5e's per-unit
// data to the SRD's bundle and corrected where that scaling demonstrably
// departs from the SRD 5.1 text (srdItemCorrections). Both directions are
// checked: a Gear() entry with no fixture backing, and a fixture entry
// Gear() is meant to ship but no longer does, so deleting an entry fails
// rather than just silently dropping its correction.
func TestGearMatchesTheSrd(t *testing.T) {
	rows := byKey(loadFixture[srdItem](t, "srd-2014-items.json"))
	got := byId(Gear())
	if len(Gear()) != len(gearSrdKey) || len(got) != len(gearSrdKey) {
		t.Errorf("Gear() has %d entries (%d distinct ids), want %d", len(Gear()), len(got), len(gearSrdKey))
	}
	for id, key := range gearSrdKey {
		raw, ok := rows[key]
		if !ok {
			t.Fatalf("%s: fixture has no %s; regenerate testdata/srd-2014-items.json", id, key)
		}
		it, ok := got[id]
		if !ok {
			t.Errorf("%s: Gear() is missing this SRD item (%s)", id, key)
			continue
		}
		if raw.Category != wantGearCategory[it.Kind] {
			t.Errorf("%s: kind %v, SRD category %s name %q", id, it.Kind, raw.Category, raw.Name)
		}
		if srdGearName(it.Name) != srdGearName(raw.Name) {
			t.Errorf("%s: name %q, SRD %q", id, it.Name, raw.Name)
		}
		wantCost, wantWeight := srdGearCostWeight(t, it, raw)
		if it.Cost != wantCost || it.Weight != wantWeight {
			t.Errorf("%s: cost/weight %s/%v, SRD %s/%v", id, it.Cost, it.Weight, wantCost, wantWeight)
		}
		if (it.Kind == ItemAmmunition) != it.Stackable {
			t.Errorf("%s: ammunition stacks, other gear does not", id)
		}
		// The SRD sells ammunition by the bundle and everything else here
		// singly, so Quantity is set for ammunition and zero otherwise.
		if (it.Kind == ItemAmmunition) != (it.Quantity > 0) {
			t.Errorf("%s: quantity %d for kind %v", id, it.Quantity, it.Kind)
		}
	}
}

// StandardItems ships every id once, with kind set, in weapons/armour/gear
// order, and RegisterStandardItems puts them all in a registry.
func TestStandardItemsAreUniqueAndRegister(t *testing.T) {
	all := StandardItems()
	seen := map[string]bool{}
	for _, it := range all {
		if it.Id == "" || it.Kind == ItemUnspecified {
			t.Errorf("%+v: id and kind are required", it)
		}
		if seen[it.Id] {
			t.Errorf("%s: duplicate id", it.Id)
		}
		seen[it.Id] = true
	}

	var wantOrder []Item
	wantOrder = append(wantOrder, Weapons()...)
	wantOrder = append(wantOrder, Armor()...)
	wantOrder = append(wantOrder, Gear()...)
	if len(all) != len(wantOrder) {
		t.Fatalf("StandardItems() has %d entries, want %d", len(all), len(wantOrder))
	}
	for i, it := range wantOrder {
		if all[i].Id != it.Id {
			t.Errorf("StandardItems()[%d] = %q, want %q (weapons, then armour, then gear)", i, all[i].Id, it.Id)
		}
	}

	r := NewRegistry[Item]()
	RegisterStandardItems(r)
	if len(r.All()) != len(all) {
		t.Errorf("registry has %d, StandardItems %d", len(r.All()), len(all))
	}
	if _, ok := r.Get("longsword"); !ok {
		t.Error("longsword not registered")
	}
}
