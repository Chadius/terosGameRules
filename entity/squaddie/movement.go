package squaddie

import "github.com/chadius/terosgamerules/entity/movement"

// Movement contains all information needed to describe a Squaddie's movement.
type Movement struct {
	movementDistance     int
	movementLogic        movement.Interface
	movementCanHitAndRun bool
}

// NewMovement creates a new Movement object.
func NewMovement(distance int, canHitAndRun bool, movementLogic movement.Interface) *Movement {
	return &Movement{
		movementDistance:     distance,
		movementCanHitAndRun: canHitAndRun,
		movementLogic:        movementLogic,
	}
}

// MovementDistance Returns the distance the Squaddie can travel.
func (movement *Movement) MovementDistance() int {
	return movement.movementDistance
}

// MovementLogic returns the Squaddie's movement logic block
func (movement *Movement) MovementLogic() movement.Interface {
	return movement.movementLogic
}

// CanHitAndRun indicates if the Squaddie can move after attacking.
func (movement *Movement) CanHitAndRun() bool {
	return movement.movementCanHitAndRun
}

// Improve improves the movement stats.
func (movement *Movement) Improve(distance int, canHitAndRun bool, movementLogic movement.Interface) {
	movement.movementDistance += distance

	if canHitAndRun {
		movement.movementCanHitAndRun = true
	}

	if movementLogic.GreaterThan(movement.movementLogic) {
		movement.movementLogic = movementLogic
	}
}
