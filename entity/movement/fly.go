package movement

import "reflect"

// Fly movement costs 1 movement per space, and can pass over (but not stop on) sky tiles.
type Fly struct{}

// Name returns a human-readable name of this logic object.
func (m *Fly) Name() string {
	return "fly"
}

// GreaterThan returns true if the other movement logic has a greater rank (aka an improvement) over this one.
func (m *Fly) GreaterThan(other Interface) bool {
	if other == nil {
		return true
	}

	if reflect.TypeOf(other).String() == "*movement.Foot" ||
		reflect.TypeOf(other).String() == "*movement.Light" {
		return true
	}

	return false
}
