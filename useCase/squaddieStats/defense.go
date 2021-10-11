package squaddiestats

import (
	"github.com/cserrant/terosbattleserver/entity/power"
	"github.com/cserrant/terosbattleserver/usecase/repositories"
)

// GetSquaddieToHitPenaltyAgainstPower returns how well the squaddie can evade the attack
func GetSquaddieToHitPenaltyAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.PowerType == power.Physical {
		return squaddie.Defense.Dodge, nil
	}

	if powerToMeasure.PowerType == power.Spell {
		return squaddie.Defense.Deflect, nil
	}

	return 0, nil
}

// GetSquaddieArmorAgainstPower returns how much armor the squaddie has to absorb the power's damage.
func GetSquaddieArmorAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	if powerToMeasure.PowerType == power.Physical {
		return squaddie.Defense.Armor, nil
	}
	return 0, nil
}

// GetSquaddieBarrierAgainstPower returns how much barrier the squaddie has to resist the power's damage.
func GetSquaddieBarrierAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, _, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.Defense.CurrentBarrier, nil
}

// GetSquaddieCurrentHitPoints returns the squaddie's current hit points.
func GetSquaddieCurrentHitPoints(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, _, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.Defense.CurrentHitPoints, nil
}
