package levelupbenefit

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/powerreference"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// Builder is used to create formula objects.
type Builder struct {
	levelID   string
	classID   string
	levelSize Size

	hitPoints int
	deflect   int
	dodge     int
	barrier   int
	armor     int

	aim      int
	strength int
	mind     int

	movementDistance     int
	movementLogic        movement.Interface
	movementCanHitAndRun bool

	powersGained []*powerreference.Reference
	powersLost   []*powerreference.Reference
}

// NewLevelUpBenefitBuilder returns a new object used to build Term objects.
func NewLevelUpBenefitBuilder() *Builder {
	return &Builder{
		levelID:   "missing ID",
		classID:   "missing class ID",
		levelSize: Small,

		hitPoints: 0,
		deflect:   0,
		dodge:     0,
		barrier:   0,
		armor:     0,

		aim:      0,
		strength: 0,
		mind:     0,

		movementDistance:     0,
		movementLogic:        movement.NewMovementLogic("foot"),
		movementCanHitAndRun: false,

		powersGained: []*powerreference.Reference{},
		powersLost:   []*powerreference.Reference{},
	}
}

// WithID sets the field.
func (b *Builder) WithID(levelID string) *Builder {
	b.levelID = levelID
	return b
}

// LevelID is an alias.
func (b *Builder) LevelID(levelID string) *Builder {
	return b.WithID(levelID)
}

// ID is an alias.
func (b *Builder) ID(levelID string) *Builder {
	return b.WithID(levelID)
}

// WithClassID sets the field.
func (b *Builder) WithClassID(classID string) *Builder {
	b.classID = classID
	return b
}

// ClassID is an alias
func (b *Builder) ClassID(classID string) *Builder {
	return b.WithClassID(classID)
}

// BigLevel means this level is a Big major milestone.
func (b *Builder) BigLevel() *Builder {
	b.levelSize = Big
	return b
}

// HitPoints increases the defensive parameter
func (b *Builder) HitPoints(hitPoints int) *Builder {
	b.hitPoints = hitPoints
	return b
}

// Deflect increases the defensive parameter
func (b *Builder) Deflect(deflect int) *Builder {
	b.deflect = deflect
	return b
}

// Dodge increases the defensive parameter
func (b *Builder) Dodge(dodge int) *Builder {
	b.dodge = dodge
	return b
}

// Barrier increases the defensive parameter
func (b *Builder) Barrier(barrier int) *Builder {
	b.barrier = barrier
	return b
}

// Armor increases the defensive parameter
func (b *Builder) Armor(armor int) *Builder {
	b.armor = armor
	return b
}

// Aim increases the offensive parameter.
func (b *Builder) Aim(aim int) *Builder {
	b.aim = aim
	return b
}

// Strength increases the offensive parameter.
func (b *Builder) Strength(strength int) *Builder {
	b.strength = strength
	return b
}

// Mind increases the offensive parameter.
func (b *Builder) Mind(mind int) *Builder {
	b.mind = mind
	return b
}

// MovementDistance increases the squaddie's movement.
func (b *Builder) MovementDistance(distance int) *Builder {
	b.movementDistance = distance
	return b
}

// MovementLogic will upgrade the squaddie movement logic using the given keyword
func (b *Builder) MovementLogic(newMovementType string) *Builder {
	b.movementLogic = movement.NewMovementLogic(newMovementType)
	return b
}

// FootMovement will change the squaddie's movement type.
func (b *Builder) FootMovement() *Builder {
	return b.MovementLogic("foot")
}

// LightMovement will change the squaddie's movement type.
func (b *Builder) LightMovement() *Builder {
	return b.MovementLogic("light")
}

// FlyMovement will change the squaddie's movement type.
func (b *Builder) FlyMovement() *Builder {
	return b.MovementLogic("fly")
}

// TeleportMovement will change the squaddie's movement type.
func (b *Builder) TeleportMovement() *Builder {
	return b.MovementLogic("teleport")
}

// CanHitAndRun means the squaddie can use hit-and-run movement.
func (b *Builder) CanHitAndRun() *Builder {
	b.movementCanHitAndRun = true
	return b
}

// GainPower will change powers.
func (b *Builder) GainPower(powerID, powerName string) *Builder {
	b.powersGained = append(
		b.powersGained,
		&powerreference.Reference{
			Name:    powerName,
			PowerID: powerID,
		},
	)
	return b
}

