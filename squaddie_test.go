package terosbattleserver_test

import (
	"encoding/json"
	"fmt"

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

	Context("Can calculate net offense and defense", func() {
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

	FContext("Marshalling Squaddies from byte streams", func() {
		It("Can create a Squaddie using JSON", func() {
			byteStream := []byte(`{
			"name": "Teros",
			"aim": 5
		}`)
			teros, err := terosbattleserver.NewSquaddieFromJSON(byteStream)
			Expect(err).To(BeNil())
			Expect(teros.Name).To(Equal("Teros"))
			Expect(teros.Aim).To(Equal(5))
		})

		It("Throws an error if Squaddie is created with wrong affiliation in JSON", func() {
			byteStream := []byte(`{
			"Name": "Teros",
			"Affiliation": "Unknown Affiliation"
		}`)
			_, err := terosbattleserver.NewSquaddieFromJSON(byteStream)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})

		It("Can Marshall a Squaddie", func() {
			teros := terosbattleserver.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1
			byteStream, err := json.Marshal(&teros)
			Expect(err).To(BeNil())
			Expect(byteStream).To(Equal([]byte(`{"name":"Teros","affiliation":"Player","current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,"attacking_powers":[]}`)))
		})

		It("Can Marshall a Squaddie with powers", func() {
			teros := terosbattleserver.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1

			attackA := terosbattleserver.NewAttackingPower("Attack Formation A")
			teros.GainInnatePower(&attackA)
			byteStream, err := json.Marshal(&teros)
			Expect(err).To(BeNil())

			attackingPowersJSON := fmt.Sprintf(`"attacking_powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
			Expect(byteStream).To(Equal([]byte(fmt.Sprintf(`{"name":"Teros","affiliation":"Player","current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,%s}`, attackingPowersJSON))))
		})

		It("Can create a Squaddie using YAML", func() {
			byteStream := []byte(`
name: Teros
aim: 5
max_barrier: 3
`)
			teros, err := terosbattleserver.NewSquaddieFromYAML(byteStream)
			Expect(err).To(BeNil())
			Expect(teros.Name).To(Equal("Teros"))
			Expect(teros.Aim).To(Equal(5))
			Expect(teros.MaxBarrier).To(Equal(3))
		})

		It("Throws an error if Squaddie is created with wrong affiliation in YAML", func() {
			byteStream := []byte(`
name: Teros
affiliation: Unknown Affiliation
`)
			_, err := terosbattleserver.NewSquaddieFromYAML(byteStream)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})
	})

	It("Can gain access to powers", func() {
		teros := terosbattleserver.NewSquaddie("Teros")
		Expect(teros.Name).To(Equal("Teros"))

		attackA := terosbattleserver.NewAttackingPower("Attack Formation A")
		teros.GainInnatePower(&attackA)

		attackIDNamePairs := teros.GetInnatePowerIDNames()
		Expect(len(attackIDNamePairs)).To(Equal(1))
		Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
		Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
	})
})
