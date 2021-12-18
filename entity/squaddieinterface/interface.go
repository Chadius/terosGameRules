package squaddieinterface

import "github.com/chadius/terosbattleserver/entity/affiliation"

// Interface will shape how healing powers work with squaddies.
type Interface interface {
	ID() string
	AffiliationLogic() affiliation.Interface

	CurrentHitPoints() int
	MaxHitPoints() int
	CurrentBarrier() int
	Dodge() int
	Deflect() int
	Armor() int

	Mind() int
	Strength() int
}
