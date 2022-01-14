package squaddiestats

import (
	"github.com/chadius/terosgamerules/usecase/repositories"
)

// CalculateSquaddieDefenseStatsStrategy describes the shape of the defense calculator.
type CalculateSquaddieDefenseStatsStrategy interface {
	GetSquaddieToHitPenaltyAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieArmorAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieBarrierAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
	GetSquaddieCurrentHitPoints(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error)
}

// CalculateSquaddieDefenseStats determines how a squaddie can evade a given attack
type CalculateSquaddieDefenseStats struct{}

// GetSquaddieToHitPenaltyAgainstPower returns how well the squaddie can evade the attack
func (c *CalculateSquaddieDefenseStats) GetSquaddieToHitPenaltyAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.PowerSourceLogic().ToHitPenalty(squaddie), nil
}

// GetSquaddieArmorAgainstPower returns how well the squaddie can evade the attack
func (c *CalculateSquaddieDefenseStats) GetSquaddieArmorAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.PowerSourceLogic().ArmorResistance(squaddie), nil
}

// GetSquaddieBarrierAgainstPower returns how much barrier the squaddie has to resist the power's damage.
func (c *CalculateSquaddieDefenseStats) GetSquaddieBarrierAgainstPower(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, powerToMeasure, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return powerToMeasure.PowerSourceLogic().BarrierResistance(squaddie), nil
}

// GetSquaddieCurrentHitPoints returns the squaddie's current hit points.
func (c *CalculateSquaddieDefenseStats) GetSquaddieCurrentHitPoints(squaddieID, powerID string, repos *repositories.RepositoryCollection) (int, error) {
	squaddie, _, err := getSquaddieAndAttackPower(squaddieID, powerID, repos)
	if err != nil {
		return 0, err
	}

	return squaddie.CurrentHitPoints(), nil
}
