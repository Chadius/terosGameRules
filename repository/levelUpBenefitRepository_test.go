package repository_test

import (
	"github.com/cserrant/terosBattleServer/entity"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cserrant/terosBattleServer/repository"
)

var _ = Describe("CRUD LevelUpBenefits", func() {
	var (
		repo *repository.LevelUpBenefitRepository
		jsonByteStream []byte
		yamlByteStream []byte
	)
	BeforeEach(func() {
		repo = repository.NewLevelUpBenefitRepository()
		jsonByteStream = []byte(`[
  {
    "squaddie_name": "Teros",
    "level_ups_by_class": [
      {
        "class_name": "Mage",
        "level_up_benefits": [
          {
            "id":"abcdefg0",
            "level_up_benefit_type": "small",
            "max_hit_points": 1,
            "aim": 0,
            "strength": 2,
            "mind": 3,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7,
            "powers": [
              {
                "name": "Scimitar",
                "id": "deadbeef"
              }
            ],
            "movement": {
              "distance": 1,
              "type": "teleport",
              "hit_and_run": true
            }
          }
        ]
      }
    ]
  }
]`)
		yamlByteStream = []byte(
`---
- squaddie_name: Teros
  level_ups_by_class:
  - class_name: Mage
    level_up_benefits:
    - id: abcdefg0
      level_up_benefit_type: small
      max_hit_points: 1
      aim: 0
      strength: 2
      mind: 3
      dodge: 4
      deflect: 5
      max_barrier: 6
      armor: 7
      powers:
      - name: Scimitar
        id: deadbeef
      movement:
        distance: 1,
        type: teleport
        hit_and_run": true
`)
	})
	Context("Create LevelUpBenefit objects from different sources", func() {
		It("Can create LevelUpBenefits from JSON", func() {
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(0))
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(1))
		})
		It("Can create LevelUpBenefits from YAML", func() {
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(0))
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfLevelUpBenefits()).To(Equal(1))
		})
	})
	Context("Can search and retrieve LevelUpBenefit objects using descriptors", func() {
		It("Can get LevelUpBenefits by squaddie and class", func() {
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			benefits, err := repo.GetLevelUpBenefitsByNameAndClass("Teros", "Mage")
			Expect(err).To(BeNil())
			Expect(len(benefits)).To(Equal(1))

			firstBenefit := benefits[0]
			Expect(firstBenefit.LevelUpBenefitType).To(Equal(entity.LevelUpBenefitTypeSmall))
			Expect(firstBenefit.SquaddieName).To(Equal("Teros"))
			Expect(firstBenefit.ClassName).To(Equal("Mage"))
			Expect(firstBenefit.MaxHitPoints).To(Equal(1))
			Expect(firstBenefit.Aim).To(Equal(0))
			Expect(firstBenefit.Strength).To(Equal(2))
			Expect(firstBenefit.Mind).To(Equal(3))
			Expect(firstBenefit.Dodge).To(Equal(4))
			Expect(firstBenefit.Deflect).To(Equal(5))
			Expect(firstBenefit.MaxBarrier).To(Equal(6))
			Expect(firstBenefit.Armor).To(Equal(7))

			Expect(len(firstBenefit.PowerIDGained)).To(Equal(1))
			Expect(firstBenefit.PowerIDGained[0].Name).To(Equal("Scimitar"))
			Expect(firstBenefit.PowerIDGained[0].ID).To(Equal("deadbeef"))
		})
		It("Raises an error if you search for wrong LevelUpBenefits", func() {
			repo.AddJSONSource(jsonByteStream)

			benefits, err := repo.GetLevelUpBenefitsByNameAndClass("Squaddie not found", "Mage")
			Expect(err.Error()).To(Equal(`no LevelUpBenefits for this squaddie: "Squaddie not found"`))
			Expect(len(benefits)).To(Equal(0))

			benefits, err = repo.GetLevelUpBenefitsByNameAndClass("Teros", "Class not found")
			Expect(err.Error()).To(Equal(`no LevelUpBenefits for this class: "Class not found"`))
			Expect(len(benefits)).To(Equal(0))
		})
	})
	It("Stops loading LevelUpBenefits upon validating the first invalid LevelUpBenefit", func() {
		byteStream := []byte(`[
  {
    "squaddie_name": "Teros",
    "level_ups_by_class": [
      {
        "class_name": "Mage",
        "level_up_benefits": [
          {
            "id":"abcdefg0",
            "level_up_benefit_type": "small",
            "max_hit_points": 1,
            "aim": 0,
            "strength": 2,
            "mind": 3,
            "dodge": 4,
            "deflect": 5,
            "max_barrier": 6,
            "armor": 7,
            "powers": [
              {
                "name": "Scimitar",
                "id": "deadbeef"
              }
            ]
          },
		  {
				"level_up_benefit_type": "unknown",
				"max_hit_points": 1,
				"aim": 0,
				"strength": 2,
				"mind": 3,
				"dodge": 4,
				"deflect": 5,
				"max_barrier": 6,
				"armor": 7,
				"powers": [{"name": "Scimitar", "id": "deadbeef"}]
          }
        ]
      }
    ]
  }
]`)
		success, err := repo.AddJSONSource(byteStream)
		Expect(success).To(BeFalse())
		Expect(err.Error()).To(Equal(`unknown level up benefit type: "unknown"`))
	})
})