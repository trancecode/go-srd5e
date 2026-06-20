// Package effect is the declarative vocabulary for what an attack, spell, or
// action applies: damage, healing, a condition, forced movement, or an ongoing
// numeric modifier, each keyed to the resolution outcome that triggers it. The
// package executes nothing. It provides the data types, the Triggered selector
// (which effects fire for an outcome), and ModifierBonus (folding active
// modifiers into a roll's bonus). The game applies the selected effects to its
// own entities, geometry, and stat assembly.
package effect
