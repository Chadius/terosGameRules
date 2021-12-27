package squaddieinterface

import (
	"github.com/chadius/terosbattleserver/entity/affiliation"
	"github.com/chadius/terosbattleserver/entity/movement"
	"github.com/chadius/terosbattleserver/entity/powerreference"
)

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	ID() string
	Name() string
	AffiliationLogic() affiliation.Interface

	ImproveMovement(int, bool, movement.Interface)

	ImproveDefense(int, int, int, int, int)
	CurrentHitPoints() int
	MaxHitPoints() int
	CurrentBarrier() int
	Dodge() int
	Deflect() int
	Armor() int

	ImproveOffense(int, int, int)
	Aim() int
	Mind() int
	Strength() int

	GetLevelCountsByClass() map[string]int
	BaseClassID() string
	CurrentClassID() string
	HasAddedClass(string) bool
	IsClassLevelAlreadyUsed(string) bool
	SetBaseClassIfNoBaseClass(string)
	MarkLevelUpBenefitAsConsumed(string, string)

	GetCopyOfPowerReferences() []*powerreference.Reference
	HasPowerWithID(string) bool
	EquipPower(string)
	ClearPowerReferences()
	AddPowerReference(*powerreference.Reference)
	AddPowerReferences([]*powerreference.Reference)
	RemovePowerReferences([]*powerreference.Reference)
}
