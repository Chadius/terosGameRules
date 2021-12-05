package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// Builder is used to define the parameters for a squaddie builder.
type Builder struct {
	identificationOptions   *IdentificationBuilderOptions
	offenseOptions          *OffenseBuilderOptions
	defenseOptions          *DefenseBuilderOptions
	movementOptions         *MovementBuilderOptions
	powerReferencesToAdd    []*power.Reference
	classReferencesToAdd    []*squaddieclass.ClassReference
	levelsConsumedByClassID map[string]*[]string
	classIDToUse            string
	baseClassID             string
}

// NewSquaddieBuilder creates a Builder with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func NewSquaddieBuilder() *Builder {
	return &Builder{
		identificationOptions:   IdentificationBuilder(),
		offenseOptions:          OffenseBuilder(),
		defenseOptions:          DefenseBuilder(),
		movementOptions:         MovementBuilder(),
		powerReferencesToAdd:    []*power.Reference{},
		classReferencesToAdd:    []*squaddieclass.ClassReference{},
		classIDToUse:            "",
		baseClassID:             "",
		levelsConsumedByClassID: map[string]*[]string{},
	}
}

// WithName delegates to the IdentificationBuilderOptions.
func (s *Builder) WithName(name string) *Builder {
	s.identificationOptions.WithName(name)
	return s
}

// WithID delegates to the IdentificationBuilderOptions.
func (s *Builder) WithID(id string) *Builder {
	s.identificationOptions.WithID(id)
	return s
}

// AsPlayer delegates to the IdentificationBuilderOptions.
func (s *Builder) AsPlayer() *Builder {
	s.identificationOptions.AsPlayer()
	return s
}

// AsEnemy delegates to the IdentificationBuilderOptions.
func (s *Builder) AsEnemy() *Builder {
	s.identificationOptions.AsEnemy()
	return s
}

// AsAlly delegates to the IdentificationBuilderOptions.
func (s *Builder) AsAlly() *Builder {
	s.identificationOptions.AsAlly()
	return s
}

// AsNeutral delegates to the IdentificationBuilderOptions.
func (s *Builder) AsNeutral() *Builder {
	s.identificationOptions.AsNeutral()
	return s
}

// Aim delegates to the OffenseBuilderOptions.
func (s *Builder) Aim(aim int) *Builder {
	s.offenseOptions.Aim(aim)
	return s
}

// Strength delegates to the OffenseBuilderOptions.
func (s *Builder) Strength(strength int) *Builder {
	s.offenseOptions.Strength(strength)
	return s
}

// Mind delegates to the OffenseBuilderOptions.
func (s *Builder) Mind(mind int) *Builder {
	s.offenseOptions.Mind(mind)
	return s
}

// HitPoints delegates to the DefenseBuilderOptions.
func (s *Builder) HitPoints(maxHitPoints int) *Builder {
	s.defenseOptions.HitPoints(maxHitPoints)
	return s
}

// Barrier delegates to the DefenseBuilderOptions.
func (s *Builder) Barrier(maxBarrier int) *Builder {
	s.defenseOptions.Barrier(maxBarrier)
	return s
}

// Armor delegates to the DefenseBuilderOptions.
func (s *Builder) Armor(armor int) *Builder {
	s.defenseOptions.Armor(armor)
	return s
}

// Dodge delegates to the DefenseBuilderOptions.
func (s *Builder) Dodge(dodge int) *Builder {
	s.defenseOptions.Dodge(dodge)
	return s
}

// Deflect delegates to the DefenseBuilderOptions.
func (s *Builder) Deflect(deflect int) *Builder {
	s.defenseOptions.Deflect(deflect)
	return s
}

// MoveDistance delegates to the MovementBuilderOptions.
func (s *Builder) MoveDistance(distance int) *Builder {
	s.movementOptions.Distance(distance)
	return s
}

// CanHitAndRun delegates to the MovementBuilderOptions.
func (s *Builder) CanHitAndRun() *Builder {
	s.movementOptions.CanHitAndRun()
	return s
}

// MovementFoot delegates to the MovementBuilderOptions.
func (s *Builder) MovementFoot() *Builder {
	s.movementOptions.Foot()
	return s
}

