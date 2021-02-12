package levelupbenefit_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
. "github.com/onsi/ginkgo"
. "github.com/onsi/gomega"
)

var _ = Describe("LevelUpBenefit error checking", func() {
	It("raises an error if there is no benefit type is set", func() {
		badLevel := levelupbenefit.LevelUpBenefit{
			ClassID:            "class0",
		}
		err := badLevel.CheckForErrors()
		Expect(err.Error()).To(Equal(`unknown level up benefit type: ""`))
	})
	It("raises an error if there is no class ID", func() {
		badLevel := levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            "",
		}
		err := badLevel.CheckForErrors()
		Expect(err.Error()).To(Equal(`no classID found for LevelUpBenefit`))
	})
	It("Filters a list", func() {
		listToTest := []*levelupbenefit.LevelUpBenefit{
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            "class0",
				ID:                 "level0",
				Aim:                1,
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            "class0",
				ID:                 "level1",
				MaxHitPoints:       1,
			},
			{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            "class0",
				ID:                 "level2",
				Aim:                1,
			},
		}

		noLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
			return false
		})
		Expect(noLevelsFound).To(HaveLen(0))

		allLevelsFound := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
			return true
		})
		Expect(allLevelsFound).To(HaveLen(3))

		onlySmallLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
			return benefit.LevelUpBenefitType == levelupbenefit.Small
		})
		Expect(onlySmallLevels).To(HaveLen(2))

		onlyBigLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
			return benefit.LevelUpBenefitType == levelupbenefit.Big
		})
		Expect(onlyBigLevels).To(HaveLen(1))

		increasesAimLevels := levelupbenefit.FilterLevelUpBenefits(listToTest, func(benefit *levelupbenefit.LevelUpBenefit) bool {
			return benefit.Aim > 0
		})
		Expect(increasesAimLevels).To(HaveLen(2))
	})
})
