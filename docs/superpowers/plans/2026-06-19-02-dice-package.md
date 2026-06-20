# dice package implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `dice` package: randomness and dice notation — the `Roller` seam, dice expressions, rolling, the d20 advantage mechanic, and fixed rollers.

**Architecture:** `dice` is a dependency-free leaf (it does not import `core`). `Roller` is the one behavioral interface; everything else is value types and pure-ish functions that consume a `Roller`. Randomness enters only through `Roller`, so resolvers in later packages stay pure.

**Tech Stack:** Go 1.25, standard library only (`math/rand/v2`, `regexp`, `strconv`, `strings`). Table-driven tests using deterministic rollers.

Plan 2 of 8 (see the design spec `docs/superpowers/specs/2026-06-18-go-srd5e-module-design.md`, the `## Package: dice` section).

## Global Constraints

- Module path `github.com/trancecode/go-srd5e`, Go `1.25`. `dice` imports nothing from this module (leaf).
- Acronyms in Java camel case: the interface method is `IntN` (the `math/rand/v2` spelling), not `Intn`. `RollD20`, not `RollD20` variants.
- `Roller` is `interface { IntN(n int) int }`; `*math/rand/v2.Rand` satisfies it directly.
- Value types (`Expr`, `Result`) have exported fields and round-trip through `encoding/json`.
- `Vantage` zero is `VantageNone` (valid neutral). Raw die faces stay `int`.
- Errors are phrased `<context>: <reason>` (per `docs/styleguide.md`), e.g. `fmt.Errorf("parsing dice %q: invalid format", s)`.
- TDD: failing test first, watch it fail, implement minimally, watch it pass, `gofmt -l`, commit. Work on `main` directly (no feature branch).

---

### Task 1: Roller, Result, and fixed rollers

**Files:**
- Create: `dice/doc.go`
- Create: `dice/roller.go`
- Test: `dice/roller_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `type Roller interface { IntN(n int) int }`
  - `type Result struct { Dice []int; Total int }`
  - `func Constant(face int) Roller` — every die shows `min(face, sides)`
  - `var Take10 Roller` (= `Constant(10)`), `var Take20 Roller` (= `Constant(20)`)
  - `func NewRoller(seed uint64) Roller` — a fair roller seeded deterministically, wrapping `*math/rand/v2.Rand`

- [ ] **Step 1: Write the failing test**

```go
package dice

import "testing"

func TestConstant(t *testing.T) {
	// Constant(face).IntN(n) yields min(face,n)-1, so a die of `sides` shows min(face,sides).
	cases := []struct {
		face, sides, wantIntN int
	}{
		{10, 20, 9}, // d20 shows 10
		{20, 20, 19}, // d20 shows 20
		{20, 6, 5},  // d6 clamps to 6
		{10, 8, 7},  // d8 clamps to 8
		{3, 6, 2},   // d6 shows 3
	}
	for _, c := range cases {
		if got := Constant(c.face).IntN(c.sides); got != c.wantIntN {
			t.Errorf("Constant(%d).IntN(%d) = %d, want %d", c.face, c.sides, got, c.wantIntN)
		}
	}
}

func TestTakeRollers(t *testing.T) {
	if got := Take10.IntN(20); got != 9 {
		t.Errorf("Take10.IntN(20) = %d, want 9 (face 10)", got)
	}
	if got := Take20.IntN(20); got != 19 {
		t.Errorf("Take20.IntN(20) = %d, want 19 (face 20)", got)
	}
}

