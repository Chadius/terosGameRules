package report

// PowerReport shows the results of using a power.
type PowerReport struct {
	AttackerID           string
	TargetID             string
	PowerID              string
	AttackingPowerResults []*AttackingPowerReport
}

// AttackingPowerReport shows the results of using a power with an AttackEffect.
type AttackingPowerReport struct {
	DamageTaken     int
	BarrierDamage   int
	WasACriticalHit bool
}