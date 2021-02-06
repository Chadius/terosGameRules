package levelUp_test

import (
	"github.com/cserrant/terosBattleServer/entity/levelUpBenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/levelUp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddie uses LevelUpBenefits", func() {
	Context("Squaddie uses level up benefit", func() {
		var teros *squaddie.Squaddie
		var statBooster levelUpBenefit.LevelUpBenefit

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
			teros.AddClass("Mage")
			teros.SetHPToMax()
			teros.SetBarrierToMax()

			statBooster = levelUpBenefit.LevelUpBenefit{
				ID: "deadbeef",
				ClassName: "Mage",
				MaxHitPoints: 0,
				Aim : 7,
				Strength : 6,
				Mind : 5,
				Dodge : 4,
				Deflect : 3,
				MaxBarrier : 2,
				Armor : 1,
			}
		})
		It("uses a LevelUpBenefit to increase Squaddie stats", func() {
			err := levelUp.LevelUpSquaddie(&statBooster, teros, nil)
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
			err := levelUp.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())
		})
		It("Raise an error if you add a level to a class that does not exist", func() {
			mushroomClassLevel := levelUpBenefit.LevelUpBenefit{
				ID: "deedbeeg",
				ClassName: "Mushroom",
				MaxHitPoints: 0,
				Aim : 7,
				Strength : 6,
				Mind : 5,
				Dodge : 4,
				Deflect : 3,
				MaxBarrier : 2,
				Armor : 1,
			}
			err := levelUp.LevelUpSquaddie(&mushroomClassLevel, teros, nil)
			Expect(err.Error()).To(Equal(`squaddie "Teros" cannot add levels to unknown class "Mushroom"`))
		})
		It("raises an error if you add a level that was already used", func() {
			err := levelUp.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())

			err = levelUp.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err.Error()).To(Equal(`Teros already consumed LevelUpBenefit - class:"Mage" id:"deadbeef"`))
		})
		Context("can increase and change movement after using a level up benefit", func() {
			var (
				improveAllMovement *levelUpBenefit.LevelUpBenefit
				upgradeToLightMovement *levelUpBenefit.LevelUpBenefit
			)
			BeforeEach(func() {
				improveAllMovement = &levelUpBenefit.LevelUpBenefit{
					ID: "aaaaaaa0",
					ClassName: "Mage",
					Movement: &squaddie.Movement{
						Distance: 1,
						Type: "fly",
						HitAndRun: true,
					},
				}

				upgradeToLightMovement = &levelUpBenefit.LevelUpBenefit{
					ID: "aaaaaaa1",
					ClassName: "Mage",
					Movement: &squaddie.Movement{
						Type: "light",
					},
				}
			})
			It("can increase and change movement from one level up benefit", func() {
				startingMovement := teros.GetMovementDistancePerRound()

				err := levelUp.LevelUpSquaddie(improveAllMovement, teros, nil)
				Expect(err).To(BeNil())

				Expect(teros.GetMovementDistancePerRound()).To(Equal(startingMovement + 1))
				Expect(teros.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Fly)))
				Expect(teros.CanHitAndRun()).To(BeTrue())
			})
			It("will not downgrade movement type", func() {
				startingMovement := teros.GetMovementDistancePerRound()
				levelUp.LevelUpSquaddie(improveAllMovement, teros, nil)

				err := levelUp.LevelUpSquaddie(upgradeToLightMovement, teros, nil)
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
			gainPower    levelUpBenefit.LevelUpBenefit
			upgradePower levelUpBenefit.LevelUpBenefit
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
			teros.AddClass("Mage")
			teros.SetHPToMax()
			teros.SetBarrierToMax()

			powerRepo = power.NewPowerRepository()

			spear = power.NewPower("Spear")
			spear.PowerType = power.Physical
			spear.ToHitBonus = 1
			spear.ID = "spearlvl1"
			teros.PowerReferences = []*power.Reference{{Name: "Spear", ID: "spearlvl1"}}

			spearLevel2 = power.NewPower("Spear")
			spearLevel2.PowerType = power.Physical
			spearLevel2.ToHitBonus = 1
			spearLevel2.ID = "spearlvl2"
			newPowers := []*power.Power{spear, spearLevel2}
			powerRepo.AddSlicePowerSource(newPowers)

			gainPower = levelUpBenefit.LevelUpBenefit{
				ID:                 "aaab1234",
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassName:          "Mage",
				PowerIDGained:      []*power.Reference{{Name: "Spear", ID: spear.ID}},
			}

			upgradePower = levelUpBenefit.LevelUpBenefit{
				ID:                 "aaaa1235",
				LevelUpBenefitType: levelUpBenefit.Big,
				ClassName:          "Mage",
				PowerIDLost:        []*power.Reference{{Name: "Spear", ID: spear.ID}},
				PowerIDGained:      []*power.Reference{{Name: "Spear", ID: spearLevel2.ID}},
			}
		})

		It("Squaddie can gain more powers based on level", func() {
			err := levelUp.LevelUpSquaddie(&gainPower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spear.ID))
		})

		It("Squaddie can lose powers due to LevelUpBenefits", func() {
			levelUp.LevelUpSquaddie(&gainPower, teros, powerRepo)
			teros.GetInnatePowerIDNames()

			err := levelUp.LevelUpSquaddie(&upgradePower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(attackIDNamePairs).To(HaveLen(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spearLevel2.ID))
		})
	})
})