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
	Foot = "foot"
	// Light Squaddies have no terrain penalties and cannot cross pits.
	Light = "light"
	// Fly Squaddies have no terrain penalties and can cross pits.
	Fly = "fly"
	// Teleport Squaddies have no terrain penalties and can cross pits.
	//   They also ignore walls and other barriers.
	Teleport = "teleport"
)

var MovementValueByType = map[MovementType]int{
	Foot:     0,
	Light:    1,
	Fly:      2,
	Teleport: 3,
}