package movement

// Foot movement is the slowest type of movement, and it is land locked.
type Foot struct{}

// Name returns a human-readable name of this logic object.
func (m *Foot) Name() string {
	return "foot"
}

// GreaterThan returns true if the other movement logic has a greater rank (aka an improvement) over this one.
func (m *Foot) GreaterThan(other Interface) bool {
	if other == nil {
		return true
	}

	return false
}
