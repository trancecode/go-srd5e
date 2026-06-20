# go-srd5e module design

## Purpose

`go-srd5e` is a reusable Go module that encodes the rules of the System
Reference Document for 5th edition (SRD 5e, CC-BY-4.0). It exists so that
multiple games can share one correct, tested implementation of the d20 core
instead of each reimplementing ability modifiers, dice, checks, combat
resolution, damage, and turn bookkeeping.

The module is a rules kernel and a set of assist utilities. By *kernel* we mean
the minimal, game-agnostic core of the rules: pure functions and value types
that turn inputs into outcomes. It deliberately owns none of the surrounding
game, neither content, game state, rendering, geometry, nor difficulty policy.
Those belong to each game.

## Consumers

Three Go games are the initial consumers. All three are full Go stacks
(including rendering), so the module is consumed as a plain Go dependency.

* A first-person fantasy dungeon crawler. It already has a mature,
  self-contained SRD5e package that this module generalizes. Lowest migration
  friction; its call sites are the reference for proven shapes.
* A second fantasy game, earlier stage, with a documented plan to adopt SRD5e
  rules. Its native model uses a different stat set that maps down to the six
  SRD abilities at the call boundary.
* A 2D turn-based tactical cyberpunk game. It keeps the d20 core verbatim but
  reskins almost everything above it (renamed races, spells, and classes) and
  adds bespoke meta-systems with no SRD analogue. It uses an
  entity-component-system (ECS) architecture.

The third game is the reason the module must stay generic: it shares the kernel
but supplies entirely its own content.

## Scope

The module ships two layers now and defers a third.

* Layer 1, the rules kernel. Pure functions and value types over the SRD's
  units: ability math, dice, checks, combat resolution, the damage pipeline, a
  declarative effects vocabulary, and turn-bookkeeping helpers. Identical across
  all three games.
* Layer 2, content shapes. Field-bearing structs (`Class`, `Race`, `Spell`,
  `Item`, `Creature`) plus a generic registry. Offered, never required. Games
  populate them with their own content.
* Layer 3, reference content data (deferred). The actual standard SRD 5e
  classes, races, spells, and items as ready-to-use data. Only the fantasy
  games would want it; the cyberpunk game would ignore it. Go links only what is used, so
  it can be added later as an opt-in subpackage with no cost to games that skip
  it.

## Design principles

These principles are binding for the implementation.

* Data flows in as values; behavior is injected as interfaces. Records such as
  ability scores and dice expressions are passed as concrete value types, not
  hidden behind interfaces. The only interfaces in the module are behavioral
  seams: `dice.Roller` (a rolling strategy) and `resource.Restorable` (recover on
  a rest), both small and behavioral, never used to hide data.
* No mandatory aggregate type. The kernel never requires a game to construct a
  `Character`, `Actor`, or `Combatant` object. Every function takes the minimal
  values it needs, so a game reads those from wherever it stores them (including
  separate ECS components) and calls in.
* Value types are designed to be composed. Games embed or alias the module's
  value types into their own structs and components, so adapting to the library
  is composition, not copying.
* Rolling and resolving are separate verbs. Only `dice` consumes a `Roller`;
  every other package interprets already-rolled numbers. This gives one
  teachable rule, lets resolvers be tested with no randomness, and leaves a
  clean seam for reroll and luck features that act between the roll and the
  outcome.
* Rollers are passed per call, never stored. No library type holds a `Roller`.
  This keeps value types serializable, keeps resolution deterministic and
  replayable, and lets a game use different rollers for different entities in
  the same tick.
* Geometry stays in the game. The kernel never computes distance, line of sight,
  adjacency, or what grants cover, and it never moves a creature. It maps
  already-determined spatial facts to mechanical consequences; the game's
  spatial layer supplies the facts and executes the movement.
* Named types for units. A distinct unit that crosses an API boundary gets a
  named type (`Distance`, `AbilityScore`, `ArmorClass`, and so on) rather than a
  bare `int` with a comment. Plain `int` is kept for transient arithmetic. See
  "Types and units."
* Enum zero values are explicit. Every enum names its zero value. Where zero is
  a genuine neutral or absent state it is named `None` (for example
  `VantageNone`, `CoverNone`) or `AbilityNone`, and is a valid value. Where the
  enum has no natural neutral, zero is a `…Unspecified` sentinel meaning the
  field was not set, which functions treat as invalid or a no-op. This stops an
  unset struct field from silently meaning a real member such as a sphere or a
  push. Enums only ever returned by functions, never struct-constructed
  (`AttackOutcome`, `Band`), keep their zero as the safe member.
* Version-neutral kernel. Every kernel function is identical between SRD 5.1 and
  SRD 5.2. Version differences live in content modeling, which layer 2 shapes
  absorb through optional fields. SRD 5.2 is documented as the reference
  baseline because the cyberpunk game is the active project. No version flags, no forks.

## Types and units

Distinct SRD units are named types so the compiler rejects argument swaps (for
example passing an ability score where a distance is wanted). They live in
`core`, the vocabulary package every other package imports.

```go
type AbilityScore int   // 1..30; method: (s AbilityScore) Modifier() Modifier
type Modifier     int    // a value added to a d20: ability modifier, proficiency, attack or skill bonus
type Level        int    // 1..20
type ArmorClass   int
type Dc           int    // difficulty class; kept separate from ArmorClass so the two cannot be swapped
type Distance     int    // feet; shared by combat.Range, content.Race.MovementSpeed, and turn movement
type HitPoints    int
type Xp           int
type Weight       int    // pounds
```

Two properties make this ergonomic rather than noisy in Go. Untyped constant
literals convert freely, so `Check{Dc: 15}` and `MeleeRange(5)` need no casts;
only a wrongly typed variable is rejected. And the mixed arithmetic that needs
conversions is almost entirely inside a handful of kernel functions, so client
code passes and receives typed values without doing that arithmetic itself.

Deliberate exceptions kept as `int`: raw die faces (`dice.Result.Dice` and
`Total`), so `dice` stays a dependency-free leaf; and damage amounts inside the
`damage` package, which are heavily arithmetic and not confusable with the
roll-and-Dc cluster.

## Package layout

