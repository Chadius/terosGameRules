package squaddie

import (
	"github.com/chadius/terosbattleserver/entity/affiliation"
	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"reflect"
)

// Squaddie represents a person, creature or thing that can take actions on a battlefield.
type Squaddie struct {
	identification  Identification
	classProgress   squaddieclass.ClassProgress
	defense         Defense
	offense         Offense
	movement        Movement
	powerCollection PowerCollection
}

// NewSquaddie returns a Squaddie object.
func NewSquaddie(
	identification *Identification,
	offense *Offense,
	defense *Defense,
	movement *Movement,
	classProgress *squaddieclass.ClassProgress,
) *Squaddie {
	newSquaddie := &Squaddie{
		identification: *identification,
		offense:        *offense,
		defense:        *defense,
		movement:       *movement,
		classProgress:  *classProgress,
	}
	return newSquaddie
}

// ID delegates.
func (s *Squaddie) ID() string {
	return s.identification.ID()
}

// AffiliationLogic delegates.
func (s *Squaddie) AffiliationLogic() affiliation.Interface {
	return s.identification.AffiliationLogic()
}

// Name delegates.
func (s *Squaddie) Name() string {
	return s.identification.Name()
}

// SetNewIDToRandom delegates.
func (s *Squaddie) SetNewIDToRandom() {
	s.identification.SetNewIDToRandom()
}

// ImproveDefense delegates.
func (s *Squaddie) ImproveDefense(maxHitPoints, dodge, deflect, maxBarrier, armor int) {
	s.defense.Improve(maxHitPoints, dodge, deflect, maxBarrier, armor)
}

// MaxHitPoints delegates.
func (s *Squaddie) MaxHitPoints() int {
	return s.defense.MaxHitPoints()
}

// Dodge delegates.
func (s *Squaddie) Dodge() int {
	return s.defense.Dodge()
}

// Deflect delegates.
func (s *Squaddie) Deflect() int {
	return s.defense.Deflect()
}

// MaxBarrier delegates.
func (s *Squaddie) MaxBarrier() int {
	return s.defense.MaxBarrier()
}

// Armor delegates.
func (s *Squaddie) Armor() int {
	return s.defense.Armor()
}

// CurrentHitPoints delegates.
func (s *Squaddie) CurrentHitPoints() int {
	return s.defense.CurrentHitPoints()
}

// CurrentBarrier delegates.
func (s *Squaddie) CurrentBarrier() int {
	return s.defense.CurrentBarrier()
}

// ReduceHitPoints delegates.
func (s *Squaddie) ReduceHitPoints(damage int) {
	s.defense.ReduceHitPoints(damage)
}

// GainHitPoints delegates.
func (s *Squaddie) GainHitPoints(healingAmount int) int {
	return s.defense.GainHitPoints(healingAmount)
}

// ReduceBarrier delegates.
func (s *Squaddie) ReduceBarrier(damage int) {
	s.defense.ReduceBarrier(damage)
}

// SetBarrierToMax delegates.
func (s *Squaddie) SetBarrierToMax() {
	s.defense.SetBarrierToMax()
}

// SetHPToMax delegates.
func (s *Squaddie) SetHPToMax() {
	s.defense.SetHPToMax()
}

// IsDead delegates.
func (s *Squaddie) IsDead() bool {
	return s.defense.IsDead()
}

// TakeDamageDistribution delegates.
func (s *Squaddie) TakeDamageDistribution(distribution *damagedistribution.DamageDistribution) {
	s.defense.TakeDamageDistribution(distribution)
}

// ImproveOffense delegates.
func (s *Squaddie) ImproveOffense(aim, strength, mind int) {
	s.offense.Improve(aim, strength, mind)
}

// Aim delegates.
func (s *Squaddie) Aim() int {
	return s.offense.Aim()
}

// Strength delegates.
func (s *Squaddie) Strength() int {
	return s.offense.Strength()
}

