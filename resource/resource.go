package resource

import "sort"

// ResourceId is a game-declared resource key, e.g. const Ki ResourceId = "ki".
type ResourceId string

// RechargeRule describes when a resource refills.
type RechargeRule int

const (
	RechargeNone RechargeRule = iota
	RechargeShortRest
	RechargeLongRest
)

// Resource is a generic limited-use pool: a count with a recharge trigger.
type Resource struct {
	Id           ResourceId
	Max, Current int
	Recharge     RechargeRule
}

// Use spends n; false if not enough (or n is negative).
func (r *Resource) Use(n int) bool {
	if n < 0 || r.Current < n {
		return false
	}
	r.Current -= n
	return true
}

// Restore refills to Max if this rest covers the resource's recharge rule.
func (r *Resource) Restore(rest RestKind) {
	if covers(rest, r.Recharge) {
		r.Current = r.Max
	}
}

// covers encodes that a long rest also triggers short-rest recharges.
func covers(rest RestKind, rule RechargeRule) bool {
	switch rule {
	case RechargeShortRest:
		return rest == RestShort || rest == RestLong
	case RechargeLongRest:
		return rest == RestLong
	default: // RechargeNone
		return false
	}
}

// ResourceSet aggregates a character's generic resources into one value (one ECS
// component), keyed by ResourceId.
type ResourceSet struct {
	Items map[ResourceId]Resource
}

// Add inserts or replaces a resource by its Id.
func (s *ResourceSet) Add(r Resource) {
	if s.Items == nil {
		s.Items = map[ResourceId]Resource{}
	}
	s.Items[r.Id] = r
}

// Use spends n of one resource (get-modify-store internally); false if absent or
// not enough.
func (s *ResourceSet) Use(id ResourceId, n int) bool {
	r, ok := s.Items[id]
	if !ok || !r.Use(n) {
		return false
	}
	s.Items[id] = r
	return true
}

// Get returns a resource by Id.
func (s *ResourceSet) Get(id ResourceId) (Resource, bool) {
	r, ok := s.Items[id]
	return r, ok
}

// Restore sweeps every resource, applying each one's recharge rule.
func (s *ResourceSet) Restore(rest RestKind) {
	for id, r := range s.Items {
		r.Restore(rest)
		s.Items[id] = r
	}
}

// All returns the resources in Id order, for rendering.
func (s *ResourceSet) All() []Resource {
	out := make([]Resource, 0, len(s.Items))
	for _, r := range s.Items {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })
	return out
}
