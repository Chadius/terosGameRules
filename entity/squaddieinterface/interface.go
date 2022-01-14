package squaddieinterface

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	"github.com/chadius/terosgamerules/entity/damagedistribution"
	"github.com/chadius/terosgamerules/entity/movement"
	"github.com/chadius/terosgamerules/entity/powerreference"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
)

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	HasSameStatsAs(Interface) bool

	ID() string
	Name() string
	AffiliationLogic() affiliation.Interface
	SetNewIDToRandom()

	ImproveMovement(int, bool, movement.Interface)
	MovementDistance() int
	MovementLogic() movement.Interface
	MovementCanHitAndRun() bool

	ImproveDefense(int, int, int, int, int)
	CurrentHitPoints() int
	MaxHitPoints() int
	SetHPToMax()
	CurrentBarrier() int
	MaxBarrier() int
	SetBarrierToMax()
	ReduceHitPoints(int)
	ReduceBarrier(int)
	Dodge() int
	Deflect() int
	Armor() int
	IsDead() bool
	TakeDamageDistribution(distribution *damagedistribution.DamageDistribution)
	GainHitPoints(healingAmount int) int

	ImproveOffense(int, int, int)
	Aim() int
	Mind() int
	Strength() int

	GetLevelCountsByClass() map[string]int
	BaseClassID() string
	AddClass(*squaddieclass.ClassReference)
	SetClass(classID string) error
	CurrentClassID() string
	HasAddedClass(string) bool
	IsClassLevelAlreadyUsed(string) bool
	SetBaseClassIfNoBaseClass(string)
	MarkLevelUpBenefitAsConsumed(string, string)
	ClassLevelsConsumed() *map[string]*squaddieclass.ClassLevelsConsumed

	GetCopyOfPowerReferences() []*powerreference.Reference
	HasPowerWithID(string) bool
	HasEquippedPower() bool
	EquipPower(string)
	ClearPowerReferences()
	AddPowerReference(*powerreference.Reference)
	AddPowerReferences([]*powerreference.Reference)
	RemovePowerReferences([]*powerreference.Reference)
	GetEquippedPowerID() string
	RemovePowerReferenceByPowerID(powerID string)
}