// MovementLight delegates to the MovementBuilderOptions.
func (s *Builder) MovementLight() *Builder {
	s.movementOptions.Light()
	return s
}

// MovementFly delegates to the MovementBuilderOptions.
func (s *Builder) MovementFly() *Builder {
	s.movementOptions.Fly()
	return s
}

// MovementTeleport delegates to the MovementBuilderOptions.
func (s *Builder) MovementTeleport() *Builder {
	s.movementOptions.Teleport()
	return s
}

// AddPowerByReference adds the power to the squaddie's collection.
func (s *Builder) AddPowerByReference(newPowerReference *power.Reference) *Builder {
	s.powerReferencesToAdd = append(s.powerReferencesToAdd, newPowerReference)
	return s
}

// AddClassByReference adds the class to the squaddie's list of possible classes.
func (s *Builder) AddClassByReference(newClassReference *squaddieclass.ClassReference) *Builder {
	s.classReferencesToAdd = append(s.classReferencesToAdd, newClassReference)
	return s
}

// AddClassLevelsConsumed adds consumed class levels.
func (s *Builder) AddClassLevelsConsumed(classID string, levelIDsConsumed *[]string) *Builder {
	s.levelsConsumedByClassID[classID] = levelIDsConsumed
	return s
}

// SetClassByID sets the squaddie's class to the given class.
func (s *Builder) SetClassByID(targetClassID string) *Builder {
	s.classIDToUse = targetClassID
	return s
}

// SetBaseClassByID sets the squaddie's class to the given class.
func (s *Builder) SetBaseClassByID(targetClassID string) *Builder {
	s.baseClassID = targetClassID
	return s
}

// Build uses the Builder to create a Movement.
func (s *Builder) Build() *Squaddie {
	// TODO Make a NewSquaddie() function that takes these smaller objects as arguments
	newSquaddie := &Squaddie{
		Identification: *s.identificationOptions.Build(),
		Offense:        *s.offenseOptions.Build(),
		Defense:        *s.defenseOptions.Build(),
		Movement:       *s.movementOptions.Build(),
		ClassProgress:  *squaddieclass.NewClassProgress("", "", nil),
	}

	for _, newPowerReference := range s.powerReferencesToAdd {
		newSquaddie.AddPowerReference(newPowerReference)
	}

	for _, newClassReference := range s.classReferencesToAdd {
		newSquaddie.AddClass(newClassReference)
	}

	if s.baseClassID != "" {
		newSquaddie.SetBaseClassIfNoBaseClass(s.baseClassID)
	}

	if s.classIDToUse != "" {
		newSquaddie.SetClass(s.classIDToUse)
	}

	for classID, levelsConsumed := range s.levelsConsumedByClassID {
		for _, levelID := range *levelsConsumed {
			newSquaddie.MarkLevelUpBenefitAsConsumed(classID, levelID)
		}
	}

	return newSquaddie
}

// Teros returns a specific squaddie build for testing.
//   Teros is a player combines physical attacks with magical attacks.
func (s *Builder) Teros() *Builder {
	teros := NewSquaddieBuilder().WithName("Teros").WithID("squaddieTeros").MovementFoot().MoveDistance(3).AsPlayer()
	return teros
}

// Bandit returns a specific squaddie build for testing.
//   Bandit is a weak enemy with an axe.
func (s *Builder) Bandit() *Builder {
	bandit := NewSquaddieBuilder().WithName("Bandit").WithID("squaddieBandit").AsEnemy()
	return bandit
}

// Lini returns a specific squaddie build for testing.
//   Lini is a player who carries a healing staff to aid her allies.
func (s *Builder) Lini() *Builder {
	lini := NewSquaddieBuilder().WithName("Lini").WithID("squaddieLini").AsPlayer()
	return lini
}

// MysticMage returns a specific squaddie build for testing.
//   MysticMage is an enemy with a potent fireball and magical defenses.
func (s *Builder) MysticMage() *Builder {
	mysticMage := NewSquaddieBuilder().WithName("Mystic Mage").WithID("squaddieMysticMage")
	return mysticMage
}

