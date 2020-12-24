package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("Calculate combination of Attacking Power and Squaddie", func() {
	It("Calculates the To Hit Bonus", func() {
		teros := terosbattleserver.NewSquaddie("Teros")
		teros.Aim = 2

		blot := terosbattleserver.NewAttackingPower("Blot")
		blot.ToHitBonus = 1

		totalToHitBonus := blot.GetTotalToHitBonus(&teros)
		Expect(totalToHitBonus).To(Equal(3))
	})

	Context("Calculate damage bonus", func() {
		var (
			teros *terosbattleserver.Squaddie
			spear *terosbattleserver.AttackingPower
			blot  *terosbattleserver.AttackingPower
		)

		BeforeEach(func() {
			tempSquaddie := terosbattleserver.NewSquaddie("Teros")
			teros = &tempSquaddie
			teros.Name = "Teros"
			teros.Strength = 2
			teros.Mind = 3

			tempPower := terosbattleserver.NewAttackingPower("Spear")
			spear = &tempPower
			spear.PowerType = terosbattleserver.PowerTypePhysical
			spear.DamageBonus = 2

			tempPower2 := terosbattleserver.NewAttackingPower("Blot")
			blot = &tempPower2
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
})
