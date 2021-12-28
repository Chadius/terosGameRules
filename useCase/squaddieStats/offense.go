package squaddiestats

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/powerinterface"
	"github.com/chadius/terosbattleserver/entity/squaddieinterface"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
)

// CalculateSquaddieOffenseStatsStrategy describes the shape of an attack calculator.
type CalculateSquaddieOffenseStatsStrategy interface {
	GetSquaddieAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieRawDamageWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieCriticalThresholdWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieCriticalRawDamageWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieCanCounterAttackWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (bool, error)
	GetSquaddieCounterAttackAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieExtraBarrierBurnWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieCanCriticallyHitWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (bool, error)
	GetHitPointsHealedWithPower(squaddieID, powerID, targetID string, repos *repositories.RepositoryCollection) (int, error)
	CanSquaddieCounterWithEquippedWeapon(squaddieID string, repos *repositories.RepositoryCollection) (bool, error)
}

// CalculateSquaddieOffenseStats returns information about a squaddie's attacks with a given power.
type CalculateSquaddieOffenseStats struct{}

// GetSquaddieAimWithPower returns the to hit bonus against a target without dodge.
func (c *CalculateSquaddieOffenseStats) GetSquaddieAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.Aim() + powerToMeasure.ToHitBonus(), nil
}

// GetSquaddieRawDamageWithPower returns the amount of damage that will be dealt to an unprotected target.
func (c *CalculateSquaddieOffenseStats) GetSquaddieRawDamageWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.PowerSourceLogic().RawDamage(squaddie) + powerToMeasure.DamageBonus(), nil
}

// GetSquaddieCriticalThresholdWithPower returns the critical hit threshold the squaddie needs to beat in order to crit.
func (c *CalculateSquaddieOffenseStats) GetSquaddieCriticalThresholdWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.CanCriticallyHit() != true {
		newError := fmt.Errorf("cannot critical hit with power, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return 0, newError
	}

	return powerToMeasure.CriticalHitThreshold(), nil
}

// GetSquaddieCriticalRawDamageWithPower returns the total critical hit damage the squaddie deals (includes raw damage.)
func (c *CalculateSquaddieOffenseStats) GetSquaddieCriticalRawDamageWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	rawDamage, rawDamageErr := c.GetSquaddieRawDamageWithPower(squaddieID, powerID, repos)
	if rawDamageErr != nil {
		return 0, rawDamageErr
	}

	if powerToMeasure.CanCriticallyHit() != true {
		newError := fmt.Errorf("cannot critical hit with power, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return 0, newError
	}

	return rawDamage + powerToMeasure.ExtraCriticalHitDamage(), nil
}

// GetSquaddieCanCounterAttackWithPower returns true if the squaddie can counterattack with this power.
func (c *CalculateSquaddieOffenseStats) GetSquaddieCanCounterAttackWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (bool, error) {
	_, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return false, err
	}

	return powerToMeasure.CanCounterAttack(), nil
}

// GetSquaddieCounterAttackAimWithPower returns the to critical hit bonus against a target without dodge.
func (c *CalculateSquaddieOffenseStats) GetSquaddieCounterAttackAimWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	counterAttackPenalty, counterAttackErr := powerToMeasure.CounterAttackPenalty()
	if counterAttackErr != nil {
		return 0, counterAttackErr
	}

	return squaddie.Aim() + powerToMeasure.ToHitBonus() + counterAttackPenalty, nil
}

// GetSquaddieExtraBarrierBurnWithPower returns the amount of extra barrier burn that will be dealt to a target with a barrier.
func (c *CalculateSquaddieOffenseStats) GetSquaddieExtraBarrierBurnWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	_, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.ExtraBarrierBurn(), nil
}

// GetSquaddieCanCriticallyHitWithPower returns true if the squaddie can crit with this power.
func (c *CalculateSquaddieOffenseStats) GetSquaddieCanCriticallyHitWithPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (bool, error) {
	_, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return false, err
	}

	return powerToMeasure.CanCritical(), nil
}

// GetHitPointsHealedWithPower returns the actual number of hit points healed.
func (c *CalculateSquaddieOffenseStats) GetHitPointsHealedWithPower(squaddieID, powerID, targetID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddieToHeal, healingPower, err := getSquaddieAndHealingPower(squaddieID, powerID, repos)
	target := repos.SquaddieRepo.GetSquaddieByID(targetID)

	if err != nil {
		return 0, nil
	}

	hitPoints := healingPower.HealingLogic().CalculateExpectedHeal(squaddieToHeal, healingPower.HitPointsHealed(), target)
	return hitPoints, nil
}

func getSquaddie(squaddieID string, repos *repositories.RepositoryCollection) (squaddieinterface.Interface, error) {
	squaddie := repos.SquaddieRepo.GetOriginalSquaddieByID(squaddieID)
	if squaddie == nil {
		newError := fmt.Errorf("squaddie could not be found, id: %s", squaddieID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}
	return squaddie, nil
}

func getHealingPower(powerID string, repos *repositories.RepositoryCollection) (powerinterface.Interface, error) {
	power := repos.PowerRepo.GetPowerByID(powerID)
	if power == nil {
		newError := fmt.Errorf("power could not be found, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}
	if !power.CanHeal() {
		newError := fmt.Errorf("cannot heal with power, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}
	return power, nil
}

func getAttackPower(powerID string, repos *repositories.RepositoryCollection) (powerinterface.Interface, error) {
	power := repos.PowerRepo.GetPowerByID(powerID)
	if power == nil {
		newError := fmt.Errorf("power could not be found, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}
	if !power.CanAttack() {
		newError := fmt.Errorf("cannot attack with power, id: %s", powerID)
		utility.Log(newError.Error(), 0, utility.Error)
		return nil, newError
	}
	return power, nil
}

func getSquaddieAndHealingPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (squaddieinterface.Interface, powerinterface.Interface, error) {
	squaddie, squaddieErr := getSquaddie(squaddieID, repos)
	if squaddieErr != nil {
		return nil, nil, squaddieErr
	}

	power, powerErr := getHealingPower(powerID, repos)
	if powerErr != nil {
		return nil, nil, powerErr
	}

	return squaddie, power, nil
}

func getSquaddieAndAttackPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (squaddieinterface.Interface, powerinterface.Interface, error) {
	squaddie, squaddieErr := getSquaddie(squaddieID, repos)
	if squaddieErr != nil {
		return nil, nil, squaddieErr
	}

	power, powerErr := getAttackPower(powerID, repos)
	if powerErr != nil {
		return nil, nil, powerErr
	}

	return squaddie, power, nil
}

// CanSquaddieCounterWithEquippedWeapon returns true if the squaddie can use the currently equipped weapon for counterattacks.
func (c *CalculateSquaddieOffenseStats) CanSquaddieCounterWithEquippedWeapon(squaddieID string, repos *repositories.RepositoryCollection) (bool, error) {
	squaddie := repos.SquaddieRepo.GetOriginalSquaddieByID(squaddieID)
	equippedPowerID := squaddie.GetEquippedPowerID()
	if equippedPowerID == "" {
		newError := fmt.Errorf("squaddie has no equipped power, %s", squaddieID)
		utility.Log(newError.Error(), 0, utility.Error)
		return false, newError
	}

	canCounter, counterErr := c.GetSquaddieCanCounterAttackWithPower(squaddieID, equippedPowerID, repos)
	return canCounter, counterErr
}
