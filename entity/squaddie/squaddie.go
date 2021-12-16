package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/affiliation"
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"reflect"
)

// Squaddie represents a person, creature or thing that can take actions on a battlefield.
type Squaddie struct {
	Identification  Identification
	ClassProgress   squaddieclass.ClassProgress
	Defense         Defense
	Offense         Offense
	Movement        Movement
	PowerCollection PowerCollection
}

// ID delegates.
func (s *Squaddie) ID() string {
	return s.Identification.ID()
}

// AffiliationLogic delegates.
func (s *Squaddie) AffiliationLogic() affiliation.Interface {
	return s.Identification.AffiliationLogic()
}

// Name delegates.
func (s *Squaddie) Name() string {
	return s.Identification.Name()
}

// SetNewIDToRandom delegates.
func (s *Squaddie) SetNewIDToRandom() {
	s.Identification.SetNewIDToRandom()
}

// ImproveDefense delegates.
func (s *Squaddie) ImproveDefense(maxHitPoints, dodge, deflect, maxBarrier, armor int) {
	s.Defense.Improve(maxHitPoints, dodge, deflect, maxBarrier, armor)
}

// MaxHitPoints delegates.
func (s *Squaddie) MaxHitPoints() int {
	return s.Defense.MaxHitPoints()
}

// Dodge delegates.
func (s *Squaddie) Dodge() int {
	return s.Defense.Dodge()
}

// Deflect delegates.
func (s *Squaddie) Deflect() int {
	return s.Defense.Deflect()
}

// MaxBarrier delegates.
func (s *Squaddie) MaxBarrier() int {
	return s.Defense.MaxBarrier()
}

// Armor delegates.
func (s *Squaddie) Armor() int {
	return s.Defense.Armor()
}

// CurrentHitPoints delegates.
func (s *Squaddie) CurrentHitPoints() int {
	return s.Defense.CurrentHitPoints()
}

// CurrentBarrier delegates.
func (s *Squaddie) CurrentBarrier() int {
	return s.Defense.CurrentBarrier()
}

// ReduceHitPoints delegates.
func (s *Squaddie) ReduceHitPoints(damage int) {
	s.Defense.ReduceHitPoints(damage)
}

// ReduceBarrier delegates.
func (s *Squaddie) ReduceBarrier(damage int) {
	s.Defense.ReduceBarrier(damage)
}

// ImproveOffense delegates.
func (s *Squaddie) ImproveOffense(aim, strength, mind int) {
	s.Offense.Improve(aim, strength, mind)
}

// Aim delegates.
func (s *Squaddie) Aim() int {
	return s.Offense.Aim()
}

// Strength delegates.
func (s *Squaddie) Strength() int {
	return s.Offense.Strength()
}

// Mind delegates.
func (s *Squaddie) Mind() int {
	return s.Offense.Mind()
}

// ImproveMovement delegates.
func (s *Squaddie) ImproveMovement(distance int, canHitAndRun bool, movementLogic movement.Interface) {
	s.Movement.Improve(distance, canHitAndRun, movementLogic)
}

//MovementDistance delegates.
func (s *Squaddie) MovementDistance() int {
	return s.Movement.MovementDistance()
}

// MovementLogic delegates.
func (s *Squaddie) MovementLogic() movement.Interface {
	return s.Movement.MovementLogic()
}

// MovementCanHitAndRun delegates.
func (s *Squaddie) MovementCanHitAndRun() bool {
	return s.Movement.CanHitAndRun()
}

// HasSameStatsAs returns true if other's stats matches this one.
//   The comparison ignores the ID.
func (s *Squaddie) HasSameStatsAs(other *Squaddie) bool {
	if s.Name() != other.Name() {
		return false
	}

	if reflect.TypeOf(s.AffiliationLogic()).String() != reflect.TypeOf(other.AffiliationLogic()).String() {
		return false
	}

	if !s.hasSameDefenseAs(other) {
		return false
	}
	if !s.hasSameOffenseAs(other) {
		return false
	}
	if !s.hasSameMovementAs(other) {
		return false
	}
	if !s.hasSamePowersAs(other) {
		return false
	}
	if !s.hasSameClassesAs(other) {
		return false
	}
	return true
}

func (s *Squaddie) hasSameMovementAs(other *Squaddie) bool {
	if reflect.TypeOf(s.MovementLogic()).String() != reflect.TypeOf(other.MovementLogic()).String() {
		return false
	}

	if s.MovementDistance() != other.MovementDistance() {
		return false
	}
	if s.MovementCanHitAndRun() != other.MovementCanHitAndRun() {
		return false
	}
	return true
}

func (s *Squaddie) hasSameOffenseAs(other *Squaddie) bool {
	if s.Aim() != other.Aim() {
		return false
	}
	if s.Strength() != other.Strength() {
		return false
	}
	if s.Mind() != other.Mind() {
		return false
	}
	return true
}

