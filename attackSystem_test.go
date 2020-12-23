package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("Calculate Squaddie stats", func() {
	It("Sets the Squaddie name.", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		Expect(teros.Name).To(Equal("T'eros"))
	})

	It("Sets current HP to max.", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		Expect(teros.MaxHitPoints).To(Equal(5))
		teros.SetHPToMax()
		Expect(teros.CurrentHitPoints).To(Equal(5))
	})

	It("Can set Barrier to Max Barrier", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		teros.MaxBarrier = 2
		teros.SetBarrierToMax()
		Expect(teros.CurrentBarrier).To(Equal(2))
	})

	It("Can calculate defenses against physical attacks", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		teros.Armor = 2
		teros.Dodge = 3
		teros.Deflect = 4
		teros.MaxBarrier = 1
		teros.SetBarrierToMax()
		evasion, barrierDamageReduction, armorDamageReduction := teros.GetDefensiveStatsAgainstPhysical()
		Expect(evasion).To(Equal(3))
		Expect(barrierDamageReduction).To(Equal(1))
		Expect(armorDamageReduction).To(Equal(2))
	})

	It("Can calculate defenses against spell attacks", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		teros.Armor = 2
		teros.Dodge = 3
		teros.Deflect = 4
		teros.MaxBarrier = 1
		teros.SetBarrierToMax()
		evasion, barrierDamageReduction, armorDamageReduction := teros.GetDefensiveStatsAgainstSpell()
		Expect(evasion).To(Equal(4))
		Expect(barrierDamageReduction).To(Equal(1))
		Expect(armorDamageReduction).To(Equal(0))
	})

	It("Can calculate offense with physical attacks", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		teros.Aim = 2
		teros.Strength = 3
		teros.Mind = 4
		toHitBonus, damageBonus := teros.GetOffensiveStatsWithPhysical()
		Expect(toHitBonus).To(Equal(2))
		Expect(damageBonus).To(Equal(3))
	})

	It("Can calculate offense with spell attacks", func() {
		teros := terosbattleserver.NewSquaddie("T'eros")
		teros.Aim = 2
		teros.Strength = 3
		teros.Mind = 4
		toHitBonus, damageBonus := teros.GetOffensiveStatsWithSpell()
		Expect(toHitBonus).To(Equal(2))
		Expect(damageBonus).To(Equal(4))
	})
})
