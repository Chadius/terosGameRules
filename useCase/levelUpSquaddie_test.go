package usecase_test

import (
	"github.com/cserrant/terosBattleServer/entity"
	"github.com/cserrant/terosBattleServer/repository"
	"github.com/cserrant/terosBattleServer/usecase"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Squaddie uses LevelUpBenefits", func() {
	Context("Squaddie uses level up benefit", func() {
		var teros *entity.Squaddie
		var statBooster entity.LevelUpBenefit

		BeforeEach(func() {
			teros = entity.NewSquaddie("Teros")
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

			statBooster = entity.LevelUpBenefit{
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
			err := usecase.LevelUpSquaddie(&statBooster, teros, nil)
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
			err := usecase.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())
		})
		It("Raise an error if you add a level to a class that does not exist", func() {
			mushroomClassLevel := entity.LevelUpBenefit{
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
			err := usecase.LevelUpSquaddie(&mushroomClassLevel, teros, nil)
			Expect(err.Error()).To(Equal(`squaddie "Teros" cannot add levels to unknown class "Mushroom"`))
		})
		It("raises an error if you add a level that was already used", func() {
			err := usecase.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err).To(BeNil())
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1}))
			Expect(teros.IsClassLevelAlreadyUsed(statBooster.ID)).To(BeTrue())

			err = usecase.LevelUpSquaddie(&statBooster, teros, nil)
			Expect(err.Error()).To(Equal(`Teros already consumed LevelUpBenefit - class:"Mage" id:"deadbeef"`))
		})
		Context("can increase and change movement after using a level up benefit", func() {
			var (
				improveAllMovement *entity.LevelUpBenefit
				upgradeToLightMovement *entity.LevelUpBenefit
			)
			BeforeEach(func() {
				improveAllMovement = &entity.LevelUpBenefit{
					ID: "aaaaaaa0",
					ClassName: "Mage",
					Movement: &entity.SquaddieMovement{
						Distance: 1,
						Type: "fly",
						HitAndRun: true,
					},
				}

				upgradeToLightMovement = &entity.LevelUpBenefit{
					ID: "aaaaaaa1",
					ClassName: "Mage",
					Movement: &entity.SquaddieMovement{
						Type: "light",
					},
				}
			})
			It("can increase and change movement from one level up benefit", func() {
				startingMovement := teros.GetMovementDistancePerRound()

				err := usecase.LevelUpSquaddie(improveAllMovement, teros, nil)
				Expect(err).To(BeNil())

				Expect(teros.GetMovementDistancePerRound()).To(Equal(startingMovement + 1))
				Expect(teros.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFly)))
				Expect(teros.CanHitAndRun()).To(BeTrue())
			})
			It("will not downgrade movement type", func() {
				startingMovement := teros.GetMovementDistancePerRound()
				usecase.LevelUpSquaddie(improveAllMovement, teros, nil)

				err := usecase.LevelUpSquaddie(upgradeToLightMovement, teros, nil)
				Expect(err).To(BeNil())

				Expect(teros.GetMovementDistancePerRound()).To(Equal(startingMovement + 1))
				Expect(teros.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFly)))
				Expect(teros.CanHitAndRun()).To(BeTrue())
			})
		})
	})
	Context("Squaddie changes powers with level up benefits", func() {
		var (
			teros *entity.Squaddie
			powerRepo *repository.PowerRepository
			gainPower entity.LevelUpBenefit
			upgradePower entity.LevelUpBenefit
			spear *entity.Power
			spearLevel2 *entity.Power
		)
		BeforeEach(func() {
			teros = entity.NewSquaddie("Teros")
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

			powerRepo = repository.NewPowerRepository()

			spear = entity.NewPower("Spear")
			spear.PowerType = entity.PowerTypePhysical
			spear.ToHitBonus = 1
			spear.ID = "spearlvl1"
			teros.TemporaryPowerReferences = []*entity.PowerReference{{Name: "Spear", ID: "spearlvl1"}}

			spearLevel2 = entity.NewPower("Spear")
			spearLevel2.PowerType = entity.PowerTypePhysical
			spearLevel2.ToHitBonus = 1
			spearLevel2.ID = "spearlvl2"
			newPowers := []*entity.Power{spear, spearLevel2}
			powerRepo.AddSlicePowerSource(newPowers)

			gainPower = entity.LevelUpBenefit{
				ID: "aaab1234",
				LevelUpBenefitType: entity.LevelUpBenefitTypeBig,
				ClassName: "Mage",
				PowerIDGained: []*entity.PowerReference{{Name: "Spear", ID: spear.ID}},
			}

			upgradePower = entity.LevelUpBenefit{
				ID: "aaaa1235",
				LevelUpBenefitType: entity.LevelUpBenefitTypeBig,
				ClassName: "Mage",
				PowerIDLost: []*entity.PowerReference{{Name: "Spear", ID: spear.ID}},
				PowerIDGained: []*entity.PowerReference{{Name: "Spear", ID: spearLevel2.ID}},
			}
		})

		It("Squaddie can gain more powers based on level", func() {
			err := usecase.LevelUpSquaddie(&gainPower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spear.ID))

			Expect(teros.TemporaryPowerReferences).To(BeEmpty())
		})

		It("Squaddie can lose powers due to LevelUpBenefits", func() {
			usecase.LevelUpSquaddie(&gainPower, teros, powerRepo)
			teros.GetInnatePowerIDNames()

			err := usecase.LevelUpSquaddie(&upgradePower, teros, powerRepo)
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(attackIDNamePairs).To(HaveLen(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spearLevel2.ID))
		})
	})
})