package repository_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cserrant/terosBattleServer/entity"
	"github.com/cserrant/terosBattleServer/repository"
)

var _ = Describe("CRUD Squaddies", func() {
	var (
		repo *repository.SquaddieRepository
	)
	BeforeEach(func() {
		repo = repository.NewSquaddieRepository()
	})
	Context("Load Squaddie using JSON sources", func() {
		It("Can add a JSON source", func() {
			Expect(repo.GetNumberOfSquaddies()).To(Equal(0))
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfSquaddies()).To(Equal(1))
		})
		It("Can get a Squaddie by name", func() {
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")
			Expect(teros).NotTo(BeNil())
			Expect(teros.Name).To(Equal("Teros"))
			Expect(teros.Aim).To(Equal(5))

			missingno := repo.GetByName("Does not exist")
			Expect(missingno).To(BeNil())
		})
		It("Stops loading Squaddies upon validating the first invalid Squaddie", func() {
			jsonByteStream := []byte(`[{
				"Name": "Teros",
				"Affiliation": "Player"
			},{
				"Name": "Teros2",
				"Affiliation": "Unknown Affiliation"
			}]`)
			success, err := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeFalse())
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})
		It("Can give Squaddie innate Powers with a repository", func() {
			powerRepo := repository.NewPowerRepository()
			powerByteStream := []byte(`[{
				"name": "Scimitar",
				"id": "deadbeef",
				"damage_bonus": 2,
				"power_type": "Physical"
			}]`)
			powerRepo.AddJSONSource(powerByteStream)
			scimitar := powerRepo.GetByName("Scimitar")

			squaddieByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player",
				"powers": [{"name": "Scimitar", "id": "deadbeef"}]
			}]`)
			success, _ := repo.AddJSONSource(squaddieByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")

			numberOfPowersAdded, err := repo.AddInnatePowersToSquaddie(teros, powerRepo)
			Expect(numberOfPowersAdded).To(Equal(1))
			Expect(err).To(BeNil())

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Scimitar"))
			Expect(attackIDNamePairs[0].ID).To(Equal(scimitar.ID))
		})
		It("Stop adding Powers to Squaddie if it doesn't exist", func() {
			powerRepo := repository.NewPowerRepository()

			squaddieByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player",
				"powers": [{"name": "Scimitar", "id": "deadbeef"}]
			}]`)
			success, _ := repo.AddJSONSource(squaddieByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")

			numberOfPowersAdded, err := repo.AddInnatePowersToSquaddie(teros, powerRepo)
			Expect(numberOfPowersAdded).To(Equal(0))
			Expect(err.Error()).To(Equal("squaddie 'Teros' tried to add Power 'Scimitar' but it does not exist"))

			attackIDNamePairs := teros.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(0))
		})
	})
	Context("Save Squaddie as JSON", func() {
		It("Can Marshall a Squaddie into JSON", func() {
			teros := entity.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1
			byteStream, err := repo.MarshalSquaddieIntoJSON(&teros)
			Expect(err).To(BeNil())
			Expect(byteStream).To(Equal([]byte(`{"name":"Teros","affiliation":"Player","current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,"powers":[]}`)))
		})
		It("Can Marshall a Squaddie with powers into JSON", func() {
			teros := entity.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1

			attackA := entity.NewAttackingPower("Attack Formation A")
			teros.GainInnatePower(&attackA)
			byteStream, err := repo.MarshalSquaddieIntoJSON(&teros)
			Expect(err).To(BeNil())

			powersJSON := fmt.Sprintf(`"powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
			Expect(byteStream).To(Equal([]byte(fmt.Sprintf(`{"name":"Teros","affiliation":"Player","current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,%s}`, powersJSON))))
		})
	})
	Context("Load Squaddie using YAML sources", func() {
		It("Can add a YAML source", func() {
			Expect(repo.GetNumberOfSquaddies()).To(Equal(0))
			yamlByteStream := []byte(`-
  name: Teros
  aim: 5
  max_barrier: 3
  affiliation: Player
`)

			repo.AddYAMLSource(yamlByteStream)
			Expect(repo.GetNumberOfSquaddies()).To(Equal(1))
		})
		It("Stops loading Squaddies upon validating the first invalid Squaddie", func() {
			yamlByteStream := []byte(`-
  name: Teros
  affiliation: Player
-
  name: Teros2
  affiliation: Unknown Affiliation`)
			success, err := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeFalse())
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal("Squaddie has unknown affiliation: 'Unknown Affiliation'"))
		})
	})
})
