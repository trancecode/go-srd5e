// Package turn assists running turns without being a turn engine. The game owns
// the loop and executes actions; this package provides initiative values, the
// action-economy budget (Economy), the pure movement and jump rules, and the
// Tracker — an initiative-ordered, looping timeline of creature turns and
// scheduled effect ticks. Everything keys off opaque Ids the game supplies, so
// the package never touches entities or geometry.
package turn
