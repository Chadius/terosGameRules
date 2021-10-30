package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// BuilderOptions is used to create healing effects.
type BuilderOptions struct {
	identificationOptions *IdentificationBuilderOptions
	offenseOptions        *OffenseBuilderOptions
	defenseOptions        *DefenseBuilderOptions
	movementOptions       *MovementBuilderOptions
	powersToAdd           []*power.Power
	classesToAdd          []*squaddieclass.Class
	classToUse            *squaddieclass.Class
}

// Builder creates a BuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func Builder() *BuilderOptions {
	return &BuilderOptions{
		identificationOptions: IdentificationBuilder(),
		offenseOptions:        OffenseBuilder(),
		defenseOptions:        DefenseBuilder(),
		movementOptions:       MovementBuilder(),
		powersToAdd:           []*power.Power{},
		classesToAdd:          []*squaddieclass.Class{},
		classToUse:            nil,
	}
}

// WithName delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) WithName(name string) *BuilderOptions {
	s.identificationOptions.WithName(name)
	return s
}

// WithID delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) WithID(id string) *BuilderOptions {
	s.identificationOptions.WithID(id)
	return s
}

// AsPlayer delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) AsPlayer() *BuilderOptions {
	s.identificationOptions.AsPlayer()
	return s
}

// AsEnemy delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) AsEnemy() *BuilderOptions {
	s.identificationOptions.AsEnemy()
	return s
}

// AsAlly delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) AsAlly() *BuilderOptions {
	s.identificationOptions.AsAlly()
	return s
}

// AsNeutral delegates to the IdentificationBuilderOptions.
func (s *BuilderOptions) AsNeutral() *BuilderOptions {
	s.identificationOptions.AsNeutral()
	return s
}

// Aim delegates to the OffenseBuilderOptions.
func (s *BuilderOptions) Aim(aim int) *BuilderOptions {
	s.offenseOptions.Aim(aim)
	return s
}

// Strength delegates to the OffenseBuilderOptions.
func (s *BuilderOptions) Strength(strength int) *BuilderOptions {
	s.offenseOptions.Strength(strength)
	return s
}

// Mind delegates to the OffenseBuilderOptions.
func (s *BuilderOptions) Mind(mind int) *BuilderOptions {
	s.offenseOptions.Mind(mind)
	return s
}

// HitPoints delegates to the DefenseBuilderOptions.
func (s *BuilderOptions) HitPoints(maxHitPoints int) *BuilderOptions {
	s.defenseOptions.HitPoints(maxHitPoints)
	return s
}

// Barrier delegates to the DefenseBuilderOptions.
func (s *BuilderOptions) Barrier(maxBarrier int) *BuilderOptions {
	s.defenseOptions.Barrier(maxBarrier)
	return s
}

// Armor delegates to the DefenseBuilderOptions.
func (s *BuilderOptions) Armor(armor int) *BuilderOptions {
	s.defenseOptions.Armor(armor)
	return s
}

// Dodge delegates to the DefenseBuilderOptions.
func (s *BuilderOptions) Dodge(dodge int) *BuilderOptions {
	s.defenseOptions.Dodge(dodge)
	return s
}

// Deflect delegates to the DefenseBuilderOptions.
func (s *BuilderOptions) Deflect(deflect int) *BuilderOptions {
	s.defenseOptions.Deflect(deflect)
	return s
}

// MoveDistance delegates to the MovementBuilderOptions.
func (s *BuilderOptions) MoveDistance(distance int) *BuilderOptions {
	s.movementOptions.Distance(distance)
	return s
}

// CanHitAndRun delegates to the MovementBuilderOptions.
func (s *BuilderOptions) CanHitAndRun() *BuilderOptions {
	s.movementOptions.CanHitAndRun()
	return s
}

// MovementFoot delegates to the MovementBuilderOptions.
func (s *BuilderOptions) MovementFoot() *BuilderOptions {
	s.movementOptions.Foot()
	return s
}

