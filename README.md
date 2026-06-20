# go-srd5e

A reusable Go module that encodes the rules of the System Reference Document for
5th edition (SRD 5e, CC-BY-4.0), so multiple games can share one correct, tested
implementation of the d20 core instead of each reimplementing it.

**[API documentation (godoc)](https://pkg.go.dev/github.com/trancecode/go-srd5e)**

## What it is

`go-srd5e` is a rules *kernel*: the minimal, game-agnostic core of the rules,
expressed as pure functions and value types. It deliberately owns none of the
surrounding game (no content, game state, rendering, geometry, or difficulty
policy); those belong to each game. It is built to be shared by several games at
once, including a fantasy dungeon crawler and a cyberpunk tactics game, so it
stays neutral about setting and supplies only the mechanics.

## Status

Early. The foundation package `core` is implemented and tested; the rest of the
module is being built in dependency order. See the design spec and the
implementation plans under `docs/superpowers/`.

Planned package layout (dependency order):

```
core/     units, open value catalogs (Ability, Skill, Condition, DamageType), stat math, vision/light   [done]
dice/     Roller, Expr, Result, Vantage, RollD20, CombineVantage, fixed rollers
check/    Check, Contest, PassiveScore: ability checks, skill checks, saving throws (pure)
combat/   ResolveAttack, Cover, Range, Attack/Setup (pure)
damage/   typed Damage, Roll, Mitigation, ApplyMitigation, ApplyToHp, Apply
effect/   declarative effects (damage, healing, condition, movement, modifier), Triggered
turn/     initiative event timeline, action-economy and movement bookkeeping
resource/ SpellSlotPool, HitDicePool, Resource, ResourceSet, Restorable
content/  Class, Race, Spell, Item, Creature, Registry (layer 2, optional)
```

## Design and plans

The authoritative design lives under `docs/`:

* `docs/superpowers/specs/2026-06-18-go-srd5e-module-design.md` — the full API
  design and the binding principles. This is the source of truth for the
  module's shape; read it before changing any API.
* `docs/superpowers/plans/` — one implementation plan per package, executed in
  dependency order with test-driven development.

## Using it

`go get github.com/trancecode/go-srd5e@latest`. See [CONSUMING.md](CONSUMING.md)
for local cross-repo development (a `replace` directive).

## Conventions

Code follows [docs/styleguide.md](docs/styleguide.md). The one deviation from
idiomatic Go worth noting: acronyms use Java camel case (`Id`, `Hp`, `Ac`, `Dc`,
`Xp`), not all-caps initialisms.

## Licensing

The module encodes SRD 5e mechanics (SRD 5.1 and 5.2 are available under
CC-BY-4.0). Mechanics are reimplemented in code; SRD text is not reproduced
verbatim. This is a non-commercial project.