// BuilderOptionMarshal is a flattened representation of all Squaddie NewSquaddieBuilder options.
type BuilderOptionMarshal struct {
	ID          string      `json:"id" yaml:"id"`
	Name        string      `json:"name" yaml:"name"`
	Affiliation Affiliation `json:"affiliation" yaml:"affiliation"`

	MaxHitPoints int `json:"max_hit_points" yaml:"max_hit_points"`
	Dodge        int `json:"dodge" yaml:"dodge"`
	Deflect      int `json:"deflect" yaml:"deflect"`
	MaxBarrier   int `json:"max_barrier" yaml:"max_barrier"`
	Armor        int `json:"armor" yaml:"armor"`

	Aim      int `json:"aim" yaml:"aim"`
	Strength int `json:"strength" yaml:"strength"`
	Mind     int `json:"mind" yaml:"mind"`

	MovementDistance     int          `json:"movement_distance" yaml:"movement_distance"`
	MovementType         MovementType `json:"movement_type" yaml:"movement_type"`
	MovementCanHitAndRun bool         `json:"hit_and_run" yaml:"hit_and_run"`

	ClassProgress   []*classProgressMarshal `json:"class_progress" yaml:"class_progress"`
	PowerReferences []*power.Reference      `json:"powers" yaml:"powers"`
}

type classProgressMarshal struct {
	BaseClass      bool     `json:"is_base_class" yaml:"is_base_class"`
	CurrentClass   bool     `json:"is_current_class" yaml:"is_current_class"`
	ClassID        string   `json:"class_id" yaml:"class_id"`
	ClassName      string   `json:"class_name" yaml:"class_name"`
	LevelsConsumed []string `json:"levels_gained" yaml:"levels_gained"`
}

// UsingYAML uses the yaml data to generate Builder.
func (s *Builder) UsingYAML(yamlData []byte) *Builder {
	return s.usingByteStream(yamlData, yaml.Unmarshal)
}

// UsingJSON uses the json data to generate Builder.
func (s *Builder) UsingJSON(jsonData []byte) *Builder {
	return s.usingByteStream(jsonData, json.Unmarshal)
}

func (s *Builder) usingByteStream(data []byte, unmarshal utility.UnmarshalFunc) *Builder {
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

	if marshaledOptions.Affiliation == Player {
		s.AsPlayer()
	}
	if marshaledOptions.Affiliation == Enemy {
		s.AsEnemy()
	}
	if marshaledOptions.Affiliation == Ally {
		s.AsAlly()
	}
	if marshaledOptions.Affiliation == Neutral {
		s.AsNeutral()
	}

	if marshaledOptions.MovementType == Foot {
		s.MovementFoot()
	}
	if marshaledOptions.MovementType == Light {
		s.MovementLight()
	}
	if marshaledOptions.MovementType == Fly {
		s.MovementFly()
	}
	if marshaledOptions.MovementType == Teleport {
		s.MovementTeleport()
	}

	if marshaledOptions.MovementCanHitAndRun == true {
		s.CanHitAndRun()
	}

	if marshaledOptions.PowerReferences != nil {
		for _, reference := range marshaledOptions.PowerReferences {
			s.AddPowerByReference(reference)
		}
	}

	if marshaledOptions.ClassProgress != nil {
		for _, progress := range marshaledOptions.ClassProgress {
			s.AddClassByReference(&squaddieclass.ClassReference{
				ID:   progress.ClassID,
				Name: progress.ClassName,
			})
			s.AddClassLevelsConsumed(progress.ClassID, &progress.LevelsConsumed)
			if progress.BaseClass {
				s.SetBaseClassByID(progress.ClassID)
			}
			if progress.CurrentClass {
				s.SetClassByID(progress.ClassID)
			}
		}
	}

	return s
}

// CloneOf modifies the Builder based on the source, except for the classID.
func (s *Builder) CloneOf(source *Squaddie) *Builder {
	s.WithName(source.Name()).
		HitPoints(source.MaxHitPoints()).Deflect(source.Deflect()).Barrier(source.MaxBarrier()).Armor(source.Armor()).Dodge(source.Dodge()).
		Aim(source.Aim()).Strength(source.Strength()).Mind(source.Mind()).
		MoveDistance(source.MovementDistance())
	s.cloneAffiliation(source)
	s.cloneMovement(source)
	s.clonePowerReferences(source)
	s.cloneClassProgress(source)
	return s
}