// MovementLight delegates to the MovementBuilderOptions.
func (s *BuilderOptions) MovementLight() *BuilderOptions {
	s.movementOptions.Light()
	return s
}

// MovementFly delegates to the MovementBuilderOptions.
func (s *BuilderOptions) MovementFly() *BuilderOptions {
	s.movementOptions.Fly()
	return s
}

// MovementTeleport delegates to the MovementBuilderOptions.
func (s *BuilderOptions) MovementTeleport() *BuilderOptions {
	s.movementOptions.Teleport()
	return s
}

// AddPower makes the squaddie able to use this power.
func (s *BuilderOptions) AddPower(newPower *power.Power) *BuilderOptions {
	s.powersToAdd = append(s.powersToAdd, newPower)
	return s
}

// AddClass adds the class to the squaddie's list of possible classes.
func (s *BuilderOptions) AddClass(newClass *squaddieclass.Class) *BuilderOptions {
	s.classesToAdd = append(s.classesToAdd, newClass)
	return s
}

// SetClass sets the squaddie's class to the given class.
func (s *BuilderOptions) SetClass(targetClass *squaddieclass.Class) *BuilderOptions {
	s.classToUse = targetClass
	return s
}

// Build uses the BuilderOptions to create a Movement.
func (s *BuilderOptions) Build() *squaddie.Squaddie {
	newSquaddie := &squaddie.Squaddie{
		Identification: *s.identificationOptions.Build(),
		Offense:        *s.offenseOptions.Build(),
		Defense:        *s.defenseOptions.Build(),
		Movement:       *s.movementOptions.Build(),
		ClassProgress: squaddie.ClassProgress{
			ClassLevelsConsumed: map[string]*squaddie.ClassLevelsConsumed{},
		},
	}

	for _, newPower := range s.powersToAdd {
		newSquaddie.PowerCollection.AddInnatePower(newPower)
	}

	for _, newClass := range s.classesToAdd {
		newSquaddie.ClassProgress.AddClass(newClass)
	}

	if s.classToUse != nil {
		newSquaddie.ClassProgress.SetClass(s.classToUse.ID)
	}

	return newSquaddie
}

// Teros returns a specific squaddie build for testing.
//   Teros is a player combines physical attacks with magical attacks.
func (s *BuilderOptions) Teros() *BuilderOptions {
	teros := Builder().WithName("Teros").WithID("squaddieTeros").MovementFoot().MoveDistance(3).AsPlayer()
	return teros
}

// Bandit returns a specific squaddie build for testing.
//   Bandit is a weak enemy with an axe.
func (s *BuilderOptions) Bandit() *BuilderOptions {
	bandit := Builder().WithName("Bandit").WithID("squaddieBandit").AsEnemy()
	return bandit
}

// Lini returns a specific squaddie build for testing.
//   Lini is a player who carries a healing staff to aid her allies.
func (s *BuilderOptions) Lini() *BuilderOptions {
	lini := Builder().WithName("Lini").WithID("squaddieLini").AsPlayer()
	return lini
}

// MysticMage returns a specific squaddie build for testing.
//   MysticMage is an enemy with a potent fireball and magical defenses.
func (s *BuilderOptions) MysticMage() *BuilderOptions {
	mysticMage := Builder().WithName("Mystic Mage").WithID("squaddieMysticMage")
	return mysticMage
}

