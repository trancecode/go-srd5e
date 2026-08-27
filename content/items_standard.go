package content

// StandardItems returns the SRD equipment this module ships as content: the
// weapons, then the armour and shield, then the gear.
func StandardItems() []Item {
	var all []Item
	all = append(all, Weapons()...)
	all = append(all, Armor()...)
	all = append(all, Gear()...)
	return all
}

// RegisterStandardItems adds every standard item to a registry by id.
func RegisterStandardItems(r *Registry[Item]) {
	for _, it := range StandardItems() {
		r.Register(it.Id, it)
	}
}
