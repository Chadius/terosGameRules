package levelupbenefit

import (
	"errors"
	"fmt"
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility"
)

// LevelUpBenefit describes how a Squaddie improves upon levelling up.
type LevelUpBenefit struct {
	identification *Identification
	defense        *Defense
	offense        *Offense
	powerChanges   *PowerChanges
	movement       *squaddie.Movement
}

// NewLevelUpBenefit returns a new LevelUpBenefit object.
func NewLevelUpBenefit(
	identification *Identification,
	defense *Defense,
	offense *Offense,
	movement *squaddie.Movement,
	changes *PowerChanges,
) *LevelUpBenefit {
	return &LevelUpBenefit{
		identification: identification,
		defense:        defense,
		offense:        offense,
		movement:       movement,
		powerChanges:   changes,
	}
}

// CheckForErrors ensures the LevelUpBenefit has valid fields
func (l *LevelUpBenefit) CheckForErrors() error {
	if l.LevelUpBenefitType() != Small && l.LevelUpBenefitType() != Big {
		newError := fmt.Errorf(`unknown level up benefit type`)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	if l.ClassID() == "" {
		newError := errors.New(`no classID found for LevelUpBenefit`)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}
	return nil
}

// FilterLevelUpBenefits filters a slice of LevelUpBenefits.
func FilterLevelUpBenefits(sliceToFilter []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) []*LevelUpBenefit {
	keptLevelUpBenefits := []*LevelUpBenefit{}
	for _, benefitToTest := range sliceToFilter {
		if condition(benefitToTest) {
			keptLevelUpBenefits = append(keptLevelUpBenefits, benefitToTest)
		}
	}
	return keptLevelUpBenefits
}

// AnyLevelUpBenefits returns true if at least LevelUpBenefit in the sliceToAnalyze satisfies the condition.
func AnyLevelUpBenefits(sliceToAnalyze []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) bool {
	for _, benefitToTest := range sliceToAnalyze {
		if condition(benefitToTest) {
			return true
		}
	}
	return false
}

// CountLevelUpBenefits returns the number of LevelUpBenefit that satisfy the given condition.
func CountLevelUpBenefits(sliceToAnalyze []*LevelUpBenefit, condition func(benefit *LevelUpBenefit) bool) int {
	count := 0
	for _, benefitToTest := range sliceToAnalyze {
		if condition(benefitToTest) {
			count = count + 1
		}
	}
	return count
}

// ID is a getter.
func (l LevelUpBenefit) ID() string {
	return l.identification.LevelID()
}

// ClassID is a getter.
func (l LevelUpBenefit) ClassID() string {
	return l.identification.ClassID()
}

// LevelUpBenefitType is a getter.
func (l LevelUpBenefit) LevelUpBenefitType() Size {
	return l.identification.LevelUpBenefitSize()
}

// MaxHitPoints is a getter.
func (l LevelUpBenefit) MaxHitPoints() int {
	if l.defense == nil {
		return 0
	}
	return l.defense.MaxHitPoints()
}

// Dodge is a getter.
func (l LevelUpBenefit) Dodge() int {
	if l.defense == nil {
		return 0
	}
	return l.defense.Dodge()
}

// Deflect is a getter.
func (l LevelUpBenefit) Deflect() int {
	if l.defense == nil {
		return 0
	}
	return l.defense.Deflect()
}

// MaxBarrier is a getter.
func (l LevelUpBenefit) MaxBarrier() int {
	if l.defense == nil {
		return 0
	}
	return l.defense.MaxBarrier()
}

// Armor is a getter.
func (l LevelUpBenefit) Armor() int {
	if l.defense == nil {
		return 0
	}
	return l.defense.Armor()
}

// Aim is a getter.
func (l LevelUpBenefit) Aim() int {
	if l.offense == nil {
		return 0
	}
	return l.offense.Aim()
}

// Strength is a getter.
func (l LevelUpBenefit) Strength() int {
	if l.offense == nil {
		return 0
	}
	return l.offense.Strength()
}

// Mind is a getter.
func (l LevelUpBenefit) Mind() int {
	if l.offense == nil {
		return 0
	}
	return l.offense.Mind()
}

// MovementDistance is a getter.
func (l LevelUpBenefit) MovementDistance() int {
	if l.movement == nil {
		return 0
	}
	return l.movement.MovementDistance()
}

// MovementLogic is a getter.
func (l LevelUpBenefit) MovementLogic() movement.Interface {
	return l.movement.MovementLogic()
}

// CanHitAndRun is a getter.
func (l LevelUpBenefit) CanHitAndRun() bool {
	if l.movement == nil {
		return false
	}
	return l.movement.CanHitAndRun()
}

// PowersGained is a getter.
func (l LevelUpBenefit) PowersGained() []*powerreference.Reference {
	return l.powerChanges.Gained()
}

// PowersLost is a getter.
func (l LevelUpBenefit) PowersLost() []*powerreference.Reference {
	return l.powerChanges.Lost()
}
