// Package content holds optional layer-2 shapes: serializable structs for
// classes, races, spells, items, and creatures that carry the fields the kernel
// pipeline reads as inputs, plus a small generic Registry. Games may use, extend,
// or ignore them; the kernel never requires them. The standard SRD content as
// data is deferred (layer 3); these are the shapes, not a complete schema.
package content
