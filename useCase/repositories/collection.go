package repositories

import (
	"github.com/cserrant/terosbattleserver/entity/levelupbenefit"
	"github.com/cserrant/terosbattleserver/entity/power"
	"github.com/cserrant/terosbattleserver/entity/squaddie"
	"github.com/cserrant/terosbattleserver/entity/squaddieclass"
)

// RepositoryCollection holds all of the repositories used in the setup.
type RepositoryCollection struct {
	SquaddieRepo *squaddie.Repository
	PowerRepo    *power.Repository
	LevelRepo    *levelupbenefit.Repository
	ClassRepo    *squaddieclass.Repository
}
