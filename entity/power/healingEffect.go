package power

// HealingEffect is a power designed to restore hit points and cure ailments.
type HealingEffect struct {
	healingAdjustmentBasedOnUserMind HealingAdjustmentBasedOnUserMind
	hitPointsHealed                  int
}

// NewHealingEffect creates a new HealingEffect object.
func NewHealingEffect(hitPointsHealed int, adjustment HealingAdjustmentBasedOnUserMind) *HealingEffect {
	return &HealingEffect{
		healingAdjustmentBasedOnUserMind: adjustment,
		hitPointsHealed:                  hitPointsHealed,
	}
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
	return h.hitPointsHealed
}

// HealingAdjustmentBasedOnUserMind returns the value.
func (h *HealingEffect) HealingAdjustmentBasedOnUserMind() HealingAdjustmentBasedOnUserMind {
	return h.healingAdjustmentBasedOnUserMind
}
