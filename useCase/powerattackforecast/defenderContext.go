package powerattackforecast

import (
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
)

// DefenderContextStrategy describes the shape of a DefenderContext.
type DefenderContextStrategy interface {
	Calculate(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error
	TargetID() string
	TotalToHitPenalty() int
	HitPoints() int
	ArmorResistance() int
	BarrierResistance() int
}

// DefenderContext lists the target's relevant information when under attack
type DefenderContext struct {
	targetID          string
	totalToHitPenalty int
	hitPoints         int
	armorResistance   int
	barrierResistance int
	defenseStrategy   squaddiestats.CalculateSquaddieDefenseStatsStrategy
}

// NewDefenderContext creates a new object.
func NewDefenderContext(targetID string, strategy squaddiestats.CalculateSquaddieDefenseStatsStrategy) *DefenderContext {
	return &DefenderContext{targetID: targetID, defenseStrategy: strategy}
}

// Calculate figures out how the Squaddie can resist the attack.
func (context *DefenderContext) Calculate(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) error {
	var err error

	context.totalToHitPenalty, err = context.calculateTotalToHitPenalty(setup, repositories)
	if err != nil {
		return err
	}

	context.armorResistance, err = context.calculateArmorResistance(setup, repositories)
	if err != nil {
		return err
	}

	context.barrierResistance, err = context.calculateBarrierResistance(setup, repositories)
	if err != nil {
		return err
	}

	context.hitPoints, err = context.calculateHitPoints(setup, repositories)
	if err != nil {
		return err
	}
	return nil
}

func (context *DefenderContext) calculateTotalToHitPenalty(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	evade, err := context.defenseStrategy.GetSquaddieToHitPenaltyAgainstPower(context.targetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return evade, nil
}

func (context *DefenderContext) calculateArmorResistance(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	armor, err := context.defenseStrategy.GetSquaddieArmorAgainstPower(context.targetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return armor, nil
}

func (context *DefenderContext) calculateBarrierResistance(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	barrier, err := context.defenseStrategy.GetSquaddieBarrierAgainstPower(context.targetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return barrier, nil
}

func (context *DefenderContext) calculateHitPoints(setup *powerusagescenario.Setup, repositories *repositories.RepositoryCollection) (int, error) {
	hitPoints, err := context.defenseStrategy.GetSquaddieCurrentHitPoints(context.targetID, setup.PowerID, repositories)
	if err != nil {
		return 0, err
	}
	return hitPoints, nil
}

// TargetID is a getter.
func (context *DefenderContext) TargetID() string {
	return context.targetID
}

// TotalToHitPenalty is a getter.
func (context *DefenderContext) TotalToHitPenalty() int {
	return context.totalToHitPenalty
}

// HitPoints is a getter.
func (context *DefenderContext) HitPoints() int {
	return context.hitPoints
}

// ArmorResistance is a getter.
func (context *DefenderContext) ArmorResistance() int {
	return context.armorResistance
}

// BarrierResistance is a getter.
func (context *DefenderContext) BarrierResistance() int {
	return context.barrierResistance
}