func (s *Builder) cloneMovement(source *Squaddie) {
	if source.MovementType() == Foot {
		s.MovementFoot()
	}
	if source.MovementType() == Light {
		s.MovementLight()
	}
	if source.MovementType() == Fly {
		s.MovementFly()
	}
	if source.MovementType() == Teleport {
		s.MovementTeleport()
	}
	if source.MovementCanHitAndRun() {
		s.CanHitAndRun()
	}
}

func (s *Builder) cloneAffiliation(source *Squaddie) {
	if source.Affiliation() == Player {
		s.AsPlayer()
	}
	if source.Affiliation() == Enemy {
		s.AsEnemy()
	}
	if source.Affiliation() == Ally {
		s.AsAlly()
	}
	if source.Affiliation() == Neutral {
		s.AsNeutral()
	}
}

func (s *Builder) clonePowerReferences(source *Squaddie) {
	for _, reference := range source.GetCopyOfPowerReferences() {
		s.AddPowerByReference(reference)
	}
}

func (s *Builder) cloneClassProgress(source *Squaddie) {
	for classID, classLevelsConsumed := range *source.ClassLevelsConsumed() {
		s.AddClassByReference(&squaddieclass.ClassReference{
			ID:   classID,
			Name: classLevelsConsumed.GetClassName(),
		})
		levelsConsumed := classLevelsConsumed.GetLevelsConsumed()
		s.AddClassLevelsConsumed(classID, &levelsConsumed)
	}

	s.SetClassByID(source.CurrentClassID())
}

// NewSquaddieFromMarshal creates a new NewSquaddieBuilder with fields based on the Marshal object
func NewSquaddieFromMarshal(builderFields BuilderOptionMarshal) *Builder {
	s := NewSquaddieBuilder().populateBuilderBasedOnMarshal(builderFields)
	return s
}

func (s *Builder) populateBuilderBasedOnMarshal(builderFields BuilderOptionMarshal) *Builder {
	s.WithName(builderFields.Name).
		WithID(builderFields.ID).
		Aim(builderFields.Aim).
		Strength(builderFields.Strength).
		Mind(builderFields.Mind).
		HitPoints(builderFields.MaxHitPoints).
		Barrier(builderFields.MaxBarrier).
		Armor(builderFields.Armor).
		Dodge(builderFields.Dodge).
		Deflect(builderFields.Deflect).
		MoveDistance(builderFields.MovementDistance)

	if builderFields.Affiliation == Player {
		s.AsPlayer()
	}
	if builderFields.Affiliation == Enemy {
		s.AsEnemy()
	}
	if builderFields.Affiliation == Ally {
		s.AsAlly()
	}
	if builderFields.Affiliation == Neutral {
		s.AsNeutral()
	}

	if builderFields.MovementCanHitAndRun {
		s.CanHitAndRun()
	}

	if builderFields.MovementType == Foot {
		s.MovementFoot()
	}
	if builderFields.MovementType == Fly {
		s.MovementFly()
	}
	if builderFields.MovementType == Light {
		s.MovementLight()
	}
	if builderFields.MovementType == Teleport {
		s.MovementTeleport()
	}

	for _, reference := range builderFields.PowerReferences {
		s.AddPowerByReference(reference)
	}

	for _, progress := range builderFields.ClassProgress {
		if progress.BaseClass {
			s.SetBaseClassByID(progress.ClassID)
		}
	}

	for _, progress := range builderFields.ClassProgress {
		if progress.CurrentClass {
			s.SetClassByID(progress.ClassID)
		}
	}

	for _, progress := range builderFields.ClassProgress {
		s.AddClassByReference(&squaddieclass.ClassReference{
			ID:   progress.ClassID,
			Name: progress.ClassName,
		})
	}

	for _, progress := range builderFields.ClassProgress {
		s.AddClassLevelsConsumed(progress.ClassID, &progress.LevelsConsumed)
	}

	return s
}