No Go files live at the repository root. The module is a set of focused
packages, each with a single purpose and its own `doc.go`.

```
github.com/trancecode/go-srd5e/
  core/     units, open value catalogs (Ability, Skill, Condition, DamageType), stat math, vision/light
  dice/     Roller, Expr, Result, Vantage, RollD20, CombineVantage, fixed rollers
  check/    Check, Contest, PassiveScore: ability checks, skill checks, saving throws (pure)
  combat/   ResolveAttack, Cover, Range, Attack/Setup (pure)
  damage/   typed Damage, Roll, Mitigation, ApplyMitigation, ApplyToHp, Apply
  effect/   declarative effects (damage, healing, condition, movement, modifier), Triggered, ModifierBonus
  turn/     initiative event timeline, action-economy and movement bookkeeping
  resource/ SpellSlotPool, HitDicePool, Resource, ResourceSet, Restorable (runtime pools)
  content/  Class, Race, Spell, Item, Creature, Registry (layer 2, optional)
```

Dependency graph (acyclic):

```
core     (leaf)
dice     (leaf)
check    -> core, dice
combat   -> core, dice
damage   -> core, dice, combat
effect   -> core, dice, damage
turn     -> core, dice
resource -> core, dice
content  -> core, dice, combat, damage, effect
```

## Package: core

The SRD vocabulary, the named units, the open value catalogs, and the pure stat
math. No randomness; fully deterministic.

```go
type Ability int   // AbilityNone (0), then AbilityStrength .. AbilityCharisma; closed set, kept by all games
type AbilityScores struct {
    Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma AbilityScore
}

func AbilityModifier(score AbilityScore) Modifier             // (score-10)/2, floored
func ProficiencyBonus(level Level) Modifier                    // 2 + (level-1)/4
func SkillBonus(score AbilityScore, proficient bool, level Level) Modifier
func SavingThrowBonus(score AbilityScore, proficient bool, level Level) Modifier
func SpellSaveDc(score AbilityScore, level Level) Dc
func SpellAttackBonus(score AbilityScore, level Level) Modifier
func MaxHp(hpRolls []int, conMod Modifier) HitPoints           // hpRolls are die faces
func XpForLevel(level Level) Xp
func XpForNextLevel(current Level) Xp
```

### Size, carrying capacity, and encumbrance

Size is a closed enum (the SRD fixes the categories, and every game keeps them).
Carrying capacity is Strength-derived stat math, scaled by size. The encumbrance
*variant* is opt-in: it reports a tier and its documented effects, and a game
that does not use the variant (for example a team-level model) simply
never calls it.

```go
type Size int  // SizeUnspecified (0), SizeTiny, SizeSmall, SizeMedium, SizeLarge, SizeHuge, SizeGargantuan

func CarryingCapacity(str AbilityScore, size Size) Weight   // STR × 15, scaled by size (Tiny halves, Large ×2, ...)
func PushDragLift(str AbilityScore, size Size) Weight        // STR × 30, scaled by size: max to push, drag, or lift

type Encumbrance int  // Unencumbered, Encumbered, HeavilyEncumbered
// EncumbranceTier reports the optional SRD encumbrance variant:
//   Encumbered (carried > STR×5): speed -10 ft
//   HeavilyEncumbered (carried > STR×10): speed -20 ft, plus disadvantage on
//   STR/DEX/CON checks, attack rolls, and saving throws
func EncumbranceTier(str AbilityScore, carried Weight) Encumbrance
```

The kernel returns per-creature capacity and the tier; the game owns how it
accumulates carried weight (per-character, team-level, or not at all) and how it
applies the consequences. A speed penalty feeds back into `turn.Economy.MovementSpeed`,
and the disadvantage is surfaced through the generic `Disadvantage` input on
checks and attacks.

### Open value catalogs

Skills, conditions, and damage types are open value types rather than closed
enums, because games genuinely diverge on these lists (a cyberpunk game adds a Hacking
skill, a Slow condition, and thermal and electromagnetic-pulse damage). Each
carries the only facts the kernel needs and ships the SRD standard set as
predefined values; a game declares more the same way, first-class with no kernel
change.

```go
type Skill struct {
    Id, Name string
    Ability  Ability   // the only fact the check math needs
}
type Condition  struct { Id, Name string }
type DamageType struct { Id, Name string }

// SRD standard sets ship as predefined values, for example:
var (
    SkillStealth = Skill{Id: "stealth", Name: "Stealth", Ability: AbilityDexterity}
    Slashing     = DamageType{Id: "slashing", Name: "Slashing"}
    Prone        = Condition{Id: "prone", Name: "Prone"}
    // ...
)
var SRDSkills []Skill   // the eighteen

// DamageAny is a wildcard used only as a key in damage.Mitigation maps (flat
// reduction or resistance that applies to every type). It is not a real damage
// type and must never be put on a damage part.
var DamageAny = DamageType{Id: "any", Name: "Any"}

// AbilityAny is a wildcard for qualifiers that apply to every ability (for
// example a buff to all saving throws). Like DamageAny it is not a concrete
// ability and must never be used where one is required (AbilityScores, check.Save).
const AbilityAny Ability = -1
```

A cyberpunk setting declares its own, governed by the same six abilities and
keyed the same way:

```go
// in the game, not in this module
var (
    SkillHacking = core.Skill{Id: "hacking", Name: "Hacking", Ability: core.AbilityIntelligence}
    Slow         = core.Condition{Id: "slow", Name: "Slow"}
    Thermal      = core.DamageType{Id: "thermal", Name: "Thermal"}
)
```

Per-entity state (which skills a creature is proficient in, which conditions it
has, which damage types it resists) is stored by Id, which is serialization- and
ECS-friendly. The full value is catalog metadata looked up by Id, through a
game map or `content.Registry[core.Skill]`. Abilities stay a closed enum
because the SRD fixes the six scores and every game keeps them.

### Vision and light

The rule-bearing part of sight, with geometry left to the game. The game
supplies the ambient light at the target and whether the target is within the
observer's sense range (both are spatial facts); the kernel applies the
darkvision and obscurement rules.

