package powerattackforecast

import (
	"github.com/chadius/terosgamerules/entity/powersource"
	"github.com/chadius/terosgamerules/entity/powerusagescenario"
	"github.com/chadius/terosgamerules/usecase/repositories"
	"github.com/chadius/terosgamerules/usecase/squaddiestats"
)

// AttackerContextStrategy describes the shape of an AttackContext.
type AttackerContextStrategy interface {
	Calculate(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error
	TotalToHitBonus() int
	IsCounterAttack() bool
	RawDamage() int
	ExtraBarrierBurn() int
	CanCritical() bool
	CriticalHitThreshold() int
	CriticalHitDamage() int
	PowerSourceLogic() powersource.Interface
}

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	isCounterAttack  bool
	totalToHitBonus  int
	rawDamage        int
	powerSourceLogic powersource.Interface

	extraBarrierBurn     int
	canCritical          bool
	criticalHitThreshold int
	criticalHitDamage    int

	offenseStrategy squaddiestats.CalculateSquaddieOffenseStatsStrategy
}

// NewAttackerContext creates a new object.
func NewAttackerContext(strategy squaddiestats.CalculateSquaddieOffenseStatsStrategy) *AttackerContext {
	return &AttackerContext{offenseStrategy: &squaddiestats.CalculateSquaddieOffenseStats{}}
}

// Calculate sets the fields according to the attacker and power they are using.
func (context *AttackerContext) Calculate(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	var err error

	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)

	context.powerSourceLogic = power.PowerSourceLogic()

	context.extraBarrierBurn = power.ExtraBarrierBurn()

	context.rawDamage, err = context.calculateRawDamage(setup, repositories)
	if err != nil {
		return err
	}

	err = context.calculateToHitBonus(setup, repositories)
	if err != nil {
		return err
	}

	context.calculateCriticalHit(setup, repositories)
	return nil
}

func (context *AttackerContext) calculateToHitBonus(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	context.isCounterAttack = setup.IsCounterAttack
	if context.isCounterAttack {
		counterAttackHitBonus, counterAttackHitBonusError := context.offenseStrategy.GetSquaddieCounterAttackAimWithPower(setup.UserID, setup.PowerID, repositories)
		if counterAttackHitBonusError != nil {
			return counterAttackHitBonusError
		}
		context.totalToHitBonus = counterAttackHitBonus
		return nil
	}

	hitBonus, hitBonusError := context.offenseStrategy.GetSquaddieAimWithPower(setup.UserID, setup.PowerID, repositories)
	if hitBonusError != nil {
		return hitBonusError
	}
	context.totalToHitBonus = hitBonus
	return nil
}

func (context *AttackerContext) calculateRawDamage(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	rawDamage, damageErr := context.offenseStrategy.GetSquaddieRawDamageWithPower(setup.UserID, setup.PowerID, repositories)
	if damageErr != nil {
		return 0, damageErr
	}
	return rawDamage, nil
}

func (context *AttackerContext) calculateCriticalHit(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)
	if !power.CanAttack() {
		return nil
	}
	context.canCritical = power.CanCriticallyHit()

	canCounter, counterErr := context.offenseStrategy.GetSquaddieCanCriticallyHitWithPower(setup.UserID, setup.PowerID, repositories)
	if counterErr != nil {
		return counterErr
	}

	context.canCritical = canCounter
	if context.canCritical == false {
		return nil
	}

	critThreshold, critThresholdError := context.offenseStrategy.GetSquaddieCriticalThresholdWithPower(setup.UserID, setup.PowerID, repositories)
	if critThresholdError != nil {
		return critThresholdError
	}
	context.criticalHitThreshold = critThreshold

	criticalDamage, criticalDamageErr := context.offenseStrategy.GetSquaddieCriticalRawDamageWithPower(setup.UserID, setup.PowerID, repositories)
	if criticalDamageErr != nil {
		return counterErr
	}
	context.criticalHitDamage = criticalDamage
	return nil
}

// TotalToHitBonus gets the field value.
func (context *AttackerContext) TotalToHitBonus() int {
	return context.totalToHitBonus
}

// IsCounterAttack gets the field value.
func (context *AttackerContext) IsCounterAttack() bool {
	return context.isCounterAttack
}

// RawDamage gets the field value.
func (context *AttackerContext) RawDamage() int {
	return context.rawDamage
}

// PowerSourceLogic gets the field value.
func (context *AttackerContext) PowerSourceLogic() powersource.Interface {
	return context.powerSourceLogic
}

// ExtraBarrierBurn gets the field value.
func (context *AttackerContext) ExtraBarrierBurn() int {
	return context.extraBarrierBurn
}

// CanCritical gets the field value.
func (context *AttackerContext) CanCritical() bool {
	return context.canCritical
}

// CriticalHitThreshold gets the field value.
func (context *AttackerContext) CriticalHitThreshold() int {
	return context.criticalHitThreshold
}

// CriticalHitDamage gets the field value.
func (context *AttackerContext) CriticalHitDamage() int {
	return context.criticalHitDamage
}
