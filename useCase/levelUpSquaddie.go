package usecase

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity"
	"github.com/cserrant/terosBattleServer/repository"
)

// improveSquaddieStats improves the Squaddie by using the LevelUpBenefit.
func improveSquaddieStats(benefit *entity.LevelUpBenefit, squaddie *entity.Squaddie) {
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
func LevelUpSquaddie(benefit *entity.LevelUpBenefit, squaddie *entity.Squaddie, powerRepo *repository.PowerRepository) error {
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

	squaddie.MarkLevelUpBenefitAsConsumed(benefit.ClassName, benefit.ID)
	return nil
}

func refreshSquaddiePowers(benefit *entity.LevelUpBenefit, squaddie *entity.Squaddie, powerRepo *repository.PowerRepository) error {
	initialSquaddiePowerReferences := squaddie.TemporaryPowerReferences
	if initialSquaddiePowerReferences == nil || len(initialSquaddiePowerReferences) == 0 {
		initialSquaddiePowerReferences = squaddie.GetInnatePowerIDNames()
	}

	powerReferencesToKeep := []*entity.PowerReference{}
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

	_, err := LoadAllOfSquaddieInnatePowers(squaddie, powerReferencesToLoad, powerRepo)
	return err
}