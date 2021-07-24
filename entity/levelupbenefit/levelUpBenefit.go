package levelupbenefit

import (
	"errors"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/utility"
)

// LevelUpBenefit describes how a Squaddie improves upon levelling up.
type LevelUpBenefit struct {
	Identification				*Identification              `json:"identification" yaml:"identification"`
	Defense						*Defense						`json:"defense" yaml:"defense"`
	Offense						*Offense						`json:"offense" yaml:"offense"`
	PowerChanges				*PowerChanges				`json:"powers" yaml:"powers"`
	Movement           *squaddie.Movement `json:"Movement" yaml:"Movement"`
}

// CheckForErrors ensures the LevelUpBenefit has valid fields
func (benefit *LevelUpBenefit) CheckForErrors() error {
	if benefit.Identification.LevelUpBenefitType != Small && benefit.Identification.LevelUpBenefitType != Big {
		newError := fmt.Errorf(`unknown level up benefit type: "%s"`, benefit.Identification.LevelUpBenefitType)
		utility.Log(newError.Error(),0, utility.Error)
		return newError
	}

	if benefit.Identification.ClassID == "" {
		newError := errors.New(`no classID found for LevelUpBenefit`)
		utility.Log(newError.Error(),0, utility.Error)
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
