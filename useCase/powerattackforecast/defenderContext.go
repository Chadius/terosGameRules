package powerattackforecast

import (
	"github.com/cserrant/terosbattleserver/entity/powerusagescenario"
	"github.com/cserrant/terosbattleserver/usecase/repositories"
	"github.com/cserrant/terosbattleserver/usecase/squaddiestats"
)

// DefenderContext lists the target's relevant information when under attack
type DefenderContext struct {
	TargetID string

	TotalToHitPenalty int

	HitPoints         int
	ArmorResistance   int
	BarrierResistance int
}

func (context *DefenderContext) calculate(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	var err error

	context.TotalToHitPenalty, err = context.calculateTotalToHitPenalty(setup, repositories)
	if err != nil {
		return err
	}

	context.ArmorResistance, err = context.calculateArmorResistance(setup, repositories)
	if err != nil {
		return err
	}

	context.BarrierResistance, err = context.calculateBarrierResistance(setup, repositories)
	if err != nil {
		return err
	}

	context.HitPoints, err = context.calculateHitPoints(setup, repositories)
	if err != nil {
		return err
	}
	return nil
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	evade, err := squaddiestats.GetSquaddieToHitPenaltyAgainstPower(context.TargetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return evade, nil
}

func (context *DefenderContext) calculateArmorResistance(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	armor, err := squaddiestats.GetSquaddieArmorAgainstPower(context.TargetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return armor, nil
}

func (context *DefenderContext) calculateBarrierResistance(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	barrier, err := squaddiestats.GetSquaddieBarrierAgainstPower(context.TargetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return barrier, nil
}

func (context *DefenderContext) calculateHitPoints(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	hitPoints, err := squaddiestats.GetSquaddieCurrentHitPoints(context.TargetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return hitPoints, nil
}
