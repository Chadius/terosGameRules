package power_test

import (
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	. "gopkg.in/check.v1"
)

type HealingEffectLoadedFromData struct {
	healStaffJSON []byte
	healStaffYAML []byte
	repo          *powerrepository.Repository
}

var _ = Suite(&HealingEffectLoadedFromData{})

func (suite *HealingEffectLoadedFromData) SetUpTest(checker *C) {
	suite.healStaffJSON = []byte(`[{
					"name": "Heal",
					"id": "power_heal",
					"power_type": "spell",
					"can_heal": true,
					"hit_points_healed": 2,
					"healing_adjustment_based_on_user_mind": "full"
				}]`)

	suite.healStaffYAML = []byte(`-
  name: Heal
  id: power_heal
  power_type: spell
  can_heal: true
  hit_points_healed: 2
  healing_adjustment_based_on_user_mind: full
`)

	suite.repo = powerrepository.NewPowerRepository()
}

func (suite *HealingEffectLoadedFromData) TestLoadFromJSON(checker *C) {
	success, err := suite.repo.AddJSONSource(suite.healStaffJSON)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)

	healStaff := suite.repo.GetPowerByID("power_heal")
	checker.Assert(healStaff.HitPointsHealed(), Equals, 2)
}

func (suite *HealingEffectLoadedFromData) TestLoadFromYAML(checker *C) {
	success, err := suite.repo.AddYAMLSource(suite.healStaffYAML)
	checker.Assert(err, IsNil)
	checker.Assert(success, Equals, true)

	healStaff := suite.repo.GetPowerByID("power_heal")
	checker.Assert(healStaff.HitPointsHealed(), Equals, 2)
}
