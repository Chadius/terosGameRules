package levelup

import (
	"fmt"
	"github.com/chadius/terosgamerules/entity/levelupbenefit"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
	"github.com/chadius/terosgamerules/utility"
)

// ImproveSquaddieStrategy describes objects that can upgrade squaddie stats.
type ImproveSquaddieStrategy interface {
	ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove squaddieinterface.Interface) error
}

// ImproveSquaddieClass describes objects that can upgrade squaddie stats.
type ImproveSquaddieClass struct{}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func (i *ImproveSquaddieClass) ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove squaddieinterface.Interface) error {
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

func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove squaddieinterface.Interface) {
	squaddieToImprove.ImproveDefense(
		benefit.MaxHitPoints(),
		benefit.Dodge(),
		benefit.Deflect(),
		benefit.MaxBarrier(),
		benefit.Armor(),
	)

	squaddieToImprove.ImproveOffense(benefit.Aim(), benefit.Strength(), benefit.Mind())
}

func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove squaddieinterface.Interface) {
	squaddieToImprove.RemovePowerReferences(benefit.PowersLost())
	squaddieToImprove.AddPowerReferences(benefit.PowersGained())
}

func improveSquaddieMovement(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove squaddieinterface.Interface) {
	squaddieToImprove.ImproveMovement(
		benefit.MovementDistance(),
		benefit.CanHitAndRun(),
		benefit.MovementLogic(),
	)
}
