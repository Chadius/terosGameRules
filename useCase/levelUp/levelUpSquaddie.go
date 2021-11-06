package levelup

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Defense != nil {
		squaddieToImprove.Defense.SquaddieMaxHitPoints = squaddieToImprove.MaxHitPoints() + benefit.Defense.MaxHitPoints
		squaddieToImprove.Defense.SquaddieDodge = squaddieToImprove.Dodge() + benefit.Defense.Dodge
		squaddieToImprove.Defense.SquaddieDeflect = squaddieToImprove.Deflect() + benefit.Defense.Deflect
		squaddieToImprove.Defense.SquaddieMaxBarrier = squaddieToImprove.MaxBarrier() + benefit.Defense.MaxBarrier
		squaddieToImprove.Defense.SquaddieArmor = squaddieToImprove.Armor() + benefit.Defense.Armor
	}

	if benefit.Offense != nil {
		squaddieToImprove.Offense.SquaddieAim = squaddieToImprove.Aim() + benefit.Offense.Aim
		squaddieToImprove.Offense.SquaddieStrength = squaddieToImprove.Strength() + benefit.Offense.Strength
		squaddieToImprove.Offense.SquaddieMind = squaddieToImprove.Mind() + benefit.Offense.Mind
	}
}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, repos *repositories.RepositoryCollection) error {
	if squaddieToImprove.HasAddedClass(benefit.Identification.ClassID) == false {
		newError := fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddieToImprove.Name(), benefit.Identification.ClassID)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}
	if squaddieToImprove.IsClassLevelAlreadyUsed(benefit.Identification.ID) {
		newError := fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddieToImprove.Name(), benefit.Identification.ClassID, benefit.Identification.ID)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	improveSquaddieStats(benefit, squaddieToImprove)
	refreshSquaddiePowers(benefit, squaddieToImprove)
	improveSquaddieMovement(benefit, squaddieToImprove)

	squaddieToImprove.SetBaseClassIfNoBaseClass(benefit.Identification.ClassID)
	squaddieToImprove.MarkLevelUpBenefitAsConsumed(benefit.Identification.ClassID, benefit.Identification.ID)
	return nil
}

func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	initialSquaddiePowerReferences := squaddieToImprove.PowerCollection.PowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddieToImprove.PowerCollection.GetCopyOfPowerReferences()
	}

	if benefit.PowerChanges != nil {
		for _, powerReferenceLost := range benefit.PowerChanges.Lost {
			squaddieToImprove.RemovePowerReferenceByPowerID(powerReferenceLost.PowerID)
		}

		for _, powerReferenceGained := range benefit.PowerChanges.Gained {
			squaddieToImprove.AddPowerReference(powerReferenceGained)
		}
	}
}

func improveSquaddieMovement(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Movement == nil {
		return
	}

	squaddieToImprove.Movement.SquaddieMovementDistance = squaddieToImprove.Movement.SquaddieMovementDistance + benefit.Movement.SquaddieMovementDistance

	if squaddie.MovementValueByType[squaddieToImprove.Movement.SquaddieMovementType] < squaddie.MovementValueByType[benefit.Movement.SquaddieMovementType] {
		squaddieToImprove.Movement.SquaddieMovementType = benefit.Movement.SquaddieMovementType
	}

	if benefit.Movement.SquaddieMovementCanHitAndRun {
		squaddieToImprove.Movement.SquaddieMovementCanHitAndRun = benefit.Movement.SquaddieMovementCanHitAndRun
	}
}
