package levelUp_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieClass"
	"github.com/cserrant/terosBattleServer/usecase/levelUp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logic for Squaddies and Classes", func() {

	var (
		classRepo *squaddieClass.Repository
		levelRepo *levelUpBenefit.Repository

		mageClass *squaddieClass.Class
		mageLevel0 *levelUpBenefit.LevelUpBenefit
		mageLevel1 *levelUpBenefit.LevelUpBenefit

		dimensionWalkerClass *squaddieClass.Class
		dimensionWalkerLevel0 *levelUpBenefit.LevelUpBenefit
		dimensionWalkerLevel1 *levelUpBenefit.LevelUpBenefit

		ancientTomeClass *squaddieClass.Class
		ancientTomeClassLevel0 *levelUpBenefit.LevelUpBenefit

		atLeastTenLevelsBaseClass *squaddieClass.Class
		lotsOfLevels []*levelUpBenefit.LevelUpBenefit

		teros *squaddie.Squaddie
	)

	BeforeEach(func() {
		mageClass = &squaddieClass.Class{
			ID:                "class1",
			Name:              "Mage",
			BaseClassRequired: false,
		}

		dimensionWalkerClass = &squaddieClass.Class{
			ID:                "class2",
			Name:              "Dimension Walker",
			BaseClassRequired: true,
		}

		ancientTomeClass = &squaddieClass.Class{
			ID:                "class3",
			Name:              "Ancient Tome",
			BaseClassRequired: true,
		}

		atLeastTenLevelsBaseClass = &squaddieClass.Class{
			ID:                "class4",
			Name:              "Base with many levels",
			BaseClassRequired: false,
		}

		classRepo = squaddieClass.NewRepository()
		classRepo.AddListOfClasses([]*squaddieClass.Class{mageClass, dimensionWalkerClass, ancientTomeClass})

		mageLevel0 = &levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Small,
			ClassID:            mageClass.ID,
			ID:                 "mageLevel0",
			MaxHitPoints:       1,
			Aim:                1,
			Strength:           1,
			Mind:               1,
			Dodge:              1,
			Deflect:            1,
			MaxBarrier:         1,
			Armor:              1,
			PowerIDGained:      nil,
			PowerIDLost:        nil,
			Movement:           nil,
		}

		mageLevel1 = &levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Big,
			ClassID:            mageClass.ID,
			ID:                 "mageLevel1",
			MaxHitPoints:       0,
			Aim:                0,
			Strength:           0,
			Mind:               0,
			Dodge:              0,
			Deflect:            0,
			MaxBarrier:         0,
			Armor:              0,
			PowerIDGained:      nil,
			PowerIDLost:        nil,
			Movement:           &squaddie.Movement{
				Distance:  1,
				Type:      "",
				HitAndRun: false,
			},
		}

		dimensionWalkerLevel0 = &levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Big,
			ClassID:            dimensionWalkerClass.ID,
			ID:                 "dwLevel0",
			Movement:           &squaddie.Movement{
				Distance:  1,
				Type:      "light",
				HitAndRun: false,
			},
		}

		dimensionWalkerLevel1 = &levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Big,
			ClassID:            dimensionWalkerClass.ID,
			ID:                 "dwLevel1",
			MaxHitPoints:       1,
			Aim:                0,
			Strength:           0,
			Mind:               0,
		}

		ancientTomeClassLevel0 = &levelUpBenefit.LevelUpBenefit{
			LevelUpBenefitType: levelUpBenefit.Small,
			ClassID:            ancientTomeClass.ID,
			ID:                 "ancientTomeLevel0",
			Mind:               1,
			Dodge:              1,
			Deflect:            1,
		}

		levelRepo = levelUpBenefit.NewLevelUpBenefitRepository()
		levelRepo.AddLevels([]*levelUpBenefit.LevelUpBenefit{
			mageLevel0, mageLevel1, dimensionWalkerLevel0, dimensionWalkerLevel1, ancientTomeClassLevel0,
		})

		lotsOfLevels = []*levelUpBenefit.LevelUpBenefit{
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall0",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall1",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall2",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall3",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall4",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall5",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall6",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall7",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall8",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall9",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall10",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig0",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig1",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig2",
			},
			{
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig3",
			},
		}
		levelRepo.AddLevels(lotsOfLevels)

		teros = squaddie.NewSquaddie("Teros")
		teros.AddClass(mageClass)
		teros.AddClass(dimensionWalkerClass)
		teros.AddClass(ancientTomeClass)
		teros.AddClass(atLeastTenLevelsBaseClass)
	})

	Context("Qualifies for switching into a new class", func() {
		It("Will allow you to switch to a base class if squaddie has not set a base class", func() {
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching to the class the squaddie is already in", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching until the squaddie finishes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will allow switching to advanced classes if the squaddie finishes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			levelUp.LevelUpSquaddie(mageLevel1, teros, nil)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
		It("will not allow switching between advanced classes if the class is under level 10", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			levelUp.LevelUpSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching to an already completed class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			levelUp.LevelUpSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			levelUp.LevelUpSquaddie(dimensionWalkerLevel0, teros, nil)
			levelUp.LevelUpSquaddie(dimensionWalkerLevel1, teros, nil)
			teros.SetClass(ancientTomeClass.ID)
			Expect(teros.CurrentClass).To(Equal(ancientTomeClass.ID))
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("Squaddie qualifies for Advanced classes if it completes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			levelUp.LevelUpSquaddie(mageLevel1, teros, nil)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
		It("keeps squaddie on the advanced class until completion", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelUp.LevelUpSquaddie(mageLevel0, teros, nil)
			levelUp.LevelUpSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			levelUp.LevelUpSquaddie(dimensionWalkerLevel0, teros, nil)
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("lets you switch to a new class if you have at least 10 small levels in this class", func() {
			teros.SetBaseClassIfNoBaseClass(atLeastTenLevelsBaseClass.ID)
			teros.SetClass(atLeastTenLevelsBaseClass.ID)
			for index, _ := range [10]int{} {
				levelUp.LevelUpSquaddie(lotsOfLevels[index], teros, nil)
			}
			Expect(teros.GetLevelCountsByClass()[atLeastTenLevelsBaseClass.ID]).To(Equal(10))
			Expect(levelUp.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelUp.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
	})
})