// Package resource holds per-character resource pools as serializable value
// types the game embeds in its components: spell slots, hit dice, and a generic
// limited-use resource for everything that is "a count with a recharge trigger".
// All pools share the Restorable interface so a whole-character rest is one sweep.
package resource
