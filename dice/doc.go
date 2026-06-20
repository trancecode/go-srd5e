// Package dice is the randomness and dice-notation layer of go-srd5e. It is a
// dependency-free leaf. Roller is the only behavioral interface; everything else
// is value types and functions that consume a Roller, so that randomness enters
// the rest of the module only here.
package dice
