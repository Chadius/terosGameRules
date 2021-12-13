package power

// HealingEffectOptions is used to create healing effects.
type HealingEffectOptions struct {
	hitPointsHealed int
}

// HealingEffectBuilder creates a HealingEffectOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func HealingEffectBuilder() *HealingEffectOptions {
	return &HealingEffectOptions{
		hitPointsHealed: 0,
	}
}

// HitPointsHealed sets the amount of healed hit points.
func (h *HealingEffectOptions) HitPointsHealed(heal int) *HealingEffectOptions {
	h.hitPointsHealed = heal
	return h
}

// Build uses the HealingEffectOptions to create a healingEffect.
func (h *HealingEffectOptions) Build() *HealingEffect {
	newHealingEffect := NewHealingEffect(
		h.hitPointsHealed,
	)
	return newHealingEffect
}
