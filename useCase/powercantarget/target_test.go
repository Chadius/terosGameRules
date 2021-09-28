package powercantarget_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercantarget"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/repositories"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TargetingCheck struct {
	teros			*squaddie.Squaddie
	lini			*squaddie.Squaddie
	bandit			*squaddie.Squaddie
	bandit2			*squaddie.Squaddie
	citizen			*squaddie.Squaddie
	mayor			*squaddie.Squaddie
	bomb			*squaddie.Squaddie
	bomb2			*squaddie.Squaddie

	meditation    	*power.Power
	axe      		*power.Power
	healingStaff	*power.Power
	selfDestruct	*power.Power

	powerRepo 		*power.Repository
	squaddieRepo 	*squaddie.Repository
	repos			*repositories.RepositoryCollection
}

var _ = Suite(&TargetingCheck{})

func (suite *TargetingCheck) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.teros.Identification.Name = "teros"
	suite.teros.Identification.ID = "teros"
	suite.teros.Identification.Affiliation = squaddie.Player

	suite.lini = squaddie.NewSquaddie("lini")
	suite.lini.Identification.Name = "lini"
	suite.lini.Identification.ID = "lini"
	suite.lini.Identification.Affiliation = squaddie.Player

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"
	suite.bandit.Identification.ID = "bandit"
	suite.bandit.Identification.Affiliation = squaddie.Enemy

	suite.bandit2 = squaddie.NewSquaddie("bandit2")
	suite.bandit2.Identification.Name = "bandit2"
	suite.bandit2.Identification.ID = "bandit2"
	suite.bandit2.Identification.Affiliation = squaddie.Enemy

	suite.citizen = squaddie.NewSquaddie("citizen")
	suite.citizen.Identification.Name = "citizen"
	suite.citizen.Identification.ID = "citizen"
	suite.citizen.Identification.Affiliation = squaddie.Ally

	suite.mayor = squaddie.NewSquaddie("mayor")
	suite.mayor.Identification.Name = "mayor"
	suite.mayor.Identification.ID = "mayor"
	suite.mayor.Identification.Affiliation = squaddie.Ally

	suite.bomb = squaddie.NewSquaddie("bomb")
	suite.bomb.Identification.Name = "bomb"
	suite.bomb.Identification.ID = "bomb"
	suite.bomb.Identification.Affiliation = squaddie.Neutral

	suite.bomb2 = squaddie.NewSquaddie("bomb2")
	suite.bomb2.Identification.Name = "bomb2"
	suite.bomb2.Identification.ID = "bomb2"
	suite.bomb2.Identification.Affiliation = squaddie.Neutral

	suite.axe = power.NewPower("axe")
	suite.axe.ID = "axe"
	suite.axe.PowerType = power.Physical
	suite.axe.Targeting = power.Targeting{
		TargetSelf:   false,
		TargetFoe:    true,
		TargetFriend: false,
	}

	suite.meditation = power.NewPower("meditation")
	suite.meditation.ID = "meditation"
	suite.meditation.PowerType = power.Physical
	suite.meditation.Targeting = power.Targeting{
		TargetSelf:   true,
		TargetFoe:    false,
		TargetFriend: false,
	}

	suite.healingStaff = power.NewPower("healingStaff")
	suite.healingStaff.ID = "healingStaff"
	suite.healingStaff.PowerType = power.Physical
	suite.healingStaff.Targeting = power.Targeting{
		TargetSelf:   false,
		TargetFoe:    false,
		TargetFriend: true,
	}

	suite.selfDestruct = power.NewPower("selfDestruct")
	suite.selfDestruct.ID = "selfDestruct"
	suite.selfDestruct.PowerType = power.Physical
	suite.selfDestruct.Targeting = power.Targeting{
		TargetSelf:   false,
		TargetFoe:    true,
		TargetFriend: false,
	}

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.lini,
		suite.bandit,
		suite.bandit2,
		suite.citizen,
		suite.mayor,
		suite.bomb,
		suite.bomb2,
	})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.meditation,
		suite.axe,
		suite.healingStaff,
		suite.selfDestruct,
	})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.meditation.GetReference(),
			suite.axe.GetReference(),
			suite.healingStaff.GetReference(),
			suite.selfDestruct.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.lini,
		[]*power.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.meditation.GetReference(),
			suite.axe.GetReference(),
			suite.healingStaff.GetReference(),
			suite.selfDestruct.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.citizen,
		[]*power.Reference{
			suite.meditation.GetReference(),
			suite.axe.GetReference(),
			suite.healingStaff.GetReference(),
			suite.selfDestruct.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.mayor,
		[]*power.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bomb,
		[]*power.Reference{
			suite.meditation.GetReference(),
			suite.axe.GetReference(),
			suite.healingStaff.GetReference(),
			suite.selfDestruct.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bomb2,
		[]*power.Reference{
			suite.selfDestruct.GetReference(),
		},
		suite.repos,
	)
}