```go
type VisionType int  // VisionNormal (0), VisionDarkvision, VisionBlindsight, VisionTruesight, VisionTremorsense
type LightLevel int  // LightUnspecified (0), LightBright, LightDim, LightDark
type Visibility int  // VisibilityClear (0), VisibilityObscured, VisibilityBlocked; function-produced

// EffectiveLight applies a creature's vision to the ambient light: darkvision in
// range treats dark as dim and dim as bright; blindsight or truesight see
// regardless of light.
func EffectiveLight(ambient LightLevel, vision VisionType, withinRange bool) LightLevel

// SightVisibility turns effective light and conditions into a visibility state:
// dim light is Obscured (disadvantage on sight Perception); darkness, blinded,
// or an invisible target is Blocked (cannot see).
func SightVisibility(light LightLevel, blinded, targetInvisible bool) Visibility
```

`VisionNormal` is the zero value because plain sight is every creature's
baseline. A creature with several senses passes the one relevant to the
situation. Detection then resolves with primitives we already have: an active
search is `check.Contest` (the searcher's Perception against the hider's
Stealth), and passive detection compares the hider's Stealth result to the
searcher's `check.PassiveScore`. The visibility result feeds the vantage seam:
attacking what you cannot see is at disadvantage, and an unseen attacker has
advantage.

### Effect duration

How long an applied condition or effect lasts. It is named `EffectDuration`
(not `Duration`) to stay clear of `time.Duration` and to leave room for other
duration concepts, and it lives in `core` so `core.EffectDuration` does not
stutter and both `effect` and `content` can share it.

```go
type DurationKind int  // DurationUnspecified (0), DurationInstant, DurationRounds,
                        // DurationEndOfNextTurn, DurationConcentration, DurationUntilRemoved
type EffectDuration struct {
    Kind        DurationKind
    Rounds      int      // when DurationRounds
    SaveEnds    bool     // the bearer repeats the save at the end of each of its turns to end it early
    SaveAbility Ability  // which save, when SaveEnds
}

func RoundsInMinutes(m int) int  // m * 10   (1 round = 6 seconds)
func RoundsInHours(h int) int     // h * 600
```

The kernel describes the duration and supplies the primitives that end it; the
game holds each entity's active conditions and counts them down. `DurationRounds`
decrements at the end of the bearer's turn; `DurationConcentration` ends when the
caster's concentration drops (`combat.ConcentrationDc` rolled as a `check.Save`,
incapacitation, or a recast); `SaveEnds` makes a `check.Save` at the end of the
bearer's turn against the caster's original Dc, which the game stores with the
active effect since the Dc is caster-dependent. `DurationEndOfNextTurn` refers to
the bearer's next turn by default; the rarer "end of the caster's next turn" case
is left to the game until a consumer needs it modeled.

## Package: dice

Randomness and dice notation, isolated. Owns the advantage mechanic because
advantage is purely "roll two d20, keep the higher," a rolling concept that both
checks and attacks need. Stays a dependency-free leaf working in raw `int` die
faces.

```go
type Roller interface { IntN(n int) int }       // *math/rand/v2.Rand satisfies it (v2 spells it IntN)

type Expr struct { Count, Sides, Modifier int }  // for example 2d6+3
func Parse(s string) (Expr, error)
func (e Expr) Roll(r Roller) Result
func (e Expr) RollCritical(r Roller) Result      // doubles dice, not modifier
func (e Expr) Min() int
func (e Expr) Max() int
func (e Expr) Average() float64

type Result struct {
    Dice  []int   // individual die values, in roll order
    Total int      // sum of Dice plus any modifier
}

type Vantage int  // VantageNone, VantageAdvantage, VantageDisadvantage
func RollD20(r Roller, v Vantage) Result          // Result.Dice holds one or both d20s

// CombineVantage applies the SRD rule: any advantage source plus any
// disadvantage source cancel to a straight roll, regardless of count.
func CombineVantage(advantage, disadvantage bool) Vantage
```

A `Result` exposes individual dice so games can render rolls and so future
reroll and keep mechanics (4d6 drop lowest, Great Weapon Fighting, Halfling
Lucky) can inspect per-die values. Those mechanics are not built now, but the
types are shaped to allow them without breaking callers.

Reference roller implementations: a fair roller wrapping `*math/rand/v2.Rand`, and
a deterministic stub for tests. Biased, karmic, or difficulty-weighted rollers
are game policy and live in the game, not here.

Two fixed rollers support assist modes and best-case analysis. `Constant` clamps
to each die's range, so it is safe across the d20 and smaller damage dice:

```go
func Constant(face int) Roller   // every die shows min(face, sides)

var (
    Take10 = Constant(10)        // d20 shows 10; dice smaller than d10 clamp to their max
    Take20 = Constant(20)        // d20 shows 20; effectively the highest face for d4..d20
)
```

These are deterministic and slot into the difficulty seam (an easy mode can hand
players `Take10` while enemies keep a fair roller). They are package-level `var`s,
not `const`s, since `Roller` is an interface. Note "take 10" and "take 20" are
not SRD 5e rules (5e folds "take 10" into passive checks); these rollers serve
house rules, difficulty modes, and tests.

## Package: check

A first-class d20 check, unifying ability checks, skill checks, and saving
throws (all "d20 plus bonus against a difficulty class"), plus the opposed
contest used by Shove and Grapple. `check` is pure: it interprets a roll the
caller already made.

```go
type Check struct {
    Bonus core.Modifier   // ability modifier plus proficiency plus situational
    Dc    core.Dc
}
type Result struct {
    Roll    dice.Result
    Total   int
    Success bool
    Margin  int            // Total minus Dc
}
func (c Check) Resolve(roll dice.Result) Result   // pure; no Roller

// convenience constructors that compute Bonus via the core math:
func Ability(scores core.AbilityScores, a core.Ability, dc core.Dc) Check
func Skill(scores core.AbilityScores, s core.Skill, proficient bool, level core.Level, dc core.Dc) Check
func Save(scores core.AbilityScores, a core.Ability, proficient bool, level core.Level, dc core.Dc) Check

// opposed contest (2014-style Shove and Grapple); ties favor the responder:
type ContestResult struct { InitiatorTotal, ResponderTotal int; InitiatorWins bool }
func Contest(initiatorTotal, responderTotal int) ContestResult

// passive check: 10 + modifier, +5 for advantage, -5 for disadvantage.
// Generic across passive Perception, Investigation, and Insight.
func PassiveScore(modifier core.Modifier, v dice.Vantage) int
```

