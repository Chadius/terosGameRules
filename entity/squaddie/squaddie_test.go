package squaddie_test

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieClass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manage Squaddie stats and Powers", func() {
	It("Sets the Squaddie name.", func() {
		teros := squaddie.NewSquaddie("Teros")
		Expect(teros.Name).To(Equal("Teros"))
	})

	It("Will get a random ID if none is given", func() {
		teros := squaddie.NewSquaddie("Teros")
		Expect(teros.ID).NotTo(BeNil())
		Expect(teros.ID).NotTo(Equal(""))
	})

	It("Sets current HP to max.", func() {
		teros := squaddie.NewSquaddie("Teros")
		Expect(teros.MaxHitPoints).To(Equal(5))
		teros.SetHPToMax()
		Expect(teros.CurrentHitPoints).To(Equal(5))
	})

	It("Can set Barrier to Max Barrier", func() {
		teros := squaddie.NewSquaddie("Teros")
		teros.MaxBarrier = 2
		teros.SetBarrierToMax()
		Expect(teros.CurrentBarrier).To(Equal(2))
	})

	Context("Default Settings", func() {
		var teros *squaddie.Squaddie
		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
		})
		It("Max Hit Points is set to 5", func() {
			Expect(teros.MaxHitPoints).To(Equal(5))
			Expect(teros.CurrentHitPoints).To(Equal(5))
		})
		It("Default movement is 3 on foot", func() {
			Expect(teros.GetMovementDistancePerRound()).To(Equal(3))
			Expect(teros.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Foot)))
		})
	})

	Context("Check Squaddies for valid data", func() {
		It("Throws an error if Squaddie is created with wrong affiliation", func() {
			newSquaddie := squaddie.NewSquaddie("Teros")
			newSquaddie.Affiliation = "Unknown Affiliation"
			err := squaddie.CheckSquaddieForErrors(newSquaddie)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})
	})

	Context("Can calculate net offense and defense", func() {
		It("Can calculate defenses against physical attacks", func() {
			teros := squaddie.NewSquaddie("Teros")
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
			teros := squaddie.NewSquaddie("Teros")
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
			teros := squaddie.NewSquaddie("Teros")
			teros.Aim = 2
			teros.Strength = 3
			teros.Mind = 4
			toHitBonus, damageBonus := teros.GetOffensiveStatsWithPhysical()
			Expect(toHitBonus).To(Equal(2))
			Expect(damageBonus).To(Equal(3))
		})

		It("Can calculate offense with spell attacks", func() {
			teros := squaddie.NewSquaddie("Teros")
			teros.Aim = 2
			teros.Strength = 3
			teros.Mind = 4
			toHitBonus, damageBonus := teros.GetOffensiveStatsWithSpell()
			Expect(toHitBonus).To(Equal(2))
			Expect(damageBonus).To(Equal(4))
		})
	})

	Context("Manage powers", func() {
		It("Can gain access to powers and report them", func() {
			teros := squaddie.NewSquaddie("Teros")
			Expect(teros.Name).To(Equal("Teros"))

			attackA := power.NewPower("Attack Formation A")
			teros.AddInnatePower(attackA)

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
			Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
		})

		It("Clears squaddie known powers", func() {
			teros := squaddie.NewSquaddie("Teros")
			Expect(teros.Name).To(Equal("Teros"))

			attackA := power.NewPower("Attack Formation A")
			teros.AddInnatePower(attackA)
			teros.ClearInnatePowers()

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(attackIDNamePairs).To(BeEmpty())
		})

		It("Clears squaddie temporary power references", func() {
			teros := squaddie.NewSquaddie("Teros")
			teros.PowerReferences = []*power.Reference{{Name: "Pow pow", ID: "Power Wheels"}}

			teros.ClearTemporaryPowerReferences()

			Expect(teros.PowerReferences).To(BeEmpty())
		})

		It("Can remove squaddie powers", func() {
			teros := squaddie.NewSquaddie("Teros")
			Expect(teros.Name).To(Equal("Teros"))

			attackA := power.NewPower("Attack Formation A")
			teros.AddInnatePower(attackA)
			teros.RemovePowerByID(attackA.ID)

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(attackIDNamePairs).To(BeEmpty())
		})

		It("Raises an error if you try to gain the same innate power", func() {
			teros := squaddie.NewSquaddie("Teros")
			Expect(teros.Name).To(Equal("Teros"))

			attackA := power.NewPower("Attack Formation A")
			err := teros.AddInnatePower(attackA)
			Expect(err).To(BeNil())
			err = teros.AddInnatePower(attackA)
			expectedErrorMessage := fmt.Sprintf(`squaddie "Teros" already has innate power with ID "%s"`, attackA.ID)
			Expect(err.Error()).To(Equal(expectedErrorMessage))

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
			Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
		})
	})

	Context("Class levels", func() {
		var (
			teros *squaddie.Squaddie
			mageClass *squaddieClass.Class
			mushroomClass *squaddieClass.Class
		)

		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			mageClass = &squaddieClass.Class{ID: "1", Name: "Mage"}
			mushroomClass = &squaddieClass.Class{ID: "2", Name: "Mushroom"}
		})

		It("Has no class and level upon creation", func() {
			Expect(teros.CurrentClass).To(Equal(""))
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{}))
		})

		It("Can add a class", func() {
			Expect(teros.GetLevelCountsByClass()).To(BeEmpty())
			teros.AddClass(mageClass)
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{mageClass.ID: 0}))
		})

		It("Can tell if a class was already added", func() {
			teros.AddClass(mageClass)
			Expect(teros.HasAddedClass(mageClass.ID)).To(BeTrue())
			Expect(teros.HasAddedClass(mushroomClass.ID)).To(BeFalse())
		})

		It("Can set the current class", func() {
			teros.AddClass(mageClass)
			Expect(teros.CurrentClass).To(Equal(""))
			err := teros.SetClass(mageClass.ID)
			Expect(err).To(BeNil())
			Expect(teros.CurrentClass).To(Equal(mageClass.ID))
		})

		It("Sets the base class", func() {
			teros.AddClass(mageClass)
			Expect(teros.BaseClassID).To(Equal(""))
			teros.SetBaseClassIfNoBaseClass(mageClass.ID)
			Expect(teros.BaseClassID).To(Equal(mageClass.ID))
		})

		It("Raise an error if you set to a class that does not exist", func() {
			teros.AddClass(mageClass)
			Expect(teros.CurrentClass).To(Equal(""))
			err := teros.SetClass(mushroomClass.ID)
			Expect(err.Error()).To(Equal(`cannot switch "Teros" to unknown class "2"`))
		})
	})
})
