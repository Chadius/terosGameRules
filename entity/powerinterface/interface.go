package powerinterface

import (
	"github.com/chadius/terosbattleserver/entity/healing"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/powersource"
)

// Interface shapes the power.
type Interface interface {
	ID() string
	HitPointsHealed() int
	Name() string
	HealingLogic() healing.Interface
	CanAttack() bool
	ToHitBonus() int
	DamageBonus() int
	ExtraBarrierBurn() int
	CounterAttackPenaltyReduction() int
	CanCritical() bool
	CriticalHitThresholdBonus() int
	ExtraCriticalHitDamage() int
	CanBeEquipped() bool
	CanCounterAttack() bool
	CanPowerTargetFriend() bool
	CanPowerTargetFoe() bool
	CanPowerTargetSelf() bool
	PowerSourceLogic() powersource.Interface
	GetReference() *powerreference.Reference
	CanHeal() bool
	CounterAttackPenalty() (int, error)
	CanCriticallyHit() bool
	CriticalHitThreshold() int
}
