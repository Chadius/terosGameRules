package squaddie

// Movement contains all of the information needed to describe a Squaddie's movement.
type Movement struct {
	Distance  int          `json:"distance" yaml:"distance"`
	Type      MovementType `json:"type" yaml:"type"`
	HitAndRun bool         `json:"hit_and_run" yaml:"hit_and_run"`
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

// GetMovementDistancePerRound Returns the distance the Squaddie can travel.
func (movement *Movement) GetMovementDistancePerRound() int {
	return movement.Distance
}

// GetMovementType returns the Squaddie's movement type
func (movement *Movement) GetMovementType() MovementType {
	return movement.Type
}

// CanHitAndRun indicates if the Squaddie can move after attacking.
func (movement *Movement) CanHitAndRun() bool {
	return movement.HitAndRun
}