`Skill` reads the governing ability from `s.Ability`, so an SRD skill and a
game-defined custom skill resolve identically. `Vantage` is not stored on
`Check`: advantage is a property of how the d20 is rolled, so it is supplied to
`dice.RollD20` at the roll site, which is also where the difficulty seam lives:

```go
c    := check.Skill(scores, SkillHacking, proficient, level, dc)
roll := dice.RollD20(r, vantage)
res  := c.Resolve(roll)
```

Shove's resolution differs by edition (2014 is a Strength contest, 2024 is a
Strength saving throw); both are expressible here, and the push-or-prone
consequence is applied as an effect afterward.

## Package: combat

Attack resolution and the targeting rules (cover, range bands) that feed it.
Pure: it interprets already-rolled numbers and already-determined spatial facts.
It never consumes a `Roller` and never computes geometry. Critical hits occur on
a natural twenty only; the SRD has no confirmation roll, so the crit consequence
(doubled dice) lands at the damage step.

```go
type AttackOutcome int  // AttackMiss, AttackHit, AttackCritical
type AttackResult struct {
    Outcome     AttackOutcome
    NaturalRoll int
    Total       int
}
func ResolveAttack(naturalRoll int, mod core.Modifier, ac core.ArmorClass) AttackResult
func ConcentrationDc(damage int) core.Dc
```

`ac` is the defender's full current armor class, including all standing and
active temporary bonuses (armor, capped Dexterity, shield, natural armor, magic,
the Shield spell). The kernel has no touch-attack concept; SRD spell attacks
target full armor class. Situational cover is added separately by `Setup`.

### Cover and range

Cover and range bands map an already-determined spatial fact to its mechanical
consequence. The game supplies the facts (distance, line of sight, what grants
cover); geometry is the game's spatial layer.

```go
type Cover int  // CoverNone, CoverHalf, CoverThreeQuarters, CoverTotal
func (c Cover) AcBonus() core.Modifier       // +0, +2, +5
func (c Cover) DexSaveBonus() core.Modifier  // +0, +2, +5; cover also aids Dexterity saves
func (c Cover) BlocksTargeting() bool         // true only for CoverTotal

type Range struct { Normal, Long core.Distance }
type Band int  // BandNormal, BandLong, BandOutOf
func (r Range) Band(distance core.Distance) Band
func MeleeRange(reach core.Distance) Range    // {reach, reach}: no long band; beyond reach is out
```

### Unified attack setup

Melee and ranged share one path: melee reach is just a `Range` with no long
band, so a single `Attack` handles both. `Setup` runs before the roll, because
it decides the vantage to roll with and whether the attack is even possible.
`ResolveAttack` stays the pure post-roll step.

```go
type Attack struct {
    Range        Range          // ranged {80, 320}; melee MeleeRange(5) or MeleeRange(10)
    Distance     core.Distance
    TargetCover  Cover
    Advantage    bool           // all advantage sources, pre-aggregated by the game
    Disadvantage bool           // all disadvantage sources; long-range disadvantage is added here
}
type AttackSetup struct {
    Possible    bool             // false if beyond range or total cover
    Vantage     dice.Vantage
    EffectiveAc core.ArmorClass
}
func (a Attack) Setup(baseAc core.ArmorClass) AttackSetup
```

There is no melee/ranged split and no ranged-only field. The "ranged attack
while a hostile is within five feet" disadvantage is not kernel-evaluable (it
depends on geometry, line of sight, and the hostile's state), so the game folds
it into the generic `Disadvantage` input. Call flow:

```go
s := atk.Setup(baseAc)
if !s.Possible { /* out of range or fully covered */ }
roll := dice.RollD20(r, s.Vantage)
res  := combat.ResolveAttack(roll.Total, attackMod, s.EffectiveAc)
```

## Package: damage

The typed damage pipeline: roll typed damage from a hit, mitigate it against a
target's defenses, then apply it to a hit-point pool. Three pure stages (plus
the roll, which uses `dice`). Damage amounts are kept as `int` here.

```go
type DamagePart struct { Amount int; Type core.DamageType; Magical bool }  // Magical overcomes "resistance to nonmagical"
type Damage     struct { Parts []DamagePart }

type PartSpec struct { Dice dice.Expr; Type core.DamageType; Magical bool }
type Spec     struct { Parts []PartSpec }

// Roll keys off the attack outcome so the caller writes no if/else: a miss
// yields zero damage, a critical doubles the dice per part. bonus is the
// attacker's flat damage modifier (an ability modifier, or 0 for an off-hand
// attack or most spells); the game decides which ability per the attack type.
// It is added once to the primary part, so it carries that part's damage type
// for resistance and is not doubled on a critical.
func Roll(spec Spec, bonus core.Modifier, outcome combat.AttackOutcome, r dice.Roller) Damage

// Mitigation is everything that reduces or modifies incoming damage. Resistance,
// vulnerability, and immunity are SRD; FlatReduction is a generic extension
// (empty for fantasy games, used by a cyberpunk game's armor-class plus damage-reduction
// hybrid). All maps are keyed by DamageType.Id.
type Mitigation struct {
    Immune, Resist, Vulnerable map[string]bool
    ResistNonmagicalPhysical   bool
    // FlatReduction subtracts after resistance, flooring at zero per part. It is
    // per damage type, plus the wildcard core.DamageAny that applies to every
    // type: a part of type T is reduced by FlatReduction[T.Id] +
    // FlatReduction[core.DamageAny.Id] (e.g. Heavy Armor Master, adamantine).
    FlatReduction map[string]int
}
type Result struct {
    Raw    int
    Final  int
    ByType map[string]int
}
func ApplyMitigation(d Damage, m Mitigation) Result

// Hp application is separate, because hit points usually live in their own
// component. delta is a signed HP change: negative is damage (floors at zero,
// and surfaces the SRD state transitions), positive is healing (caps at max).
// Damage and healing share this one application point; only damage is typed,
// mitigated, and crit-doubled upstream, so healing is applied as a positive
// delta directly, never through ApplyMitigation.
type HpOutcome struct {
    Hp            core.HitPoints
    DroppedToZero bool
    InstantDeath  bool   // massive-damage rule (damage only)
}
func ApplyToHp(current, max core.HitPoints, delta int) HpOutcome

// Apply bundles the two pure post-roll stages (mitigate, then apply to Hp) into
// one call for the damage path, applying the mitigated total as a negative
// delta. The game writes back HpOutcome.Hp and can show the Result breakdown,
// and supplies the target's Mitigation and Hp from its own components.
func Apply(d Damage, m Mitigation, cur, max core.HitPoints) (HpOutcome, Result)
```