func TestNewRollerInRange(t *testing.T) {
	r := NewRoller(42)
	for i := 0; i < 1000; i++ {
		v := r.IntN(6)
		if v < 0 || v >= 6 {
			t.Fatalf("NewRoller IntN(6) = %d, out of [0,6)", v)
		}
	}
	// Deterministic for a given seed.
	a, b := NewRoller(7), NewRoller(7)
	for i := 0; i < 10; i++ {
		if a.IntN(100) != b.IntN(100) {
			t.Fatal("NewRoller not deterministic for equal seeds")
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./dice/ -run 'TestConstant|TestTakeRollers|TestNewRoller'`
Expected: FAIL (undefined: Constant, Roller, ...).

- [ ] **Step 3: Write `dice/doc.go`**

```go
// Package dice is the randomness and dice-notation layer of go-srd5e. It is a
// dependency-free leaf. Roller is the only behavioral interface; everything else
// is value types and functions that consume a Roller, so that randomness enters
// the rest of the module only here.
package dice
```

- [ ] **Step 4: Write `dice/roller.go`**

```go
package dice

import "math/rand/v2"

// Roller is the source of randomness. *math/rand/v2.Rand satisfies it directly.
// It is always passed per call, never stored, so value types stay serializable.
type Roller interface {
	IntN(n int) int
}

// Result is the outcome of a roll: the individual die faces in roll order, and
// the total including any modifier.
type Result struct {
	Dice  []int
	Total int
}

// constRoller always reports the same face, clamped to each die's range.
type constRoller int

// IntN returns min(face, n) - 1, so a die of n sides shows min(face, n).
func (c constRoller) IntN(n int) int {
	face := int(c)
	if face > n {
		face = n
	}
	if face < 1 {
		face = 1
	}
	return face - 1
}

// Constant returns a Roller on which every die shows min(face, sides). Useful
// for assist modes (take 10, take 20), best-case analysis, and tests.
func Constant(face int) Roller { return constRoller(face) }

// Take10 and Take20 are fixed rollers: a d20 shows 10 or 20, smaller dice clamp
// to their maximum. Not SRD rules (5e uses passive checks); for house rules,
// difficulty modes, and tests.
var (
	Take10 Roller = Constant(10)
	Take20 Roller = Constant(20)
)

// NewRoller returns a fair, deterministically seeded roller backed by
// math/rand/v2 (a PCG source).
func NewRoller(seed uint64) Roller {
	return rand.New(rand.NewPCG(seed, seed))
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./dice/ -run 'TestConstant|TestTakeRollers|TestNewRoller'`
Expected: PASS.

- [ ] **Step 6: gofmt and commit**

```bash
gofmt -w dice/
git add dice/doc.go dice/roller.go dice/roller_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "dice: Roller, Result, and fixed rollers"
```

---

### Task 2: Expr, Parse, and bounds

**Files:**
- Create: `dice/expr.go`
- Test: `dice/expr_test.go`

**Interfaces:**
- Consumes: nothing (no randomness in this task).
- Produces:
  - `type Expr struct { Count, Sides, Modifier int }`
  - `func Parse(s string) (Expr, error)` — parses `"2d6+3"`, `"d8"` (count defaults to 1), `"1d20"`, `"2d6-1"`
  - `func (e Expr) Min() int`, `func (e Expr) Max() int`, `func (e Expr) Average() float64`

- [ ] **Step 1: Write the failing test**

```go
package dice

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want Expr
	}{
		{"2d6+3", Expr{2, 6, 3}},
		{"1d20", Expr{1, 20, 0}},
		{"d8", Expr{1, 8, 0}},
		{"2d6-1", Expr{2, 6, -1}},
		{" 3d4 ", Expr{3, 4, 0}},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Parse(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
	for _, bad := range []string{"", "abc", "2x6", "d", "2d", "d0", "2d6+"} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) expected error, got nil", bad)
		}
	}
}

func TestBounds(t *testing.T) {
	e := Expr{2, 6, 3} // 2d6+3
	if e.Min() != 5 {
		t.Errorf("Min = %d, want 5", e.Min())
	}
	if e.Max() != 15 {
		t.Errorf("Max = %d, want 15", e.Max())
	}
	if e.Average() != 10.0 {
		t.Errorf("Average = %v, want 10.0", e.Average())
	}
	d20 := Expr{1, 20, 0}
	if d20.Average() != 10.5 {
		t.Errorf("d20 Average = %v, want 10.5", d20.Average())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./dice/ -run 'TestParse|TestBounds'`
Expected: FAIL (undefined: Parse, Expr).

- [ ] **Step 3: Write `dice/expr.go`**

```go
package dice

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Expr is a dice expression NdS+M: Count dice of Sides faces, plus Modifier.
type Expr struct {
	Count, Sides, Modifier int
}

var diceRe = regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)

// Parse reads a dice expression such as "2d6+3", "1d20", "d8" (count defaults to
// 1), or "2d6-1".
func Parse(s string) (Expr, error) {
	m := diceRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Expr{}, fmt.Errorf("parsing dice %q: invalid format", s)
	}
	count := 1
	if m[1] != "" {
		count, _ = strconv.Atoi(m[1])
	}
	sides, _ := strconv.Atoi(m[2])
	if sides < 1 {
		return Expr{}, fmt.Errorf("parsing dice %q: sides must be at least 1", s)
	}
	mod := 0
	if m[3] != "" {
		mod, _ = strconv.Atoi(m[3])
	}
	return Expr{Count: count, Sides: sides, Modifier: mod}, nil
}

// Min is the lowest possible total (every die shows 1, plus the modifier).
func (e Expr) Min() int { return e.Count + e.Modifier }

// Max is the highest possible total (every die shows its max, plus the modifier).
func (e Expr) Max() int { return e.Count*e.Sides + e.Modifier }

// Average is the expected total.
func (e Expr) Average() float64 {
	return float64(e.Count)*(float64(e.Sides)+1)/2 + float64(e.Modifier)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./dice/ -run 'TestParse|TestBounds'`
Expected: PASS.

- [ ] **Step 5: gofmt and commit**

```bash
gofmt -w dice/
git add dice/expr.go dice/expr_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "dice: Expr, Parse, and bounds"
```

---

### Task 3: Roll and RollCritical

**Files:**
- Modify: `dice/expr.go`
- Test: `dice/roll_test.go`

**Interfaces:**
- Consumes: `Expr` (Task 2), `Roller`, `Result` (Task 1).
- Produces:
  - `func (e Expr) Roll(r Roller) Result`
  - `func (e Expr) RollCritical(r Roller) Result` — doubles the dice count, not the modifier

- [ ] **Step 1: Write the failing test**

```go
package dice

import (
	"reflect"
	"testing"
)

// seqRoller returns predetermined IntN values in order (a die face is value+1).
type seqRoller struct {
	vals []int
	i    int
}

func (s *seqRoller) IntN(int) int {
	v := s.vals[s.i]
	s.i++
	return v
}

func TestRoll(t *testing.T) {
	// 2d6+1 with a roller yielding IntN 2 then 4 -> faces 3 and 5 -> total 9.
	r := &seqRoller{vals: []int{2, 4}}
	got := Expr{2, 6, 1}.Roll(r)
	if !reflect.DeepEqual(got.Dice, []int{3, 5}) {
		t.Errorf("Dice = %v, want [3 5]", got.Dice)
	}
	if got.Total != 9 {
		t.Errorf("Total = %d, want 9", got.Total)
	}
}

func TestRollCritical(t *testing.T) {
	// 1d8+5 crit: dice doubled to 2d8, modifier added once. Constant(8) -> both 8.
	got := Expr{1, 8, 5}.RollCritical(Constant(8))
	if !reflect.DeepEqual(got.Dice, []int{8, 8}) {
		t.Errorf("crit Dice = %v, want [8 8]", got.Dice)
	}
	if got.Total != 21 { // 8 + 8 + 5
		t.Errorf("crit Total = %d, want 21", got.Total)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./dice/ -run 'TestRoll|TestRollCritical'`
Expected: FAIL (Roll/RollCritical undefined).

- [ ] **Step 3: Append to `dice/expr.go`**

```go
// Roll rolls each die and sums, adding the modifier once. Result.Dice holds the
// individual faces in roll order.
func (e Expr) Roll(r Roller) Result {
	dice := make([]int, e.Count)
	sum := 0
	for i := range dice {
		face := r.IntN(e.Sides) + 1
		dice[i] = face
		sum += face
	}
	return Result{Dice: dice, Total: sum + e.Modifier}
}

// RollCritical rolls a critical hit: the number of dice is doubled, the modifier
// is added once (SRD rule: double the dice, not the modifier).
func (e Expr) RollCritical(r Roller) Result {
	res := Expr{Count: e.Count * 2, Sides: e.Sides}.Roll(r)
	res.Total += e.Modifier
	return res
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./dice/ -run 'TestRoll|TestRollCritical'`
Expected: PASS.

- [ ] **Step 5: gofmt and commit**

```bash
gofmt -w dice/
git add dice/expr.go dice/roll_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "dice: Roll and RollCritical"
```

---

### Task 4: Vantage, RollD20, and CombineVantage

**Files:**
- Create: `dice/vantage.go`
- Test: `dice/vantage_test.go`

**Interfaces:**
- Consumes: `Roller`, `Result` (Task 1).
- Produces:
  - `type Vantage int` with `VantageNone Vantage = 0`, `VantageAdvantage`, `VantageDisadvantage`
  - `func RollD20(r Roller, v Vantage) Result` — `Result.Dice` holds one die (none) or both (advantage/disadvantage); `Total` is the kept die
  - `func CombineVantage(advantage, disadvantage bool) Vantage`

- [ ] **Step 1: Write the failing test**

```go
package dice

import (
	"reflect"
	"testing"
)

func TestRollD20(t *testing.T) {
	// none: single die.
	got := RollD20(&seqRoller{vals: []int{14}}, VantageNone) // face 15
	if !reflect.DeepEqual(got.Dice, []int{15}) || got.Total != 15 {
		t.Errorf("none = %+v, want Dice [15] Total 15", got)
	}
	// advantage: roll 15 and 10, keep 15; both dice recorded.
	got = RollD20(&seqRoller{vals: []int{14, 9}}, VantageAdvantage)
	if !reflect.DeepEqual(got.Dice, []int{15, 10}) || got.Total != 15 {
		t.Errorf("advantage = %+v, want Dice [15 10] Total 15", got)
	}
	// disadvantage: roll 15 and 10, keep 10.
	got = RollD20(&seqRoller{vals: []int{14, 9}}, VantageDisadvantage)
	if !reflect.DeepEqual(got.Dice, []int{15, 10}) || got.Total != 10 {
		t.Errorf("disadvantage = %+v, want Dice [15 10] Total 10", got)
	}
}

func TestCombineVantage(t *testing.T) {
	cases := []struct {
		adv, dis bool
		want     Vantage
	}{
		{false, false, VantageNone},
		{true, false, VantageAdvantage},
		{false, true, VantageDisadvantage},
		{true, true, VantageNone}, // any advantage + any disadvantage cancel
	}
	for _, c := range cases {
		if got := CombineVantage(c.adv, c.dis); got != c.want {
			t.Errorf("CombineVantage(%v,%v) = %v, want %v", c.adv, c.dis, got, c.want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./dice/ -run 'TestRollD20|TestCombineVantage'`
Expected: FAIL (undefined: RollD20, Vantage).

- [ ] **Step 3: Write `dice/vantage.go`**

```go
package dice

// Vantage is the advantage state of a d20 roll. Advantage rolls two d20 and
// keeps the higher; disadvantage keeps the lower.
type Vantage int

const (
	VantageNone Vantage = iota
	VantageAdvantage
	VantageDisadvantage
)

// RollD20 rolls a d20 under the given vantage. With advantage or disadvantage it
// rolls twice; Result.Dice holds both faces and Total is the kept one. There is
// no modifier on a raw d20, so Total is the kept die.
func RollD20(r Roller, v Vantage) Result {
	a := r.IntN(20) + 1
	switch v {
	case VantageAdvantage, VantageDisadvantage:
		b := r.IntN(20) + 1
		keep := a
		if (v == VantageAdvantage && b > keep) || (v == VantageDisadvantage && b < keep) {
			keep = b
		}
		return Result{Dice: []int{a, b}, Total: keep}
	default:
		return Result{Dice: []int{a}, Total: a}
	}
}

// CombineVantage applies the SRD rule: any advantage source plus any
// disadvantage source cancel to a straight roll, regardless of count.
func CombineVantage(advantage, disadvantage bool) Vantage {
	switch {
	case advantage == disadvantage:
		return VantageNone
	case advantage:
		return VantageAdvantage
	default:
		return VantageDisadvantage
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./dice/ -run 'TestRollD20|TestCombineVantage'`
Expected: PASS.

- [ ] **Step 5: Whole package + commit**

```bash
go test ./dice/
gofmt -w dice/
git add dice/vantage.go dice/vantage_test.go
git -c user.name="Claude Code" -c user.email="herve.quiroz+claude@gmail.com" commit -m "dice: Vantage, RollD20, and CombineVantage"
```

---

## Self-review notes

Spec coverage for `## Package: dice`:
- `Roller` (with `IntN`), `Result`, fixed rollers `Constant`/`Take10`/`Take20`, fair `NewRoller`: Task 1.
- `Expr`, `Parse`, `Min`/`Max`/`Average`: Task 2.
- `Roll`, `RollCritical` (crit doubles dice not modifier): Task 3.
- `Vantage`, `RollD20`, `CombineVantage`: Task 4.

`dice` imports only the standard library; no `core` dependency, preserving the leaf. The `seqRoller` test helper lives in test files only.
