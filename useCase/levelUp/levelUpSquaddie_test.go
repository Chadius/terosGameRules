package levelup_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/usecase/levelup"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddie uses LevelUpBenefits", func() {
	var mageClass *squaddieclass.Class

	BeforeEach(func() {
		mageClass = &squaddieclass.Class{
			ID:   "ffffffff",
			Name: "Mage",
		}
	})

	Context("Squaddie uses level up benefit", func() {
		var teros *squaddie.Squaddie
		var statBooster levelupbenefit.LevelUpBenefit

		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			teros.MaxHitPoints = 5
			teros.Aim = 0
			teros.Strength = 1
			teros.Mind = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 6
			teros.Armor = 7
			teros.AddClass(mageClass)
			teros.SetHPToMax()
			teros.SetBarrierToMax()

			statBooster = levelupbenefit.LevelUpBenefit{
				ID:           "deadbeef",
				ClassID:      mageClass.ID,
				MaxHitPoints: 0,
				Aim :         7,
				Strength :    6,
				Mind :        5,
				Dodge :       4,
				Deflect :     3,
				MaxBarrier :  2,
				Armor :       1,
			}
		})
		It("uses a LevelUpBenefit to increase Squaddie stats", func() {
			err := levelup.ImproveSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.MaxHitPoints).To(Equal(5))

			Expect(teros.MaxHitPoints).To(Equal(5))
			Expect(teros.Aim).To(Equal(7))
			Expect(teros.Strength).To(Equal(7))
			Expect(teros.Mind).To(Equal(7))
			Expect(teros.Dodge).To(Equal(7))
			Expect(teros.Deflect).To(Equal(7))
			Expect(teros.MaxBarrier).To(Equal(8))
			Expect(teros.Armor).To(Equal(8))
		})
		It("tells the Squaddie which levels were used", func() {
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeFalse())
			err := levelup.ImproveSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{mageClass.ID: 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())
		})
		It("Raise an error if you add a level to a class that does not exist", func() {
			mushroomClassLevel := levelupbenefit.LevelUpBenefit{
				ID:           "deedbeeg",
				ClassID:      "bad ID",
				MaxHitPoints: 0,
				Aim :         7,
				Strength :    6,
				Mind :        5,
				Dodge :       4,
				Deflect :     3,
				MaxBarrier :  2,
				Armor :       1,
			}
			err := levelup.ImproveSquaddie(&mushroomClassLevel, teros, nil)
			Expect(err.Error()).To(Equal(`squaddie "Teros" cannot add levels to unknown class "bad ID"`))
		})
		It("raises an error if you add a level that was already used", func() {
			err := levelup.ImproveSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"ffffffff": 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())

			err = levelup.ImproveSquaddie(&statBooster, teros, nil)
			Expect(err.Error()).To(Equal(`Teros already consumed LevelUpBenefit - class:"ffffffff" id:"deadbeef"`))
		})
		It("sets the squaddie's base class if it isn't already set", func() {
			Expect(teros.BaseClassID).To(Equal(""))
			levelup.ImproveSquaddie(&statBooster, teros, nil)
			Expect(teros.BaseClassID).To(Equal(mageClass.ID))
		})
		Context("can increase and change movement after using a level up benefit", func() {
			var (
				improveAllMovement *levelupbenefit.LevelUpBenefit
				upgradeToLightMovement *levelupbenefit.LevelUpBenefit
			)
			BeforeEach(func() {
				improveAllMovement = &levelupbenefit.LevelUpBenefit{
					ID:      "aaaaaaa0",
					ClassID: mageClass.ID,
					Movement: &squaddie.Movement{
						Distance: 1,
						Type: "fly",
						HitAndRun: true,
					},
				}

				upgradeToLightMovement = &levelupbenefit.LevelUpBenefit{
					ID:      "aaaaaaa1",
					ClassID: mageClass.ID,
					Movement: &squaddie.Movement{
						Type: "light",
					},
				}
			})
			It("can increase and change movement from one level up benefit", func() {
				startingMovement := teros.GetMovementDistancePerRound()

				err := levelup.ImproveSquaddie(improveAllMovement, teros, nil)
				Expect(err).To(BeNil())

				Expect(teros.GetMovementDistancePerRound()).To(Equal(startingMovement + 1))
				Expect(teros.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Fly)))
				Expect(teros.CanHitAndRun()).To(BeTrue())
			})
			It("will not downgrade movement type", func() {
				startingMovement := teros.GetMovementDistancePerRound()
				levelup.ImproveSquaddie(improveAllMovement, teros, nil)

				err := levelup.ImproveSquaddie(upgradeToLightMovement, teros, nil)
				Expect(err).To(BeNil())

				Expect(teros.GetMovementDistancePerRound()).To(Equal(startingMovement + 1))
				Expect(teros.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Fly)))
				Expect(teros.CanHitAndRun()).To(BeTrue())
			})
		})
	})
	Context("Squaddie changes powers with level up benefits", func() {
		var (
			teros        *squaddie.Squaddie
			powerRepo    *power.Repository
			gainPower    levelupbenefit.LevelUpBenefit
			upgradePower levelupbenefit.LevelUpBenefit
			spear        *power.Power
			spearLevel2  *power.Power
		)
		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			teros.MaxHitPoints = 5
			teros.Aim = 0
			teros.Strength = 1
			teros.Mind = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 6
			teros.Armor = 7
			teros.AddClass(&squaddieclass.Class{
				ID:   mageClass.ID,
				Name: "Mage",
			})
			teros.SetHPToMax()
			teros.SetBarrierToMax()

			powerRepo = power.NewPowerRepository()

			spear = power.NewPower("Spear")
			spear.PowerType = power.Physical
			spear.AttackEffect.ToHitBonus = 1
			spear.ID = "spearlvl1"
			teros.PowerReferences = []*power.Reference{{Name: "Spear", ID: "spearlvl1"}}

			spearLevel2 = power.NewPower("Spear")
			spearLevel2.PowerType = power.Physical
			spearLevel2.AttackEffect.ToHitBonus = 1
			spearLevel2.ID = "spearlvl2"
			newPowers := []*power.Power{spear, spearLevel2}
			powerRepo.AddSlicePowerSource(newPowers)

			gainPower = levelupbenefit.LevelUpBenefit{
				ID:                 "aaab1234",
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            mageClass.ID,
				PowerIDGained:      []*power.Reference{{Name: "Spear", ID: spear.ID}},
			}

			upgradePower = levelupbenefit.LevelUpBenefit{
				ID:                 "aaaa1235",
				LevelUpBenefitType: levelupbenefit.Big,
				ClassID:            mageClass.ID,
				PowerIDLost:        []*power.Reference{{Name: "Spear", ID: spear.ID}},
				PowerIDGained:      []*power.Reference{{Name: "Spear", ID: spearLevel2.ID}},
			}
		})

		It("Squaddie can gain more powers based on level", func() {
			err := levelup.ImproveSquaddie(&gainPower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spear.ID))
		})

		It("Squaddie can lose powers due to LevelUpBenefits", func() {
			levelup.ImproveSquaddie(&gainPower, teros, powerRepo)
			teros.GetInnatePowerIDNames()

			err := levelup.ImproveSquaddie(&upgradePower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(attackIDNamePairs).To(HaveLen(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spearLevel2.ID))
		})
	})
})