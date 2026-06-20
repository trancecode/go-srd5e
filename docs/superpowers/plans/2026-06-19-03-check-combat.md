# check and combat packages implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Build the `check` package (ability checks, skill checks, saving throws, opposed contests, passive scores) and the `combat` package (attack resolution, cover, range, the unified melee/ranged attack setup).

**Architecture:** Both are pure resolvers over already-rolled numbers — only `dice` consumes a `Roller`. `check` and `combat` depend on `core` and `dice`. A small `core.AbilityScores.Score` accessor is added first (Task 1) since the check constructors need it.

**Tech Stack:** Go 1.25, standard library plus `core` and `dice`.

Plan 3 of 8 (design spec `docs/superpowers/specs/2026-06-18-go-srd5e-module-design.md`, the `## Package: check` and `## Package: combat` sections).

## Global Constraints

- Module `github.com/trancecode/go-srd5e`, Go `1.25`. Acronyms in Java camel case (`Dc`, `Ac`, `Hp`).
- `check` and `combat` import only `core`, `dice`, and the standard library. They never consume a `Roller` (pure; they take an already-rolled `dice.Result` or natural roll).
- Enum zeros: `AttackMiss` is the (function-produced) zero for `AttackOutcome`; `CoverNone`, `BandNormal` are valid zeros.
- Value types have exported fields. Errors (none expected here) follow `<context>: <reason>`.
- Crit on natural 20 only; natural 1 auto-misses. Cover gives +2 (half) / +5 (three-quarters) to AC and Dexterity saves; total cover blocks targeting.
- TDD; `gofmt -l`; commit on `main` directly with author "Claude Code" and NO `Co-Authored-By:` line.

---

### Task 1: core.AbilityScores.Score accessor

**Files:**
- Modify: `core/ability.go`
- Test: `core/ability_score_test.go`

**Interfaces:**
- Consumes: `AbilityScores`, `Ability`, `AbilityScore` (existing in `core`).
- Produces: `func (s AbilityScores) Score(a Ability) AbilityScore`

- [ ] **Step 1: Write the failing test**

```go
package core

import "testing"

func TestAbilityScoresScore(t *testing.T) {
	s := AbilityScores{Strength: 15, Dexterity: 12, Constitution: 14, Intelligence: 8, Wisdom: 10, Charisma: 13}
	cases := []struct {
		a    Ability
		want AbilityScore
	}{
		{AbilityStrength, 15}, {AbilityDexterity, 12}, {AbilityConstitution, 14},
		{AbilityIntelligence, 8}, {AbilityWisdom, 10}, {AbilityCharisma, 13},
	}
	for _, c := range cases {
		if got := s.Score(c.a); got != c.want {
			t.Errorf("Score(%v) = %d, want %d", c.a, got, c.want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("Score(AbilityNone) should panic")
		}
	}()
	s.Score(AbilityNone)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./core/ -run TestAbilityScoresScore`
Expected: FAIL (Score undefined).

- [ ] **Step 3: Append to `core/ability.go`**

```go
// Score returns the score for a concrete ability. It panics on AbilityNone or
// AbilityAny, which are not concrete abilities.
func (s AbilityScores) Score(a Ability) AbilityScore {
	switch a {
	case AbilityStrength:
		return s.Strength
	case AbilityDexterity:
		return s.Dexterity
	case AbilityConstitution:
		return s.Constitution
	case AbilityIntelligence:
		return s.Intelligence
	case AbilityWisdom:
		return s.Wisdom
	case AbilityCharisma:
		return s.Charisma
	default:
		panic("core: Score requires a concrete ability")
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./core/ -run TestAbilityScoresScore`
Expected: PASS.

- [ ] **Step 5: gofmt and commit**

```bash
gofmt -w core/
git add core/ability.go core/ability_score_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "core: AbilityScores.Score accessor"
```

---

### Task 2: check package

**Files:**
- Create: `check/doc.go`
- Create: `check/check.go`
- Test: `check/check_test.go`

