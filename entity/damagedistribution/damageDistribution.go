package damagedistribution

// DamageDistribution tracks how a Squaddie handles damage across their Barrier, Armor and Health.
type DamageDistribution struct {
	DamageAbsorbedByArmor   int
	DamageAbsorbedByBarrier int
	DamageDealt             int
	ExtraBarrierBurnt       int
	TotalBarrierBurnt       int
}
