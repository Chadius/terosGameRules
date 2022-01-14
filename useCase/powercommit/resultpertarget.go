package powercommit

import (
	"github.com/chadius/terosgamerules/entity/powerinterface"
	"github.com/chadius/terosgamerules/entity/squaddieinterface"
)

// ResultPerTarget shows what happened to each target.
type ResultPerTarget struct {
	userID   string
	powerID  string
	targetID string
	attack   *AttackResult
	healing  *HealResult
}

// UserID is a getter.
func (r *ResultPerTarget) UserID() string {
	return r.userID
}

// PowerID is a getter.
func (r *ResultPerTarget) PowerID() string {
	return r.powerID
}

// TargetID is a getter.
func (r *ResultPerTarget) TargetID() string {
	return r.targetID
}

// Attack is a getter.
func (r *ResultPerTarget) Attack() *AttackResult {
	return r.attack
}

// Healing is a getter.
func (r *ResultPerTarget) Healing() *HealResult {
	return r.healing
}

// ResultPerTargetBuilder can be used to create ResultPerTarget objects
type ResultPerTargetBuilder struct {
	userID   string
	powerID  string
	targetID string
	attack   *AttackResult
	healing  *HealResult
}

// NewResultPerTargetBuilder creates a new ResultPerTargetBuilder.
func NewResultPerTargetBuilder() *ResultPerTargetBuilder {
	return &ResultPerTargetBuilder{
		userID:   "",
		powerID:  "",
		targetID: "",
		attack:   nil,
		healing:  nil,
	}
}

// User sets the field.
func (rp *ResultPerTargetBuilder) User(user squaddieinterface.Interface) *ResultPerTargetBuilder {
	rp.userID = user.ID()
	return rp
}

// Power sets the field.
func (rp *ResultPerTargetBuilder) Power(pow powerinterface.Interface) *ResultPerTargetBuilder {
	rp.powerID = pow.ID()
	return rp
}

// Target sets the field.
func (rp *ResultPerTargetBuilder) Target(target squaddieinterface.Interface) *ResultPerTargetBuilder {
	rp.targetID = target.ID()
	return rp
}

// AttackResult sets the field.
func (rp *ResultPerTargetBuilder) AttackResult(result *AttackResult) *ResultPerTargetBuilder {
	rp.attack = result
	return rp
}

// HealResult sets the field.
func (rp *ResultPerTargetBuilder) HealResult(result *HealResult) *ResultPerTargetBuilder {
	rp.healing = result
	return rp
}

// Build creates the ResultPerTarget.
func (rp *ResultPerTargetBuilder) Build() *ResultPerTarget {
	return &ResultPerTarget{
		userID:   rp.userID,
		powerID:  rp.powerID,
		targetID: rp.targetID,
		attack:   rp.attack,
		healing:  rp.healing,
	}
}
