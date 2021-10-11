package levelupbenefit

import "github.com/cserrant/terosbattleserver/entity/power"

// PowerChanges tracks changes to the squaddie's innate powers.
type PowerChanges struct {
	Gained []*power.Reference `json:"gained" yaml:"gained"`
	Lost   []*power.Reference `json:"lost" yaml:"lost"`
}
