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
			totalHealthDamage, _ := spear.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(4))
		})

		It("Armor reduces damage against physical attacks", func() {
			bandit.Armor = 3
			totalHealthDamage, _ := spear.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(1))
		})

		It("Barrier absorbs damage against physical attacks and is depleted first", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 1
			totalHealthDamage, totalBarrierDamage := spear.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(3))
			Expect(totalBarrierDamage).To(Equal(1))
		})

		It("Will deal no damage if barrier is strong enough", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 4
			totalHealthDamage, totalBarrierDamage := spear.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(totalBarrierDamage).To(Equal(4))
		})

		It("May deal no damage if armor is strong enough", func() {
			bandit.Armor = 4
			totalHealthDamage, totalBarrierDamage := spear.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(totalBarrierDamage).To(Equal(0))
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
			totalHealthDamage, _ := blot.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(6))
		})

		It("Ignores Armor when using spell attacks", func() {
			bandit.Armor = 9001
			totalHealthDamage, _ := blot.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(6))
		})

		It("Barrier absorbs damage against spell attacks and is depleted first", func() {
			bandit.MaxBarrier = 4
			bandit.CurrentBarrier = 1
			totalHealthDamage, totalBarrierDamage := blot.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(5))
			Expect(totalBarrierDamage).To(Equal(1))
		})

		It("Will deal no damage if barrier is strong enough", func() {
			bandit.MaxBarrier = 9001
			bandit.CurrentBarrier = 9001
			totalHealthDamage, totalBarrierDamage := blot.GetDamageAgainstTarget(teros, bandit)
			Expect(totalHealthDamage).To(Equal(0))
			Expect(totalBarrierDamage).To(Equal(6))
		})
	})

})
