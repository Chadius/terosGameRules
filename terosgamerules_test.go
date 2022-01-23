package terosgamerules_test

import (
	"bytes"
	"github.com/chadius/terosgamerules"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func useValidSquaddieData() *bytes.Buffer {
	squaddieData := []byte(`
-
  name: Teros
  id: squaddieTeros
  affiliation: player
  aim: 2
  strength: 1
  max_hit_points: 5
  max_barrier: 3
  armor: 2
  dodge: 3
  deflect: 4
  powers:
    -
      name: Spear
      id: powerSpear
-
  name: Bandit
  id: squaddieBandit0
  affiliation: enemy
  aim: 0
  strength: 1
  max_hit_points: 5
  max_barrier: 0
  armor: 0
  dodge: 0
  deflect: 0
  powers:
    -
      name: Axe
      id: powerAxe
`)
	return bytes.NewBuffer(squaddieData)
}

func useValidPowerData() *bytes.Buffer {
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
	return bytes.NewBuffer(powerData)
}

func useValidScriptData() *bytes.Buffer {
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
	return bytes.NewBuffer(scriptData)
}

func TestReplayScriptExpectedOutput(t *testing.T) {
	suite.Run(t, new(ReplayScriptExpectedOutput))
}

type ReplayScriptExpectedOutput struct {
	suite.Suite
}

func (suite *ReplayScriptExpectedOutput) TestWhenCustomerUsesValidDataIsUsed_GetExpectedResponse() {
	// Setup
	var output strings.Builder
	gameRunner := terosgamerules.GameRules{}

	// Run
	err := gameRunner.ReplayBattleScript(
		useValidScriptData(),
		useValidSquaddieData(),
		useValidPowerData(),
		&output,
	)

	// Require
	require := require.New(suite.T())
	require.Nil(err, "no errors should have been found")

	expectedOutput := "Teros (Spear) vs Bandit: +2 (30/36), for 3 damage\n crit: 3/36, FATAL\nBandit (Axe) counters Teros: -5 (1/36) for NO DAMAGE + 2 barrier burn\nTeros (Spear) hits Bandit, for 3 damage\n   Bandit: 2/5 HP\nBandit (Axe) misses Teros\n   Teros: 5/5 HP, 3 barrier\n---\nBandit (Axe) vs Teros: -3 (6/36) for NO DAMAGE + 2 barrier burn\nTeros (Spear) counters Bandit: +0 (21/36), FATAL\nBandit (Axe) misses Teros\n   Teros: 5/5 HP, 3 barrier\nTeros (Spear) counters Bandit, felling\n   Bandit: 0/5 HP\n---\n"
	require.Equal(expectedOutput, output.String())
}

func TestReplayScriptErrorsSuite(t *testing.T) {
	suite.Run(t, new(ReplayScriptErrorsSuite))
}

type ReplayScriptErrorsSuite struct {
	suite.Suite
	byteOutput strings.Builder
	gameRunner terosgamerules.GameRules
}

func (suite *ReplayScriptErrorsSuite) TestWhenAllDataExists_ThenNoErrors() {
	// Setup
	squaddieDataBuffer := useValidSquaddieData()
	powerDataBuffer := useValidPowerData()
	scriptDataBuffer := useValidScriptData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		squaddieDataBuffer,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Nil(err)
}

func (suite *ReplayScriptErrorsSuite) TestWhenSquaddieDataIsMissing_ThenSquaddieDataErrors() {
	powerDataBuffer := useValidPowerData()
	scriptDataBuffer := useValidScriptData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		nil,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report squaddie reader error")
	require.Containsf(err.Error(), "no squaddie data found", "Error message does not match.")
}

func (suite *ReplayScriptErrorsSuite) TestWhenSquaddieDataIsInvalid_ThenReportNoSquaddieInformation() {
	squaddieData := []byte(`Not Valid JSON/YAML`)
	squaddieDataBuffer := bytes.NewBuffer(squaddieData)
	powerDataBuffer := useValidPowerData()
	scriptDataBuffer := useValidScriptData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		squaddieDataBuffer,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report squaddie data error")
	require.Containsf(err.Error(), "squaddie data is invalid", "Error message does not match.")
}

func (suite *ReplayScriptErrorsSuite) TestWhenPowerDataIsMissing_ThenPowerDataErrors() {
	squaddieDataBuffer := useValidSquaddieData()
	scriptDataBuffer := useValidScriptData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		squaddieDataBuffer,
		nil,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report power reader error")
	require.Containsf(err.Error(), "no power data found", "Error message does not match.")
}

func (suite *ReplayScriptErrorsSuite) TestWhenPowerDataIsInvalid_ThenReportNoPowerInformation() {
	powerData := []byte(`Not Valid JSON/YAML`)
	powerDataBuffer := bytes.NewBuffer(powerData)
	squaddieDataBuffer := useValidSquaddieData()
	scriptDataBuffer := useValidScriptData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		squaddieDataBuffer,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report power data error")
	require.Containsf(err.Error(), "power data is invalid", "Error message does not match.")
}

func (suite *ReplayScriptErrorsSuite) TestWhenScriptDataIsMissing_ThenScriptDataErrors() {
	squaddieDataBuffer := useValidSquaddieData()
	powerDataBuffer := useValidPowerData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		nil,
		squaddieDataBuffer,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report script reader error")
	require.Containsf(err.Error(), "no script data found", "Error message does not match.")
}

func (suite *ReplayScriptErrorsSuite) TestWhenScriptDataIsInvalid_ThenReportNoScriptInformation() {
	scriptData := []byte(`Not Valid JSON/YAML`)
	scriptDataBuffer := bytes.NewBuffer(scriptData)
	squaddieDataBuffer := useValidSquaddieData()
	powerDataBuffer := useValidPowerData()

	// Run
	err := suite.gameRunner.ReplayBattleScript(
		scriptDataBuffer,
		squaddieDataBuffer,
		powerDataBuffer,
		&suite.byteOutput,
	)

	// Require
	require := require.New(suite.T())
	require.Error(err, "Did not report script data error")
	require.Containsf(err.Error(), "script data is invalid", "Error message does not match.")
}
