package squaddie

import "github.com/cserrant/terosbattleserver/utility"

// Identification is akin to an ID card for each Squaddie. Each Squaddie carries a unique Identification.
type Identification struct {
	ID          string      `json:"id" yaml:"id"`
	Name        string      `json:"name" yaml:"name"`
	Affiliation Affiliation `json:"affiliation" yaml:"affiliation"`
}

// SetNewIDToRandom changes the ID to a random value.
func (identification *Identification) SetNewIDToRandom() {
	identification.ID = utility.StringWithCharset(8, "abcdefgh0123456789")
}
