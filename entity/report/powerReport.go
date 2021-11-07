package report

// PowerReport shows what happened after committing to using the power.
type PowerReport struct {
	AttackerID            string
	PowerID               string
	AttackingPowerReports []*AttackingPowerReport
}

// AttackingPowerReport shows what happened after using a power with an attackEffect.
type AttackingPowerReport struct {
	AttackerID      string
	TargetID        string
	PowerID         string
	DamageTaken     int
	BarrierDamage   int
	WasAHit         bool
	WasACriticalHit bool
}
