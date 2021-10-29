package power

// HealingEffect is a power designed to restore hit points and cure ailments.
type HealingEffect struct {
	HealingHealingAdjustmentBasedOnUserMind HealingAdjustmentBasedOnUserMind `json:"healing_adjustment_based_on_user_mind" yaml:"healing_adjustment_based_on_user_mind"`
	HealingHitPointsHealed                  int                              `json:"hit_points_healed" yaml:"hit_points_healed"`
}

// HealingAdjustmentBasedOnUserMind indicates how much the user's SquaddieMind should be adjusted.
type HealingAdjustmentBasedOnUserMind string

// User's SquaddieMind is added to most healing abilities (Full). But it may be at a Half bonus or doesn't affect (Zero)
const (
	Full HealingAdjustmentBasedOnUserMind = "full"
	Half HealingAdjustmentBasedOnUserMind = "half"
	Zero HealingAdjustmentBasedOnUserMind = "zero"
)

// HitPointsHealed returns the value.
func (h *HealingEffect) HitPointsHealed() int {
	return h.HealingHitPointsHealed
}

// HealingAdjustmentBasedOnUserMind returns the value.
func (h *HealingEffect) HealingAdjustmentBasedOnUserMind() HealingAdjustmentBasedOnUserMind {
	return h.HealingHealingAdjustmentBasedOnUserMind
}
