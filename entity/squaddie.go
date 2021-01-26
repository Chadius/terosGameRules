package entity

import (
	"fmt"
)

// Squaddie is the base unit you can deploy and control on a field.
type Squaddie struct {
	Name                string               `json:"name" yaml:"name"`
	Affiliation         string               `json:"affiliation" yaml:"affiliation"`
	CurrentClass        string               `json:"current_class" yaml:"current_class"`
	ClassLevels         map[string][]string  `json:"class_levels" yaml:"class_levels"`
	CurrentHitPoints         int               `json:"current_hit_points" yaml:"current_hit_points"`
	MaxHitPoints             int               `json:"max_hit_points" yaml:"max_hit_points"`
	Aim                      int               `json:"aim" yaml:"aim"`
	Strength                 int               `json:"strength" yaml:"strength"`
	Mind                     int               `json:"mind" yaml:"mind"`
	Dodge                    int               `json:"dodge" yaml:"dodge"`
	Deflect                  int               `json:"deflect" yaml:"deflect"`
	CurrentBarrier           int               `json:"current_barrier" yaml:"current_barrier"`
	MaxBarrier               int               `json:"max_barrier" yaml:"max_barrier"`
	Armor                    int               `json:"armor" yaml:"armor"`
	TemporaryPowerReferences []*PowerReference `json:"powers" yaml:"powers"`
	innatePowers             []*Power
}

// NewSquaddie generates a squaddie with maxed out health.
func NewSquaddie(name string) *Squaddie {
	newSquaddie := Squaddie{
		Name:             name,
		Affiliation:      "Player",
		CurrentHitPoints:    0,
		MaxHitPoints:        5,
		Aim:                 0,
		Strength:            0,
		Mind:                0,
		Dodge:               0,
		Deflect:             0,
		CurrentBarrier:      0,
		MaxBarrier:          0,
		Armor:               0,
		ClassLevels: map[string][]string{},
	}
	newSquaddie.SetHPToMax()
	return &newSquaddie
}

// CheckSquaddieForErrors makes sure the created squaddie doesn't have an error.
func CheckSquaddieForErrors(newSquaddie *Squaddie) (newError error) {
	if newSquaddie.Affiliation != "Player" {
		return fmt.Errorf("Squaddie has unknown affiliation: '%s'", newSquaddie.Affiliation)
	}

	return nil
}

// SetHPToMax restores the Squaddie's HitPoints.
func (squaddie *Squaddie) SetHPToMax() {
	squaddie.CurrentHitPoints = squaddie.MaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (squaddie *Squaddie) SetBarrierToMax() {
	squaddie.CurrentBarrier = squaddie.MaxBarrier
}

// GetDefensiveStatsAgainstPhysical calculates how this squaddie can defend against physical attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstPhysical() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Dodge, squaddie.CurrentBarrier, squaddie.Armor
}

// GetDefensiveStatsAgainstSpell calculates how this squaddie can defend against spell attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstSpell() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Deflect, squaddie.CurrentBarrier, 0
}

// GetOffensiveStatsWithPhysical calculates the squaddie's bonuses with physical attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithPhysical() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Strength
}

// GetOffensiveStatsWithSpell calculates the squaddie's bonuses with Spell attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithSpell() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Mind
}

// AddInnatePower gives the Squaddie access to the power.
//  Raises an error if the squaddie already has the power.
func (squaddie *Squaddie) AddInnatePower(newPower *Power) error {
	for _, innatePower := range squaddie.innatePowers {
		if newPower.ID == innatePower.ID {
			return fmt.Errorf(`squaddie "%s" already has innate power with ID "%s"`, squaddie.Name, innatePower.ID)
		}
	}

	squaddie.innatePowers = append(squaddie.innatePowers, newPower)
	return nil
}

// GetInnatePowerIDNames returns a list of all the powers the squaddie has access to.
func (squaddie *Squaddie) GetInnatePowerIDNames() []*PowerReference {
	powerIDNames := []*PowerReference{}
	for _, power := range squaddie.innatePowers {
		powerIDNames = append(powerIDNames, &PowerReference{Name: power.Name, ID: power.ID})
	}
	return powerIDNames
}

// AddClass gives the Squaddie a new class it can gain levels in.
func (squaddie *Squaddie) AddClass(className string) {
	squaddie.ClassLevels[className] = []string{}
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (squaddie *Squaddie) GetLevelCountsByClass() (map[string]int) {
	count := map[string]int{}
	for className, benefitIDs := range squaddie.ClassLevels {
		count[className] = len(benefitIDs)
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the Squaddie remember it used this benefit to level up already.
func (squaddie *Squaddie) MarkLevelUpBenefitAsConsumed(benefitClassName, benefitID string)  {
	squaddie.ClassLevels[benefitClassName] = append(squaddie.ClassLevels[benefitClassName], benefitID)
}

// SetClass changes the Squaddie's CurrentClass to the given className.
//   Raises an error if className has not been added to the squaddie yet.
func (squaddie *Squaddie) SetClass(className string) error {
	if _, exists := squaddie.ClassLevels[className]; !exists {
		return fmt.Errorf(`cannot switch "%s" to unknown class "%s"`, squaddie.Name, className)
	}

	squaddie.CurrentClass = className
	return nil
}

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given ID has already been used.
func (squaddie *Squaddie) IsClassLevelAlreadyUsed(benefitID string) bool {
	for _, levels := range squaddie.ClassLevels {
		for _, levelID := range levels {
			if levelID == benefitID {
				return true
			}
		}
	}
	return false
}

// HasAddedClass returns true if the Squaddie has already added a class with the name classNameToFind
func (squaddie *Squaddie) HasAddedClass(classNameToFind string) bool {
	for className, _ := range squaddie.ClassLevels {
		if className == classNameToFind {
			return true
		}
	}
	return false
}

// ClearInnatePowers removes all of the squaddie's powers.
func (squaddie *Squaddie) ClearInnatePowers() {
	squaddie.innatePowers = []*Power{}
}

// RemovePowerByID removes the power with the given ID from the squaddie's
func (squaddie *Squaddie) RemovePowerByID(powerToRemoveID string) {
	powerFound := false
	powerIndexToDelete := 0
	for index, power := range squaddie.innatePowers {
		if power.ID == powerToRemoveID {
			powerIndexToDelete = index
			powerFound = true
		}
	}
	if powerFound == false {
		return
	}

	squaddie.innatePowers = append(squaddie.innatePowers[:powerIndexToDelete], squaddie.innatePowers[powerIndexToDelete+1:]...)
}

// ClearTemporaryPowerReferences empties the temporary references to powers.
func (squaddie *Squaddie) ClearTemporaryPowerReferences() {
	squaddie.TemporaryPowerReferences = []*PowerReference{}
}