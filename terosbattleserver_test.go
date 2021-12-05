package terosbattleserver_test

import (
	"bytes"
	"github.com/chadius/terosbattleserver"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ReplayScriptSuite struct{}

var _ = Suite(&ReplayScriptSuite{})

func (suite *ReplayScriptSuite) Test(checker *C) {
	squaddieData := []byte(`
-
  identification:
    name: Teros
    id: squaddieTeros
    affiliation: player
  offense:
    aim: 2
    strength: 1
  defense:
    max_hit_points: 5
    max_barrier: 3
    armor: 2
    dodge: 3
    deflect: 4
  powers:
    powers:
    -
      name: Spear
      id: powerSpear
-
  identification:
    name: Bandit
    id: squaddieBandit0
    affiliation: enemy
  offense:
    aim: 0
    strength: 1
  defense:
    max_hit_points: 5
    max_barrier: 0
    armor: 0
    dodge: 0
    deflect: 0
  powers:
    powers:
    -
      name: Axe
      id: powerAxe
`)
	squaddieByteStream := bytes.NewBuffer(squaddieData)

	powerData := []byte(`
-
  name: Spear
  id: powerSpear
  power_type: physical
  target_foe: true
  can_attack: true
  damage_bonus: 2
  can_be_equipped: true
  can_counter_attack: true
  counter_attack_penalty_reduction: 0
  can_critical: true
  critical_damage: 2
-
  name: Axe
  id: powerAxe
  power_type: physical
  target_foe: true
  can_attack: true
  damage_bonus: 1
  can_be_equipped: true
  can_counter_attack: true
`)
	powerByteStream := bytes.NewBuffer(powerData)

	scriptData := []byte(`---
version: 0.1F
actions:
  -
    random_seed: 1000
    user_id: squaddieTeros
    power_id: powerSpear
    target_ids:
      - squaddieBandit0
  -
    random_seed: 2
    user_id: squaddieBandit0
    power_id: powerAxe
    target_ids:
      - squaddieTeros
`)
	scriptByteStream := bytes.NewBuffer(scriptData)

	var output strings.Builder
	err := terosbattleserver.ReplayBattleScript(scriptByteStream, squaddieByteStream, powerByteStream, &output)

	checker.Assert(err, IsNil)

	expectedOutput := "Teros (Spear) vs Bandit: +2 (30/36), for 3 damage\n crit: 3/36, FATAL\nBandit (Axe) counters Teros: -5 (1/36) for NO DAMAGE + 2 barrier burn\nTeros (Spear) hits Bandit, for 3 damage\n   Bandit: 2/5 HP\nBandit (Axe) misses Teros\n   Teros: 5/5 HP, 3 barrier\n---\nBandit (Axe) vs Teros: -3 (6/36) for NO DAMAGE + 2 barrier burn\nTeros (Spear) counters Bandit: +0 (21/36), FATAL\nBandit (Axe) misses Teros\n   Teros: 5/5 HP, 3 barrier\nTeros (Spear) counters Bandit, felling\n   Bandit: 0/5 HP\n---\n"
	checker.Assert(expectedOutput, Equals, output.String())
}
