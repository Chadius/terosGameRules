package report

// PowerReport shows the results of using a power.
type PowerReport struct {
	AttackerID           string
	PowerID              string
	AttackingPowerResults []*AttackingPowerReport
}

// AttackingPowerReport shows the results of using a power with an AttackEffect.
type AttackingPowerReport struct {
	TargetID		string
	DamageTaken		int
	BarrierDamage	int
	WasAHit			bool
	WasACriticalHit	bool
}