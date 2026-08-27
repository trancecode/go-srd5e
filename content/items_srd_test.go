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

// srdProperty is one weapon property as Open5e's v2 API serves it.
type srdProperty struct {
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

type srdWeapon struct {
	Key        string        `json:"key"`
	Name       string        `json:"name"`
	Cost       string        `json:"cost"`
	Weight     string        `json:"weight"`
	DamageDice string        `json:"damage_dice"`
	DamageType string        `json:"damage_type"`
	IsSimple   bool          `json:"is_simple"`
	IsMartial  bool          `json:"is_martial"`
	Range      float64       `json:"range"`
	LongRange  float64       `json:"long_range"`
	Properties []srdProperty `json:"properties"`
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

// idOf turns an Open5e key ("srd_light-crossbow") into the content id.
func idOf(key string) string { return strings.TrimPrefix(key, "srd_") }

func byId(items []Item) map[string]Item {
	m := make(map[string]Item, len(items))
	for _, it := range items {
		m[it.Id] = it
	}
	return m
}

// srdCorrection is one place the SRD 5.1 weapon table, as Open5e's v2 API
// currently serves it, demonstrably departs from the SRD 5.1 text itself.
// The SRD text is the authority; Open5e only serves it. apply corrects one
// fetched fixture row in place before it is compared against production, so
// production (content/items_weapons.go) carries the SRD's own value and this
// table is the single place documenting why it differs from the raw fetch
// kept, untouched, in testdata/srd-2014-weapons.json.
type srdCorrection struct {
	key    string
	reason string
	apply  func(*srdWeapon)
}

var srdCorrections = []srdCorrection{
	{
		key:    "srd_crossbow-heavy",
		reason: `heavy-crossbow: add property Ammunition (SRD 5.1: "Ammunition (range 100/400), heavy, loading, two-handed")`,
		apply: func(w *srdWeapon) {
			w.Properties = append(w.Properties, srdProperty{Name: "Ammunition", Detail: "range 100/400"})
		},
	},
	{
		key:    "srd_handaxe",
		reason: `handaxe: add property Thrown (SRD 5.1: "Light, thrown (range 20/60)")`,
		apply: func(w *srdWeapon) {
			w.Properties = append(w.Properties, srdProperty{Name: "Thrown", Detail: "range 20/60"})
		},
	},
	{
		key:    "srd_greatsword",
		reason: `greatsword: weight 6 (SRD 5.1: 6 lb.)`,
		apply: func(w *srdWeapon) {
			w.Weight = "6.000"
		},
	},
}

// applyCorrections returns w with any known Open5e departure from the SRD
// 5.1 text corrected, so the rest of TestWeaponsMatchTheSrd compares
// production against the SRD, not against a gap in how Open5e serves it.
func applyCorrections(w srdWeapon) srdWeapon {
	for _, c := range srdCorrections {
		if c.key == w.Key {
			c.apply(&w)
		}
	}
	return w
}

// srdCost parses Open5e's decimal gold amount ("15.00", "0.10", "0.02") and
// renders it the way core.Coins does ("15 gp", "1 sp", "2 cp").
func srdCost(gp string) string {
	f, err := strconv.ParseFloat(gp, 64)
	if err != nil {
		panic(fmt.Sprintf("parsing SRD cost %q: %v", gp, err))
	}
	// Round rather than truncate: gold amounts like 0.05 are not exactly
	// representable in binary floating point.
	return core.Coins(math.Round(f * 100)).String()
}

// srdWeight parses Open5e's decimal pound weight ("3.000", "0.250") into
// core.Weight.
func srdWeight(w string) core.Weight {
	f, err := strconv.ParseFloat(w, 64)
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

// propertyByName maps an SRD weapon property's display name to its enum
// value. Open5e suffixes "Special" with the weapon it names ("Special
// (Lance)", "Special (Net)"); the part before the parenthesis is the
// property. Properties outside the SRD 5.1 table (weapon mastery and
// later additions) report false.
func propertyByName(name string) (WeaponProperty, bool) {
	if i := strings.Index(name, " ("); i >= 0 {
		name = name[:i]
	}
	switch name {
	case "Ammunition":
		return PropertyAmmunition, true
	case "Finesse":
		return PropertyFinesse, true
	case "Heavy":
		return PropertyHeavy, true
	case "Light":
		return PropertyLight, true
	case "Loading":
		return PropertyLoading, true
	case "Reach":
		return PropertyReach, true
	case "Special":
		return PropertySpecial, true
	case "Thrown":
		return PropertyThrown, true
	case "Two-Handed":
		return PropertyTwoHanded, true
	case "Versatile":
		return PropertyVersatile, true
	default:
		return PropertyNone, false
	}
}

// countKnown counts the fixture properties propertyByName recognizes.
func countKnown(props []srdProperty) int {
	n := 0
	for _, p := range props {
		if _, ok := propertyByName(p.Name); ok {
			n++
		}
	}
	return n
}

func TestSrdCost(t *testing.T) {
	cases := []struct{ in, want string }{
		{"15.00", "15 gp"},
		{"0.10", "1 sp"},
		{"0.02", "2 cp"},
	}
	for _, c := range cases {
		if got := srdCost(c.in); got != c.want {
			t.Errorf("srdCost(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// Every SRD 5.1 weapon is present, once, and every pinned field matches the
// SRD 5.1 text, corrected (via srdCorrections, above) for the few places
// Open5e's fetched data demonstrably departs from that text.
func TestWeaponsMatchTheSrd(t *testing.T) {
	fixture := loadFixture[srdWeapon](t, "srd-2014-weapons.json")
	got := byId(Weapons())
	if len(Weapons()) != len(fixture) || len(got) != len(fixture) {
		t.Errorf("Weapons() has %d entries (%d distinct ids), fixture has %d", len(Weapons()), len(got), len(fixture))
	}
	for _, raw := range fixture {
		w := applyCorrections(raw)
		it, ok := got[idOf(w.Key)]
		if !ok {
			t.Errorf("%s: missing", w.Key)
			continue
		}
		if it.Kind != ItemWeapon || it.Name != w.Name {
			t.Errorf("%s: kind/name %v %q", w.Key, it.Kind, it.Name)
		}
		if it.Cost.String() != srdCost(w.Cost) {
			t.Errorf("%s: cost %s, SRD %s", w.Key, it.Cost, srdCost(w.Cost))
		}
		if it.Weight != srdWeight(w.Weight) {
			t.Errorf("%s: weight %v, SRD %v", w.Key, it.Weight, srdWeight(w.Weight))
		}
		if !srdHasDamage(w.DamageDice) {
			if it.Damage != nil {
				t.Errorf("%s: has damage, SRD has none", w.Key)
			}
		} else if it.Damage == nil || len(it.Damage.Parts) != 1 || it.Damage.Parts[0].Dice.String() != srdDiceExpr(w.DamageDice) || !strings.EqualFold(it.Damage.Parts[0].Type.Id, w.DamageType) {
			t.Errorf("%s: damage %+v, SRD %s %s", w.Key, it.Damage, w.DamageDice, w.DamageType)
		}
		if (it.WeaponCategory == WeaponCategorySimple) != w.IsSimple || (it.WeaponCategory == WeaponCategoryMartial) != w.IsMartial {
			t.Errorf("%s: category %v, SRD simple=%v martial=%v", w.Key, it.WeaponCategory, w.IsSimple, w.IsMartial)
		}
		if float64(it.Range) != w.Range || float64(it.LongRange) != w.LongRange {
			t.Errorf("%s: range %v/%v, SRD %v/%v", w.Key, it.Range, it.LongRange, w.Range, w.LongRange)
		}
		for _, p := range w.Properties {
			prop, ok := propertyByName(p.Name)
			if !ok {
				continue // mastery and other non-5.1 properties do not appear in the 5.1 fixture, but guard anyway
			}
			if !it.HasProperty(prop) {
				t.Errorf("%s: missing property %s", w.Key, p.Name)
			}
			if prop == PropertyVersatile && (it.VersatileDamage == nil || it.VersatileDamage.Parts[0].Dice.String() != p.Detail) {
				t.Errorf("%s: versatile damage %+v, SRD %s", w.Key, it.VersatileDamage, p.Detail)
			}
		}
		if len(it.Properties) != countKnown(w.Properties) {
			t.Errorf("%s: %d properties, SRD %d", w.Key, len(it.Properties), countKnown(w.Properties))
		}
		// The SRD's rule: a weapon with the Ammunition property is not a
		// melee weapon; every other weapon (including a thrown one) is.
		if it.Melee == it.HasProperty(PropertyAmmunition) {
			t.Errorf("%s: melee=%v is inconsistent with ammunition=%v", w.Key, it.Melee, it.HasProperty(PropertyAmmunition))
		}
	}
}
