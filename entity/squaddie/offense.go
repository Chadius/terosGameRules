package squaddie

// Offense helps the Squaddie fight better and be more effective.
type Offense struct {
	SquaddieAim      int `json:"aim" yaml:"aim"`
	SquaddieStrength int `json:"strength" yaml:"strength"`
	SquaddieMind     int `json:"mind" yaml:"mind"`
}

// NewOffense returns a new Offense object.
func NewOffense(aim, strength, mind int) *Offense {
	return &Offense{
		SquaddieAim:      aim,
		SquaddieStrength: strength,
		SquaddieMind:     mind,
	}
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

// Improve improves the offensive stats.
func (o *Offense) Improve(aim, strength, mind int) {
	o.SquaddieAim += aim
	o.SquaddieStrength += strength
	o.SquaddieMind += mind
}