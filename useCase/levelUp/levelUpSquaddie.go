package levelup

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerusage"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie) {
	squaddieToImprove.MaxHitPoints = squaddieToImprove.MaxHitPoints + benefit.MaxHitPoints
	squaddieToImprove.Aim = squaddieToImprove.Aim + benefit.Aim
	squaddieToImprove.Strength = squaddieToImprove.Strength + benefit.Strength
	squaddieToImprove.Mind = squaddieToImprove.Mind + benefit.Mind
	squaddieToImprove.Dodge = squaddieToImprove.Dodge + benefit.Dodge
	squaddieToImprove.Deflect = squaddieToImprove.Deflect + benefit.Deflect
	squaddieToImprove.MaxBarrier = squaddieToImprove.MaxBarrier + benefit.MaxBarrier
	squaddieToImprove.Armor = squaddieToImprove.Armor + benefit.Armor
}

// ImproveSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func ImproveSquaddie(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, powerRepo *power.Repository) error {
	if squaddieToImprove.HasAddedClass(benefit.ClassID) == false {
		return fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddieToImprove.Name, benefit.ClassID)
	}
	if squaddieToImprove.IsClassLevelAlreadyUsed(benefit.ID) {
		return fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddieToImprove.Name, benefit.ClassID, benefit.ID)
	}

	improveSquaddieStats(benefit, squaddieToImprove)
	err := refreshSquaddiePowers(benefit, squaddieToImprove, powerRepo)
	if err != nil {
		return err
	}

	improveSquaddieMovement(benefit, squaddieToImprove)

	squaddieToImprove.SetBaseClassIfNoBaseClass(benefit.ClassID)
	squaddieToImprove.MarkLevelUpBenefitAsConsumed(benefit.ClassID, benefit.ID)
	return nil
}

func refreshSquaddiePowers(benefit *levelupbenefit.LevelUpBenefit, squaddieToImprove *squaddie.Squaddie, powerRepo *power.Repository) error {
	initialSquaddiePowerReferences := squaddieToImprove.PowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddieToImprove.GetInnatePowerIDNames()
	}

	powerReferencesToKeep := squaddie.FilterPowerID(initialSquaddiePowerReferences, func(existingPower *power.Reference) bool {
		return squaddie.ContainsPowerID(benefit.PowerIDLost, existingPower.ID) == false
	})
	powerReferencesToLoad := append(powerReferencesToKeep, benefit.PowerIDGained...)

	_, err := powerusage.LoadAllOfSquaddieInnatePowers(squaddieToImprove, powerReferencesToLoad, powerRepo)
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
