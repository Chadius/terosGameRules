package power_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CRUD Powers", func() {
	var (
		repo *power.Repository
	)
	BeforeEach(func() {
		repo = power.NewPowerRepository()
	})
	It("Can add powers directly", func() {
		Expect(repo.GetNumberOfPowers()).To(Equal(0))
		spear := power.NewPower("Spear")
		spear.PowerType = power.Physical
		newPowers := []*power.Power{spear}
		success, _ := repo.AddSlicePowerSource(newPowers)
		Expect(success).To(BeTrue())
		Expect(repo.GetNumberOfPowers()).To(Equal(1))
	})
	Context("Getting Powers from repo", func() {
		var (
			spear *power.Power
			spear2 *power.Power
			repo *power.Repository
		)

		BeforeEach(func() {
			spear = power.NewPower("Spear")
			spear.PowerType = power.Physical
			spear.ID = "spearLevel1"
			spear.AttackEffect.ToHitBonus = 1

			spear2 = power.NewPower("Spear")
			spear2.PowerType = power.Physical
			spear2.ID = "spearLevel2"
			spear2.AttackEffect.ToHitBonus = 2

			newPowers := []*power.Power{spear, spear2}

			repo = power.NewPowerRepository()
			repo.AddSlicePowerSource(newPowers)
		})
		It("Tracks powers by ID even if they have same name", func() {
			Expect(repo.GetNumberOfPowers()).To(Equal(2))
		})
		It("Can get powers by ID", func() {
			spearLevel1FromRepo := repo.GetPowerByID(spear.ID)
			Expect(spearLevel1FromRepo.Name).To(Equal("Spear"))
			Expect(spearLevel1FromRepo.ID).To(Equal(spear.ID))
			Expect(spearLevel1FromRepo.AttackEffect.ToHitBonus).To(Equal(spear.AttackEffect.ToHitBonus))

			spearLevel2FromRepo := repo.GetPowerByID(spear2.ID)
			Expect(spearLevel2FromRepo.Name).To(Equal("Spear"))
			Expect(spearLevel2FromRepo.ID).To(Equal(spear2.ID))
			Expect(spearLevel2FromRepo.AttackEffect.ToHitBonus).To(Equal(spear2.AttackEffect.ToHitBonus))
		})
		It("Returns nil if power does not exist", func() {
			nonExistentPower := repo.GetPowerByID("Nope")
			Expect(nonExistentPower).To(BeNil())
		})
		It("Get all of the powers in repo by name", func() {
			allSpearPowers := repo.GetAllPowersByName("Spear")
			Expect(len(allSpearPowers)).To(Equal(2))
			Expect(allSpearPowers).To(ContainElements([]*power.Power{spear, spear2}))
		})
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
				"power_type": "Physical",
				"attack_effect": {
					"damage_bonus": 2
				}
			}]`)
			success, _ := repo.AddJSONSource(jsonByteStream)
			Expect(success).To(BeTrue())

			scimitar := repo.GetPowerByID("deadbeef")
			Expect(scimitar.Name).To(Equal("Scimitar"))
			Expect(scimitar.ID).To(Equal("deadbeef"))
			Expect(scimitar.AttackEffect.DamageBonus).To(Equal(2))

			missingno := repo.GetPowerByID(("Does not exist"))
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
  id: deadbeef
  power_type: Physical
  attack_effect:
    damage_bonus: 2
`)
			success, _ := repo.AddYAMLSource(yamlByteStream)
			Expect(success).To(BeTrue())
			Expect(repo.GetNumberOfPowers()).To(Equal(1))

			scimitar := repo.GetPowerByID("deadbeef")
			Expect(scimitar.Name).To(Equal("Scimitar"))
			Expect(scimitar.ID).To(Equal("deadbeef"))
			Expect(scimitar.AttackEffect.DamageBonus).To(Equal(2))
		})
	})
})
