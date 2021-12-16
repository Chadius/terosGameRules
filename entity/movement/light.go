package movement

import "reflect"

// Light movement is land locked, but movement penalties are reduced (2 on rough terrain, 1 elsewhere)
type Light struct{}

// Name returns a human-readable name of this logic object.
func (m *Light) Name() string {
	return "light"
}

// GreaterThan returns true if the other movement logic has a greater rank (aka an improvement) over this one.
func (m *Light) GreaterThan(other Interface) bool {
	if other == nil {
		return true
	}

	if reflect.TypeOf(other).String() == "*movement.Foot" {
		return true
	}

	return false
}
