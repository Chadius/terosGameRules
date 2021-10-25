package squaddie

// Offense helps the Squaddie fight better and be more effective.
type Offense struct {
	SquaddieAim      int `json:"aim" yaml:"aim"`
	SquaddieStrength int `json:"strength" yaml:"strength"`
	SquaddieMind     int `json:"mind" yaml:"mind"`
}

// Aim returns the value.
func (o *Offense) Aim() int {
	return o.SquaddieAim
}

// Strength returns the value.
func (o *Offense) Strength() int {
	return o.SquaddieStrength
}

// Mind returns the value.
func (o *Offense) Mind() int {
	return o.SquaddieMind
}
