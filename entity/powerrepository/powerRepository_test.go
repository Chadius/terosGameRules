package powerrepository_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerCreationSuite struct {
	spear  *power.Power
	spear2 *power.Power
	repo   *powerrepository.Repository
}

var _ = Suite(&PowerCreationSuite{})

func (suite *PowerCreationSuite) SetUpTest(checker *C) {
	suite.spear = power.Builder().Spear().WithName("Spear").WithID("spearLevel1").Build()
	suite.spear2 = power.Builder().Spear().WithName("Spear").WithID("spearLevel2").Build()

	newPowers := []*power.Power{suite.spear, suite.spear2}

	suite.repo = powerrepository.NewPowerRepository()
	suite.repo.AddSlicePowerSource(newPowers)
}

func (suite *PowerCreationSuite) TestAddPowersToNewRepository(checker *C) {
	newRepo := powerrepository.NewPowerRepository()
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 0)
	spear := power.Builder().Spear().Build()
	newPowers := []*power.Power{spear}
	success, _ := newRepo.AddSlicePowerSource(newPowers)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)
}

func (suite *PowerCreationSuite) TestUsesIDToRetrievePowers(checker *C) {
	checker.Assert(suite.repo.GetNumberOfPowers(), Equals, 2)

	spearLevel1FromRepo := suite.repo.GetPowerByID(suite.spear.ID())
	checker.Assert(spearLevel1FromRepo.Name(), Equals, "Spear")
	checker.Assert(spearLevel1FromRepo.ID(), Equals, suite.spear.ID())
	checker.Assert(spearLevel1FromRepo.ToHitBonus(), Equals, suite.spear.ToHitBonus())

	spearLevel2FromRepo := suite.repo.GetPowerByID(suite.spear2.ID())
	checker.Assert(spearLevel2FromRepo.Name(), Equals, "Spear")
	checker.Assert(spearLevel2FromRepo.ID(), Equals, suite.spear2.ID())
	checker.Assert(spearLevel2FromRepo.ToHitBonus(), Equals, suite.spear2.ToHitBonus())
}

func (suite *PowerCreationSuite) TestReturnNilIfIDDoesNotExist(checker *C) {
	nonExistentPower := suite.repo.GetPowerByID("Nope")
	checker.Assert(nonExistentPower, IsNil)
}

func (suite *PowerCreationSuite) TestSearchForPowerByName(checker *C) {
	allSpearPowers := suite.repo.GetAllPowersByName("Spear")
	checker.Assert(allSpearPowers, HasLen, 2)

	hasSpearPower := false
	hasSpear2Power := false
	for _, power := range allSpearPowers {
		if power.PowerID == suite.spear.PowerID {
			hasSpearPower = true
		}
		if power.PowerID == suite.spear2.PowerID {
			hasSpear2Power = true
		}
	}

	checker.Assert(hasSpearPower, Equals, true)
	checker.Assert(hasSpear2Power, Equals, true)
}

func (suite *PowerCreationSuite) TestLoadPowersWithJSON(checker *C) {
	jsonByteStream := []byte(`[{
					"id": "deadbeef",
					"name": "Scimitar",
					"can_attack": true,
					"damage_bonus": 2,
					"power_type": "physical",
					"target_foe": true,
					"can_counter_attack": true,
					"counter_attack_penalty_reduction": -2
				}]`)
	newRepo := powerrepository.NewPowerRepository()
	success, err := newRepo.AddJSONSource(jsonByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)
	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.Name(), Equals, "Scimitar")
	checker.Assert(scimitar.ID(), Equals, "deadbeef")
	checker.Assert(scimitar.DamageBonus(), Equals, 2)
	checker.Assert(scimitar.CanCounterAttack(), Equals, true)
}

func (suite *PowerCreationSuite) TestLoadPowersWithYAML(checker *C) {
	yamlByteStream := []byte(`-
  id: deadbeef
  name: Scimitar
  power_type: physical
  target_foe: true
  can_attack: true
  damage_bonus: 2
  can_counter_attack: true
  counter_attack_penalty_reduction: -2
`)
	newRepo := powerrepository.NewPowerRepository()
	success, _ := newRepo.AddYAMLSource(yamlByteStream)
	checker.Assert(success, Equals, true)
	checker.Assert(newRepo.GetNumberOfPowers(), Equals, 1)

	scimitar := newRepo.GetPowerByID("deadbeef")
	checker.Assert(scimitar.Name(), Equals, "Scimitar")
	checker.Assert(scimitar.ID(), Equals, "deadbeef")
	checker.Assert(scimitar.DamageBonus(), Equals, 2)
	checker.Assert(scimitar.CanCounterAttack(), Equals, true)
	checker.Assert(scimitar.CounterAttackPenaltyReduction(), Equals, -2)
	checker.Assert(scimitar.CanPowerTargetFoe(), Equals, true)
}
