package movement

import "reflect"

// Teleport movement costs 1 movement per space, and can pass over (but not stop on) sky and wall tiles.
type Teleport struct{}

// Name returns a human-readable name of this logic object.
func (m *Teleport) Name() string {
	return "teleport"
}

// GreaterThan returns true if the other movement logic has a greater rank (aka an improvement) over this one.
func (m *Teleport) GreaterThan(other Interface) bool {
	if other == nil {
		return true
	}

	if reflect.TypeOf(other).String() == "*movement.Teleport" {
		return false
	}

	return true
}
