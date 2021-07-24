package levelup

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	"github.com/cserrant/terosBattleServer/utility"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Defense != nil {
		squaddieToImprove.Defense.MaxHitPoints = squaddieToImprove.Defense.MaxHitPoints + benefit.Defense.MaxHitPoints
		squaddieToImprove.Defense.Dodge = squaddieToImprove.Defense.Dodge + benefit.Defense.Dodge
		squaddieToImprove.Defense.Deflect = squaddieToImprove.Defense.Deflect + benefit.Defense.Deflect
		squaddieToImprove.Defense.MaxBarrier = squaddieToImprove.Defense.MaxBarrier + benefit.Defense.MaxBarrier
		squaddieToImprove.Defense.Armor = squaddieToImprove.Defense.Armor + benefit.Defense.Armor
	}

	if benefit.Offense != nil {
		squaddieToImprove.Offense.Aim = squaddieToImprove.Offense.Aim + benefit.Offense.Aim
		squaddieToImprove.Offense.Strength = squaddieToImprove.Offense.Strength + benefit.Offense.Strength
		squaddieToImprove.Offense.Mind = squaddieToImprove.Offense.Mind + benefit.Offense.Mind
	}
}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, repos *repositories.RepositoryCollection) error {
	if squaddieToImprove.ClassProgress.HasAddedClass(benefit.Identification.ClassID) == false {
		newError := fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddieToImprove.Identification.Name, benefit.Identification.ClassID)
		utility.Log(newError.Error(),0, utility.Error)
		return newError
	}
	if squaddieToImprove.ClassProgress.IsClassLevelAlreadyUsed(benefit.Identification.ID) {
		newError := fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddieToImprove.Identification.Name, benefit.Identification.ClassID, benefit.Identification.ID)
		utility.Log(newError.Error(),0, utility.Error)
		return newError
	}

	improveSquaddieStats(benefit, squaddieToImprove)
	err := refreshSquaddiePowers(benefit, squaddieToImprove, repos)
	if err != nil {
		return err
	}

	improveSquaddieMovement(benefit, squaddieToImprove)

	squaddieToImprove.ClassProgress.SetBaseClassIfNoBaseClass(benefit.Identification.ClassID)
	squaddieToImprove.ClassProgress.MarkLevelUpBenefitAsConsumed(benefit.Identification.ClassID, benefit.Identification.ID)
	return nil
}

func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, repos *repositories.RepositoryCollection) error {
	initialSquaddiePowerReferences := squaddieToImprove.PowerCollection.PowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddieToImprove.PowerCollection.GetInnatePowerIDNames()
	}

	var powerReferencesToLoad []*power.Reference
	if benefit.PowerChanges != nil {
		powerReferencesToKeep := squaddie.FilterPowerID(initialSquaddiePowerReferences, func(existingPower *power.Reference) bool {
			return squaddie.ContainsPowerID(benefit.PowerChanges.Lost, existingPower.ID) == false
		})
		powerReferencesToLoad = append(powerReferencesToKeep, benefit.PowerChanges.Gained...)
	}

	_, err := powerequip.LoadAllOfSquaddieInnatePowers(squaddieToImprove, powerReferencesToLoad, repos)
	return err
}

func improveSquaddieMovement(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	if benefit.Movement == nil {
		return
	}

	squaddieToImprove.Movement.Distance = squaddieToImprove.Movement.Distance + benefit.Movement.Distance

	if squaddie.MovementValueByType[squaddieToImprove.Movement.Type] < squaddie.MovementValueByType[benefit.Movement.Type] {
		squaddieToImprove.Movement.Type = benefit.Movement.Type
	}

	if benefit.Movement.HitAndRun {
		squaddieToImprove.Movement.HitAndRun = benefit.Movement.HitAndRun
	}
}