`Mitigation` is the bundle your design called "Defense"; it is named for what it
does (reduce or modify incoming damage), to avoid reading like armor class.

## Package: effect

A declarative vocabulary for the consequences an attack, spell, or action
applies, keyed to the outcome that triggers them. The kernel owns the numeric
consequences (damage via `damage`, healing as a signed HP delta, and the fold of
ongoing modifiers via `ModifierBonus`); condition, forced-movement, and the
storage of active modifiers are *described* here as data and *applied* by the
game, because they touch entity state, geometry, and per-game stat assembly. The
package executes nothing by itself.

```go
type Trigger int  // TriggerUnspecified (0), then OnHit, OnCrit, OnMiss, OnSaveFail, OnSaveSuccess, OnSave (either outcome), OnContestWin, OnContestLose, Always

type ConditionSpec struct { Condition core.Condition; Duration core.EffectDuration }
type MovementKind  int                                                // MovementUnspecified (0), MovementPush, MovementPull
type MovementSpec  struct { Distance core.Distance; Kind MovementKind } // relative to the source; game executes

// ModifierSpec is an ongoing numeric buff or debuff (Bless +1d4 to attacks and
// saves, Shield of Faith +2 AC, Hunter's Mark +1d6 damage). It is declarative:
// the game stores active modifiers (durations ticked like conditions) and folds
// them into the Bonus/AC it passes to check/combat via ModifierBonus. Vantage is
// deliberately not here; advantage/disadvantage flows through conditions and the
// game's vantage aggregation.
type ModifierTarget int     // ModTargetUnspecified (0), ModTargetAc, ModTargetAttack, ModTargetSave, ModTargetCheck, ModTargetDamage, ModTargetSpeed
type ModifierSource string  // game-declared, e.g. const SourceBless ModifierSource = "bless"
type ModifierSpec struct {
    Source   ModifierSource      // dedup key (same source does not stack, SRD 5e rule); also keys removal
    Targets  []ModifierTarget
    Ability  core.Ability        // core.AbilityAny = all; a specific ability narrows save/check targets
    Flat     int                 // +2 AC
    Dice     dice.Expr           // rolled per use, e.g. Bless 1d4; zero = none
    Duration core.EffectDuration
}
// ModifierBonus folds the active modifiers that apply to one roll, deduped by
// Source: each Source counts at most once per target (identical same-source
// entries collapse; highest wins on conflict). Distinct sources sum (5e: they
// stack). The game rolls the returned dice and adds flat to the roll's Bonus.
func ModifierBonus(active []ModifierSpec, target ModifierTarget, ability core.Ability) (flat int, dice []dice.Expr)

type Effect struct {
    Damage     *damage.Spec   // kernel rolls and applies via the damage pipeline
    HalfOnSave bool           // on a save-resolved effect: full damage on a failed save, half on a success (rolled once)
    Healing    *dice.Expr     // untyped HP gain; game adds any caster modifier, applies as a positive ApplyToHp delta (skips mitigation)
    Condition  *ConditionSpec // game applies
    Movement   *MovementSpec  // game executes (geometry)
    Modifier   *ModifierSpec  // ongoing numeric buff/debuff; game stores it active and folds via ModifierBonus
}
type ConditionalEffect struct { Trigger Trigger; Effect Effect }

// Outcome is a normalized resolution result, so the selector needs no
// dependency on combat or check. The game fills it from whichever resolution ran.
type Outcome struct {
    Hit, Crit  bool
    Saved      *bool   // nil if no save was involved
    ContestWon *bool   // nil if no contest was involved
}
func Triggered(effs []ConditionalEffect, outcome Outcome) []Effect
```

So a spell or weapon mastery is data, and the game runs a uniform loop: resolve,
call `Triggered`, then for each returned effect send damage through the `damage`
pipeline, apply healing as a positive `ApplyToHp` delta, store a modifier in its
active set (folded later via `ModifierBonus`), and apply conditions and movement
to its own entities.

```go
// Fireball:                {OnSave: fire damage, HalfOnSave: true}   (one effect, halved on a successful save)
// Sword + Push mastery:    {OnHit: slashing}, {OnHit: Movement(10 ft push)}
// Shove:                   {OnContestWin: Movement(5 ft)} or {OnContestWin: Condition(Prone)}
// Damage + slow hex:       {OnSaveFail: damage}, {OnSaveFail: Condition(Slow, 2 rounds)}
// Cure Wounds:             {Always: Healing(1d8)}                    (game adds the caster's modifier)
// Bless:                   {Always: Modifier{Source: bless, Targets:[AttackRoll,SaveThrow], Ability: AbilityAny, Dice: 1d4}}
```

The library deliberately stops here (the "lean" altitude): it provides the
vocabulary and the selector but no action engine that drives a turn end to end,
because such an engine would have to reach into entity state and geometry.

## Package: turn

Utilities that assist running turns, not a turn engine. The game owns the loop
and executes actions; this package provides an initiative timeline, the "is this
legal right now" bookkeeping, and the movement-cost rules. Everything keys off
opaque Ids the game supplies, so the package never touches entities or geometry.

The `Tracker` is an event timeline, not just a turn cursor. Creature turns and
recurring effect ticks are entries on one initiative-ordered, looping queue, so
`Current`/`Next` return a tagged `Event` rather than a bare Id. This lets a
persisting effect (a save-ends condition, a lasting area like a wall of fire)
fire at exactly its point in the order each round.

