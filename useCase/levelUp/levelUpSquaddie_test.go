package levelup_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/usecase/levelup"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
)

type SquaddieUsesLevelUpBenefitSuite struct {
	mageClass              *squaddieclass.Class
	statBooster            levelupbenefit.LevelUpBenefit
	teros                  *squaddie.Squaddie
	improveAllMovement     *levelupbenefit.LevelUpBenefit
	upgradeToLightMovement *levelupbenefit.LevelUpBenefit
}

var _ = Suite(&SquaddieUsesLevelUpBenefitSuite{})

func (suite *SquaddieUsesLevelUpBenefitSuite) SetUpTest(checker *C) {
	suite.mageClass = &squaddieclass.Class{
		ID:   "ffffffff",
		Name: "Mage",
	}
	suite.teros = squaddieBuilder.Builder().Teros().WithName("teros").Strength(1).Mind(2).Dodge(3).Deflect(4).Barrier(6).Armor(7).AddClass(suite.mageClass).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.statBooster = levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			ID:      "deadbeef",
			ClassID: suite.mageClass.ID,
		},
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
		Identification: &levelupbenefit.Identification{
			ID:      "aaaaaaa0",
			ClassID: suite.mageClass.ID,
		},
		Movement: squaddieBuilder.MovementBuilder().Fly().CanHitAndRun().Distance(1).Build(),
	}

	suite.upgradeToLightMovement = &levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			ID:      "aaaaaaa1",
			ClassID: suite.mageClass.ID,
		},
		Movement: squaddieBuilder.MovementBuilder().Light().Build(),
	}
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestIncreaseStats(checker *C) {
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
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
	checker.Assert(suite.teros.ClassProgress.IsClassLevelAlreadyUsed(suite.statBooster.Identification.ID), Equals, false)
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID: 1})
	checker.Assert(suite.teros.ClassProgress.IsClassLevelAlreadyUsed(suite.statBooster.Identification.ID), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorForNonexistentClass(checker *C) {
	mushroomClassLevel := levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			ID:      "deedbeeg",
			ClassID: "bad SquaddieID",
		},
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
	err := levelup.ImproveSquaddie(&mushroomClassLevel, suite.teros, nil)
	checker.Assert(err.Error(), Equals, `squaddie "teros" cannot add levels to unknown class "bad SquaddieID"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorIfReusingLevel(checker *C) {
	err := levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.ClassProgress.GetLevelCountsByClass(), DeepEquals, map[string]int{"ffffffff": 1})
	checker.Assert(suite.teros.ClassProgress.IsClassLevelAlreadyUsed(suite.statBooster.Identification.ID), Equals, true)

	err = levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(err.Error(), Equals, `teros already consumed LevelUpBenefit - class:"ffffffff" id:"deadbeef"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestUsingLevelSetsBaseClassIfBaseClassIsUnset(checker *C) {
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, "")
	levelup.ImproveSquaddie(&suite.statBooster, suite.teros, nil)
	checker.Assert(suite.teros.ClassProgress.BaseClassID, Equals, suite.mageClass.ID)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieChangeMovement(checker *C) {
	startingMovement := suite.teros.Movement.GetMovementDistancePerRound()

	err := levelup.ImproveSquaddie(suite.improveAllMovement, suite.teros, nil)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.GetMovementDistancePerRound(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieCannotDowngradeMovement(checker *C) {
	startingMovement := suite.teros.Movement.GetMovementDistancePerRound()
	levelup.ImproveSquaddie(suite.improveAllMovement, suite.teros, nil)

	err := levelup.ImproveSquaddie(suite.upgradeToLightMovement, suite.teros, nil)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.GetMovementDistancePerRound(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.GetMovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

type SquaddieChangePowersWithLevelUpBenefitsSuite struct {
	mageClass    *squaddieclass.Class
	teros        *squaddie.Squaddie
	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
	gainPower    levelupbenefit.LevelUpBenefit
	upgradePower levelupbenefit.LevelUpBenefit
	spear        *power.Power
	spearLevel2  *power.Power
}

var _ = Suite(&SquaddieChangePowersWithLevelUpBenefitsSuite{})

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) SetUpTest(checker *C) {
	suite.mageClass = &squaddieclass.Class{
		ID:   "ffffffff",
		Name: "Mage",
	}
	suite.teros = squaddieBuilder.Builder().Teros().AddClass(&squaddieclass.Class{
		ID:   suite.mageClass.ID,
		Name: "Mage",
	}).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.powerRepo = power.NewPowerRepository()

	suite.spear = powerBuilder.Builder().Spear().WithID("spearlvl1").Build()

	suite.teros.PowerCollection.PowerReferences = []*power.Reference{{Name: "spear", ID: "spearlvl1"}}

	suite.spearLevel2 = powerBuilder.Builder().Spear().WithID("spearlvl2").Build()
	newPowers := []*power.Power{suite.spear, suite.spearLevel2}
	suite.powerRepo.AddSlicePowerSource(newPowers)

	suite.gainPower = levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			ID:                 "aaab1234",
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            suite.mageClass.ID,
		},
		PowerChanges: &levelupbenefit.PowerChanges{
			Gained: []*power.Reference{{Name: "spear", ID: suite.spear.ID}},
		},
	}

	suite.upgradePower = levelupbenefit.LevelUpBenefit{
		Identification: &levelupbenefit.Identification{
			ID:                 "aaaa1235",
			LevelUpBenefitType: levelupbenefit.Big,
			ClassID:            suite.mageClass.ID,
		},
		PowerChanges: &levelupbenefit.PowerChanges{
			Lost:   []*power.Reference{{Name: "spear", ID: suite.spear.ID}},
			Gained: []*power.Reference{{Name: "spear", ID: suite.spearLevel2.ID}},
		},
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieGainPowers(checker *C) {
	err := levelup.ImproveSquaddie(&suite.gainPower, suite.teros, suite.repos)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spear.ID)
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieLosePowers(checker *C) {
	levelup.ImproveSquaddie(&suite.gainPower, suite.teros, suite.repos)
	suite.teros.PowerCollection.GetInnatePowerIDNames()

	err := levelup.ImproveSquaddie(&suite.upgradePower, suite.teros, suite.repos)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetInnatePowerIDNames()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].ID, Equals, suite.spearLevel2.ID)
}
