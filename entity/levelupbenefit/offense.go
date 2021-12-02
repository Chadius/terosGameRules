package levelupbenefit

// Offense describes the offensive benefits of leveling up.
type Offense struct {
	aim      int
	strength int
	mind     int
}

// NewOffense returns a new Offense object.
func NewOffense(aim, strength, mind int) *Offense {
	return &Offense{
		aim:      aim,
		strength: strength,
		mind:     mind,
	}
}

// Aim is a getter.
func (o *Offense) Aim() int {
	return o.aim
}

// Strength is a getter.
func (o *Offense) Strength() int {
	return o.strength
}

// Mind is a getter.
func (o *Offense) Mind() int {
	return o.mind
}
