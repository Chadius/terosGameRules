package levelup

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility"
)

// ImproveSquaddieStrategy describes objects that can upgrade squaddie stats.
type ImproveSquaddieStrategy interface {
	ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) error
}

// ImproveSquaddieClass describes objects that can upgrade squaddie stats.
type ImproveSquaddieClass struct{}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func (i *ImproveSquaddieClass) ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) error {
	if squaddieToImprove.HasAddedClass(benefit.ClassID()) == false {
		newError := fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddieToImprove.Name(), benefit.ClassID())
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}
	if squaddieToImprove.IsClassLevelAlreadyUsed(benefit.ID()) {
		newError := fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddieToImprove.Name(), benefit.ClassID(), benefit.ID())
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}
	squaddieToImprove.SetBaseClassIfNoBaseClass(benefit.ClassID())
	squaddieToImprove.MarkLevelUpBenefitAsConsumed(benefit.ClassID(), benefit.ID())

	improveSquaddieStats(benefit, squaddieToImprove)
	refreshSquaddiePowers(benefit, squaddieToImprove)
	improveSquaddieMovement(benefit, squaddieToImprove)
	return nil
}

func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Defense != nil {
		squaddieToImprove.ImproveDefense(
			benefit.MaxHitPoints(),
			benefit.Dodge(),
			benefit.Deflect(),
			benefit.MaxBarrier(),
			benefit.Armor(),
		)
	}

	if benefit.Offense != nil {
		squaddieToImprove.ImproveOffense(benefit.Aim(), benefit.Strength(), benefit.Mind())
	}
}
func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.PowerChanges != nil {
		squaddieToImprove.RemovePowerReferences(benefit.PowersLost())
		squaddieToImprove.AddPowerReferences(benefit.PowersGained())
	}
}
func improveSquaddieMovement(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Movement == nil {
		return
	}

	squaddieToImprove.ImproveMovement(
		benefit.Movement.MovementDistance(),
		benefit.Movement.MovementType(),
		benefit.Movement.CanHitAndRun(),
	)
}
