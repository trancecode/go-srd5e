// Package combat resolves attacks and the targeting rules (cover, range bands)
// that feed them. It is pure: it interprets already-rolled numbers and
// already-determined spatial facts, never consuming a Roller or computing
// geometry. Critical hits occur on a natural 20 only.
package combat