func (s *Squaddie) hasSameDefenseAs(other *Squaddie) bool {
	if s.MaxHitPoints() != other.MaxHitPoints() {
		return false
	}
	if s.Dodge() != other.Dodge() {
		return false
	}
	if s.Deflect() != other.Deflect() {
		return false
	}
	if s.MaxBarrier() != other.MaxBarrier() {
		return false
	}
	if s.Armor() != other.Armor() {
		return false
	}
	if s.CurrentHitPoints() != other.CurrentHitPoints() {
		return false
	}
	if s.CurrentBarrier() != other.CurrentBarrier() {
		return false
	}
	return true
}

func (s *Squaddie) hasSamePowersAs(other *Squaddie) bool {
	return s.PowerCollection.HasSamePowersAs(&other.PowerCollection)
}

func (s *Squaddie) hasSameClassesAs(other *Squaddie) bool {
	return s.ClassProgress.HasSameClassesAs(&other.ClassProgress)
}

// ClearPowerReferences delegates.
func (s *Squaddie) ClearPowerReferences() {
	s.PowerCollection.ClearPowerReferences()
}

// HasPowerWithID delegates.
func (s *Squaddie) HasPowerWithID(powerID string) bool {
	return s.PowerCollection.HasPowerWithID(powerID)
}

// HasEquippedPower delegates.
func (s *Squaddie) HasEquippedPower() bool {
	return s.PowerCollection.HasEquippedPower()
}

// EquipPower delegates.
func (s *Squaddie) EquipPower(powerID string) {
	s.PowerCollection.EquipPower(powerID)
}

// GetEquippedPowerID delegates.
func (s *Squaddie) GetEquippedPowerID() string {
	return s.PowerCollection.GetEquippedPowerID()
}

// AddPowerReference delegates.
func (s *Squaddie) AddPowerReference(reference *power.Reference) {
	s.PowerCollection.AddPowerReference(reference)
}

// RemovePowerReferences removes multiple powers.
func (s *Squaddie) RemovePowerReferences(powersToRemove []*power.Reference) {
	for _, powerReferenceLost := range powersToRemove {
		s.RemovePowerReferenceByPowerID(powerReferenceLost.PowerID)
	}
}

// AddPowerReferences adds multiple powers.
func (s *Squaddie) AddPowerReferences(powersToAdd []*power.Reference) {
	for _, powerReferenceGained := range powersToAdd {
		s.AddPowerReference(powerReferenceGained)
	}
}

// RemovePowerReferenceByPowerID delegates.
func (s *Squaddie) RemovePowerReferenceByPowerID(powerID string) {
	s.PowerCollection.RemovePowerReferenceByPowerID(powerID)
}

// GetLevelCountsByClass delegates.
func (s *Squaddie) GetLevelCountsByClass() map[string]int {
	return s.ClassProgress.GetLevelCountsByClass()
}

// MarkLevelUpBenefitAsConsumed delegates.
func (s *Squaddie) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	s.ClassProgress.MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID)
}

// SetClass delegates.
func (s *Squaddie) SetClass(classID string) error {
	return s.ClassProgress.SetClass(classID)
}

// SetBaseClassIfNoBaseClass delegates.
func (s *Squaddie) SetBaseClassIfNoBaseClass(classID string) {
	s.ClassProgress.SetBaseClassIfNoBaseClass(classID)
}

// IsClassLevelAlreadyUsed delegates.
func (s *Squaddie) IsClassLevelAlreadyUsed(benefitID string) bool {
	return s.ClassProgress.IsClassLevelAlreadyUsed(benefitID)
}

// HasAddedClass delegates.
func (s *Squaddie) HasAddedClass(classIDToFind string) bool {
	return s.ClassProgress.HasAddedClass(classIDToFind)
}

// AddClass delegates.
func (s *Squaddie) AddClass(classReference *squaddieclass.ClassReference) {
	s.ClassProgress.AddClass(classReference)
}

// GetCopyOfPowerReferences delegates.
func (s *Squaddie) GetCopyOfPowerReferences() []*power.Reference {
	return s.PowerCollection.GetCopyOfPowerReferences()
}

// CurrentClassID delegates
func (s *Squaddie) CurrentClassID() string {
	return s.ClassProgress.CurrentClassID()
}

// BaseClassID delegates
func (s *Squaddie) BaseClassID() string {
	return s.ClassProgress.BaseClassID()
}

// ClassLevelsConsumed delegates
func (s *Squaddie) ClassLevelsConsumed() *map[string]*squaddieclass.ClassLevelsConsumed {
	return s.ClassProgress.ClassLevelsConsumed()
}
