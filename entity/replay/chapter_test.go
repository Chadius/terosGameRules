package replay_test

import (
	"github.com/chadius/terosbattleserver/entity/replay"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MapReplayTest struct{}

var _ = Suite(&MapReplayTest{})

func (suite *MapReplayTest) TestConsumeYAML(checker *C) {
	yamlByteStream := []byte(`---
version: 0.1F
actions:
  -
    random_seed: 1000
    user_id: squaddie_teros
    power_id: power_blot
    target_ids:
    - squaddie_bandit_0
    - squaddie_bandit_1
  -
    random_seed: 2000
    user_id: squaddie_bandit
    power_id: power_axe
    target_ids:
    - squaddie_teros
`)
	replayCommands, err := replay.NewCreateMapReplayFromYAML(yamlByteStream)
	checker.Assert(err, IsNil)
	checker.Assert(replayCommands.Version, Equals, "0.1F")
	checker.Assert(replayCommands.Actions, HasLen, 2)
	checker.Assert(replayCommands.Actions[0].RandomSeed, Equals, (int64)(1000))
	checker.Assert(replayCommands.Actions[0].UserID, Equals, "squaddie_teros")
	checker.Assert(replayCommands.Actions[0].PowerID, Equals, "power_blot")
	checker.Assert(replayCommands.Actions[0].TargetIDs, HasLen, 2)
	checker.Assert(replayCommands.Actions[0].TargetIDs[0], Equals, "squaddie_bandit_0")
	checker.Assert(replayCommands.Actions[0].TargetIDs[1], Equals, "squaddie_bandit_1")
}
