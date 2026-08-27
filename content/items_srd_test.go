package content

import (
	"encoding/json"
	"fmt"
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
