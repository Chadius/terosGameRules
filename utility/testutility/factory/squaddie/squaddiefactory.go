package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
)

// SquaddieBuilderOptions is used to create healing effects.
type SquaddieBuilderOptions struct {
	identificationOptions *IdentificationBuilderOptions
	offenseOptions *OffenseBuilderOptions
	defenseOptions *DefenseBuilderOptions
	movementOptions *MovementBuilderOptions
	powersToAdd []*power.Power
	classesToAdd []*squaddieclass.Class
	classToUse *squaddieclass.Class
}

// SquaddieFactory creates a SquaddieBuilderOptions with default values.
//   Can be chained with other class functions. Call Build() to create the
//   final object.
func SquaddieFactory() *SquaddieBuilderOptions {
	return &SquaddieBuilderOptions{
		identificationOptions: IdentificationFactory(),
		offenseOptions: OffenseFactory(),
		defenseOptions: DefenseFactory(),
		movementOptions: MovementFactory(),
		powersToAdd: []*power.Power{},
		classesToAdd: []*squaddieclass.Class{},
		classToUse: nil,
	}
}

// WithName delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) WithName(name string) *SquaddieBuilderOptions {
	s.identificationOptions.WithName(name)
	return s
}

// WithID delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) WithID(id string) *SquaddieBuilderOptions {
	s.identificationOptions.WithID(id)
	return s
}

// AsPlayer delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) AsPlayer() *SquaddieBuilderOptions {
	s.identificationOptions.AsPlayer()
	return s
}

// AsEnemy delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) AsEnemy() *SquaddieBuilderOptions {
	s.identificationOptions.AsEnemy()
	return s
}

// AsAlly delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) AsAlly() *SquaddieBuilderOptions {
	s.identificationOptions.AsAlly()
	return s
}

// AsNeutral delegates to the IdentificationBuilderOptions.
func (s *SquaddieBuilderOptions) AsNeutral() *SquaddieBuilderOptions {
	s.identificationOptions.AsNeutral()
	return s
}


// Aim delegates to the OffenseBuilderOptions.
func (s *SquaddieBuilderOptions) Aim(aim int) *SquaddieBuilderOptions {
	s.offenseOptions.Aim(aim)
	return s
}

// Strength delegates to the OffenseBuilderOptions.
func (s *SquaddieBuilderOptions) Strength(strength int) *SquaddieBuilderOptions {
	s.offenseOptions.Strength(strength)
	return s
}

// Mind delegates to the OffenseBuilderOptions.
func (s *SquaddieBuilderOptions) Mind(mind int) *SquaddieBuilderOptions {
	s.offenseOptions.Mind(mind)
	return s
}


// HitPoints delegates to the DefenseBuilderOptions.
func (s *SquaddieBuilderOptions) HitPoints(maxHitPoints int) *SquaddieBuilderOptions {
	s.defenseOptions.HitPoints(maxHitPoints)
	return s
}

// Barrier delegates to the DefenseBuilderOptions.
func (s *SquaddieBuilderOptions) Barrier(maxBarrier int) *SquaddieBuilderOptions {
	s.defenseOptions.Barrier(maxBarrier)
	return s
}

// Armor delegates to the DefenseBuilderOptions.
func (s *SquaddieBuilderOptions) Armor(armor int) *SquaddieBuilderOptions {
	s.defenseOptions.Armor(armor)
	return s
}

// Dodge delegates to the DefenseBuilderOptions.
func (s *SquaddieBuilderOptions) Dodge(dodge int) *SquaddieBuilderOptions {
	s.defenseOptions.Dodge(dodge)
	return s
}

// Deflect delegates to the DefenseBuilderOptions.
func (s *SquaddieBuilderOptions) Deflect(deflect int) *SquaddieBuilderOptions {
	s.defenseOptions.Deflect(deflect)
	return s
}


// MoveDistance delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) MoveDistance(distance int) *SquaddieBuilderOptions {
	s.movementOptions.Distance(distance)
	return s
}

// CanHitAndRun delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) CanHitAndRun() *SquaddieBuilderOptions {
	s.movementOptions.CanHitAndRun()
	return s
}

// MovementFoot delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) MovementFoot() *SquaddieBuilderOptions {
	s.movementOptions.Foot()
	return s
}

// MovementLight delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) MovementLight() *SquaddieBuilderOptions {
	s.movementOptions.Light()
	return s
}

// MovementFly delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) MovementFly() *SquaddieBuilderOptions {
	s.movementOptions.Fly()
	return s
}

// MovementTeleport delegates to the MovementBuilderOptions.
func (s *SquaddieBuilderOptions) MovementTeleport() *SquaddieBuilderOptions {
	s.movementOptions.Teleport()
	return s
}


// AddPower makes the squaddie able to use this power.
func (s *SquaddieBuilderOptions) AddPower(newPower *power.Power) *SquaddieBuilderOptions {
	s.powersToAdd = append(s.powersToAdd, newPower)
	return s
}


// AddClass adds the class to the squaddie's list of possible classes.
func (s *SquaddieBuilderOptions) AddClass(newClass *squaddieclass.Class) *SquaddieBuilderOptions {
	s.classesToAdd = append(s.classesToAdd, newClass)
	return s
}

// SetClass sets the squaddie's class to the given class.
func (s *SquaddieBuilderOptions) SetClass(targetClass *squaddieclass.Class) *SquaddieBuilderOptions {
	s.classToUse = targetClass
	return s
}


// Build uses the SquaddieBuilderOptions to create a Movement.
func (s *SquaddieBuilderOptions) Build() *squaddie.Squaddie {
	newSquaddie := &squaddie.Squaddie{
		Identification: *s.identificationOptions.Build(),
		Offense: *s.offenseOptions.Build(),
		Defense: *s.defenseOptions.Build(),
		Movement: *s.movementOptions.Build(),
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
func (s *SquaddieBuilderOptions) Teros() *SquaddieBuilderOptions {
	teros := SquaddieFactory().WithName("Teros").WithID("squaddieTeros").MovementFoot().MoveDistance(3).AsPlayer()
	return teros
}

// Bandit returns a specific squaddie build for testing.
//   Bandit is a weak enemy with an axe.
func (s *SquaddieBuilderOptions) Bandit() *SquaddieBuilderOptions {
	bandit := SquaddieFactory().WithName("Bandit").WithID("squaddieBandit").AsEnemy()
	return bandit
}

// Lini returns a specific squaddie build for testing.
//   Lini is a player who carries a healing staff to aid her allies.
func (s *SquaddieBuilderOptions) Lini() *SquaddieBuilderOptions {
	lini := SquaddieFactory().WithName("Lini").WithID("squaddieLini").AsPlayer()
	return lini
}

// MysticMage returns a specific squaddie build for testing.
//   MysticMage is an enemy with a potent fireball and magical defenses.
func (s *SquaddieBuilderOptions) MysticMage() *SquaddieBuilderOptions {
	mysticMage := SquaddieFactory().WithName("Mystic Mage").WithID("squaddieMysticMage")
	return mysticMage
}