// BuilderOptionMarshal is a flattened representation of all Squaddie Builder options.
type BuilderOptionMarshal struct {
	ID          string               `json:"id" yaml:"id"`
	Name        string               `json:"name" yaml:"name"`
	Affiliation squaddie.Affiliation `json:"affiliation" yaml:"affiliation"`

	MaxHitPoints int `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge        int `json:"dodge" yaml:"dodge"`
	Deflect      int `json:"deflect" yaml:"deflect"`
	MaxBarrier   int `json:"max_barrier" yaml:"max_barrier"`
	Armor        int `json:"armor" yaml:"armor"`

	Aim      int `json:"aim" yaml:"aim"`
	Strength int `json:"strength" yaml:"strength"`
	Mind     int `json:"mind" yaml:"mind"`

	MovementDistance     int                   `json:"movement_distance" yaml:"movement_distance"`
	MovementType         squaddie.MovementType `json:"movement_type" yaml:"movement_type"`
	MovementCanHitAndRun bool                  `json:"hit_and_run" yaml:"hit_and_run"`
}

// UsingYAML uses the yaml data to generate BuilderOptions.
func (s *BuilderOptions) UsingYAML(yamlData []byte) *BuilderOptions {
	return s.usingByteStream(yamlData, yaml.Unmarshal)
}

// UsingJSON uses the json data to generate BuilderOptions.
func (s *BuilderOptions) UsingJSON(jsonData []byte) *BuilderOptions {
	return s.usingByteStream(jsonData, json.Unmarshal)
}

func (s *BuilderOptions) usingByteStream(data []byte, unmarshal utility.UnmarshalFunc) *BuilderOptions {
	var unmarshalError error
	var marshaledOptions BuilderOptionMarshal

	unmarshalError = unmarshal(data, &marshaledOptions)

	if unmarshalError != nil {
		return s
	}

	s.WithID(marshaledOptions.ID).WithName(marshaledOptions.Name).
		HitPoints(marshaledOptions.MaxHitPoints).Dodge(marshaledOptions.Dodge).Deflect(marshaledOptions.Deflect).Barrier(marshaledOptions.MaxBarrier).Armor(marshaledOptions.Armor).
		Aim(marshaledOptions.Aim).Strength(marshaledOptions.Strength).Mind(marshaledOptions.Mind).
		MoveDistance(marshaledOptions.MovementDistance)

	if marshaledOptions.Affiliation == squaddie.Player {
		s.AsPlayer()
	}
	if marshaledOptions.Affiliation == squaddie.Enemy {
		s.AsEnemy()
	}
	if marshaledOptions.Affiliation == squaddie.Ally {
		s.AsAlly()
	}
	if marshaledOptions.Affiliation == squaddie.Neutral {
		s.AsNeutral()
	}

	if marshaledOptions.MovementType == squaddie.Foot {
		s.MovementFoot()
	}
	if marshaledOptions.MovementType == squaddie.Light {
		s.MovementLight()
	}
	if marshaledOptions.MovementType == squaddie.Fly {
		s.MovementFly()
	}
	if marshaledOptions.MovementType == squaddie.Teleport {
		s.MovementTeleport()
	}

	if marshaledOptions.MovementCanHitAndRun == true {
		s.CanHitAndRun()
	}
	return s
}

// CloneOf modifies the BuilderOptions based on the source, except for the ID.
func (s *BuilderOptions) CloneOf(source *squaddie.Squaddie) *BuilderOptions {
	s.WithName(source.Name()).
		HitPoints(source.MaxHitPoints()).Deflect(source.Deflect()).Barrier(source.MaxBarrier()).Armor(source.Armor()).Dodge(source.Dodge()).
		Aim(source.Aim()).Strength(source.Strength()).Mind(source.Mind()).
		MoveDistance(source.MovementDistance())

	if source.Affiliation() == squaddie.Player {
		s.AsPlayer()
	}
	if source.Affiliation() == squaddie.Enemy {
		s.AsEnemy()
	}
	if source.Affiliation() == squaddie.Ally {
		s.AsAlly()
	}
	if source.Affiliation() == squaddie.Neutral {
		s.AsNeutral()
	}

	if source.MovementType() == squaddie.Foot {
		s.MovementFoot()
	}
	if source.MovementType() == squaddie.Light {
		s.MovementLight()
	}
	if source.MovementType() == squaddie.Fly {
		s.MovementFly()
	}
	if source.MovementType() == squaddie.Teleport {
		s.MovementTeleport()
	}
	if source.MovementCanHitAndRun() {
		s.CanHitAndRun()
	}
	return s
}
