// Package content holds optional layer-2 shapes: serializable structs for
// classes, races, spells, items, and creatures that carry the fields the kernel
// pipeline reads as inputs, plus a small generic Registry. Games may use, extend,
// or ignore them; the kernel never requires them.
//
// The SRD 5.1 equipment ships as data, not just shapes: Weapons, Armor
// (armour and the shield), and Gear (ammunition and selected adventuring
// gear) each return their slice of Item, and StandardItems returns the
// union of the three; RegisterStandardItems loads that union into a
// Registry[Item]. Every value is pinned to the SRD text by fixtures under
// testdata/. Creatures, spells, and classes remain shapes only; standard
// content for those layers is still deferred.
package content
