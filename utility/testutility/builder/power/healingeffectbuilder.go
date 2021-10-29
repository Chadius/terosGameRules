package power

import "github.com/chadius/terosbattleserver/entity/power"

// HealingEffectOptions is used to create healing effects.
type HealingEffectOptions struct {
	hitPointsHealed                  int
	healingAdjustmentBasedOnUserMind power.HealingAdjustmentBasedOnUserMind
}

// HealingEffectBuilder creates a HealingEffectOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func HealingEffectBuilder() *HealingEffectOptions {
	return &HealingEffectOptions{
		hitPointsHealed:                  0,
		healingAdjustmentBasedOnUserMind: power.Full,
	}
}

// HitPointsHealed sets the amount of healed hit points.
func (h *HealingEffectOptions) HitPointsHealed(heal int) *HealingEffectOptions {
	h.hitPointsHealed = heal
	return h
}

// HealingAdjustmentBasedOnUserMindFull applies the user's Full SquaddieMind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindFull() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = power.Full
	return h
}

// HealingAdjustmentBasedOnUserMindHalf applies Half of the user's SquaddieMind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindHalf() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = power.Half
	return h
}

// HealingAdjustmentBasedOnUserMindZero applies None of the user's SquaddieMind bonus to healing effects.
func (h *HealingEffectOptions) HealingAdjustmentBasedOnUserMindZero() *HealingEffectOptions {
	h.healingAdjustmentBasedOnUserMind = power.Zero
	return h
}

// Build uses the HealingEffectOptions to create a HealingEffect.
func (h *HealingEffectOptions) Build() *power.HealingEffect {
	newHealingEffect := &power.HealingEffect{
		HealingHitPointsHealed:                  h.hitPointsHealed,
		HealingHealingAdjustmentBasedOnUserMind: h.healingAdjustmentBasedOnUserMind,
	}
	return newHealingEffect
}
