package repositories

import (
	"github.com/chadius/terosgamerules/entity/levelupbenefit"
	"github.com/chadius/terosgamerules/entity/powerrepository"
	"github.com/chadius/terosgamerules/entity/squaddie"
	"github.com/chadius/terosgamerules/entity/squaddieclass"
)

// RepositoryCollection holds all of the repositories used in the setup.
type RepositoryCollection struct {
	SquaddieRepo *squaddie.Repository
	PowerRepo    *powerrepository.Repository
	LevelRepo    *levelupbenefit.Repository
	ClassRepo    *squaddieclass.Repository
}
