package squaddie

// Movement contains all information needed to describe a Squaddie's movement.
type Movement struct {
	movementDistance     int
	movementType         MovementType
	movementCanHitAndRun bool
}

// NewMovement creates a new Movement object.
func NewMovement(distance int, movementType MovementType, canHitAndRun bool) *Movement {
	return &Movement{
		movementDistance:     distance,
		movementType:         movementType,
		movementCanHitAndRun: canHitAndRun,
	}
}

// MovementType describes how Squaddies traverse terrain. This affects
//   movement costs and crossing pits
type MovementType string

const (
	// Foot Squaddies take full terrain penalties and cannot cross pits.
	Foot MovementType = "foot"
	// Light Squaddies have no terrain penalties and cannot cross pits.
	Light MovementType = "light"
	// Fly Squaddies have no terrain penalties and can cross pits.
	Fly MovementType = "fly"
	// Teleport Squaddies have no terrain penalties and can cross pits.
	//   They also ignore walls and other barriers.
	Teleport MovementType = "teleport"
)

// MovementValueByType orders movement types by priority (highest number is most desired)
var MovementValueByType = map[MovementType]int{
	Foot:     0,
	Light:    1,
	Fly:      2,
	Teleport: 3,
}

// MovementDistance Returns the distance the Squaddie can travel.
func (movement *Movement) MovementDistance() int {
	return movement.movementDistance
}

// MovementType returns the Squaddie's movement type
func (movement *Movement) MovementType() MovementType {
	return movement.movementType
}

// CanHitAndRun indicates if the Squaddie can move after attacking.
func (movement *Movement) CanHitAndRun() bool {
	return movement.movementCanHitAndRun
}

// Improve improves the movement stats.
func (movement *Movement) Improve(distance int, moveType MovementType, canHitAndRun bool) {
	movement.movementDistance += distance

	if canHitAndRun {
		movement.movementCanHitAndRun = true
	}

	if moveType == Teleport {
		movement.movementType = Teleport
	} else if moveType == Fly {
		if movement.movementType != Teleport {
			movement.movementType = Fly
		}
	} else if moveType == Light {
		if movement.movementType != Teleport && movement.movementType != Fly {
			movement.movementType = Light
		}
	}
}
