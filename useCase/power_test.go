package usecase_test

import (
	"github.com/cserrant/terosBattleServer/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cserrant/terosBattleServer/entity"
	"github.com/cserrant/terosBattleServer/usecase"
)

var _ = Describe("Power uses with other Entities", func() {

	Context("Calculate combination of Attacking Power and Squaddie", func() {
		var (
			teros *entity.Squaddie
			spear *entity.Power
			blot  *entity.Power
		)

		BeforeEach(func() {
			teros = entity.NewSquaddie("Teros")
			teros.Name = "Teros"

			spear = entity.NewAttackingPower("Spear")
			spear.PowerType = entity.PowerTypePhysical

			blot = entity.NewAttackingPower("Blot")
			blot.PowerType = entity.PowerTypeSpell
		})

		It("Calculates the To Hit Bonus", func() {
			teros.Aim = 2
			blot.AttackingEffect.ToHitBonus = 1

			totalToHitBonus := usecase.GetPowerToHitBonusWhenUsedBySquaddie(blot, teros)
			Expect(totalToHitBonus).To(Equal(3))
		})

		Context("Calculate damage bonus", func() {
			BeforeEach(func() {
				teros.Strength = 2
				teros.Mind = 3

				spear.PowerType = entity.PowerTypePhysical
				spear.AttackingEffect.DamageBonus = 2

				blot.PowerType = entity.PowerTypeSpell
				blot.AttackingEffect.DamageBonus = 6
			})

			It("Calculates the Damage bonus of physical attacks", func() {
				totalDamageBonus := usecase.GetPowerDamageBonusWhenUsedBySquaddie(spear, teros)
				Expect(totalDamageBonus).To(Equal(4))
			})

			It("Calculates the Damage bonus of spell attacks", func() {
				totalDamageBonus := usecase.GetPowerDamageBonusWhenUsedBySquaddie(blot, teros)
				Expect(totalDamageBonus).To(Equal(9))
			})

			It("Calculates the Critical Damage bonus of physical attacks", func() {
				totalDamageBonus := usecase.GetPowerCriticalDamageBonusWhenUsedBySquaddie(spear, teros)
				Expect(totalDamageBonus).To(Equal(8))
			})
		})

		Context("Calculate to hit penalties values against powers", func() {
			BeforeEach(func() {
				teros.Dodge = 2
				teros.Deflect = 9001

				spear.PowerType = entity.PowerTypePhysical

				blot.PowerType = entity.PowerTypeSpell
			})

			It("Calculates the to hit reduction against physical attacks", func() {
				toHitPenalty := usecase.GetPowerToHitPenaltyAgainstSquaddie(spear, teros)
				Expect(toHitPenalty).To(Equal(2))
			})

			It("Calculates the to hit reduction against spell attacks", func() {
				toHitPenalty := usecase.GetPowerToHitPenaltyAgainstSquaddie(blot, teros)
				Expect(toHitPenalty).To(Equal(9001))
			})
		})

		Context("Calculate damage if the attacker hits the target with the power", func() {
			var (
				bandit *entity.Squaddie
			)

			BeforeEach(func() {
				bandit = entity.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Strength = 1
				spear.DamageBonus = 3
			})

			It("Does full damage against targets without armor or barrier", func() {
				totalHealthDamage, _, _ := usecase.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(4))
			})

			It("Armor reduces damage against physical attacks", func() {
				bandit.Armor = 3
				totalHealthDamage, _, _ := usecase.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(1))
			})

			It("Barrier absorbs damage against physical attacks and is depleted first", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1
				totalHealthDamage, initialBarrierDamage, _ := usecase.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(3))
				Expect(initialBarrierDamage).To(Equal(1))
			})

			It("Will deal no damage if barrier is strong enough", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 4
				totalHealthDamage, initialBarrierDamage, _ := usecase.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(4))
			})

			It("May deal no damage if armor is strong enough", func() {
				bandit.Armor = 4
				totalHealthDamage, initialBarrierDamage, _ := usecase.GetHowTargetDistributesDamage(spear, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(0))
			})
		})

		Context("Calculate damage if the attacker hits the target with the power", func() {
			var (
				bandit *entity.Squaddie
			)

			BeforeEach(func() {
				bandit = entity.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Mind = 2
				blot.DamageBonus = 4
			})

			It("Does full damage against targets without armor or barrier", func() {
				totalHealthDamage, _, _ := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(6))
			})

			It("Ignores Armor when using spell attacks", func() {
				bandit.Armor = 9001
				totalHealthDamage, _, _ := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(6))
			})

			It("Barrier absorbs damage against spell attacks and is depleted first", func() {
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1
				totalHealthDamage, initialBarrierDamage, _ := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(5))
				Expect(initialBarrierDamage).To(Equal(1))
			})

			It("Will deal no damage if barrier is strong enough", func() {
				bandit.MaxBarrier = 9001
				bandit.CurrentBarrier = 9001
				totalHealthDamage, initialBarrierDamage, _ := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
			})

			It("Can deal extra Barrier damage if the barrier absorbs the attack", func() {
				bandit.MaxBarrier = 8
				bandit.CurrentBarrier = 8
				blot.ExtraBarrierDamage = 2

				totalHealthDamage, initialBarrierDamage, extraBarrierDamage := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
				Expect(extraBarrierDamage).To(Equal(2))
			})

			It("Knows extra Barrier damage is reduced if the barrier is depleted", func() {
				bandit.MaxBarrier = 8
				bandit.CurrentBarrier = 7
				blot.ExtraBarrierDamage = 2

				totalHealthDamage, initialBarrierDamage, extraBarrierDamage := usecase.GetHowTargetDistributesDamage(blot, teros, bandit)
				Expect(totalHealthDamage).To(Equal(0))
				Expect(initialBarrierDamage).To(Equal(6))
				Expect(extraBarrierDamage).To(Equal(1))
			})
		})

		Context("Calculate expected damage", func() {
			var (
				bandit *entity.Squaddie
			)

			BeforeEach(func() {
				bandit = entity.NewSquaddie("Bandit")
				bandit.Name = "Bandit"
			})

			It("Give summary of the physical attack", func() {
				bandit.Armor = 1
				bandit.Dodge = 1
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 1

				teros.Strength = 1
				spear.DamageBonus = 3

				teros.Mind = 2
				blot.DamageBonus = 4

				attackingPowerSummary := usecase.GetExpectedDamage(spear, teros, bandit)
				Expect(attackingPowerSummary.ChanceToHit).To(Equal(15))
				Expect(attackingPowerSummary.DamageTaken).To(Equal(2))
				Expect(attackingPowerSummary.ExpectedDamage).To(Equal(30))
				Expect(attackingPowerSummary.BarrierDamageTaken).To(Equal(1))
				Expect(attackingPowerSummary.ExpectedBarrierDamage).To(Equal(15))
			})

			It("Give summary of the spell attack with barrier burn", func() {
				bandit.Armor = 1
				bandit.Dodge = 1
				bandit.MaxBarrier = 10
				bandit.CurrentBarrier = 10

				teros.Aim = 3
				teros.Mind = 2
				blot.DamageBonus = 4
				blot.ExtraBarrierDamage = 3
				attackingPowerSummary := usecase.GetExpectedDamage(blot, teros, bandit)
				Expect(attackingPowerSummary.ChanceToHit).To(Equal(33))
				Expect(attackingPowerSummary.DamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.ExpectedDamage).To(Equal(0))
				Expect(attackingPowerSummary.BarrierDamageTaken).To(Equal(9))
				Expect(attackingPowerSummary.ExpectedBarrierDamage).To(Equal(9 * 33))
			})
		})

		Context("Critical Hits", func() {
			var (
				bandit *entity.Squaddie
			)

			BeforeEach(func() {
				bandit = entity.NewSquaddie("Bandit")
				bandit.Name = "Bandit"

				teros.Strength = 1
				spear.DamageBonus = 3
			})

			It("Adds the chance to crit to the attack summary", func() {
				spear.CriticalHitThreshold = 4
				attackingPowerSummary := usecase.GetExpectedDamage(spear, teros, bandit)
				Expect(attackingPowerSummary.ChanceToCrit).To(Equal(6))
			})

			It("Doubles the damage before applying armor and barrier to the attack summary", func() {
				bandit.Armor = 1
				bandit.MaxBarrier = 4
				bandit.CurrentBarrier = 4
				spear.CriticalHitThreshold = 4
				attackingPowerSummary := usecase.GetExpectedDamage(spear, teros, bandit)
				Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(3))
				Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(4))
				Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(3 * 21))
				Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(4 * 21))
			})

			It("Does not factor critcal effects if the attack cannot crit", func() {
				spear.CriticalHitThreshold = 0
				attackingPowerSummary := usecase.GetExpectedDamage(spear, teros, bandit)
				Expect(attackingPowerSummary.ChanceToCrit).To(Equal(0))
				Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(0))
				Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(0))
				Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(0))
			})
		})
	})
	Context("Give squaddie powers", func() {
		var (
			teros *entity.Squaddie
			powerRepository *repository.PowerRepository
			spear *entity.Power
		)
		BeforeEach(func() {
			powerRepository = repository.NewPowerRepository()

			spear = entity.NewAttackingPower("Spear")
			spear.PowerType = entity.PowerTypePhysical
			spear.ID = "deadbeef"
			newPowers := []*entity.Power{spear}
			powerRepository.AddSlicePowerSource(newPowers)

			teros = entity.NewSquaddie("Teros")
			teros.Name = "Teros"
		})
		It("Can give Squaddie innate Powers with a repository", func() {
			temporaryPowerReferences := []*entity.PowerReference{{Name: "Spear", ID: spear.ID}}
			numberOfPowersAdded, err := usecase.LoadAllOfSquaddieInnatePowers(teros, temporaryPowerReferences, powerRepository)
			Expect(numberOfPowersAdded).To(Equal(1))
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Spear"))
			Expect(attackIDNamePairs[0].ID).To(Equal(spear.ID))
		})
		It("Stop adding Powers to Squaddie if it doesn't exist", func() {
			scimitar := entity.NewAttackingPower("Scimitar")
			scimitar.PowerType = entity.PowerTypePhysical

			temporaryPowerReferences := []*entity.PowerReference{{Name: "Scimitar", ID: scimitar.ID}}
			numberOfPowersAdded, err := usecase.LoadAllOfSquaddieInnatePowers(teros, temporaryPowerReferences, powerRepository)
			Expect(numberOfPowersAdded).To(Equal(0))
			Expect(err.Error()).To(Equal("squaddie 'Teros' tried to add Power 'Scimitar' but it does not exist"))

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(0))
		})
	})
})