// LosePower will change powers.
func (b *Builder) LosePower(powerID string) *Builder {
	b.powersLost = append(
		b.powersLost,
		&powerreference.Reference{
			Name:    "power name does not matter",
			PowerID: powerID,
		},
	)
	return b
}

// Build creates a new LevelUpBenefit object.
func (b *Builder) Build() (*LevelUpBenefit, error) {
	return NewLevelUpBenefit(
		NewIdentification(
			b.levelID,
			b.classID,
			b.levelSize,
		),
		NewDefense(
			b.hitPoints,
			b.dodge,
			b.deflect,
			b.barrier,
			b.armor,
		),
		NewOffense(
			b.aim,
			b.strength,
			b.mind,
		),
		squaddie.NewMovement(
			b.movementDistance,
			b.movementCanHitAndRun,
			b.movementLogic,
		),
		NewPowerChanges(
			b.powersGained,
			b.powersLost,
		),
	), nil
}

// NewLevelUpBenefitBuilderFromYAML creates a new Builder using YAML.
func NewLevelUpBenefitBuilderFromYAML(data []byte) *Builder {
	builder := NewLevelUpBenefitBuilder()
	builder.unmarshalAndApplyDataStream(data, yaml.Unmarshal)
	return builder
}

// NewLevelUpBenefitBuilderFromJSON creates a new Builder using JSON.
func NewLevelUpBenefitBuilderFromJSON(data []byte) *Builder {
	builder := NewLevelUpBenefitBuilder()
	builder.unmarshalAndApplyDataStream(data, json.Unmarshal)
	return builder
}

// BuilderMarshal is used to convert to and from data streams.
type BuilderMarshal struct {
	LevelID  string `json:"id" yaml:"id"`
	ClassID  string `json:"class_id" yaml:"class_id"`
	BigLevel bool   `json:"is_a_big_level" yaml:"is_a_big_level"`

	HitPoints int `json:"hit_points" yaml:"hit_points"`
	Deflect   int `json:"deflect" yaml:"deflect"`
	Dodge     int `json:"dodge" yaml:"dodge"`
	Barrier   int `json:"barrier" yaml:"barrier"`
	Armor     int `json:"armor" yaml:"armor"`

	Aim      int `json:"aim" yaml:"aim"`
	Strength int `json:"strength" yaml:"strength"`
	Mind     int `json:"mind" yaml:"mind"`

	MovementDistance     int    `json:"movement_distance" yaml:"movement_distance"`
	MovementLogic        string `json:"movement_type" yaml:"movement_type"`
	MovementCanHitAndRun bool   `json:"can_hit_and_run" yaml:"can_hit_and_run"`

	PowersGained []*powerreference.Reference `json:"powers_gained" yaml:"powers_gained"`
	PowersLost   []string                    `json:"powers_lost" yaml:"powers_lost"`
}

// unmarshalAndApplyDataStream consumes a given bytestream of the given sourceType and tries to analyze it.
func (b *Builder) unmarshalAndApplyDataStream(data []byte, unmarshal utility.UnmarshalFunc) (bool, error) {
	var unmarshalError error

	var builderFields BuilderMarshal
	unmarshalError = unmarshal(data, &builderFields)

	if unmarshalError != nil {
		return false, unmarshalError
	}

	b.populateBuilderBasedOnMarshal(builderFields)
	return true, nil
}

func (b *Builder) populateBuilderBasedOnMarshal(builderFields BuilderMarshal) *Builder {
	b.LevelID(builderFields.LevelID).
		ClassID(builderFields.ClassID).
		HitPoints(builderFields.HitPoints).
		Dodge(builderFields.Dodge).
		Deflect(builderFields.Deflect).
		Barrier(builderFields.Barrier).
		Armor(builderFields.Armor).
		Aim(builderFields.Aim).
		Strength(builderFields.Strength).
		Mind(builderFields.Mind).
		MovementDistance(builderFields.MovementDistance).
		MovementLogic(builderFields.MovementLogic)

	if builderFields.MovementCanHitAndRun {
		b.CanHitAndRun()
	}

	if builderFields.BigLevel {
		b.BigLevel()
	}

	for _, reference := range builderFields.PowersGained {
		b.GainPower(reference.PowerID, reference.Name)
	}

	for _, powerID := range builderFields.PowersLost {
		b.LosePower(powerID)
	}

	return b
}

// NewBuilderFromMarshal creates a new Builder with fields based on the Marshal object
func NewBuilderFromMarshal(builderFields BuilderMarshal) *Builder {
	b := NewLevelUpBenefitBuilder().populateBuilderBasedOnMarshal(builderFields)
	return b
}
