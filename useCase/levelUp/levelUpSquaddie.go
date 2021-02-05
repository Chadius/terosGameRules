package levelUp

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	squaddie2 "github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powerAttack"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *levelUpBenefit.LevelUpBenefit, squaddie *squaddie2.Squaddie) {
	squaddie.MaxHitPoints = squaddie.MaxHitPoints + benefit.MaxHitPoints
	squaddie.Aim = squaddie.Aim + benefit.Aim
	squaddie.Strength = squaddie.Strength + benefit.Strength
	squaddie.Mind = squaddie.Mind + benefit.Mind
	squaddie.Dodge = squaddie.Dodge + benefit.Dodge
	squaddie.Deflect = squaddie.Deflect + benefit.Deflect
	squaddie.MaxBarrier = squaddie.MaxBarrier + benefit.MaxBarrier
	squaddie.Armor = squaddie.Armor + benefit.Armor
}

// LevelUpSquaddie uses the LevelUpBenefit to improve the squaddie.
//   Raises an error if the Squaddie does not have that class.
//   Raises an error if the Squaddie marked the LevelUpBenefit as consumed.
func LevelUpSquaddie(benefit *levelUpBenefit.LevelUpBenefit, squaddie *squaddie2.Squaddie, powerRepo *power.PowerRepository) error {
	if squaddie.HasAddedClass(benefit.ClassName) == false {
		return fmt.Errorf(`squaddie "%s" cannot add levels to unknown class "%s"`, squaddie.Name, benefit.ClassName)
	}
	if squaddie.IsClassLevelAlreadyUsed(benefit.ID) {
		return fmt.Errorf(`%s already consumed LevelUpBenefit - class:"%s" id:"%s"`, squaddie.Name, benefit.ClassName, benefit.ID)
	}

	improveSquaddieStats(benefit, squaddie)
	err := refreshSquaddiePowers(benefit, squaddie, powerRepo)
	if err != nil {
		return err
	}

	improveSquaddieMovement(benefit, squaddie)

	squaddie.MarkLevelUpBenefitAsConsumed(benefit.ClassName, benefit.ID)
	return nil
}

func refreshSquaddiePowers(benefit *levelUpBenefit.LevelUpBenefit, squaddie *squaddie2.Squaddie, powerRepo *power.PowerRepository) error {
	initialSquaddiePowerReferences := squaddie.TemporaryPowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddie.GetInnatePowerIDNames()
	}

	powerReferencesToKeep := []*power.PowerReference{}
	for _, existingPower := range initialSquaddiePowerReferences {
		powerFound := false
		for _, powerToRemove := range benefit.PowerIDLost {
			if existingPower.ID == powerToRemove.ID {
				powerFound = true
			}
		}
		if powerFound == false {
			powerReferencesToKeep = append(powerReferencesToKeep, existingPower)
		}
	}

	powerReferencesToLoad := append(powerReferencesToKeep, benefit.PowerIDGained...)

	_, err := powerAttack.LoadAllOfSquaddieInnatePowers(squaddie, powerReferencesToLoad, powerRepo)
	return err
}

func improveSquaddieMovement(benefit *levelUpBenefit.LevelUpBenefit, squaddie *squaddie2.Squaddie) {
	if benefit.Movement == nil {
		return
	}

	squaddie.Movement.Distance = squaddie.Movement.Distance + benefit.Movement.Distance

	if squaddie2.MovementValueByType[squaddie.Movement.Type] < squaddie2.MovementValueByType[benefit.Movement.Type] {
		squaddie.Movement.Type = benefit.Movement.Type
	}

	if benefit.Movement.HitAndRun {
		squaddie.Movement.HitAndRun = benefit.Movement.HitAndRun
	}
}