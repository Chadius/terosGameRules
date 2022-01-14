package replay

import (
	"github.com/chadius/terosgamerules/utility"
	"gopkg.in/yaml.v2"
)

// SquaddieAction records everything a squaddie could have performed in a single turn.
type SquaddieAction struct {
	RandomSeed int64    `json:"random_seed" yaml:"random_seed"`
	UserID     string   `json:"user_id" yaml:"user_id"`
	PowerID    string   `json:"power_id" yaml:"power_id"`
	TargetIDs  []string `json:"target_ids" yaml:"target_ids"`
}

// ChapterReplay contains the information needed to recreate a replay of one chapter in a game.
type ChapterReplay struct {
	Version string            `json:"version" yaml:"version"`
	Actions []*SquaddieAction `json:"actions" yaml:"actions"`
}

// NewCreateMapReplayFromYAML reads the YAML data and returns a list of Map objects.
func NewCreateMapReplayFromYAML(data []byte) (*ChapterReplay, error) {
	return newCreateMapReplayFromDatastream(data, yaml.Unmarshal)
}

// newCreateMapReplayFromDatastream consumes a given bytestream and tries to create multiple objects from it
func newCreateMapReplayFromDatastream(data []byte, unmarshal utility.UnmarshalFunc) (*ChapterReplay, error) {
	var unmarshalError error
	var chapterReplay ChapterReplay
	unmarshalError = unmarshal(data, &chapterReplay)

	if unmarshalError != nil {
		return nil, unmarshalError
	}

	return &chapterReplay, nil
}
