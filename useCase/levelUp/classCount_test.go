package levelUp_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieClass"
	"github.com/cserrant/terosBattleServer/usecase/levelUp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddie choosing and gaining levels", func() {
	Context("Squaddie knows about big and small level ups", func() {
		var (
			teros *squaddie.Squaddie
			mageClass *squaddieClass.Class
			onlySmallLevelsClass *squaddieClass.Class
			lotsOfSmallLevels []*levelUpBenefit.LevelUpBenefit
			lotsOfBigLevels []*levelUpBenefit.LevelUpBenefit
			classRepo *squaddieClass.Repository
			levelRepo *levelUpBenefit.Repository
		)
		BeforeEach(func() {
			mageClass = &squaddieClass.Class{
				ID:                "class0",
				Name:              "Mage",
				BaseClassRequired: false,
			}

			onlySmallLevelsClass = &squaddieClass.Class{
				ID:                "class1",
				Name:              "SmallLevels",
				BaseClassRequired: false,
			}

			classRepo = squaddieClass.NewRepository()
			classRepo.AddListOfClasses([]*squaddieClass.Class{mageClass, onlySmallLevelsClass})

			lotsOfSmallLevels = []*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall1",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall2",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall3",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall4",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall5",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall6",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall7",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall8",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall9",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsSmall10",
				},
			}

			lotsOfBigLevels = []*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig1",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig2",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Big,
					ClassID:            mageClass.ID,
					ID:                 "lotsLevelsBig3",
				},
			}

			levelRepo = levelUpBenefit.NewLevelUpBenefitRepository()
			levelRepo.AddLevels(lotsOfSmallLevels)
			levelRepo.AddLevels(lotsOfBigLevels)
			levelRepo.AddLevels([]*levelUpBenefit.LevelUpBenefit{
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            onlySmallLevelsClass.ID,
					ID:                 "smallLevel0",
				},
				{
					LevelUpBenefitType: levelUpBenefit.Small,
					ClassID:            onlySmallLevelsClass.ID,
					ID:                 "smallLevel1",
				},
			})

			teros = squaddie.NewSquaddie("Teros")
			teros.AddClass(mageClass)
		})
		It("Only considers small levels when calculating its class level", func() {
			teros.AddClass(mageClass)
			for index, _ := range [5]int{} {
				levelUp.LevelUpSquaddie(lotsOfSmallLevels[index], teros, nil)
			}

			levelUp.LevelUpSquaddie(lotsOfBigLevels[0], teros, nil)
			levelUp.LevelUpSquaddie(lotsOfBigLevels[1], teros, nil)

			classLevels := levelUp.GetSquaddieClassLevels(teros, levelRepo)
			Expect(classLevels[mageClass.ID]).To(Equal(5))
		})
		Context("Choosing levels when leveling up a squaddie", func() {
			It("gets one small and one big level if the squaddie has an odd class level", func() {
				teros.AddClass(mageClass)
				teros.SetClass(mageClass.ID)
				err := levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelUp.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[mageClass.ID]).To(Equal(1))
				Expect(teros.ClassLevelsConsumed[mageClass.ID].LevelsConsumed).To(HaveLen(2))

				hasSmallLevel := false
				hasBigLevel := false
				for _, consumedLevelID := range teros.ClassLevelsConsumed[mageClass.ID].LevelsConsumed {
					for _, level := range lotsOfSmallLevels {
						if level.ID == consumedLevelID {
							hasSmallLevel = true
						}
					}
					for _, level := range lotsOfBigLevels {
						if level.ID == consumedLevelID {
							hasBigLevel = true
						}
					}
				}

				Expect(hasSmallLevel).To(BeTrue())
				Expect(hasBigLevel).To(BeTrue())
			})
			It("raises an error if the class cannot be found", func() {
				err := levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err.Error()).To(Equal(`class repository: No class found with ID: ""`))
			})
			It("does not add big levels if there are none available", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				err := levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, lotsOfBigLevels[0].ID, levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelUp.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(1))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(1))
			})
			It("chooses any small level at most once", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, "", levelRepo, classRepo, nil)

				err := levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, "", levelRepo, classRepo, nil)
				Expect(err).To(BeNil())
				classLevels := levelUp.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(2))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(2))
			})
			It("does not add small levels if there are none available", func() {
				teros.AddClass(onlySmallLevelsClass)
				teros.SetClass(onlySmallLevelsClass.ID)
				levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, "", levelRepo, classRepo, nil)
				levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, "", levelRepo, classRepo, nil)
				err := levelUp.LevelUpSquaddieBasedOnSquaddieLevel(teros, "", levelRepo, classRepo, nil)
				Expect(err).To(BeNil())

				classLevels := levelUp.GetSquaddieClassLevels(teros, levelRepo)
				Expect(classLevels[onlySmallLevelsClass.ID]).To(Equal(2))
				Expect(teros.ClassLevelsConsumed[onlySmallLevelsClass.ID].LevelsConsumed).To(HaveLen(2))
			})
		})
	})
})
