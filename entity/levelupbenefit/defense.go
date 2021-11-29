package levelupbenefit

// Defense describes the defensive benefits of leveling up.
type Defense struct {
	PrivatizeMeMaxHitPoints int `json:"max_hit_points" yaml:"max_hit_points"`
	PrivatizeMeDodge   int `json:"dodge" yaml:"dodge"`
	PrivatizeMeDeflect    int `json:"deflect" yaml:"deflect"`
	PrivatizeMeMaxBarrier int `json:"max_barrier" yaml:"max_barrier"`
	PrivatizeMeArmor      int `json:"armor" yaml:"armor"`
}

// NewDefense returns a new Defense object.
func NewDefense (maxHitPoints, dodge, deflect, maxBarrier, armor int) *Defense {
	return &Defense{
		PrivatizeMeMaxHitPoints: maxHitPoints,
		PrivatizeMeDodge:        dodge,
		PrivatizeMeDeflect:      deflect,
		PrivatizeMeMaxBarrier:   maxBarrier,
		PrivatizeMeArmor:        armor,
	}
}

// MaxHitPoints is a getter.
func (d *Defense) MaxHitPoints() int {
	return d.PrivatizeMeMaxHitPoints
}

// Dodge is a getter.
func (d *Defense) Dodge() int {
	return d.PrivatizeMeDodge
}

// Deflect is a getter.
func (d *Defense) Deflect() int {
	return d.PrivatizeMeDeflect
}

// MaxBarrier is a getter.
func (d *Defense) MaxBarrier() int {
	return d.PrivatizeMeMaxBarrier
}

// Armor is a getter.
func (d *Defense) Armor() int {
	return d.PrivatizeMeArmor
}
