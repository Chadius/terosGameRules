package levelUpBenefit_test

import (
"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
)

var _ = Describe("LevelUpBenefit error checking", func() {
	It("raises an error if there is no benefit type is set", func() {
		badLevel := levelUpBenefit.LevelUpBenefit{
			ClassID:            "class0",
		}
		err := badLevel.CheckForErrors()
		Expect(err.Error()).To(Equal(`unknown level up benefit type: ""`))
	})
	It("raises an error if there is no class ID", func() {
		badLevel := levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Small,
			ClassID:            "",
		}
		err := badLevel.CheckForErrors()
		Expect(err.Error()).To(Equal(`no classID found for LevelUpBenefit`))
	})
	It("Filters a list", func() {
		listToTest := []*levelUpBenefit.LevelUpBenefit{
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID: "class0",
				ID: "level0",
				Aim: 1,
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID: "class0",
				ID: "level1",
				MaxHitPoints: 1,
			},
			{
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassID: "class0",
				ID: "level2",
				Aim: 1,
			},
		}

		noLevelsFound := levelUpBenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelUpBenefit.LevelUpBenefit) bool {
			return false
		})
		Expect(noLevelsFound).To(HaveLen(0))

		allLevelsFound := levelUpBenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelUpBenefit.LevelUpBenefit) bool {
			return true
		})
		Expect(allLevelsFound).To(HaveLen(3))

		onlySmallLevels := levelUpBenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelUpBenefit.LevelUpBenefit) bool {
			return benefit.LevelUpBenefitType == levelUpBenefit.Small
		})
		Expect(onlySmallLevels).To(HaveLen(2))

		onlyBigLevels := levelUpBenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelUpBenefit.LevelUpBenefit) bool {
			return benefit.LevelUpBenefitType == levelUpBenefit.Big
		})
		Expect(onlyBigLevels).To(HaveLen(1))

		increasesAimLevels := levelUpBenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelUpBenefit.LevelUpBenefit) bool {
			return benefit.Aim > 0
		})
		Expect(increasesAimLevels).To(HaveLen(2))
	})
})
