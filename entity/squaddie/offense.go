package squaddie

// Offense helps the Squaddie fight better and be more effective.
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

// Aim returns the value.
func (o *Offense) Aim() int {
	return o.aim
}

// Strength returns the value.
func (o *Offense) Strength() int {
	return o.strength
}

// Mind returns the value.
func (o *Offense) Mind() int {
	return o.mind
}

// Improve improves the offensive stats.
func (o *Offense) Improve(aim, strength, mind int) {
	o.aim += aim
	o.strength += strength
	o.mind += mind
}
