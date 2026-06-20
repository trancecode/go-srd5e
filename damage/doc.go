// Package damage is the typed-damage pipeline: roll typed damage from a hit,
// mitigate it against a target's defenses, then apply it to a hit-point pool as
// a signed delta. Rolling consumes a dice.Roller; mitigation and HP application
// are pure. Damage amounts are plain ints.
package damage