**Interfaces:**
- Consumes: `core.Modifier`, `core.Dc`, `core.AbilityScores`, `core.Ability`, `core.Skill`, `core.Level`, `core.AbilityModifier`, `core.SkillBonus`, `core.SavingThrowBonus`, `core.AbilityScores.Score` (Task 1); `dice.Result`, `dice.Vantage`.
- Produces:
  - `type Check struct { Bonus core.Modifier; Dc core.Dc }`
  - `type Result struct { Roll dice.Result; Total int; Success bool; Margin int }`
  - `func (c Check) Resolve(roll dice.Result) Result`
  - `func Ability(scores core.AbilityScores, a core.Ability, dc core.Dc) Check`
  - `func Skill(scores core.AbilityScores, s core.Skill, proficient bool, level core.Level, dc core.Dc) Check`
  - `func Save(scores core.AbilityScores, a core.Ability, proficient bool, level core.Level, dc core.Dc) Check`
  - `type ContestResult struct { InitiatorTotal, ResponderTotal int; InitiatorWins bool }`
  - `func Contest(initiatorTotal, responderTotal int) ContestResult`
  - `func PassiveScore(modifier core.Modifier, v dice.Vantage) int`

- [ ] **Step 1: Write the failing test**

```go
package check

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestResolve(t *testing.T) {
	c := Check{Bonus: 5, Dc: 15}
	// rolled a 12 -> total 17 -> success, margin 2.
	got := c.Resolve(dice.Result{Dice: []int{12}, Total: 12})
	if !got.Success || got.Total != 17 || got.Margin != 2 {
		t.Errorf("success case = %+v, want Total 17 Success true Margin 2", got)
	}
	// rolled a 9 -> total 14 -> fail, margin -1.
	got = c.Resolve(dice.Result{Dice: []int{9}, Total: 9})
	if got.Success || got.Total != 14 || got.Margin != -1 {
		t.Errorf("fail case = %+v, want Total 14 Success false Margin -1", got)
	}
}

func TestConstructors(t *testing.T) {
	sc := core.AbilityScores{Strength: 16, Dexterity: 16, Intelligence: 16}
	// STR 16 (+3), level 5 (prof +3).
	if Ability(sc, core.AbilityStrength, 12).Bonus != 3 {
		t.Errorf("Ability bonus = %d, want 3", Ability(sc, core.AbilityStrength, 12).Bonus)
	}
	if Skill(sc, core.SkillAthletics, true, 5, 12).Bonus != 6 {
		t.Errorf("Skill proficient bonus = %d, want 6", Skill(sc, core.SkillAthletics, true, 5, 12).Bonus)
	}
	if Save(sc, core.AbilityDexterity, false, 5, 12).Bonus != 3 {
		t.Errorf("Save non-proficient bonus = %d, want 3", Save(sc, core.AbilityDexterity, false, 5, 12).Bonus)
	}
}

func TestContest(t *testing.T) {
	if !Contest(18, 12).InitiatorWins {
		t.Error("18 vs 12 should be initiator win")
	}
	if Contest(12, 18).InitiatorWins {
		t.Error("12 vs 18 should be responder win")
	}
	if Contest(15, 15).InitiatorWins {
		t.Error("tie should favor responder (InitiatorWins false)")
	}
}

func TestPassiveScore(t *testing.T) {
	if PassiveScore(3, dice.VantageNone) != 13 {
		t.Errorf("passive none = %d, want 13", PassiveScore(3, dice.VantageNone))
	}
	if PassiveScore(3, dice.VantageAdvantage) != 18 {
		t.Errorf("passive advantage = %d, want 18", PassiveScore(3, dice.VantageAdvantage))
	}
	if PassiveScore(3, dice.VantageDisadvantage) != 8 {
		t.Errorf("passive disadvantage = %d, want 8", PassiveScore(3, dice.VantageDisadvantage))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./check/`
Expected: FAIL (package/symbols undefined).

- [ ] **Step 3: Write `check/doc.go`**

```go
// Package check resolves d20 checks against a difficulty class: ability checks,
// skill checks, and saving throws, plus the opposed contest used by Shove and
// Grapple and passive scores. It is pure: it interprets a roll the caller already
// made; the roll (with vantage) happens in dice.
package check
```

- [ ] **Step 4: Write `check/check.go`**

