package levelupbenefit

// Offense describes the offensive benefits of leveling up.
type Offense struct {
	PrivatizeMeAim      int `json:"aim" yaml:"aim"`
	PrivatizeMeStrength int `json:"strength" yaml:"strength"`
	PrivatizeMeMind     int `json:"mind" yaml:"mind"`
}

// NewOffense returns a new Offense object.
func NewOffense(aim, strength, mind int) *Offense {
	return &Offense{
		PrivatizeMeAim:      aim,
		PrivatizeMeStrength: strength,
		PrivatizeMeMind:     mind,
	}
}

// Aim is a getter.
func (o *Offense) Aim() int {
	return o.PrivatizeMeAim
}

// Strength is a getter.
func (o *Offense) Strength() int {
	return o.PrivatizeMeStrength
}

// Mind is a getter.
func (o *Offense) Mind() int {
	return o.PrivatizeMeMind
}