// Mind delegates.
func (s *Squaddie) Mind() int {
	return s.offense.Mind()
}

// ImproveMovement delegates.
func (s *Squaddie) ImproveMovement(distance int, canHitAndRun bool, movementLogic movement.Interface) {
	s.movement.Improve(distance, canHitAndRun, movementLogic)
}

//MovementDistance delegates.
func (s *Squaddie) MovementDistance() int {
	return s.movement.MovementDistance()
}

// MovementLogic delegates.
func (s *Squaddie) MovementLogic() movement.Interface {
	return s.movement.MovementLogic()
}

// MovementCanHitAndRun delegates.
func (s *Squaddie) MovementCanHitAndRun() bool {
	return s.movement.CanHitAndRun()
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
	return s.powerCollection.HasSamePowersAs(&other.powerCollection)
}

func (s *Squaddie) hasSameClassesAs(other *Squaddie) bool {
	return s.classProgress.HasSameClassesAs(&other.classProgress)
}

// ClearPowerReferences delegates.
func (s *Squaddie) ClearPowerReferences() {
	s.powerCollection.ClearPowerReferences()
}

// HasPowerWithID delegates.
func (s *Squaddie) HasPowerWithID(powerID string) bool {
	return s.powerCollection.HasPowerWithID(powerID)
}

// HasEquippedPower delegates.
func (s *Squaddie) HasEquippedPower() bool {
	return s.powerCollection.HasEquippedPower()
}

// EquipPower delegates.
func (s *Squaddie) EquipPower(powerID string) {
	s.powerCollection.EquipPower(powerID)
}

// GetEquippedPowerID delegates.
func (s *Squaddie) GetEquippedPowerID() string {
	return s.powerCollection.GetEquippedPowerID()
}

// AddPowerReference delegates.
func (s *Squaddie) AddPowerReference(reference *power.Reference) {
	s.powerCollection.AddPowerReference(reference)
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
	s.powerCollection.RemovePowerReferenceByPowerID(powerID)
}

// GetLevelCountsByClass delegates.
func (s *Squaddie) GetLevelCountsByClass() map[string]int {
	return s.classProgress.GetLevelCountsByClass()
}

// MarkLevelUpBenefitAsConsumed delegates.
func (s *Squaddie) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	s.classProgress.MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID)
}

// SetClass delegates.
func (s *Squaddie) SetClass(classID string) error {
	return s.classProgress.SetClass(classID)
}

// SetBaseClassIfNoBaseClass delegates.
func (s *Squaddie) SetBaseClassIfNoBaseClass(classID string) {
	s.classProgress.SetBaseClassIfNoBaseClass(classID)
}

// IsClassLevelAlreadyUsed delegates.
func (s *Squaddie) IsClassLevelAlreadyUsed(benefitID string) bool {
	return s.classProgress.IsClassLevelAlreadyUsed(benefitID)
}

// HasAddedClass delegates.
func (s *Squaddie) HasAddedClass(classIDToFind string) bool {
	return s.classProgress.HasAddedClass(classIDToFind)
}

// AddClass delegates.
func (s *Squaddie) AddClass(classReference *squaddieclass.ClassReference) {
	s.classProgress.AddClass(classReference)
}

// GetCopyOfPowerReferences delegates.
func (s *Squaddie) GetCopyOfPowerReferences() []*power.Reference {
	return s.powerCollection.GetCopyOfPowerReferences()
}

// CurrentClassID delegates
func (s *Squaddie) CurrentClassID() string {
	return s.classProgress.CurrentClassID()
}

// BaseClassID delegates
func (s *Squaddie) BaseClassID() string {
	return s.classProgress.BaseClassID()
}

// ClassLevelsConsumed delegates
func (s *Squaddie) ClassLevelsConsumed() *map[string]*squaddieclass.ClassLevelsConsumed {
	return s.classProgress.ClassLevelsConsumed()
}
