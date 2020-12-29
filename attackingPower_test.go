package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("Calculate combination of Attacking Power and Squaddie", func() {
	var (
		teros *terosbattleserver.Squaddie
		spear *terosbattleserver.AttackingPower
		blot  *terosbattleserver.AttackingPower
	)

	BeforeEach(func() {
		tempSquaddie := terosbattleserver.NewSquaddie("Teros")
		teros = &tempSquaddie
		teros.Name = "Teros"

		tempPower := terosbattleserver.NewAttackingPower("Spear")
		spear = &tempPower
		spear.PowerType = terosbattleserver.PowerTypePhysical

		tempPower2 := terosbattleserver.NewAttackingPower("Blot")
		blot = &tempPower2
		blot.PowerType = terosbattleserver.PowerTypeSpell
	})

	It("Calculates the To Hit Bonus", func() {
		teros.Aim = 2
		blot.ToHitBonus = 1

		totalToHitBonus := blot.GetTotalToHitBonus(teros)
		Expect(totalToHitBonus).To(Equal(3))
	})

	It("Calculates the crit chance based on To Hit Bonus", func() {
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(9001)).To(Equal(36))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(5)).To(Equal(36))

		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-6)).To(Equal(0))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-9001)).To(Equal(0))

		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(4)).To(Equal(35))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(3)).To(Equal(33))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(2)).To(Equal(30))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(1)).To(Equal(26))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(0)).To(Equal(21))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-1)).To(Equal(15))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-2)).To(Equal(10))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-3)).To(Equal(6))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-4)).To(Equal(3))
		Expect(terosbattleserver.GetChanceToHitBasedOnHitRate(-5)).To(Equal(1))
	})

	Context("Calculate damage bonus", func() {
		BeforeEach(func() {
			teros.Strength = 2
			teros.Mind = 3

			spear.PowerType = terosbattleserver.PowerTypePhysical
			spear.DamageBonus = 2

			blot.PowerType = terosbattleserver.PowerTypeSpell
			blot.DamageBonus = 6
		})

		It("Calculates the Damage bonus of physical attacks", func() {
			totalDamageBonus := spear.GetTotalDamageBonus(teros)
			Expect(totalDamageBonus).To(Equal(4))
		})

		It("Calculates the Damage bonus of spell attacks", func() {
			totalDamageBonus := blot.GetTotalDamageBonus(teros)
			Expect(totalDamageBonus).To(Equal(9))
		})

		It("Calculates the Critical Damage bonus of physical attacks", func() {
			totalDamageBonus := spear.GetCriticalDamageBonus(teros)
			Expect(totalDamageBonus).To(Equal(8))
		})

	})

	Context("Calculate to hit penalties values against powers", func() {
		BeforeEach(func() {
			teros.Dodge = 2
			teros.Deflect = 9001

			spear.PowerType = terosbattleserver.PowerTypePhysical

			blot.PowerType = terosbattleserver.PowerTypeSpell
		})

		It("Calculates the to hit reduction against physical attacks", func() {
			toHitPenalty := spear.GetToHitPenalty(teros)
			Expect(toHitPenalty).To(Equal(2))
		})

		It("Calculates the to hit reduction against spell attacks", func() {
			toHitPenalty := blot.GetToHitPenalty(teros)
			Expect(toHitPenalty).To(Equal(9001))
		})
	})

	Context("Calculate damage if the attacker hits the target with the power", func() {
		var (
			bandit *terosbattleserver.Squaddie
		)

		BeforeEach(func() {
			tempSquaddie := terosbattleserver.NewSquaddie("Bandit")
			bandit = &tempSquaddie
			bandit.Name = "Bandit"

			teros.Strength = 1
			spear.DamageBonus = 3
		})

		It("Does full damage against targets without armor or barrier", func() {
			totalHealthDamage, _, _ := spear.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(4))
		})

		It("Armor reduces damage against physical attacks", func() {
			bandit.Armor = 3
			totalHealthDamage, _, _ := spear.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(1))
		})

		It("Barrier absorbs damage against physical attacks and is depleted first", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 1
			totalHealthDamage, initialBarrierDamage, _ := spear.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(3))
			Expect(initialBarrierDamage).To(Equal(1))
		})

		It("Will deal no damage if barrier is strong enough", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 4
			totalHealthDamage, initialBarrierDamage, _ := spear.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(initialBarrierDamage).To(Equal(4))
		})

		It("May deal no damage if armor is strong enough", func() {
			bandit.Armor = 4
			totalHealthDamage, initialBarrierDamage, _ := spear.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(initialBarrierDamage).To(Equal(0))
		})
	})

	Context("Calculate damage if the attacker hits the target with the power", func() {
		var (
			bandit *terosbattleserver.Squaddie
		)

		BeforeEach(func() {
			tempSquaddie := terosbattleserver.NewSquaddie("Bandit")
			bandit = &tempSquaddie
			bandit.Name = "Bandit"

			teros.Mind = 2
			blot.DamageBonus = 4
		})

		It("Does full damage against targets without armor or barrier", func() {
			totalHealthDamage, _, _ := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(6))
		})

		It("Ignores Armor when using spell attacks", func() {
			bandit.Armor = 9001
			totalHealthDamage, _, _ := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(6))
		})

		It("Barrier absorbs damage against spell attacks and is depleted first", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 1
			totalHealthDamage, initialBarrierDamage, _ := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(5))
			Expect(initialBarrierDamage).To(Equal(1))
		})

		It("Will deal no damage if barrier is strong enough", func() {
			bandit.MaxBarrier = 9001
			bandit.CurrentBarrier = 9001
			totalHealthDamage, initialBarrierDamage, _ := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(initialBarrierDamage).To(Equal(6))
		})

		It("Can deal extra Barrier damage if the barrier absorbs the attack", func() {
			bandit.MaxBarrier = 8
			bandit.CurrentBarrier = 8
			blot.ExtraBarrierDamage = 2

			totalHealthDamage, initialBarrierDamage, extraBarrierDamage := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(initialBarrierDamage).To(Equal(6))
			Expect(extraBarrierDamage).To(Equal(2))
		})

		It("Knows extra Barrier damage is reduced if the barrier is depleted", func() {
			bandit.MaxBarrier = 8
			bandit.CurrentBarrier = 7
			blot.ExtraBarrierDamage = 2

			totalHealthDamage, initialBarrierDamage, extraBarrierDamage := blot.GetHowTargetDistributesDamage(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(initialBarrierDamage).To(Equal(6))
			Expect(extraBarrierDamage).To(Equal(1))
		})
	})

	Context("Calculate expected damage", func() {
		var (
			bandit *terosbattleserver.Squaddie
		)

		BeforeEach(func() {
			tempSquaddie := terosbattleserver.NewSquaddie("Bandit")
			bandit = &tempSquaddie
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

			attackingPowerSummary := spear.GetExpectedDamage(teros, bandit)
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
			attackingPowerSummary := blot.GetExpectedDamage(teros, bandit)
			Expect(attackingPowerSummary.ChanceToHit).To(Equal(33))
			Expect(attackingPowerSummary.DamageTaken).To(Equal(0))
			Expect(attackingPowerSummary.ExpectedDamage).To(Equal(0))
			Expect(attackingPowerSummary.BarrierDamageTaken).To(Equal(9))
			Expect(attackingPowerSummary.ExpectedBarrierDamage).To(Equal(9 * 33))
		})
	})

	Context("Critical Hits", func() {
		var (
			bandit *terosbattleserver.Squaddie
		)

		BeforeEach(func() {
			tempSquaddie := terosbattleserver.NewSquaddie("Bandit")
			bandit = &tempSquaddie
			bandit.Name = "Bandit"

			teros.Strength = 1
			spear.DamageBonus = 3
		})

		It("Can look up the chance of a critical hit", func() {
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(9001)).To(Equal(36))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(12)).To(Equal(36))

			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(-6)).To(Equal(0))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(0)).To(Equal(0))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(1)).To(Equal(0))

			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(2)).To(Equal(1))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(3)).To(Equal(3))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(4)).To(Equal(6))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(5)).To(Equal(10))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(6)).To(Equal(15))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(7)).To(Equal(21))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(8)).To(Equal(26))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(9)).To(Equal(30))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(10)).To(Equal(33))
			Expect(terosbattleserver.GetChanceToCritBasedOnThreshold(11)).To(Equal(35))
		})

		It("Adds the chance to crit to the attack summary", func() {
			spear.CriticalHitThreshold = 4
			attackingPowerSummary := spear.GetExpectedDamage(teros, bandit)
			Expect(attackingPowerSummary.ChanceToCrit).To(Equal(6))
		})

		It("Doubles the damage before applying armor and barrier to the attack summary", func() {
			bandit.Armor = 1
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 4
			spear.CriticalHitThreshold = 4
			attackingPowerSummary := spear.GetExpectedDamage(teros, bandit)
			Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(3))
			Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(4))
			Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(3 * 21))
			Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(4 * 21))
		})

		It("Does not factor critcal effects if the attack cannot crit", func() {
			spear.CriticalHitThreshold = 0
			attackingPowerSummary := spear.GetExpectedDamage(teros, bandit)
			Expect(attackingPowerSummary.ChanceToCrit).To(Equal(0))
			Expect(attackingPowerSummary.CriticalDamageTaken).To(Equal(0))
			Expect(attackingPowerSummary.CriticalBarrierDamageTaken).To(Equal(0))
			Expect(attackingPowerSummary.CriticalExpectedDamage).To(Equal(0))
			Expect(attackingPowerSummary.CriticalExpectedBarrierDamage).To(Equal(0))
		})
	})

	Context("Creating powers from datastreams", func() {
		It("Can create powers from JSON", func() {
			byteStream := []byte(`{
				"name": "Scimitar",
				"damage_bonus": 2
			}`)
			scimitar, err := terosbattleserver.NewAttackingPowerFromJSON(byteStream)
			Expect(err).To(BeNil())
			Expect(scimitar.Name).To(Equal("Scimitar"))
			Expect(scimitar.DamageBonus).To(Equal(2))
		})

		It("Throws an error if Power is created with wrong power type in JSON", func() {
			byteStream := []byte(`{
				"name": "Scimitar",
				"power_type": "mystery"
			}`)
			_, err := terosbattleserver.NewAttackingPowerFromJSON(byteStream)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Power 'Scimitar' has unknown power_type: 'mystery'"))
		})

		It("Can create a Power using YAML", func() {
			byteStream := []byte(`
name: Scimitar
damage_bonus: 2
`)
			scimitar, err := terosbattleserver.NewAttackingPowerFromYAML(byteStream)
			Expect(err).To(BeNil())
			Expect(scimitar.Name).To(Equal("Scimitar"))
			Expect(scimitar.DamageBonus).To(Equal(2))
		})
	})
})