func (suite *TargetingCheck) TestTargetingSelfCanOnlyTargetUser(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.meditation.ID, suite.teros.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.meditation.ID, suite.bandit.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.meditation.ID, suite.citizen.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.meditation.ID, suite.bomb.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.meditation.ID, suite.lini.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.meditation.ID, suite.bandit.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.meditation.ID, suite.bomb.Identification.ID, suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestPlayerAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.Identification.ID, suite.healingStaff.ID, suite.teros.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.Identification.ID, suite.healingStaff.ID, suite.citizen.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.axe.ID, suite.bomb.Identification.ID, suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestPlayerAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.Identification.ID, suite.healingStaff.ID, suite.bandit.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.Identification.ID, suite.healingStaff.ID, suite.bomb.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.axe.ID, suite.lini.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.Identification.ID, suite.axe.ID, suite.citizen.Identification.ID, suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestEnemyAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.healingStaff.ID, suite.bandit.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.axe.ID, suite.teros.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.axe.ID, suite.citizen.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.axe.ID, suite.bomb.Identification.ID, suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestEnemyAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.healingStaff.ID, suite.teros.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.healingStaff.ID, suite.citizen.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.healingStaff.ID, suite.bomb.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.Identification.ID, suite.axe.ID, suite.bandit2.Identification.ID, suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestAllyAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.healingStaff.ID, suite.teros.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.healingStaff.ID, suite.mayor.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.axe.ID, suite.bomb.Identification.ID, suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestAllyAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.healingStaff.ID, suite.bandit.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.healingStaff.ID, suite.bomb.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.axe.ID, suite.teros.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.Identification.ID, suite.axe.ID, suite.mayor.Identification.ID, suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestNeutralAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.selfDestruct.ID, suite.teros.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.selfDestruct.ID, suite.bandit.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.selfDestruct.ID, suite.citizen.Identification.ID, suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.selfDestruct.ID, suite.bomb2.Identification.ID, suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestNeutralAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.healingStaff.ID, suite.teros.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.healingStaff.ID, suite.bandit.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.healingStaff.ID, suite.citizen.Identification.ID, suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.Identification.ID, suite.healingStaff.ID, suite.bomb2.Identification.ID, suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestTargetGivesAffiliationReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.meditation.ID, suite.teros.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.meditation.ID, suite.lini.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.PowerCannotTargetAffiliation)
}

func (suite *TargetingCheck) TestTargetGivesTargetIsDeadReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	suite.bandit.Defense.ReduceHitPoints(suite.bandit.Defense.MaxHitPoints)
	checker.Assert(suite.bandit.Defense.IsDead(), Equals, true)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsDead)
}

func (suite *TargetingCheck) TestTargetGivesUserIsDeadReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	suite.teros.Defense.ReduceHitPoints(suite.teros.Defense.MaxHitPoints)
	checker.Assert(suite.teros.Defense.IsDead(), Equals, true)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.Identification.ID, suite.axe.ID, suite.bandit.Identification.ID, suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.UserIsDead)
}
