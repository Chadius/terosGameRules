package power

// HealingEffectOptions is used to create healing effects.
type HealingEffectOptions struct {
	hitPointsHealed                  int
	healingAdjustmentBasedOnUserMind HealingAdjustmentBasedOnUserMind
}

// HealingEffectBuilder creates a HealingEffectOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func HealingEffectBuilder() *HealingEffectOptions {
	return &HealingEffectOptions{
		hitPointsHealed:                  0,
		healingAdjustmentBasedOnUserMind: Full,
	}
}

// HitPointsHealed sets the amount of healed hit points.
func (h *HealingEffectOptions) HitPointsHealed(heal int) *HealingEffectOptions {
	h.hitPointsHealed = heal
	return h
}

// HealingAdjustmentBasedOnUserMindFull applies the user's Full mind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindFull() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = Full
	return h
}

// HealingAdjustmentBasedOnUserMindHalf applies Half of the user's mind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindHalf() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = Half
	return h
}

// HealingAdjustmentBasedOnUserMindZero applies None of the user's mind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindZero() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = Zero
	return h
}

// Build uses the HealingEffectOptions to create a healingEffect.
func (h *HealingEffectOptions) Build() *HealingEffect {
	newHealingEffect := NewHealingEffect(
		h.hitPointsHealed,
		h.healingAdjustmentBasedOnUserMind,
	)
	return newHealingEffect
}
