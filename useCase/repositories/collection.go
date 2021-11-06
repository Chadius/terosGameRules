package repositories

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
)

// RepositoryCollection holds all of the repositories used in the setup.
type RepositoryCollection struct {
	SquaddieRepo *squaddie.Repository
	PowerRepo    *powerrepository.Repository
	LevelRepo    *levelupbenefit.Repository
	ClassRepo    *squaddieclass.Repository
}
