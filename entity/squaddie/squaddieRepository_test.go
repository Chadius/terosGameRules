package squaddie_test

import (
	"bytes"
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CRUD Squaddies", func() {
	var (
		repo *squaddie.Repository
	)
	BeforeEach(func() {
		repo = squaddie.NewSquaddieRepository()
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
		It("gets a clone of the squaddie by name", func() {
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player"
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			teros0 := repo.GetByName("Teros")
			teros1 := repo.GetByName("Teros")
			Expect(teros0).ToNot(BeIdenticalTo(teros1))
		})
		It("Can load class levels", func() {
			jsonByteStream := []byte(`[{
				"name": "Teros",
				"aim": 5,
				"affiliation": "Player",
				"class_levels": {
                  "1": {
                    "id": "1",
                    "name": "Mage",
                    "levels_gained": ["123"]
                  },
                  "2": {
                    "id": "2",
                    "name": "Dimension Walker"
                  }
               }
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"1": 1, "2": 0}))
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
			Expect(soldier.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Foot)))
			Expect(soldier.CanHitAndRun()).To(BeFalse())

			scout := repo.GetByName("Scout")
			Expect(scout.Name).To(Equal("Scout"))
			Expect(scout.GetMovementDistancePerRound()).To(Equal(4))
			Expect(scout.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Light)))
			Expect(scout.CanHitAndRun()).To(BeFalse())

			bird := repo.GetByName("Bird")
			Expect(bird.Name).To(Equal("Bird"))
			Expect(bird.GetMovementDistancePerRound()).To(Equal(3))
			Expect(bird.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Fly)))
			Expect(bird.CanHitAndRun()).To(BeTrue())

			teleporter := repo.GetByName("Teleporter")
			Expect(teleporter.Name).To(Equal("Teleporter"))
			Expect(teleporter.GetMovementDistancePerRound()).To(Equal(2))
			Expect(teleporter.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Teleport)))
			Expect(teleporter.CanHitAndRun()).To(BeFalse())
		})
	})
	Context("Save Squaddie as JSON", func() {
		var teros *squaddie.Squaddie
		BeforeEach(func() {
			teros = squaddie.NewSquaddie("Teros")
			teros.Armor = 2
			teros.Dodge = 3
			teros.Deflect = 4
			teros.MaxBarrier = 1
		})
		It("Can Marshall a Squaddie into JSON", func() {
			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())

			hasIDNameAffiliationInJSON := bytes.Contains(byteStream, []byte(fmt.Sprintf(`"id":"%s","name":"Teros","affiliation":"Player"`, teros.ID)))
			Expect(hasIDNameAffiliationInJSON).To(BeTrue())

			hasDefaultHitPointsInJSON := bytes.Contains(byteStream, []byte(`"current_hit_points":5,"max_hit_points":5`))
			Expect(hasDefaultHitPointsInJSON).To(BeTrue())

			hasOffensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"aim":0,"strength":0,"mind":0`))
			Expect(hasOffensiveStatsInJSON).To(BeTrue())

			hasDefensiveStatsInJSON := bytes.Contains(byteStream, []byte(`"dodge":3,"deflect":4,"current_barrier":0,"max_barrier":1,"armor":2`))
			Expect(hasDefensiveStatsInJSON).To(BeTrue())

			hasDefaultMovementInJSON := bytes.Contains(byteStream, []byte(`"movement":{"distance":3,"type":"foot","hit_and_run":false}`))
			Expect(hasDefaultMovementInJSON).To(BeTrue())

			hasNoPowersInJSON := bytes.Contains(byteStream, []byte(`"powers":[]`))
			Expect(hasNoPowersInJSON).To(BeTrue())

			hasNoClassesInJSON := bytes.Contains(byteStream, []byte(`"base_class":"","current_class":"","class_levels":{}`))
			Expect(hasNoClassesInJSON).To(BeTrue())
		})
		It("Can Marshall a Squaddie with extraordinary movement into JSON", func() {
			teros.Movement.Type = squaddie.Teleport
			teros.Movement.Distance = 8
			teros.Movement.HitAndRun = true

			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())
			movementJSON := `"movement":{"distance":8,"type":"teleport","hit_and_run":true}`
			containsPowersJson := bytes.Contains(byteStream, []byte(movementJSON))
			Expect(containsPowersJson).To(BeTrue())
		})
		It("Can Marshall a Squaddie with powers into JSON", func() {
			attackA := power.NewPower("Attack Formation A")
			teros.AddInnatePower(attackA)
			byteStream, err := repo.MarshalSquaddieIntoJSON(teros)
			Expect(err).To(BeNil())

			powersJSON := fmt.Sprintf(`"powers":[{"name":"Attack Formation A","id":"%s"}]`, attackA.ID)
			containsPowersJson := bytes.Contains(byteStream, []byte(powersJSON))
			Expect(containsPowersJson).To(BeTrue())
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
    '1':
      id: '1'
      name: Mage
      levels_gained:
      - '123'
    '2':
      id: '2'
      name: Dimension Walker
`)
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())

			teros := repo.GetByName("Teros")
			Expect(teros.GetLevelCountsByClass()).To(Equal(map[string]int{"1": 1, "2": 0}))
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
			Expect(soldier.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Foot)))
			Expect(soldier.CanHitAndRun()).To(BeFalse())

			scout := repo.GetByName("Scout")
			Expect(scout.Name).To(Equal("Scout"))
			Expect(scout.GetMovementDistancePerRound()).To(Equal(4))
			Expect(scout.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Light)))
			Expect(scout.CanHitAndRun()).To(BeFalse())

			bird := repo.GetByName("Bird")
			Expect(bird.Name).To(Equal("Bird"))
			Expect(bird.GetMovementDistancePerRound()).To(Equal(3))
			Expect(bird.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Fly)))
			Expect(bird.CanHitAndRun()).To(BeTrue())

			teleporter := repo.GetByName("Teleporter")
			Expect(teleporter.Name).To(Equal("Teleporter"))
			Expect(teleporter.GetMovementDistancePerRound()).To(Equal(2))
			Expect(teleporter.GetMovementType()).To(Equal(squaddie.MovementType(squaddie.Teleport)))
			Expect(teleporter.CanHitAndRun()).To(BeFalse())
		})
	})
	It("Add Squaddie directly", func() {
		success, err := repo.AddSquaddies([]*squaddie.Squaddie{squaddie.NewSquaddie("Generic")})
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
		Expect(repo.GetNumberOfSquaddies()).To(Equal(1))
	})
	Context("Cloning an existing squaddie", func() {
		var base *squaddie.Squaddie

		BeforeEach(func() {
			base = squaddie.NewSquaddie("Base")
			repo.AddSquaddies([]*squaddie.Squaddie{base})
		})
		It("copies name and affiliation but not ID", func() {
			base.Affiliation = "Enemy"
			clone, err := repo.CloneSquaddie(base, "")
			Expect(err).To(BeNil())
			Expect(clone.Name).To(Equal(base.Name))
			Expect(clone.Affiliation).To(Equal(base.Affiliation))
			Expect(clone.ID).ToNot(Equal(base.ID))
		})
		It("will set the clone ID to the given ID", func() {
			clone, _ := repo.CloneSquaddie(base, "12345")
			Expect(clone.ID).To(Equal("12345"))
		})
		It("will copy basic stats", func() {
			base.CurrentHitPoints = 1
			base.MaxHitPoints += 5
			base.CurrentBarrier = 2
			base.MaxBarrier += 5

			base.Aim = 2
			base.Strength = 3
			base.Mind = 4
			base.Dodge = 5
			base.Deflect = 6
			base.Armor = 7

			clone, _ := repo.CloneSquaddie(base, "")
			Expect(clone.CurrentHitPoints).To(Equal(base.CurrentHitPoints))
			Expect(clone.MaxHitPoints).To(Equal(base.MaxHitPoints))
			Expect(clone.Aim).To(Equal(base.Aim))
			Expect(clone.Strength).To(Equal(base.Strength))
			Expect(clone.Mind).To(Equal(base.Mind))
			Expect(clone.Dodge).To(Equal(base.Dodge))
			Expect(clone.Deflect).To(Equal(base.Deflect))
			Expect(clone.CurrentBarrier).To(Equal(base.CurrentBarrier))
			Expect(clone.MaxBarrier).To(Equal(base.MaxBarrier))
			Expect(clone.Armor).To(Equal(base.Armor))
		})
		It("will copy Movement", func() {
			base.Movement = squaddie.Movement{
				Distance:  base.Movement.Distance + 2,
				Type:      squaddie.Fly,
				HitAndRun: true,
			}

			clone, _ := repo.CloneSquaddie(base, "")
			Expect(clone.Movement.Distance).To(Equal(base.Movement.Distance))
			Expect(clone.Movement.Type).To(Equal(base.Movement.Type))
			Expect(clone.Movement.HitAndRun).To(Equal(base.Movement.HitAndRun))
		})
		It("will copy PowerReferences", func() {
			attackA := power.NewPower("Attack Formation A")
			base.AddInnatePower(attackA)
			clone, _ := repo.CloneSquaddie(base, "")

			attackIDNamePairs := clone.GetInnatePowerIDNames()
			Expect(len(attackIDNamePairs)).To(Equal(1))
			Expect(attackIDNamePairs[0].Name).To(Equal("Attack Formation A"))
			Expect(attackIDNamePairs[0].ID).To(Equal(attackA.ID))
		})
		It("will copy class information", func() {
			initialClass := &squaddieclass.Class{
				ID: "initial",
				Name: "Initial Class",
				BaseClassRequired: false,
				InitialBigLevelID: "",
			}
			advancedClass := &squaddieclass.Class{
				ID: "advanced",
				Name: "Advanced Class",
				BaseClassRequired: true,
				InitialBigLevelID: "advanceLevel0",
			}

			base.AddClass(initialClass)
			base.AddClass(advancedClass)
			base.SetBaseClassIfNoBaseClass(initialClass.ID)
			base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel0")
			base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel1")
			base.MarkLevelUpBenefitAsConsumed(initialClass.ID, "initialLevel2")

			clone, _ := repo.CloneSquaddie(base, "")
			Expect(clone.BaseClassID).To(Equal(base.BaseClassID))
			Expect(clone.CurrentClass).To(Equal(base.CurrentClass))
			for classID, levelsConsumed := range base.ClassLevelsConsumed {
				Expect(clone.ClassLevelsConsumed).To(HaveKey(classID))

				cloneLevelsConsumed := clone.ClassLevelsConsumed[classID]
				Expect(cloneLevelsConsumed).NotTo(BeIdenticalTo(levelsConsumed))
				Expect(cloneLevelsConsumed).To(Equal(levelsConsumed))
			}
		})
	})
})
