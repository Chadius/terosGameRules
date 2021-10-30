package power

// CriticalEffect records the various extras that affect the target once the power crits.
type CriticalEffect struct {
	CriticalCriticalHitThresholdBonus int `json:"critical_hit_threshold_bonus" yaml:"critical_hit_threshold_bonus"`
	Damage                            int `json:"damage" yaml:"damage"`
}

// NewCriticalEffect returns a new CriticalEffect.
func NewCriticalEffect (criticalHitThresholdBonus, damage int) *CriticalEffect {
	return &CriticalEffect{
		CriticalCriticalHitThresholdBonus: criticalHitThresholdBonus,
		Damage:                            damage,
	}
}

// CriticalHitThresholdInitialValue is the default value the attack must exceed the defense by.
const CriticalHitThresholdInitialValue = 6

// CriticalHitThreshold returns how far the attacker must exceed the defender in order to crit.
func (criticalEffect *CriticalEffect) CriticalHitThreshold() int {
	return CriticalHitThresholdInitialValue - criticalEffect.CriticalHitThresholdBonus()
}

//CriticalHitThresholdBonus returns the raw bonus.
func (criticalEffect *CriticalEffect) CriticalHitThresholdBonus() int {
	return criticalEffect.CriticalCriticalHitThresholdBonus
}

// ExtraCriticalHitDamage returns the extra Damage dealt upon a critical hit.
func (criticalEffect *CriticalEffect) ExtraCriticalHitDamage() int {
	return criticalEffect.Damage
}
