package powerattackforecast

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
)

// AttackerContext lists the attacker's relevant information when attacking
type AttackerContext struct {
	IsCounterAttack bool
	TotalToHitBonus int

	RawDamage  int
	DamageType power.DamageType

	ExtraBarrierBurn int

	CanCritical          bool
	CriticalHitThreshold int
	CriticalHitDamage    int
}

func (context *AttackerContext) calculate(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	var err error

	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)

	context.DamageType = power.PowerType
	context.ExtraBarrierBurn = power.AttackEffect.ExtraBarrierBurn

	context.RawDamage, err = context.calculateRawDamage(setup, repositories)
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
	context.IsCounterAttack = setup.IsCounterAttack
	if context.IsCounterAttack {
		counterAttackHitBonus, counterAttackHitBonusError := squaddiestats.GetSquaddieCounterAttackAimWithPower(setup.UserID, setup.PowerID, repositories)
		if counterAttackHitBonusError != nil {
			return counterAttackHitBonusError
		}
		context.TotalToHitBonus = counterAttackHitBonus
		return nil
	}

	hitBonus, hitBonusError := squaddiestats.GetSquaddieAimWithPower(setup.UserID, setup.PowerID, repositories)
	if hitBonusError != nil {
		return hitBonusError
	}
	context.TotalToHitBonus = hitBonus
	return nil
}

func (context *AttackerContext) calculateRawDamage(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	rawDamage, damageErr := squaddiestats.GetSquaddieRawDamageWithPower(setup.UserID, setup.PowerID, repositories)
	if damageErr != nil {
		return 0, damageErr
	}
	return rawDamage, nil
}

func (context *AttackerContext) calculateCriticalHit(setup powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	power := repositories.PowerRepo.GetPowerByID(setup.PowerID)
	if power.AttackEffect == nil {
		return nil
	}
	context.CanCritical = power.AttackEffect.CanCriticallyHit()

	canCounter, counterErr := squaddiestats.GetSquaddieCanCriticallyHitWithPower(setup.UserID, setup.PowerID, repositories)
	if counterErr != nil {
		return counterErr
	}

	context.CanCritical = canCounter
	if context.CanCritical == false {
		return nil
	}

	critThreshold, critThresholdError := squaddiestats.GetSquaddieCriticalThresholdWithPower(setup.UserID, setup.PowerID, repositories)
	if critThresholdError != nil {
		return critThresholdError
	}
	context.CriticalHitThreshold = critThreshold

	counterDamage, counterDamageErr := squaddiestats.GetSquaddieCriticalRawDamageWithPower(setup.UserID, setup.PowerID, repositories)
	if counterDamageErr != nil {
		return counterErr
	}
	context.CriticalHitDamage = counterDamage
	return nil
}
