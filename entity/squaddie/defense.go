package squaddie

// Defense holds everything needed to prevent the squaddie from getting hindered.
type Defense struct {
	CurrentHitPoints    		int                            `json:"current_hit_points" yaml:"current_hit_points"`
	MaxHitPoints        		int                            `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge               		int                               `json:"dodge" yaml:"dodge"`
	Deflect             		int                             `json:"deflect" yaml:"deflect"`
	CurrentBarrier      		int                              `json:"current_barrier" yaml:"current_barrier"`
	MaxBarrier          		int                              `json:"max_barrier" yaml:"max_barrier"`
	Armor               		int                               `json:"armor" yaml:"armor"`
}

// SetHPToMax restores the Squaddie's HitPoints.
func (defense *Defense) SetHPToMax() {
	defense.CurrentHitPoints = defense.MaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (defense *Defense) SetBarrierToMax() {
	defense.CurrentBarrier = defense.MaxBarrier
}