```go
package check

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// Check is a configured d20 check: a total bonus against a difficulty class.
type Check struct {
	Bonus core.Modifier
	Dc    core.Dc
}

// Result is the outcome of resolving a Check against a d20 roll.
type Result struct {
	Roll    dice.Result
	Total   int
	Success bool
	Margin  int // Total minus Dc
}

// Resolve interprets an already-made d20 roll against the check.
func (c Check) Resolve(roll dice.Result) Result {
	total := roll.Total + int(c.Bonus)
	margin := total - int(c.Dc)
	return Result{Roll: roll, Total: total, Success: margin >= 0, Margin: margin}
}

// Ability builds an ability check: the ability modifier against the DC.
func Ability(scores core.AbilityScores, a core.Ability, dc core.Dc) Check {
	return Check{Bonus: core.AbilityModifier(scores.Score(a)), Dc: dc}
}

// Skill builds a skill check, reading the governing ability from the skill.
func Skill(scores core.AbilityScores, s core.Skill, proficient bool, level core.Level, dc core.Dc) Check {
	return Check{Bonus: core.SkillBonus(scores.Score(s.Ability), proficient, level), Dc: dc}
}

// Save builds a saving throw.
func Save(scores core.AbilityScores, a core.Ability, proficient bool, level core.Level, dc core.Dc) Check {
	return Check{Bonus: core.SavingThrowBonus(scores.Score(a), proficient, level), Dc: dc}
}

// ContestResult is the outcome of an opposed check. Ties favor the responder.
type ContestResult struct {
	InitiatorTotal, ResponderTotal int
	InitiatorWins                  bool
}

// Contest resolves an opposed check from the two already-computed totals.
func Contest(initiatorTotal, responderTotal int) ContestResult {
	return ContestResult{
		InitiatorTotal: initiatorTotal,
		ResponderTotal: responderTotal,
		InitiatorWins:  initiatorTotal > responderTotal,
	}
}

// PassiveScore is 10 + modifier, +5 for advantage, -5 for disadvantage. Generic
// across passive Perception, Investigation, and Insight.
func PassiveScore(modifier core.Modifier, v dice.Vantage) int {
	p := 10 + int(modifier)
	switch v {
	case dice.VantageAdvantage:
		p += 5
	case dice.VantageDisadvantage:
		p -= 5
	}
	return p
}
```

- [ ] **Step 5: Run tests; gofmt; commit**

Run: `go test ./check/` (expect PASS), then:

```bash
gofmt -w check/
git add check/
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "check: Check, Resolve, constructors, Contest, PassiveScore"
```

---

### Task 3: combat resolution

**Files:**
- Create: `combat/doc.go`
- Create: `combat/combat.go`
- Test: `combat/combat_test.go`

**Interfaces:**
- Consumes: `core.Modifier`, `core.ArmorClass`, `core.Dc`.
- Produces:
  - `type AttackOutcome int` with `AttackMiss AttackOutcome = 0`, `AttackHit`, `AttackCritical`
  - `type AttackResult struct { Outcome AttackOutcome; NaturalRoll int; Total int }`
  - `func ResolveAttack(naturalRoll int, mod core.Modifier, ac core.ArmorClass) AttackResult`
  - `func ConcentrationDc(damage int) core.Dc`

- [ ] **Step 1: Write the failing test**

```go
package combat

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
)

func TestResolveAttack(t *testing.T) {
	// natural 20 -> critical regardless of AC.
	if r := ResolveAttack(20, 0, 99); r.Outcome != AttackCritical {
		t.Errorf("nat 20 = %v, want Critical", r.Outcome)
	}
	// natural 1 -> miss regardless of bonus.
	if r := ResolveAttack(1, 100, 5); r.Outcome != AttackMiss {
		t.Errorf("nat 1 = %v, want Miss", r.Outcome)
	}
	// 15 + 4 = 19 vs AC 18 -> hit; Total recorded.
	r := ResolveAttack(15, 4, 18)
	if r.Outcome != AttackHit || r.Total != 19 || r.NaturalRoll != 15 {
		t.Errorf("hit = %+v, want Hit Total 19 NaturalRoll 15", r)
	}
	// 10 + 2 = 12 vs AC 18 -> miss.
	if r := ResolveAttack(10, 2, 18); r.Outcome != AttackMiss {
		t.Errorf("low total = %v, want Miss", r.Outcome)
	}
}

func TestConcentrationDc(t *testing.T) {
	cases := []struct {
		dmg  int
		want core.Dc
	}{{8, 10}, {19, 10}, {22, 11}, {30, 15}}
	for _, c := range cases {
		if got := ConcentrationDc(c.dmg); got != c.want {
			t.Errorf("ConcentrationDc(%d) = %d, want %d", c.dmg, got, c.want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./combat/`
Expected: FAIL (package/symbols undefined).

- [ ] **Step 3: Write `combat/doc.go`**

```go
// Package combat resolves attacks and the targeting rules (cover, range bands)
// that feed them. It is pure: it interprets already-rolled numbers and
// already-determined spatial facts, never consuming a Roller or computing
// geometry. Critical hits occur on a natural 20 only.
package combat
```

- [ ] **Step 4: Write `combat/combat.go`**

