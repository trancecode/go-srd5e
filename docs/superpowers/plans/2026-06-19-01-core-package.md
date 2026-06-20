# core package implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `core` package: the SRD vocabulary, named units, open value catalogs, and pure stat math that every other go-srd5e package imports.

**Architecture:** `core` is a dependency-free leaf. Pure functions and value types over named-int units, fully deterministic (no randomness, no I/O). Each concern lives in its own file with table-driven tests.

**Tech Stack:** Go 1.25, standard library only, `go test` with table-driven tests.

This is plan 1 of 8 (see `2026-06-19_srd-module-adoption-findings.md` and the design spec `docs/superpowers/specs/2026-06-18-go-srd5e-module-design.md`). It also does the one-time repo scaffolding.

## Global Constraints

Copied verbatim from the spec; every task inherits these.

- Module path: `github.com/trancecode/go-srd5e`. Go version `1.25`.
- No Go files at the repository root; every package is its own directory with a `doc.go`.
- Acronyms use Java camel case, not Go all-caps initialisms: `Id`, `Hp`, `Ac`, `Dc`, `Xp` (per `~/src/nrg/doc/styleguide.md`).
- Enum zero values are explicit: a genuine neutral is named `None` (e.g. `VisionNormal`, `AbilityNone`) and is valid; a must-set enum reserves a `…Unspecified` zero; function-only enums keep their zero as the safe member.
- Value types have exported fields and round-trip through `encoding/json` with no custom marshaler.
- Ability modifier rounds **down** (floor), not toward zero: `AbilityModifier(7) == -2`.
- TDD: write the failing test first, watch it fail, implement minimally, watch it pass, commit. Table-driven tests. Frequent commits. Document every exported type and func starting with its name.

---

### Task 1: Repository scaffolding

**Files:**
- Create: `go.mod`
- Create: `core/doc.go`
- Create: `CONSUMING.md`

**Interfaces:**
- Consumes: nothing.
- Produces: a buildable module rooted at `github.com/trancecode/go-srd5e` with an empty `core` package.

- [ ] **Step 1: Create `go.mod`**

```
module github.com/trancecode/go-srd5e

go 1.25
```

- [ ] **Step 2: Create `core/doc.go`**

```go
// Package core holds the SRD 5e vocabulary, the named-unit types, the open
// value catalogs (Ability, Skill, Condition, DamageType), and the pure stat
// math that every other go-srd5e package builds on. It is a dependency-free
// leaf with no randomness and no I/O.
package core
```

- [ ] **Step 3: Verify it builds**

Run: `go build ./...`
Expected: no output, exit 0.

- [ ] **Step 4: Create `CONSUMING.md`**

