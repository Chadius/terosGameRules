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
		It("Can load class levels", func() {
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player",
				"class_levels": {"Mage":["123"],"Dimension Walker":[]}
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1, "Dimension Walker": 0}))
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
		It("Can create Squaddies with different movement", func() {
			jsonByteStream := []byte(`[
			{
				"name": "Soldier",
				"affiliation": "Player",
				"movement": { "distance": 5, "type": "foot"}
			},
			{
				"name": "Scout",
				"affiliation": "Player",
				"movement": { "distance": 4, "type": "light"}
			},
			{
				"name": "Bird",
				"affiliation": "Player",
				"movement": { "distance": 3, "type": "fly", "hit_and_run": true}
			},
			{
				"name": "Teleporter",
				"affiliation": "Player",
				"movement": { "distance": 2, "type": "teleport"}
			}
			]`)

			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfSquaddies()).To(Equal(4))

			soldier := repo.GetByName("Soldier")
			Expect(soldier.Name).To(Equal("Soldier"))
			Expect(soldier.GetMovementDistancePerRound()).To(Equal(5))
			Expect(soldier.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFoot)))
			Expect(soldier.CanHitAndRun()).To(BeFalse())

			scout := repo.GetByName("Scout")
			Expect(scout.Name).To(Equal("Scout"))
			Expect(scout.GetMovementDistancePerRound()).To(Equal(4))
			Expect(scout.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeLight)))
			Expect(scout.CanHitAndRun()).To(BeFalse())

			bird := repo.GetByName("Bird")
			Expect(bird.Name).To(Equal("Bird"))
			Expect(bird.GetMovementDistancePerRound()).To(Equal(3))
			Expect(bird.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFly)))
			Expect(bird.CanHitAndRun()).To(BeTrue())

			teleporter := repo.GetByName("Teleporter")
			Expect(teleporter.Name).To(Equal("Teleporter"))
			Expect(teleporter.GetMovementDistancePerRound()).To(Equal(2))
			Expect(teleporter.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeTeleport)))
			Expect(teleporter.CanHitAndRun()).To(BeFalse())
		})
	})
	Context("Save Squaddie as JSON", func() {
		var teros *entity.Squaddie
		BeforeEach(func() {
			teros = entity.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1
		})
		It("Can Marshall a Squaddie into JSON", func() {
			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())
			Expect(byteStream).To(Equal([]byte(`{"name":"Teros","affiliation":"Player","current_class":"","class_levels":{},"current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,"movement":{"distance":3,"type":"foot","hit_and_run":false},"powers":[]}`)))
		})
		It("Can Marshall a Squaddie with extraordinary movement into JSON", func() {
			teros.Movement.Type = entity.SquaddieMovementTypeTeleport
			teros.Movement.Distance = 8
			teros.Movement.HitAndRun = true

			movementJSON := `"movement":{"distance":8,"type":"teleport","hit_and_run":true}`
			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())
			Expect(byteStream).To(Equal([]byte(fmt.Sprintf(`{"name":"Teros","affiliation":"Player","current_class":"","class_levels":{},"current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,%s,"powers":[]}`, movementJSON))))
		})
		It("Can Marshall a Squaddie with powers into JSON", func() {
			attackA := entity.NewPower("Attack Formation A")
			teros.AddInnatePower(attackA)
			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())

			powersJSON := fmt.Sprintf(`"powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
			Expect(byteStream).To(Equal([]byte(fmt.Sprintf(`{"name":"Teros","affiliation":"Player","current_class":"","class_levels":{},"current_hit_points":5,"max_hit_points":5,"aim":0,"strength":0,"mind":0,"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2,"movement":{"distance":3,"type":"foot","hit_and_run":false},%s}`, powersJSON))))
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
		It("Can load class levels", func() {
			yamlByteStream := []byte(`-
  name: Teros
  aim: 5
  max_barrier: 3
  affiliation: Player
  class_levels:
    Mage: ["hi"]
    Dimension Walker: []
`)
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"Mage": 1, "Dimension Walker": 0}))
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
		It("Can create Squaddies with different movement", func() {
			yamlByteStream := []byte(`-
  name: Soldier
  affiliation: Player
  movement:
    distance: 5
    type: foot
-
  name: Scout
  affiliation: Player
  movement:
    distance: 4
    type: light
-
  name: Bird
  affiliation: Player
  movement:
    distance: 3
    type: fly
    hit_and_run: true
-
  name: Teleporter
  affiliation: Player
  movement:
    distance: 2
    type: teleport
    hit_and_run: false
`)

			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfSquaddies()).To(Equal(4))

			soldier := repo.GetByName("Soldier")
			Expect(soldier.Name).To(Equal("Soldier"))
			Expect(soldier.GetMovementDistancePerRound()).To(Equal(5))
			Expect(soldier.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFoot)))
			Expect(soldier.CanHitAndRun()).To(BeFalse())

			scout := repo.GetByName("Scout")
			Expect(scout.Name).To(Equal("Scout"))
			Expect(scout.GetMovementDistancePerRound()).To(Equal(4))
			Expect(scout.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeLight)))
			Expect(scout.CanHitAndRun()).To(BeFalse())

			bird := repo.GetByName("Bird")
			Expect(bird.Name).To(Equal("Bird"))
			Expect(bird.GetMovementDistancePerRound()).To(Equal(3))
			Expect(bird.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeFly)))
			Expect(bird.CanHitAndRun()).To(BeTrue())

			teleporter := repo.GetByName("Teleporter")
			Expect(teleporter.Name).To(Equal("Teleporter"))
			Expect(teleporter.GetMovementDistancePerRound()).To(Equal(2))
			Expect(teleporter.GetMovementType()).To(Equal(entity.SquaddieMovementType(entity.SquaddieMovementTypeTeleport)))
			Expect(teleporter.CanHitAndRun()).To(BeFalse())
		})
	})
})
