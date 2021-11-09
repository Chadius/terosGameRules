package powercommit

import "github.com/chadius/terosbattleserver/entity/damagedistribution"

// AttackResult shows what happens when the power was an attack.
type AttackResult struct {
	attackRoll           int
	defendRoll           int
	attackerToHitBonus   int
	defenderToHitPenalty int
	attackerTotal        int
	defenderTotal        int
	hitTarget            bool
	criticallyHitTarget  bool
	damage               *damagedistribution.DamageDistribution
	isCounterAttack      bool
}

// NewAttackResult generates a new object.
func NewAttackResult(
	attackRoll int,
	defendRoll int,
	attackerToHitBonus int,
	defenderToHitPenalty int,
	attackerTotal int,
	defenderTotal int,
	hitTarget bool,
	criticallyHitTarget bool,
	damage *damagedistribution.DamageDistribution,
	isCounterAttack bool,
) *AttackResult {
	return &AttackResult{
		attackRoll:           attackRoll,
		defendRoll:           defendRoll,
		attackerToHitBonus:   attackerToHitBonus,
		defenderToHitPenalty: defenderToHitPenalty,
		attackerTotal:        attackerTotal,
		defenderTotal:        defenderTotal,
		hitTarget:            hitTarget,
		criticallyHitTarget:  criticallyHitTarget,
		damage:               damage,
		isCounterAttack:      isCounterAttack,
	}
}

// AttackRoll is a getter.
func (a *AttackResult) AttackRoll() int {
	return a.attackRoll
}

// DefendRoll is a getter.
func (a *AttackResult) DefendRoll() int {
	return a.defendRoll
}

// AttackerToHitBonus is a getter.
func (a *AttackResult) AttackerToHitBonus() int {
	return a.attackerToHitBonus
}

// DefenderToHitPenalty is a getter.
func (a *AttackResult) DefenderToHitPenalty() int {
	return a.defenderToHitPenalty
}

// AttackerTotal is a getter.
func (a *AttackResult) AttackerTotal() int {
	return a.attackerTotal
}

// DefenderTotal is a getter.
func (a *AttackResult) DefenderTotal() int {
	return a.defenderTotal
}

// HitTarget is a getter.
func (a *AttackResult) HitTarget() bool {
	return a.hitTarget
}

// CriticallyHitTarget is a getter.
func (a *AttackResult) CriticallyHitTarget() bool {
	return a.criticallyHitTarget
}

// Damage is a getter.
func (a *AttackResult) Damage() *damagedistribution.DamageDistribution {
	return a.damage
}

// IsCounterAttack is a getter.
func (a *AttackResult) IsCounterAttack() bool {
	return a.isCounterAttack
}

// AttackResultBuilder saves instructions so you can create an AttackResult.
type AttackResultBuilder struct {
	attackRoll           int
	defendRoll           int
	attackerToHitBonus   int
	defenderToHitPenalty int
	attackerTotal        int
	defenderTotal        int
	hitTarget            bool
	criticallyHitTarget  bool
	damage               *damagedistribution.DamageDistribution
	isCounterAttack      bool
}

// NewAttackResultBuilder returns a new builder object.
func NewAttackResultBuilder() *AttackResultBuilder {
	return &AttackResultBuilder{
		0,
		0,
		0,
		0,
		0,
		0,
		false,
		false,
		&damagedistribution.DamageDistribution{},
		false,
	}
}

// DamageDistribution assumes the attack hit.
func (ar *AttackResultBuilder) DamageDistribution(damage *damagedistribution.DamageDistribution) *AttackResultBuilder {
	ar.HitTarget()
	ar.damage = damage
	return ar
}

// HitTarget means the attack hit.
func (ar *AttackResultBuilder) HitTarget() *AttackResultBuilder {
	ar.hitTarget = true
	return ar
}

// CriticallyHit sets the attack to a critical hit.
func (ar *AttackResultBuilder) CriticallyHit() *AttackResultBuilder {
	ar.criticallyHitTarget = true
	return ar
}

// CounterAttack marks this attack as a counter.
func (ar *AttackResultBuilder) CounterAttack() *AttackResultBuilder {
	ar.isCounterAttack = true
	return ar
}

// AttackRoll sets the field.
func (ar *AttackResultBuilder) AttackRoll(attackRoll           int) *AttackResultBuilder {
	ar.attackRoll = attackRoll
	return ar
}

// DefendRoll sets the field.
func (ar *AttackResultBuilder) DefendRoll(defendRoll           int) *AttackResultBuilder {
	ar.defendRoll = defendRoll
	return ar
}

// AttackerToHitBonus sets the field.
func (ar *AttackResultBuilder) AttackerToHitBonus(attackerToHitBonus   int) *AttackResultBuilder {
	ar.attackerToHitBonus = attackerToHitBonus
	return ar
}

// DefenderToHitPenalty sets the field.
func (ar *AttackResultBuilder) DefenderToHitPenalty(defenderToHitPenalty int) *AttackResultBuilder {
	ar.defenderToHitPenalty = defenderToHitPenalty
	return ar
}

// AttackerTotal sets the field.
func (ar *AttackResultBuilder) AttackerTotal(attackerTotal        int) *AttackResultBuilder {
	ar.attackerTotal = attackerTotal
	return ar
}

// DefenderTotal sets the field.
func (ar *AttackResultBuilder) DefenderTotal(defenderTotal        int) *AttackResultBuilder {
	ar.defenderTotal = defenderTotal
	return ar
}

// Build constructs an AttackResult.
func (ar *AttackResultBuilder) Build() *AttackResult {
	return NewAttackResult(
		ar.attackRoll,
		ar.defendRoll,
		ar.attackerToHitBonus,
		ar.defenderToHitPenalty,
		ar.attackerTotal,
		ar.defenderTotal,
		ar.hitTarget,
		ar.criticallyHitTarget,
		ar.damage,
		ar.isCounterAttack,
	)
}