package levelup_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/usecase/levelup"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	squaddieClassBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieUsesLevelUpBenefitSuite struct {
	mageClass              *squaddieclass.Class
	statBooster            levelupbenefit.LevelUpBenefit
	teros                  *squaddie.Squaddie
	improveAllMovement     *levelupbenefit.LevelUpBenefit
	upgradeToLightMovement *levelupbenefit.LevelUpBenefit

	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
}

var _ = Suite(&SquaddieUsesLevelUpBenefitSuite{})

func (suite *SquaddieUsesLevelUpBenefitSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("ffffffff").WithName("Mage").Build()
	suite.teros = squaddieBuilder.Builder().Teros().WithName("teros").Strength(1).Mind(2).Dodge(3).Deflect(4).Barrier(6).Armor(7).AddClassByReference(suite.mageClass.GetReference()).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.statBooster = levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("deadbeef", suite.mageClass.ID(), levelupbenefit.Small),
		Defense: &levelupbenefit.Defense{
			MaxHitPoints: 0,
			Dodge:        4,
			Deflect:      3,
			MaxBarrier:   2,
			Armor:        1,
		},
		Offense: &levelupbenefit.Offense{
			Aim:      7,
			Strength: 6,
			Mind:     5,
		},
	}

	suite.improveAllMovement = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("aaaaaaa0", suite.mageClass.ID(), levelupbenefit.Small),
		Movement:       squaddieBuilder.MovementBuilder().Fly().CanHitAndRun().Distance(1).Build(),
	}

	suite.upgradeToLightMovement = &levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("aaaaaaa1", suite.mageClass.ID(), levelupbenefit.Small),
		Movement:       squaddieBuilder.MovementBuilder().Light().Build(),
	}

	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestIncreaseStats(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.MaxHitPoints(), Equals, 5)
	checker.Assert(suite.teros.Aim(), Equals, 7)
	checker.Assert(suite.teros.Strength(), Equals, 7)
	checker.Assert(suite.teros.Mind(), Equals, 7)
	checker.Assert(suite.teros.Dodge(), Equals, 7)
	checker.Assert(suite.teros.Deflect(), Equals, 7)
	checker.Assert(suite.teros.MaxBarrier(), Equals, 8)
	checker.Assert(suite.teros.Armor(), Equals, 8)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieRecordsLevel(checker *C) {
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, false)
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID(): 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorForNonexistentClass(checker *C) {
	mushroomClassLevel := levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("deedbeeg", "bad SquaddieID", levelupbenefit.Small),
		Defense: &levelupbenefit.Defense{
			MaxHitPoints: 0,
			Dodge:        4,
			Deflect:      3,
			MaxBarrier:   2,
			Armor:        1,
		},
		Offense: &levelupbenefit.Offense{
			Aim:      7,
			Strength: 6,
			Mind:     5,
		},
	}
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&mushroomClassLevel, suite.teros)
	checker.Assert(err.Error(), Equals, `squaddie "teros" cannot add levels to unknown class "bad SquaddieID"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorIfReusingLevel(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"ffffffff": 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, true)

	err = suite.improveSquaddieStrategy.ImproveSquaddie(&suite.statBooster, suite.teros)
	checker.Assert(err.Error(), Equals, `teros already consumed LevelUpBenefit - class:"ffffffff" id:"deadbeef"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestUsingLevelSetsBaseClassIfBaseClassIsUnset(checker *C) {
	checker.Assert(suite.teros.BaseClassID(), Equals, "")
	suite.improveSquaddieStrategy.ImproveSquaddie(&suite.statBooster, suite.teros)
	checker.Assert(suite.teros.BaseClassID(), Equals, suite.mageClass.ID())
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieChangeMovement(checker *C) {
	startingMovement := suite.teros.Movement.MovementDistance()

	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.improveAllMovement, suite.teros)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.MovementDistance(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieCannotDowngradeMovement(checker *C) {
	startingMovement := suite.teros.Movement.MovementDistance()
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.improveAllMovement, suite.teros)

	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.upgradeToLightMovement, suite.teros)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.MovementDistance(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

type SquaddieChangePowersWithLevelUpBenefitsSuite struct {
	mageClass               *squaddieclass.Class
	teros                   *squaddie.Squaddie
	powerRepo               *powerrepository.Repository
	squaddieRepo            *squaddie.Repository
	repos                   *repositories.RepositoryCollection
	gainPower               levelupbenefit.LevelUpBenefit
	upgradePower            levelupbenefit.LevelUpBenefit
	spear                   *power.Power
	spearLevel2             *power.Power
	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
}

var _ = Suite(&SquaddieChangePowersWithLevelUpBenefitsSuite{})

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieClassBuilder.ClassBuilder().WithID("ffffffff").WithName("Mage").Build()
	suite.teros = squaddieBuilder.Builder().Teros().AddClassByReference(suite.mageClass.GetReference()).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.powerRepo = powerrepository.NewPowerRepository()

	suite.spear = powerBuilder.Builder().Spear().WithID("spearlvl1").Build()

	suite.teros.AddPowerReference(&power.Reference{
		Name:    "spear",
		PowerID: "spearlvl1",
	})

	suite.spearLevel2 = powerBuilder.Builder().Spear().WithID("spearlvl2").Build()
	newPowers := []*power.Power{suite.spear, suite.spearLevel2}
	suite.powerRepo.AddSlicePowerSource(newPowers)

	suite.gainPower = levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("aaab1234", suite.mageClass.ID(), levelupbenefit.Big),
		PowerChanges: &levelupbenefit.PowerChanges{
			Gained: []*power.Reference{{Name: "spear", PowerID: suite.spear.PowerID}},
		},
	}

	suite.upgradePower = levelupbenefit.LevelUpBenefit{
		Identification: levelupbenefit.NewIdentification("aaab1235", suite.mageClass.ID(), levelupbenefit.Big),
		PowerChanges: &levelupbenefit.PowerChanges{
			Lost:   []*power.Reference{{Name: "spear", PowerID: suite.spear.PowerID}},
			Gained: []*power.Reference{{Name: "spear", PowerID: suite.spearLevel2.PowerID}},
		},
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieGainPowers(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.gainPower, suite.teros)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetCopyOfPowerReferences()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.spear.PowerID)
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieLosePowers(checker *C) {
	suite.improveSquaddieStrategy.ImproveSquaddie(&suite.gainPower, suite.teros)
	suite.teros.PowerCollection.GetCopyOfPowerReferences()

	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.upgradePower, suite.teros)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.spearLevel2.PowerID)
}