Granularity is the game's choice, because the combatant Id is opaque. A
party-as-one-turn game registers a single `"party"` combatant and lets any
member act during its `EventTurn` (per-member `Economy` is still individual; only
the party's *place in the order* is shared); a tactical game registers each
creature separately. A game that schedules no effects sees only `EventTurn`
events, so the same `Tracker` doubles as a plain turn cursor with no extra API.

```go
// initiative values to feed AddCombatant:
func StaticInitiative(dexScore core.AbilityScore) int            // 10 + Dexterity modifier
func RollInitiative(dexMod core.Modifier, r dice.Roller, v dice.Vantage) int

type EventKind int  // EventTurn, EventEffect
type Event struct {
    Kind   EventKind
    Round  int
    Actor  string    // creature whose turn it is, on EventTurn
    Effect *Active    // the effect that is due, on EventEffect
}

type TargetKind int  // TargetCreature, TargetArea
type Active struct {
    Id         string              // for Cancel, and for the game to map back to what the effect does
    Target     string              // creature or area Id; the game resolves it (and checks it is still alive)
    TargetKind TargetKind
    Duration   core.EffectDuration
    Remaining  int                 // rounds left, when DurationRounds
    SaveDc     core.Dc             // stored at apply time, for save-ends
}

type Timing int  // Before, After (relative to the anchor creature's turn)

type Tracker struct { /* combatant turns and effect ticks on one initiative-ordered, looping queue */ }
func NewTracker() *Tracker
func (t *Tracker) AddCombatant(id string, initiative int, dex core.AbilityScore)  // Dex breaks ties
func (t *Tracker) RemoveCombatant(id string)
func (t *Tracker) ScheduleEffect(a Active, anchor string, when Timing)            // recurs in that slot each round until Remaining hits 0
func (t *Tracker) Cancel(effectId string)                                         // concentration broke, dispelled, target died
func (t *Tracker) Current() Event
func (t *Tracker) Next() Event                                                    // turn and effect events interleaved, in order, looping by round
func (t *Tracker) Round() int
func (t *Tracker) Upcoming() []Event                                              // a copy of the upcoming timeline, for logging or rendering

// per-participant action economy for one round
type Economy struct {
    ActionUsed, BonusUsed, ReactionUsed bool
    MovementUsed  core.Distance
    MovementSpeed core.Distance
}
func (e *Economy) ResetTurn()                                     // at the start of the creature's turn
func (e Economy) CanReact() bool
func (e Economy) MovementRemaining() core.Distance
func (e Economy) CanMove(cost core.Distance) bool                 // cost <= MovementRemaining()

// movement cost and reach rules (pure; the game measures the path, the kernel
// turns a measured distance into a budget cost or a reach)
func DifficultTerrainCost(distance core.Distance) core.Distance  // 1 extra foot per foot, i.e. ×2
func DashDistance(speed core.Distance) core.Distance              // the Dash action adds your speed again
func StandUpCost(speed core.Distance) core.Distance               // standing from prone costs half your speed
func LongJump(str core.AbilityScore, runningStart bool) core.Distance  // STR feet with a run-up, half without
func HighJump(strMod core.Modifier, runningStart bool) core.Distance   // 3 + STR modifier feet with a run-up
```

The game's loop is then a clean `switch tracker.Next().Kind`: an `EventTurn`
hands control to a creature, an `EventEffect` says a scheduled effect is due now.
Two boundary points, both the game's job: on an `EventEffect` the game checks the
target is still alive (the kernel cannot know entity state) and, for a
`TargetArea`, finds who is currently in the area (geometry); and the effect's
payload (what it does) lives in the game, which maps `Active.Id` back to its own
data, so the tracker carries only identity, target, and timing. `ScheduleEffect`
anchors an effect to a creature's turn with `Before` or `After`, which is also
how caster-versus-bearer timing is expressed (anchor to whichever creature's turn
it should ride). Save-ends effects surface as an `EventEffect` whose `SaveDc` the
game rolls against with `check.Save`, calling `Cancel` on success; concentration
breaks and deaths are out-of-band `Cancel` calls. If an anchor creature is
removed, its effects keep their slot unless the game cancels them.

SRD 5e has no separate "move action": a creature has a movement allowance equal
to its `MovementSpeed` that it splits around its action, and the Dash action
grants more.
`Economy` models this as a budget rather than an action slot. The game measures
distances and terrain (geometry); the kernel converts a measured distance into a
budget cost and computes jump reach from Strength.

Reactions do not involve the `Tracker`: an opportunity attack or a readied
action fires inline on another creature's turn without advancing the order, so
it is purely `Economy` bookkeeping (one reaction per round, reset at the start of
your turn). A readied action is a reaction with a game-defined trigger. The game
decides what triggers a reaction and executes it.

## Package: content

Optional layer 2 shapes plus a small generic registry. Plain, pointer-free,
serializable structs that games may use, extend, or ignore. The kernel never
requires them.

```go
// slots available per spell level; index is the spell level (1..9), index 0 unused
type SpellSlots []int
// class level (1..20) -> slots at that level; an absent level means no casting yet,
// keeping late and partial casters free of filler zero rows
type SpellSlotProgression map[int]SpellSlots

type Class struct {
    Id, Name            string
    HitDie              int
    ProficientSaves     []core.Ability   // SRD classes have two; left open for custom classes
    SkillChoiceCount    int
    AvailableSkills     []core.Skill
    SpellcastingAbility core.Ability        // zero value means non-caster
    Slots               SpellSlotProgression
    // further fields as needed
}

type Race struct {
    Id, Name       string
    AbilityBonuses map[core.Ability]int       // optional; empty models SRD 5.2 and reskinned-race settings
    MovementSpeed  core.Distance              // walking speed; fly/swim/climb would be separate fields
    Traits         []string
}

// how a spell is resolved against its targets
type SpellResolution int  // ResolveUnspecified (0), ResolveSpellAttack, ResolveSave, ResolveAuto

// targeting descriptor: data the game interprets spatially (geometry stays in the game)
type RangeKind  int  // RangeUnspecified (0), RangeSelf, RangeTouch, RangeRanged
type TargetKind int  // TargetUnspecified (0), TargetSelf, TargetSingle, TargetMultiple, TargetArea
type AreaShape  int  // AreaNone (0) for non-area spells, then AreaSphere, AreaCone, AreaLine, AreaCube
type Targeting struct {
    Range      RangeKind
    Distance   core.Distance  // when RangeRanged
    Target     TargetKind
    MaxTargets int            // when TargetMultiple
    Shape      AreaShape      // when TargetArea
    AreaSize   core.Distance  // radius, length, or side
}

type Spell struct {
    Id, Name    string
    Level       int
    Resolution  SpellResolution
    SaveAbility core.Ability               // used when Resolution == ResolveSave; the Dc is the caster's, via core.SpellSaveDc
    Targeting   Targeting
    Effects     []effect.ConditionalEffect // damage and riders, keyed to the resolution outcome
    // further fields as needed
}

type Item struct { /* weapons, armor, gear; weapons carry a damage.Spec */ }

// Creature carries the stat-block fields the kernel pipeline reads; games extend
// it with their own (sprites, loot, behavior) or define their own type entirely.
type Creature struct {
    Id, Name      string
    Size          core.Size
    Abilities     core.AbilityScores
    ArmorClass    core.ArmorClass
    HitDice       dice.Expr          // rolled or averaged (dice.Expr.Average) for HP
    Mitigation    damage.Mitigation  // resistances/immunities the damage pipeline reads
    MovementSpeed core.Distance
    // attacks (each a damage.Spec), traits, and game data as needed
}

type Registry[T any] struct { /* Register, Get, All */ }
```

`Race.AbilityBonuses` being optional is what lets the same struct model SRD 5.1
(the race grants bonuses) and both SRD 5.2 and reskinned-race settings (bonuses come
from elsewhere) without a version fork.

The content shapes carry the fields the kernel pipeline reads as inputs (most
notably `Creature.Mitigation`, which `damage.ApplyMitigation` consumes, and the
hit-dice data that `NewHitDicePool` and `Creature.HitDice` feed), and otherwise
stay thin: flavor and game-specific data (sprites, loot, item tags, attachment
slots, the full class pick-lists) are added by each game, which embeds these
shapes or defines its own. Two consequences: cantrip and upcast scaling are not
modeled here, since the game computes the effective `damage.Spec` by level; and
the actual standard SRD content as data (specific classes, spells, monsters)
remains the deferred layer 3, not these shapes. The shapes are the home for what
the kernel consumes, not a complete content schema.

A `Spell` is self-describing, so the game casts it uniformly: read `Resolution`
to pick a spell attack roll (`core.SpellAttackBonus` vs armor class) or a target
saving throw (`check.Save` against the caster's `core.SpellSaveDc`), fill an
`effect.Outcome`, then `effect.Triggered` and apply. Fire Bolt is
`ResolveSpellAttack`; Fireball is `ResolveSave` with `SaveAbility =
AbilityDexterity`; Magic Missile is `ResolveAuto`. SRD 5e has no spell
resistance; the analogue, Magic Resistance, is advantage on saves against spells,
which the game supplies through the `Vantage` seam, so it needs no field here.

## Package: resource

Runtime per-character resource pools, as value types the game embeds in its
components (the way it composes `core.AbilityScores`). Spell slots and hit dice
have real structure and get specific types; everything else that is "a count
with a recharge trigger" (feature uses, ki or sorcery points, or a reskinned
tech resource and its expendable-use mechanic) shares one generic `Resource`. How they refill belongs to rest and
recovery (next section); this package is the resources themselves.

```go
type RestKind int  // RestShort, RestLong

// All three pools recover on a rest, so they share one behavioral interface.
type Restorable interface { Restore(rest RestKind) }

// spell slots: a per-level pool, built from a class's static slot counts.
// Cantrips cost no slot and never touch this pool. Index is the spell level
// (1..9); [0] is unused.
type SpellSlotPool struct {
    Max     [10]int
    Current [10]int
}
func NewSpellSlotPool(perLevel []int) SpellSlotPool   // from content.SpellSlotProgression at a class level
func (p *SpellSlotPool) Available(level int) int
func (p *SpellSlotPool) Expend(level int) bool        // spend one slot of that level; the caller upcasts by passing a higher level
func (p *SpellSlotPool) RestoreAll()                  // direct, e.g. Arcane Recovery
func (p *SpellSlotPool) Restore(rest RestKind)        // RestLong restores all; RestShort is a no-op

// hit dice: a pool keyed by die size (mixed sizes for a multiclass character)
type HitDicePool struct {
    Dice  map[int]int  // die size (e.g. 8 for d8) -> total count
    Spent map[int]int
}
func NewHitDicePool(level, hitDie int) HitDicePool    // a single-class character: level dice of size hitDie
func (p *HitDicePool) Available(die int) int
func (p *HitDicePool) Spend(die int) bool             // spent by the player to heal on a short rest
func (p *HitDicePool) Recover(n int)                  // direct: regain n spent dice
func (p HitDicePool) Total() int                       // sum of Dice across sizes
func (p *HitDicePool) Restore(rest RestKind)          // RestLong recovers up to half the pool; RestShort is a no-op
func HitDieHeal(die int, scores core.AbilityScores, r dice.Roller) core.HitPoints  // die roll + Constitution modifier; reads scores.Constitution

// generic limited-use resource
type ResourceId  string  // the game declares its own consts, e.g. const Ki ResourceId = "ki"
type RechargeRule int    // RechargeNone (0; never), RechargeShortRest, RechargeLongRest
type Resource struct {
    Id           ResourceId
    Max, Current int
    Recharge     RechargeRule  // describes WHEN it refills
}
func (r *Resource) Use(n int) bool        // spend; false if not enough
func (r *Resource) Restore(rest RestKind) // refill to Max iff this rest covers r.Recharge

// ResourceSet aggregates a character's generic resources into one value, keyed
// by the typed ResourceId. An entity-component-system holds one component per
// type, so the several generic resources (ki, Channel Divinity, a reskinned resource) cannot each
// be their own component; this bag is the single component they live in.
type ResourceSet struct {
    Items map[ResourceId]Resource
}
func (s *ResourceSet) Add(r Resource)
func (s *ResourceSet) Use(id ResourceId, n int) bool   // get-modify-store handled internally
func (s *ResourceSet) Get(id ResourceId) (Resource, bool)
func (s *ResourceSet) Restore(rest RestKind)           // sweep all, applying each one's rule
func (s *ResourceSet) All() []Resource                 // for rendering
```

Each `Restore` applies that pool's own rest rule: `SpellSlotPool` refills on
`RestLong`, `HitDicePool` recovers up to half the pool on `RestLong` (the
fraction to verify against SRD 5.2), and a `Resource` refills only when the rest
covers its `Recharge` rule, encoding that a long rest also triggers short-rest
recharges (`RestLong` covers both `RechargeShortRest` and `RechargeLongRest`,
`RestShort` covers only `RechargeShortRest`, `RechargeNone` never). `ResourceSet`
sweeps each `Resource`. Spending hit dice to heal is a separate player action
(`Spend` + `HitDieHeal`), not part of the rest sweep, and `HitDieHeal` reads
Constitution from `core.AbilityScores` so the wrong modifier cannot be passed.
`NewSpellSlotPool` takes a plain `[]int` (the slot counts the game pulls from
`content.SpellSlotProgression`), so `resource` depends only on `core` and `dice`.
Roll-based recharge (a monster's "Recharge 5-6") is turn-based, not rest-based,
with no consumer yet, so it is deferred; adding it later is additive.

Because the three pools share `Restorable`, a whole-character rest is one sweep
over them (then HP, which is a game component, not a resource):

```go
for _, r := range []resource.Restorable{&slots, &hitDice, &generics} {
    r.Restore(kind)
}
if kind == resource.RestLong { hp.Current = hp.Max }  // HP is game-owned
```

All four types (`SpellSlotPool`, `HitDicePool`, `Resource`, `ResourceSet`) have
exported fields and round-trip through `encoding/json` with no custom marshaler,
so a game serializes its whole party (these pools included) into a save without
extra accessors. This is consistent with the rest of the module's transparent
value types; the methods are conveniences that maintain invariants, not
encapsulation the game must go through. (`turn.Tracker` is mid-mission state, not
part of a between-mission save, so its serialization is left until a consumer
needs it.)

`ResourceId` is a typed string in the same spirit as the other named units: it
catches a wrong-typed key and groups the game's constants (`const Ki
resource.ResourceId = "ki"`), though untyped literals still convert, so the
safety comes from using the declared consts. `ResourceSet` bundles only the
homogeneous generic resources, because that is where the one-component-per-type
limit bites; the `SpellSlotPool` and `HitDicePool` are single per character and
stay their own components. The cross-component "all of a character's resources"
view that a rest needs (slots plus hit dice plus the set plus Hp) is a game-side
system, since it spans the game's component layout and Hp, which the kernel does
not model.

## Orchestration pattern

The module deliberately ships no action engine: resolving an attack or spell and
applying its effects interleaves kernel-pure math with reads and writes of the
game's own entities (armor class, hit points, conditions) and geometry, which
belong to the game. Instead the pieces are each a single call, and this is the
canonical loop a game writes once. `damage.Apply` bundles the pure damage stages;
`tracker.ScheduleEffect` handles the persisting effect; the `Damage` and
`Movement` arms are where the game reads and writes its own components.

```go
// resolve (attack shown; a save spell uses check.Save against the caster's Dc)
roll := dice.RollD20(r, setup.Vantage)
atk  := combat.ResolveAttack(roll.Total, mod, setup.EffectiveAc)
out  := effect.Outcome{Hit: atk.Outcome != combat.AttackMiss, Crit: atk.Outcome == combat.AttackCritical}

for _, e := range effect.Triggered(spell.Effects, out) {
    switch {
    case e.Damage != nil:
        dmg   := damage.Roll(*e.Damage, bonus, atk.Outcome, r)
        hp, _ := damage.Apply(dmg, tgt.Mitigation, tgt.Hp, tgt.MaxHp)  // game reads its component
        tgt.Hp = hp.Hp                                                 // game writes it back
    case e.Condition != nil:
        tracker.ScheduleEffect(turn.Active{
            Id: id, Target: tgtID, TargetKind: turn.TargetCreature,
            Duration: e.Condition.Duration, SaveDc: dc,
        }, casterID, turn.After)
    case e.Movement != nil:
        game.Push(tgtID, *e.Movement)   // geometry: game-owned
    }
}
```

A game keeps this loop in one place; the kernel keeps each step a pure,
testable call.

## SRD version stance

The kernel is version-neutral. The math in every kernel package is identical
between SRD 5.1 and SRD 5.2. Where a modeling choice depends on edition, it is
confined to layer 2 content fields. SRD 5.2 is the documented reference
baseline.

## Module consumption

The module is public at module path `github.com/trancecode/go-srd5e`, consumed
with a plain `go get github.com/trancecode/go-srd5e@latest`; pkg.go.dev serves
the API docs. For day-to-day local development across several repos, a consumer's
`go.mod` uses a replace directive pointing at the local checkout, so changes do
not require tagging on every edit:

```
replace github.com/trancecode/go-srd5e => ../go-srd5e
```

Once the module is stable, releases are tagged with semantic version tags (for
example `v0.1.0`), consumers pin a version, and the replace directive is dropped.
`CONSUMING.md` documents this for client code.

## Testing

Each package has table-driven unit tests. The math has known-correct outputs, so
tests assert exact values. Check, combat, damage, and effect resolution are pure
and test with no randomness (pass a `dice.Result` with a chosen total, or the
deterministic `Roller` stub from `dice` for the roll step). The kernel is pure
logic with no input or output and no mocking required.

## Out of scope

* Layer 3 reference content data. Deferred to a later opt-in subpackage.
* A turn or action engine that drives play end to end, game state, save and load,
  and rendering. These belong to each game; `turn` only assists.
* Geometry: distance, line of sight, adjacency, and movement execution.
* Difficulty policy. The module exposes the seams (roller injection and the
  vantage parameter) and stays ignorant of difficulty.
* Meta-systems with no SRD analogue, such as a cyberpunk game's escalation and
  detection systems. These are game-specific and built on top of the kernel.
