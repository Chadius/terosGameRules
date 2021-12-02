package levelupbenefit

// Defense describes the defensive benefits of leveling up.
type Defense struct {
	maxHitPoints int
	dodge        int
	deflect      int
	maxBarrier   int
	armor        int
}

// NewDefense returns a new Defense object.
func NewDefense(maxHitPoints, dodge, deflect, maxBarrier, armor int) *Defense {
	return &Defense{
		maxHitPoints: maxHitPoints,
		dodge:        dodge,
		deflect:      deflect,
		maxBarrier:   maxBarrier,
		armor:        armor,
	}
}

// MaxHitPoints is a getter.
func (d *Defense) MaxHitPoints() int {
	return d.maxHitPoints
}

// Dodge is a getter.
func (d *Defense) Dodge() int {
	return d.dodge
}

// Deflect is a getter.
func (d *Defense) Deflect() int {
	return d.deflect
}

// MaxBarrier is a getter.
func (d *Defense) MaxBarrier() int {
	return d.maxBarrier
}

// Armor is a getter.
func (d *Defense) Armor() int {
	return d.armor
}
