package repositories

import (
	"github.com/cserrant/terosBattleServer/entity/levelupbenefit"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
)

// RepositoryCollection holds all of the repositories used in the setup.
type RepositoryCollection struct {
	SquaddieRepo *squaddie.Repository
	PowerRepo    *power.Repository
	LevelRepo *levelupbenefit.Repository
	ClassRepo *squaddieclass.Repository
}