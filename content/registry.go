package content

// Registry is a small generic store of content keyed by string Id. The zero
// value is usable. All() returns items in registration order.
type Registry[T any] struct {
	items map[string]T
	order []string
}

// NewRegistry returns an empty registry.
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{items: map[string]T{}}
}

// Register adds or replaces an item by Id (replacing keeps its original order
// position).
func (r *Registry[T]) Register(id string, v T) {
	if r.items == nil {
		r.items = map[string]T{}
	}
	if _, ok := r.items[id]; !ok {
		r.order = append(r.order, id)
	}
	r.items[id] = v
}

// Get returns the item for an Id.
func (r *Registry[T]) Get(id string) (T, bool) {
	v, ok := r.items[id]
	return v, ok
}

// All returns the items in registration order.
func (r *Registry[T]) All() []T {
	out := make([]T, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.items[id])
	}
	return out
}
