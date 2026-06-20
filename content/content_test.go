package content

import (
	"encoding/json"
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/damage"
	"github.com/trancecode/go-srd5e/dice"
	"github.com/trancecode/go-srd5e/effect"
)

func TestEnumZeroValues(t *testing.T) {
	if ResolveUnspecified != 0 || RangeUnspecified != 0 || TargetUnspecified != 0 || AreaNone != 0 || ItemUnspecified != 0 {
		t.Error("zero-value enum sentinels must be 0")
	}
}

func TestClassAndRace(t *testing.T) {
	c := Class{
		Id: "wizard", Name: "Wizard", HitDie: 6,
		ProficientSaves:     []core.Ability{core.AbilityIntelligence, core.AbilityWisdom},
		SpellcastingAbility: core.AbilityIntelligence,
		Slots:               SpellSlotProgression{1: SpellSlots{0, 2}, 2: SpellSlots{0, 3}},
	}
	if c.Slots[1][1] != 2 || c.SpellcastingAbility != core.AbilityIntelligence {
		t.Errorf("class wiring wrong: %+v", c)
	}
	// non-caster: zero SpellcastingAbility.
	fighter := Class{Id: "fighter", HitDie: 10}
	if fighter.SpellcastingAbility != core.AbilityNone {
		t.Error("non-caster should have zero (AbilityNone) spellcasting ability")
	}
	// race with optional bonuses absent (SRD 5.2 style).
	r := Race{Id: "human", MovementSpeed: 30}
	if r.AbilityBonuses != nil || r.MovementSpeed != 30 {
		t.Error("race optional bonuses should be nil-able")
	}
}

func TestSpellJSONRoundTrip(t *testing.T) {
	fire := damage.Single(dice.Expr{Count: 8, Sides: 6}, core.Fire)
	half := true
	sp := Spell{
		Id: "fireball", Name: "Fireball", Level: 3,
		Resolution:  ResolveSave,
		SaveAbility: core.AbilityDexterity,
		Targeting:   Targeting{Range: RangeRanged, Distance: 150, Target: TargetArea, Shape: AreaSphere, AreaSize: 20},
		Effects: []effect.ConditionalEffect{
			{Trigger: effect.OnSave, Effect: effect.Effect{Damage: &fire, HalfOnSave: half}},
		},
	}
	data, err := json.Marshal(sp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back Spell
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.Level != 3 || back.Resolution != ResolveSave || back.Targeting.Shape != AreaSphere {
		t.Errorf("round-trip scalar fields wrong: %+v", back)
	}
	if len(back.Effects) != 1 || !back.Effects[0].Effect.HalfOnSave || back.Effects[0].Effect.Damage == nil {
		t.Fatalf("round-trip effects wrong: %+v", back.Effects)
	}
	if back.Effects[0].Effect.Damage.Parts[0].Type != core.Fire {
		t.Error("round-trip nested damage type wrong")
	}
}

func TestCreatureAndItem(t *testing.T) {
	cr := Creature{
		Id: "goblin", Name: "Goblin", Size: core.SizeSmall,
		Abilities:     core.AbilityScores{Dexterity: 14},
		ArmorClass:    15,
		HitDice:       dice.Expr{Count: 2, Sides: 6},
		Mitigation:    damage.Mitigation{Resist: map[string]bool{core.Fire.Id: true}},
		MovementSpeed: 30,
	}
	if cr.HitDice.Average() != 7 || !cr.Mitigation.Resist[core.Fire.Id] {
		t.Errorf("creature wiring wrong: %+v", cr)
	}
	sword := damage.Single(dice.Expr{Count: 1, Sides: 8}, core.Slashing)
	it := Item{Id: "longsword", Name: "Longsword", Kind: ItemWeapon, Damage: &sword}
	if it.Kind != ItemWeapon || it.Damage == nil {
		t.Error("weapon item wiring wrong")
	}
}
