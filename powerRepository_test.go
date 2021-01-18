package terosbattleserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	terosbattleserver "github.com/cserrant/terosBattleServer"
)

var _ = Describe("CRUD Powers", func() {
	var (
		repo *terosbattleserver.PowerRepository
	)
	BeforeEach(func() {
		repo = terosbattleserver.NewPowerRepository()
	})
	Context("Load Power using JSON sources", func() {
		It("Can create powers from JSON", func() {
			Expect(repo.GetNumberOfPowers()).To(Equal(0))
			jsonByteStream := []byte(`[{
					"name": "Scimitar",
					"id": "deadbeef",
					"damage_bonus": 2,
					"power_type": "Physical"
				}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfPowers()).To(Equal(1))
		})
		It("Can get a Power by name", func() {
			jsonByteStream := []byte(`[{
				"name": "Scimitar",
				"id": "deadbeef",
				"damage_bonus": 2,
				"power_type": "Physical"
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			scimitar := repo.GetByName("Scimitar")
			Expect(scimitar.Name).To(Equal("Scimitar"))
			Expect(scimitar.ID).To(Equal("deadbeef"))
			Expect(scimitar.DamageBonus).To(Equal(2))

			missingno := repo.GetByName("Does not exist")
			Expect(missingno).To(BeNil())
		})
		It("Stops loading Powers upon validating the first invalid Power", func() {
			jsonByteStream := []byte(`[{
				"name": "Scimitar",
				"id": "deadbeef",
				"power_type": "Physical"
			},{
				"name": "Scimitar2",
				"id": "deadbeee",
				"power_type": "mystery"
			}]`)
			success, err := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeFalse())
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("AttackingPower 'Scimitar2' has unknown power_type: 'mystery'"))
		})
	})
	Context("Load Power using YAML sources", func() {
		It("Can create a AttackingPower using YAML", func() {
			Expect(repo.GetNumberOfPowers()).To(Equal(0))
			yamlByteStream := []byte(`-
  name: Scimitar
  damage_bonus: 2,
  power_type: Physical
`)
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfPowers()).To(Equal(1))
		})
	})
})
