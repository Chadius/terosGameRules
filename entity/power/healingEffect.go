package power

// HealingEffect is a power designed to restore hit points and cure ailments.
type HealingEffect struct {
	hitPointsHealed int
}

// NewHealingEffect creates a new HealingEffect object.
func NewHealingEffect(hitPointsHealed int) *HealingEffect {
	return &HealingEffect{
		hitPointsHealed: hitPointsHealed,
	}
}

// HitPointsHealed returns the value.
func (h *HealingEffect) HitPointsHealed() int {
	return h.hitPointsHealed
}
