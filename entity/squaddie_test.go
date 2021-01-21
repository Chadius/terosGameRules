package entity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cserrant/terosBattleServer/entity"
)

var _ = Describe("Calculate Squaddie stats", func() {
	It("Sets the Squaddie name.", func() {
		teros := entity.NewSquaddie("T'eros")
		Expect(teros.Name).To(Equal("T'eros"))
	})

	It("Sets current HP to max.", func() {
		teros := entity.NewSquaddie("T'eros")
		Expect(teros.MaxHitPoints).To(Equal(5))
		teros.SetHPToMax()
		Expect(teros.CurrentHitPoints).To(Equal(5))
	})

	It("Can set Barrier to Max Barrier", func() {
		teros := entity.NewSquaddie("T'eros")
		teros.MaxBarrier = 2
		teros.SetBarrierToMax()
		Expect(teros.CurrentBarrier).To(Equal(2))
	})

	Context("Check Squaddies for valid data", func() {
		It("Throws an error if Squaddie is created with wrong affiliation", func() {
			newSquaddie := entity.NewSquaddie("T'eros")
			newSquaddie.Affiliation = "Unknown Affiliation"
			err := entity.CheckSquaddieForErrors(&newSquaddie)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})
	})

	Context("Can calculate net offense and defense", func() {
		It("Can calculate defenses against physical attacks", func() {
			teros := entity.NewSquaddie("T'eros")
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
			teros := entity.NewSquaddie("T'eros")
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
			teros := entity.NewSquaddie("T'eros")
			teros.Aim = 2
			teros.Strength = 3
			teros.Mind = 4
			toHitBonus, damageBonus := teros.GetOffensiveStatsWithPhysical()
			Expect(toHitBonus).To(Equal(2))
			Expect(damageBonus).To(Equal(3))
		})

		It("Can calculate offense with spell attacks", func() {
			teros := entity.NewSquaddie("T'eros")
			teros.Aim = 2
			teros.Strength = 3
			teros.Mind = 4
			toHitBonus, damageBonus := teros.GetOffensiveStatsWithSpell()
			Expect(toHitBonus).To(Equal(2))
			Expect(damageBonus).To(Equal(4))
		})
	})

	It("Squaddie can match Power names to load them from a repository", func() {
		attackA := entity.NewAttackingPower("Attack Formation A")
		teros := entity.NewSquaddie("Teros")
		powerNames := []string{"Attack Formation A"}

		teros.GetInnatePowersFromRepository(powerNames, []*entity.Power{&attackA})
		attackIDNamePairs := teros.GetInnatePowerIDNames()
		Expect(len(attackIDNamePairs)).To(Equal(1))
		Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
		Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
	})

	It("Can gain access to powers and report them", func() {
		teros := entity.NewSquaddie("Teros")
		Expect(teros.Name).To(Equal("Teros"))

		attackA := entity.NewAttackingPower("Attack Formation A")
		teros.GainInnatePower(&attackA)

		attackIDNamePairs := teros.GetInnatePowerIDNames()
		Expect(len(attackIDNamePairs)).To(Equal(1))
		Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
		Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
	})
})
