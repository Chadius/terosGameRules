package powercommit

// HealResult shows the effects of recovery abilities.
type HealResult struct {
	hitPointsRestored int
}

// HitPointsRestored is a getter.
func (h *HealResult) HitPointsRestored() int {
	return h.hitPointsRestored
}

// HealResultBuilder is used to build heal results.
type HealResultBuilder struct {
	hitPointsRestored int
}

// NewHealResultBuilder creates a new HealResultBuilder object.
func NewHealResultBuilder() *HealResultBuilder {
	return &HealResultBuilder{
		hitPointsRestored: 0,
	}
}

// HitPointsRestored sets the field
func (hr *HealResultBuilder) HitPointsRestored(hitPoints int) *HealResultBuilder {
	hr.hitPointsRestored = hitPoints
	return hr
}

// Build returns a HealResult
func (hr *HealResultBuilder) Build() *HealResult {
	return &HealResult{
		hr.hitPointsRestored,
	}
}
