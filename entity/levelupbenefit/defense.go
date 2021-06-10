package levelupbenefit

// Defense describes all of the defensive benefits of leveling up.
type Defense struct {
	MaxHitPoints       int                `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge              int                `json:"dodge" yaml:"dodge"`
	Deflect            int                `json:"deflect" yaml:"deflect"`
	MaxBarrier         int                `json:"max_barrier" yaml:"max_barrier"`
	Armor              int                `json:"armor" yaml:"armor"`
}

