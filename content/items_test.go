package content

import (
	"encoding/json"
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
)

func TestItemEnumZeroValues(t *testing.T) {
	if ItemUnspecified != 0 || PropertyNone != 0 || WeaponCategoryNone != 0 || ArmorCategoryNone != 0 {
		t.Error("item enums must have explicit zero values")
	}
}

func TestItemRoundTripsThroughJson(t *testing.T) {
	d8, _ := dice.Parse("1d8")
	d10, _ := dice.Parse("1d10")
	in := Item{
		Id: "longsword", Name: "Longsword", Kind: ItemWeapon,
		Weight: 3, Cost: core.Gp(15),
		Damage:          &damage.Spec{Parts: []damage.PartSpec{{Dice: d8, Type: core.Slashing}}},
		VersatileDamage: &damage.Spec{Parts: []damage.PartSpec{{Dice: d10, Type: core.Slashing}}},
		Properties:      []WeaponProperty{PropertyVersatile},
		WeaponCategory:  WeaponCategoryMartial, Melee: true,
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out Item
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Id != in.Id || out.Kind != in.Kind || out.Cost != in.Cost || len(out.Properties) != 1 || out.Damage == nil || out.VersatileDamage == nil {
		t.Errorf("round trip lost fields: %+v", out)
	}
}

func TestItemHasProperty(t *testing.T) {
	it := Item{Properties: []WeaponProperty{PropertyLight, PropertyFinesse}}
	if !it.HasProperty(PropertyFinesse) || it.HasProperty(PropertyHeavy) {
		t.Error("HasProperty wrong")
	}
}
