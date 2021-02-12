package levelup_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/usecase/levelup"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logic for Squaddies and Classes", func() {

	var (
		classRepo *squaddieclass.Repository
		levelRepo *levelupbenefit.Repository

		mageClass *squaddieclass.Class
		mageLevel0 *levelupbenefit.LevelUpBenefit
		mageLevel1 *levelupbenefit.LevelUpBenefit

		dimensionWalkerClass *squaddieclass.Class
		dimensionWalkerLevel0 *levelupbenefit.LevelUpBenefit
		dimensionWalkerLevel1 *levelupbenefit.LevelUpBenefit

		ancientTomeClass *squaddieclass.Class
		ancientTomeClassLevel0 *levelupbenefit.LevelUpBenefit

		atLeastTenLevelsBaseClass *squaddieclass.Class
		lotsOfLevels []*levelupbenefit.LevelUpBenefit

		teros *squaddie.Squaddie
	)

	BeforeEach(func() {
		mageClass = &squaddieclass.Class{
			ID:                "class1",
			Name:              "Mage",
			BaseClassRequired: false,
		}

		dimensionWalkerClass = &squaddieclass.Class{
			ID:                "class2",
			Name:              "Dimension Walker",
			BaseClassRequired: true,
		}

		ancientTomeClass = &squaddieclass.Class{
			ID:                "class3",
			Name:              "Ancient Tome",
			BaseClassRequired: true,
		}

		atLeastTenLevelsBaseClass = &squaddieclass.Class{
			ID:                "class4",
			Name:              "Base with many levels",
			BaseClassRequired: false,
		}

		classRepo = squaddieclass.NewRepository()
		classRepo.AddListOfClasses([]*squaddieclass.Class{mageClass, dimensionWalkerClass, ancientTomeClass})

		mageLevel0 = &levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Small,
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

		mageLevel1 = &levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Big,
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

		dimensionWalkerLevel0 = &levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            dimensionWalkerClass.ID,
			ID:                 "dwLevel0",
			Movement:           &squaddie.Movement{
				Distance:  1,
				Type:      "light",
				HitAndRun: false,
			},
		}

		dimensionWalkerLevel1 = &levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            dimensionWalkerClass.ID,
			ID:                 "dwLevel1",
			MaxHitPoints:       1,
			Aim:                0,
			Strength:           0,
			Mind:               0,
		}

		ancientTomeClassLevel0 = &levelupbenefit.LevelUpBenefit{
			LevelUpBenefitType: levelupbenefit.Small,
			ClassID:            ancientTomeClass.ID,
			ID:                 "ancientTomeLevel0",
			Mind:               1,
			Dodge:              1,
			Deflect:            1,
		}

		levelRepo = levelupbenefit.NewLevelUpBenefitRepository()
		levelRepo.AddLevels([]*levelupbenefit.LevelUpBenefit{
			mageLevel0, mageLevel1, dimensionWalkerLevel0, dimensionWalkerLevel1, ancientTomeClassLevel0,
		})

		lotsOfLevels = []*levelupbenefit.LevelUpBenefit{
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall0",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall1",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall2",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall3",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall4",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall5",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall6",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall7",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall8",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall9",
			},
			{
				LevelUpBenefitType: levelupbenefit.Small,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsSmall10",
			},
			{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig0",
			},
			{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig1",
			},
			{
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            atLeastTenLevelsBaseClass.ID,
				ID:                 "lotsLevelsBig2",
			},
			{
				LevelUpBenefitType: levelupbenefit.Big,
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
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching to the class the squaddie is already in", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching until the squaddie finishes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will allow switching to advanced classes if the squaddie finishes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			levelup.ImproveSquaddie(mageLevel1, teros, nil)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
		It("will not allow switching between advanced classes if the class is under level 10", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			levelup.ImproveSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("will not allow switching to an already completed class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			levelup.ImproveSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			levelup.ImproveSquaddie(dimensionWalkerLevel0, teros, nil)
			levelup.ImproveSquaddie(dimensionWalkerLevel1, teros, nil)
			teros.SetClass(ancientTomeClass.ID)
			Expect(teros.CurrentClass).To(Equal(ancientTomeClass.ID))
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("Squaddie qualifies for Advanced classes if it completes the base class", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			levelup.ImproveSquaddie(mageLevel1, teros, nil)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
		It("keeps squaddie on the advanced class until completion", func() {
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			teros.SetClass(mageClass.ID)
			levelup.ImproveSquaddie(mageLevel0, teros, nil)
			levelup.ImproveSquaddie(mageLevel1, teros, nil)
			teros.SetClass(dimensionWalkerClass.ID)
			levelup.ImproveSquaddie(dimensionWalkerLevel0, teros, nil)
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeFalse())
			Expect(levelup.SquaddieCanSwitchToClass(teros, ancientTomeClass.ID, classRepo, levelRepo)).To(BeFalse())
		})
		It("lets you switch to a new class if you have at least 10 small levels in this class", func() {
			teros.SetBaseClassIfNoBaseClass(atLeastTenLevelsBaseClass.ID)
			teros.SetClass(atLeastTenLevelsBaseClass.ID)
			for index, _ := range [10]int{} {
				levelup.ImproveSquaddie(lotsOfLevels[index], teros, nil)
			}
			Expect(teros.GetLevelCountsByClass()[atLeastTenLevelsBaseClass.ID]).To(Equal(10))
			Expect(levelup.SquaddieCanSwitchToClass(teros, mageClass.ID, classRepo, levelRepo)).To(BeTrue())
			Expect(levelup.SquaddieCanSwitchToClass(teros, dimensionWalkerClass.ID, classRepo, levelRepo)).To(BeTrue())
		})
	})
})