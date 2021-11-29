package levelupbenefit

import (
	"errors"
	"fmt"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility"
)

// TODO privatize all fields

// LevelUpBenefit describes how a Squaddie improves upon levelling up.
type LevelUpBenefit struct {
	Identification *Identification    `json:"identification" yaml:"identification"`
	Defense        *Defense           `json:"defense" yaml:"defense"`
	Offense        *Offense           `json:"offense" yaml:"offense"`
	PowerChanges   *PowerChanges      `json:"powers" yaml:"powers"`
	Movement       *squaddie.Movement `json:"Movement" yaml:"Movement"`
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
		Identification: identification,
		Defense: defense,
		Offense: offense,
		Movement: movement,
		PowerChanges: changes,
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
	return l.Identification.LevelID()
}

// ClassID is a getter.
func (l LevelUpBenefit) ClassID() string {
	return l.Identification.ClassID()
}

// LevelUpBenefitType is a getter.
func (l LevelUpBenefit) LevelUpBenefitType() Size {
	return l.Identification.LevelUpBenefitSize()
}

// MaxHitPoints is a getter.
func (l LevelUpBenefit) MaxHitPoints() int {
	return l.Defense.MaxHitPoints()
}

// Dodge is a getter.
func (l LevelUpBenefit) Dodge() int {
	return l.Defense.Dodge()
}

// Deflect is a getter.
func (l LevelUpBenefit) Deflect() int {
	return l.Defense.Deflect()
}

// MaxBarrier is a getter.
func (l LevelUpBenefit) MaxBarrier() int {
	return l.Defense.MaxBarrier()
}

// Armor is a getter.
func (l LevelUpBenefit) Armor() int {
	return l.Defense.Armor()
}

// Aim is a getter.
func (l LevelUpBenefit) Aim() int {
	return l.Offense.Aim()
}

// Strength is a getter.
func (l LevelUpBenefit) Strength() int {
	return l.Offense.Strength()
}

// Mind is a getter.
func (l LevelUpBenefit) Mind() int {
	return l.Offense.Mind()
}

// MovementDistance is a getter.
func (l LevelUpBenefit) MovementDistance() int {
	return l.Movement.MovementDistance()
}

// MovementType is a getter.
func (l LevelUpBenefit) MovementType() squaddie.MovementType {
	return l.Movement.MovementType()
}

// CanHitAndRun is a getter.
func (l LevelUpBenefit) CanHitAndRun() bool {
	return l.Movement.CanHitAndRun()
}

// PowersGained is a getter.
func (l LevelUpBenefit) PowersGained() []*power.Reference {
	return l.PowerChanges.Gained()
}

// PowersLost is a getter.
func (l LevelUpBenefit) PowersLost() []*power.Reference {
	return l.PowerChanges.Lost()
}