```go
package combat

import "github.com/trancecode/go-srd5e/core"

// AttackOutcome is the result kind of an attack roll.
type AttackOutcome int

const (
	AttackMiss AttackOutcome = iota
	AttackHit
	AttackCritical
)

// AttackResult is the outcome of resolving an attack roll against armor class.
type AttackResult struct {
	Outcome     AttackOutcome
	NaturalRoll int
	Total       int
}

// ResolveAttack resolves a d20 attack: natural 20 is a critical hit and natural
// 1 an automatic miss, regardless of modifiers; otherwise the total (natural +
// modifier) is compared to the target's armor class.
func ResolveAttack(naturalRoll int, mod core.Modifier, ac core.ArmorClass) AttackResult {
	total := naturalRoll + int(mod)
	switch {
	case naturalRoll == 20:
		return AttackResult{Outcome: AttackCritical, NaturalRoll: naturalRoll, Total: total}
	case naturalRoll == 1:
		return AttackResult{Outcome: AttackMiss, NaturalRoll: naturalRoll, Total: total}
	case total >= int(ac):
		return AttackResult{Outcome: AttackHit, NaturalRoll: naturalRoll, Total: total}
	default:
		return AttackResult{Outcome: AttackMiss, NaturalRoll: naturalRoll, Total: total}
	}
}

// ConcentrationDc is the save DC to maintain concentration after taking damage:
// the greater of 10 and half the damage taken.
func ConcentrationDc(damage int) core.Dc {
	dc := damage / 2
	if dc < 10 {
		dc = 10
	}
	return core.Dc(dc)
}
```

- [ ] **Step 5: Run tests; gofmt; commit**

Run: `go test ./combat/` (expect PASS), then:

```bash
gofmt -w combat/
git add combat/doc.go combat/combat.go combat/combat_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "combat: attack resolution and concentration DC"
```

---

### Task 4: combat targeting (cover, range, attack setup)

**Files:**
- Create: `combat/targeting.go`
- Test: `combat/targeting_test.go`

**Interfaces:**
- Consumes: `core.Modifier`, `core.ArmorClass`, `core.Distance`; `dice.Vantage`, `dice.CombineVantage`.
- Produces:
  - `type Cover int` with `CoverNone CoverHalf CoverThreeQuarters CoverTotal`; methods `AcBonus() core.Modifier`, `DexSaveBonus() core.Modifier`, `BlocksTargeting() bool`
  - `type Range struct { Normal, Long core.Distance }`; `type Band int` with `BandNormal BandLong BandOutOf`; `func MeleeRange(reach core.Distance) Range`; `func (r Range) Band(distance core.Distance) Band`
  - `type Attack struct { Range Range; Distance core.Distance; TargetCover Cover; Advantage, Disadvantage bool }`
  - `type AttackSetup struct { Possible bool; Vantage dice.Vantage; EffectiveAc core.ArmorClass }`
  - `func (a Attack) Setup(baseAc core.ArmorClass) AttackSetup`

- [ ] **Step 1: Write the failing test**

```go
package combat

import (
	"testing"

	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

func TestCover(t *testing.T) {
	if CoverNone.AcBonus() != 0 || CoverHalf.AcBonus() != 2 || CoverThreeQuarters.AcBonus() != 5 {
		t.Error("cover AC bonuses wrong")
	}
	if CoverHalf.DexSaveBonus() != 2 {
		t.Error("cover dex save bonus should match AC bonus")
	}
	if !CoverTotal.BlocksTargeting() || CoverHalf.BlocksTargeting() {
		t.Error("only total cover blocks targeting")
	}
}

func TestRangeBand(t *testing.T) {
	r := Range{Normal: 80, Long: 320}
	if r.Band(40) != BandNormal || r.Band(200) != BandLong || r.Band(400) != BandOutOf {
		t.Error("ranged bands wrong")
	}
	// melee: reach with no long band.
	m := MeleeRange(5)
	if m.Band(5) != BandNormal || m.Band(10) != BandOutOf {
		t.Error("melee bands wrong")
	}
}

func TestAttackSetup(t *testing.T) {
	// total cover -> not possible.
	if (Attack{TargetCover: CoverTotal, Range: MeleeRange(5)}).Setup(15).Possible {
		t.Error("total cover should block")
	}
	// out of range -> not possible.
	if (Attack{Range: Range{Normal: 80, Long: 320}, Distance: 400}).Setup(15).Possible {
		t.Error("beyond long range should block")
	}
	// long range -> disadvantage; half cover -> +2 AC.
	s := Attack{Range: Range{80, 320}, Distance: 200, TargetCover: CoverHalf}.Setup(15)
	if !s.Possible || s.Vantage != dice.VantageDisadvantage || s.EffectiveAc != 17 {
		t.Errorf("long+halfcover = %+v, want Possible Vantage Disadvantage EffectiveAc 17", s)
	}
	// long range disadvantage cancels with a supplied advantage.
	s = Attack{Range: Range{80, 320}, Distance: 200, Advantage: true}.Setup(15)
	if s.Vantage != dice.VantageNone {
		t.Errorf("advantage + long-range disadvantage = %v, want None", s.Vantage)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./combat/ -run 'TestCover|TestRangeBand|TestAttackSetup'`