```markdown
# Consuming go-srd5e

This module is private under the trancecode organization.

```bash
go env -w GOPRIVATE=github.com/trancecode/*
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

For local development across games, add to the consuming game's `go.mod`:

```
replace github.com/trancecode/go-srd5e => ../go-srd5e
```

Drop the replace and pin a tagged version once the module is stable.
```

- [ ] **Step 5: Commit**

```bash
git add go.mod core/doc.go CONSUMING.md
git commit -m "scaffold module and core package"
```

---

### Task 2: Units and abilities

**Files:**
- Create: `core/ability.go`
- Test: `core/ability_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `type AbilityScore int`, `Modifier int`, `Level int`, `ArmorClass int`, `Dc int`, `Distance int`, `HitPoints int`, `Xp int`, `Weight int`
  - `type Ability int` with `AbilityNone Ability = 0`, `AbilityStrength`, `AbilityDexterity`, `AbilityConstitution`, `AbilityIntelligence`, `AbilityWisdom`, `AbilityCharisma`, and `AbilityAny Ability = -1`
  - `type AbilityScores struct { Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma AbilityScore }`
  - `func AbilityModifier(score AbilityScore) Modifier` (floored)
  - `func (s AbilityScore) Modifier() Modifier`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestAbilityModifier(t *testing.T) {
	cases := []struct {
		score AbilityScore
		want  Modifier
	}{
		{1, -5}, {6, -2}, {7, -2}, {8, -1}, {9, -1}, {10, 0},
		{11, 0}, {12, 1}, {15, 2}, {20, 5}, {30, 10},
	}
	for _, c := range cases {
		if got := AbilityModifier(c.score); got != c.want {
			t.Errorf("AbilityModifier(%d) = %d, want %d", c.score, got, c.want)
		}
		if got := c.score.Modifier(); got != c.want {
			t.Errorf("AbilityScore(%d).Modifier() = %d, want %d", c.score, got, c.want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run TestAbilityModifier`
Expected: FAIL (undefined: AbilityModifier, AbilityScore).

- [ ] **Step 3: Write minimal implementation**

```go
package core

// Named units. Distinct SRD quantities are named types so the compiler rejects
// argument swaps; untyped constant literals still convert freely.
type (
	AbilityScore int // 1..30
	Modifier     int // a value added to a d20
	Level        int // 1..20
	ArmorClass   int
	Dc           int // difficulty class
	Distance     int // feet
	HitPoints    int
	Xp           int
	Weight       int // pounds
)

// Ability is the closed set of the six SRD abilities. AbilityAny is a wildcard
// for qualifiers only (e.g. a buff to all saves); never use it where a concrete
// ability is required.
type Ability int

const (
	AbilityNone Ability = iota
	AbilityStrength
	AbilityDexterity
	AbilityConstitution
	AbilityIntelligence
	AbilityWisdom
	AbilityCharisma
)

const AbilityAny Ability = -1

// AbilityScores is the six ability scores.
type AbilityScores struct {
	Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma AbilityScore
}

// AbilityModifier is (score-10)/2 rounded down.
func AbilityModifier(score AbilityScore) Modifier {
	d := int(score) - 10
	if d >= 0 {
		return Modifier(d / 2)
	}
	return Modifier(-((-d + 1) / 2))
}

// Modifier returns the score's ability modifier.
func (s AbilityScore) Modifier() Modifier { return AbilityModifier(s) }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run TestAbilityModifier`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add core/ability.go core/ability_test.go
git commit -m "core: named units and ability modifier"
```

---

### Task 3: Stat math

**Files:**
- Create: `core/stats.go`
- Create: `core/xp.go`
- Test: `core/stats_test.go`

**Interfaces:**
- Consumes: `AbilityScore`, `Modifier`, `Level`, `Dc`, `HitPoints`, `Xp`, `AbilityModifier` (Task 2).
- Produces:
  - `func ProficiencyBonus(level Level) Modifier`
  - `func SkillBonus(score AbilityScore, proficient bool, level Level) Modifier`
  - `func SavingThrowBonus(score AbilityScore, proficient bool, level Level) Modifier`
  - `func SpellSaveDc(score AbilityScore, level Level) Dc`
  - `func SpellAttackBonus(score AbilityScore, level Level) Modifier`
  - `func MaxHp(hpRolls []int, conMod Modifier) HitPoints`
  - `func XpForLevel(level Level) Xp`
  - `func XpForNextLevel(current Level) Xp`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestProficiencyBonus(t *testing.T) {
	cases := []struct {
		level Level
		want  Modifier
	}{{1, 2}, {4, 2}, {5, 3}, {8, 3}, {9, 4}, {13, 5}, {17, 6}, {20, 6}}
	for _, c := range cases {
		if got := ProficiencyBonus(c.level); got != c.want {
			t.Errorf("ProficiencyBonus(%d) = %d, want %d", c.level, got, c.want)
		}
	}
}

func TestBonuses(t *testing.T) {
	// STR 16 (+3) at level 5 (prof +3): skill proficient = 6, non-proficient = 3.
	if got := SkillBonus(16, true, 5); got != 6 {
		t.Errorf("SkillBonus proficient = %d, want 6", got)
	}
	if got := SkillBonus(16, false, 5); got != 3 {
		t.Errorf("SkillBonus non-proficient = %d, want 3", got)
	}
	if got := SavingThrowBonus(16, true, 5); got != 6 {
		t.Errorf("SavingThrowBonus = %d, want 6", got)
	}
	// INT 16 (+3) at level 5 (prof +3): spell save DC = 8+3+3 = 14, attack = 6.
	if got := SpellSaveDc(16, 5); got != 14 {
		t.Errorf("SpellSaveDc = %d, want 14", got)
	}
	if got := SpellAttackBonus(16, 5); got != 6 {
		t.Errorf("SpellAttackBonus = %d, want 6", got)
	}
}

func TestMaxHp(t *testing.T) {
	// three levels of rolls 6,5,5 with CON +2 = 16 + 6 = 22.
	if got := MaxHp([]int{6, 5, 5}, 2); got != 22 {
		t.Errorf("MaxHp = %d, want 22", got)
	}
}

func TestXp(t *testing.T) {
	if got := XpForLevel(1); got != 0 {
		t.Errorf("XpForLevel(1) = %d, want 0", got)
	}
	if got := XpForLevel(5); got != 6500 {
		t.Errorf("XpForLevel(5) = %d, want 6500", got)
	}
	if got := XpForLevel(20); got != 355000 {
		t.Errorf("XpForLevel(20) = %d, want 355000", got)
	}
	if got := XpForNextLevel(4); got != 6500 {
		t.Errorf("XpForNextLevel(4) = %d, want 6500", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run 'TestProficiencyBonus|TestBonuses|TestMaxHp|TestXp'`
Expected: FAIL (undefined functions).

- [ ] **Step 3: Write `core/stats.go`**

```go
package core

// ProficiencyBonus is 2 + (level-1)/4.
func ProficiencyBonus(level Level) Modifier { return Modifier(2 + (int(level)-1)/4) }

// SkillBonus is the ability modifier plus the proficiency bonus when proficient.
func SkillBonus(score AbilityScore, proficient bool, level Level) Modifier {
	return abilityCheckBonus(score, proficient, level)
}

// SavingThrowBonus is the ability modifier plus the proficiency bonus when proficient.
func SavingThrowBonus(score AbilityScore, proficient bool, level Level) Modifier {
	return abilityCheckBonus(score, proficient, level)
}

func abilityCheckBonus(score AbilityScore, proficient bool, level Level) Modifier {
	b := AbilityModifier(score)
	if proficient {
		b += ProficiencyBonus(level)
	}
	return b
}

// SpellSaveDc is 8 + proficiency bonus + spellcasting ability modifier.
func SpellSaveDc(score AbilityScore, level Level) Dc {
	return Dc(8 + int(ProficiencyBonus(level)) + int(AbilityModifier(score)))
}

// SpellAttackBonus is proficiency bonus + spellcasting ability modifier.
func SpellAttackBonus(score AbilityScore, level Level) Modifier {
	return ProficiencyBonus(level) + AbilityModifier(score)
}

// MaxHp sums the per-level hit-point rolls and adds the Constitution modifier
// once per level.
func MaxHp(hpRolls []int, conMod Modifier) HitPoints {
	total := 0
	for _, r := range hpRolls {
		total += r
	}
	return HitPoints(total + int(conMod)*len(hpRolls))
}
```

- [ ] **Step 4: Write `core/xp.go`**

```go
package core

// xpThresholds is the cumulative XP required to reach each level (index 1..20).
var xpThresholds = [21]Xp{
	0, 0, 300, 900, 2700, 6500, 14000, 23000, 34000, 48000, 64000,
	85000, 100000, 120000, 140000, 165000, 195000, 225000, 265000, 305000, 355000,
}

// XpForLevel is the cumulative XP needed to reach the given level (1..20).
func XpForLevel(level Level) Xp { return xpThresholds[level] }

// XpForNextLevel is the cumulative XP needed to reach the level after the
// current one. At level 20 it returns the level-20 threshold.
func XpForNextLevel(current Level) Xp {
	if current >= 20 {
		return xpThresholds[20]
	}
	return xpThresholds[current+1]
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./core/ -run 'TestProficiencyBonus|TestBonuses|TestMaxHp|TestXp'`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add core/stats.go core/xp.go core/stats_test.go
git commit -m "core: stat math and XP thresholds"
```

---

### Task 4: Open value catalogs

**Files:**
- Create: `core/catalog.go`
- Test: `core/catalog_test.go`

**Interfaces:**
- Consumes: `Ability` and its constants (Task 2).
- Produces:
  - `type Skill struct { Id, Name string; Ability Ability }`
  - `type Condition struct { Id, Name string }`
  - `type DamageType struct { Id, Name string }`
  - predefined `Skill` values for the eighteen SRD skills (e.g. `SkillStealth`, `SkillArcana`), and `var SRDSkills []Skill`
  - predefined `DamageType` values (e.g. `Slashing`, `Piercing`, `Bludgeoning`, `Fire`) and `var DamageAny = DamageType{Id: "any", Name: "Any"}`
  - predefined `Condition` values (e.g. `Prone`, `Poisoned`)

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestSRDSkills(t *testing.T) {
	if len(SRDSkills) != 18 {
		t.Fatalf("SRDSkills has %d entries, want 18", len(SRDSkills))
	}
	if SkillStealth.Ability != AbilityDexterity {
		t.Errorf("Stealth ability = %v, want Dexterity", SkillStealth.Ability)
	}
	if SkillArcana.Ability != AbilityIntelligence {
		t.Errorf("Arcana ability = %v, want Intelligence", SkillArcana.Ability)
	}
}

func TestDamageAny(t *testing.T) {
	if DamageAny.Id != "any" {
		t.Errorf("DamageAny.Id = %q, want \"any\"", DamageAny.Id)
	}
	if Slashing.Id != "slashing" {
		t.Errorf("Slashing.Id = %q, want \"slashing\"", Slashing.Id)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run 'TestSRDSkills|TestDamageAny'`
Expected: FAIL (undefined: SRDSkills, SkillStealth, ...).

- [ ] **Step 3: Write minimal implementation**

```go
package core

// Skill is a named proficiency governed by an ability. It is an open value:
// the eighteen SRD skills ship as predefined values and a game declares more.
type Skill struct {
	Id, Name string
	Ability  Ability
}

// Condition is an open named state (prone, poisoned, plus game-defined ones).
type Condition struct{ Id, Name string }

// DamageType is an open damage category. Games add their own (e.g. thermal).
type DamageType struct{ Id, Name string }

// The eighteen SRD skills.
var (
	SkillAcrobatics     = Skill{"acrobatics", "Acrobatics", AbilityDexterity}
	SkillAnimalHandling = Skill{"animal_handling", "Animal Handling", AbilityWisdom}
	SkillArcana         = Skill{"arcana", "Arcana", AbilityIntelligence}
	SkillAthletics      = Skill{"athletics", "Athletics", AbilityStrength}
	SkillDeception      = Skill{"deception", "Deception", AbilityCharisma}
	SkillHistory        = Skill{"history", "History", AbilityIntelligence}
	SkillInsight        = Skill{"insight", "Insight", AbilityWisdom}
	SkillIntimidation   = Skill{"intimidation", "Intimidation", AbilityCharisma}
	SkillInvestigation  = Skill{"investigation", "Investigation", AbilityIntelligence}
	SkillMedicine       = Skill{"medicine", "Medicine", AbilityWisdom}
	SkillNature         = Skill{"nature", "Nature", AbilityIntelligence}
	SkillPerception     = Skill{"perception", "Perception", AbilityWisdom}
	SkillPerformance    = Skill{"performance", "Performance", AbilityCharisma}
	SkillPersuasion     = Skill{"persuasion", "Persuasion", AbilityCharisma}
	SkillReligion       = Skill{"religion", "Religion", AbilityIntelligence}
	SkillSleightOfHand  = Skill{"sleight_of_hand", "Sleight of Hand", AbilityDexterity}
	SkillStealth        = Skill{"stealth", "Stealth", AbilityDexterity}
	SkillSurvival       = Skill{"survival", "Survival", AbilityWisdom}
)

// SRDSkills is the eighteen SRD skills.
var SRDSkills = []Skill{
	SkillAcrobatics, SkillAnimalHandling, SkillArcana, SkillAthletics,
	SkillDeception, SkillHistory, SkillInsight, SkillIntimidation,
	SkillInvestigation, SkillMedicine, SkillNature, SkillPerception,
	SkillPerformance, SkillPersuasion, SkillReligion, SkillSleightOfHand,
	SkillStealth, SkillSurvival,
}

// SRD physical and common elemental damage types.
var (
	Bludgeoning = DamageType{"bludgeoning", "Bludgeoning"}
	Piercing    = DamageType{"piercing", "Piercing"}
	Slashing    = DamageType{"slashing", "Slashing"}
	Fire        = DamageType{"fire", "Fire"}
	Cold        = DamageType{"cold", "Cold"}
	Lightning   = DamageType{"lightning", "Lightning"}
	Acid        = DamageType{"acid", "Acid"}
	Poison      = DamageType{"poison", "Poison"}
	Necrotic    = DamageType{"necrotic", "Necrotic"}
	Radiant     = DamageType{"radiant", "Radiant"}
	Psychic     = DamageType{"psychic", "Psychic"}
	Thunder     = DamageType{"thunder", "Thunder"}
	Force       = DamageType{"force", "Force"}
)

// DamageAny is a wildcard used only as a key in damage.Mitigation maps. It is
// not a real damage type and must never be put on a damage part.
var DamageAny = DamageType{"any", "Any"}

// A few common SRD conditions; games declare more as open values.
var (
	Prone     = Condition{"prone", "Prone"}
	Poisoned  = Condition{"poisoned", "Poisoned"}
	Restrained = Condition{"restrained", "Restrained"}
	Stunned   = Condition{"stunned", "Stunned"}
	Blinded   = Condition{"blinded", "Blinded"}
	Invisible = Condition{"invisible", "Invisible"}
)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run 'TestSRDSkills|TestDamageAny'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add core/catalog.go core/catalog_test.go
git commit -m "core: open value catalogs (skills, conditions, damage types)"
```

---

### Task 5: Size, carrying capacity, and encumbrance

**Files:**
- Create: `core/encumbrance.go`
- Test: `core/encumbrance_test.go`

**Interfaces:**
- Consumes: `AbilityScore`, `Weight` (Task 2).
- Produces:
  - `type Size int` with `SizeUnspecified Size = 0`, `SizeTiny`, `SizeSmall`, `SizeMedium`, `SizeLarge`, `SizeHuge`, `SizeGargantuan`
  - `func CarryingCapacity(str AbilityScore, size Size) Weight`
  - `func PushDragLift(str AbilityScore, size Size) Weight`
  - `type Encumbrance int` with `Unencumbered Encumbrance = 0`, `Encumbered`, `HeavilyEncumbered`
  - `func EncumbranceTier(str AbilityScore, carried Weight) Encumbrance`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestCarryingCapacity(t *testing.T) {
	// STR 15, Medium: 15*15 = 225; push/drag/lift = 450.
	if got := CarryingCapacity(15, SizeMedium); got != 225 {
		t.Errorf("CarryingCapacity Medium = %d, want 225", got)
	}
	if got := PushDragLift(15, SizeMedium); got != 450 {
		t.Errorf("PushDragLift Medium = %d, want 450", got)
	}
	// Large doubles, Tiny halves.
	if got := CarryingCapacity(15, SizeLarge); got != 450 {
		t.Errorf("CarryingCapacity Large = %d, want 450", got)
	}
	if got := CarryingCapacity(10, SizeTiny); got != 75 {
		t.Errorf("CarryingCapacity Tiny = %d, want 75", got)
	}
}

func TestEncumbranceTier(t *testing.T) {
	// STR 10: encumbered > 50 (STR*5), heavily > 100 (STR*10).
	if got := EncumbranceTier(10, 40); got != Unencumbered {
		t.Errorf("tier(40) = %v, want Unencumbered", got)
	}
	if got := EncumbranceTier(10, 60); got != Encumbered {
		t.Errorf("tier(60) = %v, want Encumbered", got)
	}
	if got := EncumbranceTier(10, 120); got != HeavilyEncumbered {
		t.Errorf("tier(120) = %v, want HeavilyEncumbered", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run 'TestCarryingCapacity|TestEncumbranceTier'`
Expected: FAIL (undefined: CarryingCapacity, ...).

- [ ] **Step 3: Write minimal implementation**

```go
package core

// Size is the closed set of creature sizes.
type Size int

const (
	SizeUnspecified Size = iota
	SizeTiny
	SizeSmall
	SizeMedium
	SizeLarge
	SizeHuge
	SizeGargantuan
)

// sizeFactor gives the carrying-capacity multiplier as a fraction (num/den):
// Tiny halves, Small and Medium x1, Large x2, Huge x4, Gargantuan x8. It panics
// on SizeUnspecified because Size is a must-set sentinel.
func sizeFactor(size Size) (num, den int) {
	switch size {
	case SizeTiny:
		return 1, 2
	case SizeSmall, SizeMedium:
		return 1, 1
	case SizeLarge:
		return 2, 1
	case SizeHuge:
		return 4, 1
	case SizeGargantuan:
		return 8, 1
	default: // SizeUnspecified or out of range
		panic("core: carrying capacity requires a Size")
	}
}

// CarryingCapacity is Strength x 15, scaled by size.
func CarryingCapacity(str AbilityScore, size Size) Weight {
	num, den := sizeFactor(size)
	return Weight(int(str) * 15 * num / den)
}

// PushDragLift is Strength x 30, scaled by size: the maximum to push, drag, or lift.
func PushDragLift(str AbilityScore, size Size) Weight {
	num, den := sizeFactor(size)
	return Weight(int(str) * 30 * num / den)
}

// Encumbrance is the optional SRD encumbrance variant tier.
type Encumbrance int

const (
	Unencumbered Encumbrance = iota
	Encumbered
	HeavilyEncumbered
)

// EncumbranceTier reports the optional variant: Encumbered above Strength x 5,
// HeavilyEncumbered above Strength x 10. Games that do not use the variant never
// call it.
func EncumbranceTier(str AbilityScore, carried Weight) Encumbrance {
	switch {
	case int(carried) > int(str)*10:
		return HeavilyEncumbered
	case int(carried) > int(str)*5:
		return Encumbered
	default:
		return Unencumbered
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run 'TestCarryingCapacity|TestEncumbranceTier'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add core/encumbrance.go core/encumbrance_test.go
git commit -m "core: size, carrying capacity, encumbrance variant"
```

---

### Task 6: Vision and light

**Files:**
- Create: `core/vision.go`
- Test: `core/vision_test.go`

**Interfaces:**
- Consumes: nothing from earlier core tasks.
- Produces:
  - `type VisionType int` with `VisionNormal VisionType = 0`, `VisionDarkvision`, `VisionBlindsight`, `VisionTruesight`, `VisionTremorsense`
  - `type LightLevel int` with `LightUnspecified LightLevel = 0`, `LightBright`, `LightDim`, `LightDark`
  - `type Visibility int` with `VisibilityClear Visibility = 0`, `VisibilityObscured`, `VisibilityBlocked`
  - `func EffectiveLight(ambient LightLevel, vision VisionType, withinRange bool) LightLevel`
  - `func SightVisibility(light LightLevel, blinded, targetInvisible bool) Visibility`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestEffectiveLight(t *testing.T) {
	// Darkvision in range: dark -> dim, dim -> bright.
	if got := EffectiveLight(LightDark, VisionDarkvision, true); got != LightDim {
		t.Errorf("darkvision dark = %v, want Dim", got)
	}
	if got := EffectiveLight(LightDim, VisionDarkvision, true); got != LightBright {
		t.Errorf("darkvision dim = %v, want Bright", got)
	}
	// Out of range: unchanged.
	if got := EffectiveLight(LightDark, VisionDarkvision, false); got != LightDark {
		t.Errorf("darkvision out of range = %v, want Dark", got)
	}
	// Blindsight sees regardless: treat as bright.
	if got := EffectiveLight(LightDark, VisionBlindsight, true); got != LightBright {
		t.Errorf("blindsight = %v, want Bright", got)
	}
}

func TestSightVisibility(t *testing.T) {
	if got := SightVisibility(LightBright, false, false); got != VisibilityClear {
		t.Errorf("bright = %v, want Clear", got)
	}
	if got := SightVisibility(LightDim, false, false); got != VisibilityObscured {
		t.Errorf("dim = %v, want Obscured", got)
	}
	if got := SightVisibility(LightDark, false, false); got != VisibilityBlocked {
		t.Errorf("dark = %v, want Blocked", got)
	}
	if got := SightVisibility(LightBright, true, false); got != VisibilityBlocked {
		t.Errorf("blinded = %v, want Blocked", got)
	}
	if got := SightVisibility(LightBright, false, true); got != VisibilityBlocked {
		t.Errorf("invisible target = %v, want Blocked", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run 'TestEffectiveLight|TestSightVisibility'`
Expected: FAIL (undefined: EffectiveLight, ...).

- [ ] **Step 3: Write minimal implementation**

```go
package core

// VisionType is a creature's relevant sense for a situation. VisionNormal is the
// baseline every creature has.
type VisionType int

const (
	VisionNormal VisionType = iota
	VisionDarkvision
	VisionBlindsight
	VisionTruesight
	VisionTremorsense
)

// LightLevel is the ambient light. LightUnspecified is the must-set zero.
type LightLevel int

const (
	LightUnspecified LightLevel = iota
	LightBright
	LightDim
	LightDark
)

// Visibility is the result of SightVisibility. VisibilityClear is the safe zero.
type Visibility int

const (
	VisibilityClear Visibility = iota
	VisibilityObscured
	VisibilityBlocked
)

// EffectiveLight applies a creature's vision to the ambient light. Darkvision in
// range treats dark as dim and dim as bright; blindsight and truesight see
// regardless of light.
func EffectiveLight(ambient LightLevel, vision VisionType, withinRange bool) LightLevel {
	switch vision {
	case VisionBlindsight, VisionTruesight, VisionTremorsense:
		if withinRange {
			return LightBright
		}
		return ambient
	case VisionDarkvision:
		if !withinRange {
			return ambient
		}
		switch ambient {
		case LightDark:
			return LightDim
		case LightDim:
			return LightBright
		default:
			return ambient
		}
	default:
		return ambient
	}
}

// SightVisibility turns effective light and conditions into a visibility state:
// dim is Obscured (disadvantage on sight Perception); darkness, blinded, or an
// invisible target is Blocked (cannot see).
func SightVisibility(light LightLevel, blinded, targetInvisible bool) Visibility {
	if blinded || targetInvisible || light == LightDark {
		return VisibilityBlocked
	}
	if light == LightDim {
		return VisibilityObscured
	}
	return VisibilityClear
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run 'TestEffectiveLight|TestSightVisibility'`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add core/vision.go core/vision_test.go
git commit -m "core: vision and light"
```

---

### Task 7: Effect duration

**Files:**
- Create: `core/duration.go`
- Test: `core/duration_test.go`

**Interfaces:**
- Consumes: `Ability` (Task 2).
- Produces:
  - `type DurationKind int` with `DurationUnspecified DurationKind = 0`, `DurationInstant`, `DurationRounds`, `DurationEndOfNextTurn`, `DurationConcentration`, `DurationUntilRemoved`
  - `type EffectDuration struct { Kind DurationKind; Rounds int; SaveEnds bool; SaveAbility Ability }`
  - `func RoundsInMinutes(m int) int`
  - `func RoundsInHours(h int) int`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestRoundsConversion(t *testing.T) {
	if got := RoundsInMinutes(1); got != 10 {
		t.Errorf("RoundsInMinutes(1) = %d, want 10", got)
	}
	if got := RoundsInHours(1); got != 600 {
		t.Errorf("RoundsInHours(1) = %d, want 600", got)
	}
}

func TestEffectDurationZero(t *testing.T) {
	var d EffectDuration
	if d.Kind != DurationUnspecified {
		t.Errorf("zero EffectDuration.Kind = %v, want DurationUnspecified", d.Kind)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run 'TestRoundsConversion|TestEffectDurationZero'`
Expected: FAIL (undefined: RoundsInMinutes, ...).

- [ ] **Step 3: Write minimal implementation**

```go
package core

// DurationKind classifies how long an applied effect or condition lasts.
type DurationKind int

const (
	DurationUnspecified DurationKind = iota
	DurationInstant
	DurationRounds
	DurationEndOfNextTurn
	DurationConcentration
	DurationUntilRemoved
)

// EffectDuration describes how long an applied effect lasts. Named EffectDuration
// (not Duration) to stay clear of time.Duration. The game ticks the countdown;
// SaveEnds means the bearer repeats the save at the end of each of its turns.
type EffectDuration struct {
	Kind        DurationKind
	Rounds      int // when DurationRounds
	SaveEnds    bool
	SaveAbility Ability // which save, when SaveEnds
}

// RoundsInMinutes converts minutes to rounds (1 round = 6 seconds).
func RoundsInMinutes(m int) int { return m * 10 }

// RoundsInHours converts hours to rounds.
func RoundsInHours(h int) int { return h * 600 }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run 'TestRoundsConversion|TestEffectDurationZero'`
Expected: PASS.

- [ ] **Step 5: Run the whole package and commit**

```bash
go test ./core/
git add core/duration.go core/duration_test.go
git commit -m "core: effect duration and round conversions"
```

---

## Self-review notes

Spec coverage for the `core` section of the design spec:
- Units (`AbilityScore`, `Modifier`, `Level`, `ArmorClass`, `Dc`, `Distance`, `HitPoints`, `Xp`, `Weight`): Task 2.
- `Ability`/`AbilityScores`/`AbilityModifier`/`AbilityAny`: Task 2.
- Stat math (`ProficiencyBonus`, `SkillBonus`, `SavingThrowBonus`, `SpellSaveDc`, `SpellAttackBonus`, `MaxHp`, `XpForLevel`, `XpForNextLevel`): Task 3.
- Open catalogs (`Skill`, `Condition`, `DamageType`, `SRDSkills`, `DamageAny`): Task 4.
- `Size`, `CarryingCapacity`, `PushDragLift`, `Encumbrance`, `EncumbranceTier`: Task 5.
- Vision/light (`VisionType`, `LightLevel`, `Visibility`, `EffectiveLight`, `SightVisibility`): Task 6.
- `DurationKind`, `EffectDuration`, `RoundsInMinutes`, `RoundsInHours`: Task 7.

`Distance` is declared here but consumed by `combat`/`turn`/`content` in later plans. The XP-threshold long-rest fraction and the exact hit-dice recovery are in the `resource` plan, not here.
