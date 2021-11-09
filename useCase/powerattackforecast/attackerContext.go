package powerattackforecast

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
)

// AttackerContextStrategy describes the shape of an AttackContext.
type AttackerContextStrategy interface {
	Calculate(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error
	TotalToHitBonus() int
	IsCounterAttack() bool
	RawDamage() int
	DamageType() power.DamageType
	ExtraBarrierBurn() int
	CanCritical() bool
	CriticalHitThreshold() int
	CriticalHitDamage() int
}

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	isCounterAttack      bool
	totalToHitBonus      int
	rawDamage            int
	damageType           power.DamageType
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

	context.damageType = power.Type()
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

// DamageType gets the field value.
func (context *AttackerContext) DamageType() power.DamageType {
	return context.damageType
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
