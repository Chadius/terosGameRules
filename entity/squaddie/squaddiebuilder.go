package squaddie

import (
	"encoding/json"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
	"gopkg.in/yaml.v2"
)

// BuilderOptions is used to create healing effects.
type BuilderOptions struct {
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

// Builder creates a BuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func Builder() *BuilderOptions {
	return &BuilderOptions{
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

// AddPowerByReference adds the power to the squaddie's collection.
func (s *BuilderOptions) AddPowerByReference(newPowerReference *power.Reference) *BuilderOptions {
	s.powerReferencesToAdd = append(s.powerReferencesToAdd, newPowerReference)
	return s
}

// AddClassByReference adds the class to the squaddie's list of possible classes.
func (s *BuilderOptions) AddClassByReference(newClassReference *squaddieclass.ClassReference) *BuilderOptions {
	s.classReferencesToAdd = append(s.classReferencesToAdd, newClassReference)
	return s
}

// AddClassLevelsConsumed adds consumed class levels.
func (s *BuilderOptions) AddClassLevelsConsumed(classID string, levelIDsConsumed *[]string) *BuilderOptions {
	s.levelsConsumedByClassID[classID] = levelIDsConsumed
	return s
}

// SetClassByID sets the squaddie's class to the given class.
func (s *BuilderOptions) SetClassByID(targetClassID string) *BuilderOptions {
	s.classIDToUse = targetClassID
	return s
}

// SetBaseClassByID sets the squaddie's class to the given class.
func (s *BuilderOptions) SetBaseClassByID(targetClassID string) *BuilderOptions {
	s.baseClassID = targetClassID
	return s
}

// Build uses the BuilderOptions to create a Movement.
func (s *BuilderOptions) Build() *Squaddie {
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

// CloneOf modifies the BuilderOptions based on the source, except for the classID.
func (s *BuilderOptions) CloneOf(source *Squaddie) *BuilderOptions {
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

func (s *BuilderOptions) cloneMovement(source *Squaddie) {
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

func (s *BuilderOptions) cloneAffiliation(source *Squaddie) {
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

func (s *BuilderOptions) clonePowerReferences(source *Squaddie) {
	for _, reference := range source.GetCopyOfPowerReferences() {
		s.AddPowerByReference(reference)
	}
}

func (s *BuilderOptions) cloneClassProgress(source *Squaddie) {
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