Expected: FAIL (undefined: Cover, Range, Attack, ...).

- [ ] **Step 3: Write `combat/targeting.go`**

```go
package combat

import (
	"github.com/trancecode/go-srd5e/core"
	"github.com/trancecode/go-srd5e/dice"
)

// Cover is the degree of cover a target has.
type Cover int

const (
	CoverNone Cover = iota
	CoverHalf
	CoverThreeQuarters
	CoverTotal
)

// AcBonus is the bonus cover grants to armor class (+0/+2/+5).
func (c Cover) AcBonus() core.Modifier {
	switch c {
	case CoverHalf:
		return 2
	case CoverThreeQuarters:
		return 5
	default:
		return 0
	}
}

// DexSaveBonus is the bonus cover grants to Dexterity saving throws (same as AC).
func (c Cover) DexSaveBonus() core.Modifier { return c.AcBonus() }

// BlocksTargeting reports whether the target cannot be targeted directly.
func (c Cover) BlocksTargeting() bool { return c == CoverTotal }

// Range is a weapon's normal and long range, in feet.
type Range struct {
	Normal, Long core.Distance
}

// Band classifies a distance relative to a range.
type Band int

const (
	BandNormal Band = iota
	BandLong
	BandOutOf
)

// MeleeRange builds a Range for a melee weapon: reach with no long band, so any
// distance beyond reach is out of range.
func MeleeRange(reach core.Distance) Range { return Range{Normal: reach, Long: reach} }

// Band reports the band a distance falls in: within normal, in the long-range
// (disadvantage) window, or out of range.
func (r Range) Band(distance core.Distance) Band {
	switch {
	case distance <= r.Normal:
		return BandNormal
	case distance <= r.Long:
		return BandLong
	default:
		return BandOutOf
	}
}

// Attack bundles the pre-roll facts that determine whether an attack is possible
// and at what vantage. The game supplies Distance, cover, and any pre-aggregated
// advantage/disadvantage; long-range disadvantage is added by Setup.
type Attack struct {
	Range        Range
	Distance     core.Distance
	TargetCover  Cover
	Advantage    bool
	Disadvantage bool
}

// AttackSetup is the pre-roll result: whether the attack is possible, the vantage
// to roll the d20 with, and the cover-adjusted armor class.
type AttackSetup struct {
	Possible    bool
	Vantage     dice.Vantage
	EffectiveAc core.ArmorClass
}

// Setup computes the pre-roll setup against the target's full armor class. Total
// cover or a distance beyond long range makes the attack impossible.
func (a Attack) Setup(baseAc core.ArmorClass) AttackSetup {
	if a.TargetCover.BlocksTargeting() {
		return AttackSetup{Possible: false}
	}
	band := a.Range.Band(a.Distance)
	if band == BandOutOf {
		return AttackSetup{Possible: false}
	}
	disadvantage := a.Disadvantage || band == BandLong
	return AttackSetup{
		Possible:    true,
		Vantage:     dice.CombineVantage(a.Advantage, disadvantage),
		EffectiveAc: baseAc + core.ArmorClass(a.TargetCover.AcBonus()),
	}
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./combat/`
Expected: PASS.

- [ ] **Step 5: Whole-module gate; gofmt; commit**

Run: `go build ./... && go vet ./... && go test ./...`, then:

```bash
gofmt -w combat/
git add combat/targeting.go combat/targeting_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "combat: cover, range, and unified attack setup"
```

---

## Self-review notes

Spec coverage:
- `check`: `Check`/`Result`/`Resolve`, `Ability`/`Skill`/`Save`, `Contest`/`ContestResult`, `PassiveScore` — Task 2. The `Score` accessor it relies on — Task 1.
- `combat`: `AttackOutcome`/`AttackResult`/`ResolveAttack`/`ConcentrationDc` — Task 3; `Cover`/`Range`/`Band`/`MeleeRange`/`Attack`/`AttackSetup`/`Setup` — Task 4.

`Distance` (declared in `core` in Plan 1) is consumed here. `combat` → `core`, `dice`; `check` → `core`, `dice`. Both acyclic.
