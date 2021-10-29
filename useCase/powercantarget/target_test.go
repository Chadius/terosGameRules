package powercantarget_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powercantarget"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TargetingCheck struct {
	teros   *squaddie.Squaddie
	lini    *squaddie.Squaddie
	bandit  *squaddie.Squaddie
	bandit2 *squaddie.Squaddie
	citizen *squaddie.Squaddie
	mayor   *squaddie.Squaddie
	bomb    *squaddie.Squaddie
	bomb2   *squaddie.Squaddie

	meditation   *power.Power
	axe          *power.Power
	healingStaff *power.Power
	selfDestruct *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection
}

var _ = Suite(&TargetingCheck{})

func (suite *TargetingCheck) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.lini = squaddieBuilder.Builder().Lini().Build()
	suite.bandit = squaddieBuilder.Builder().Bandit().Build()
	suite.bandit2 = squaddieBuilder.Builder().Bandit().WithID("bandit2").WithName("bandit2").Build()
	suite.citizen = squaddieBuilder.Builder().WithName("citizen").AsAlly().Build()
	suite.mayor = squaddieBuilder.Builder().WithName("mayor").AsAlly().Build()
	suite.bomb = squaddieBuilder.Builder().WithName("bomb").AsNeutral().Build()
	suite.bomb2 = squaddieBuilder.Builder().WithName("bomb2").AsNeutral().Build()

	suite.axe = powerBuilder.Builder().Axe().Build()
	suite.meditation = powerBuilder.Builder().TargetsSelf().Build()
	suite.healingStaff = powerBuilder.Builder().HealingStaff().Build()
	suite.selfDestruct = powerBuilder.Builder().TargetsFoe().Build()

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
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.meditation.ID(), suite.teros.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.meditation.ID(), suite.bandit.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.meditation.ID(), suite.citizen.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.meditation.PowerID, suite.bomb.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.meditation.PowerID, suite.lini.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.meditation.PowerID, suite.bandit.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.meditation.PowerID, suite.bomb.ID(), suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestPlayerAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.ID(), suite.healingStaff.PowerID, suite.teros.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.ID(), suite.healingStaff.PowerID, suite.citizen.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.axe.PowerID, suite.bomb.ID(), suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestPlayerAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.ID(), suite.healingStaff.PowerID, suite.bandit.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.lini.ID(), suite.healingStaff.PowerID, suite.bomb.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.axe.PowerID, suite.lini.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.teros.ID(), suite.axe.PowerID, suite.citizen.ID(), suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestEnemyAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.healingStaff.PowerID, suite.bandit.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.axe.PowerID, suite.teros.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.axe.PowerID, suite.citizen.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.axe.PowerID, suite.bomb.ID(), suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestEnemyAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.healingStaff.PowerID, suite.teros.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.healingStaff.PowerID, suite.citizen.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.healingStaff.PowerID, suite.bomb.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bandit.ID(), suite.axe.PowerID, suite.bandit2.ID(), suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestAllyAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.healingStaff.PowerID, suite.teros.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.healingStaff.PowerID, suite.mayor.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.axe.PowerID, suite.bomb.ID(), suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestAllyAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.healingStaff.PowerID, suite.bandit.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.healingStaff.PowerID, suite.bomb.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.axe.PowerID, suite.teros.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.citizen.ID(), suite.axe.PowerID, suite.mayor.ID(), suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestNeutralAffiliationValidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.selfDestruct.PowerID, suite.teros.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.selfDestruct.PowerID, suite.bandit.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.selfDestruct.PowerID, suite.citizen.ID(), suite.repos), Equals, true)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.selfDestruct.PowerID, suite.bomb2.ID(), suite.repos), Equals, true)
}

func (suite *TargetingCheck) TestNeutralAffiliationInvalidAttacks(checker *C) {
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.healingStaff.PowerID, suite.teros.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.healingStaff.PowerID, suite.bandit.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.healingStaff.PowerID, suite.citizen.ID(), suite.repos), Equals, false)
	checker.Assert(powercantarget.CanTargetTargetAffiliationWithPower(suite.bomb.ID(), suite.healingStaff.PowerID, suite.bomb2.ID(), suite.repos), Equals, false)
}

func (suite *TargetingCheck) TestTargetGivesAffiliationReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.ID(), suite.meditation.PowerID, suite.teros.ID(), suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.ID(), suite.meditation.PowerID, suite.lini.ID(), suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.PowerCannotTargetAffiliation)
}

func (suite *TargetingCheck) TestTargetGivesTargetIsDeadReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	suite.bandit.Defense.ReduceHitPoints(suite.bandit.MaxHitPoints())
	checker.Assert(suite.bandit.Defense.IsDead(), Equals, true)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsDead)
}

func (suite *TargetingCheck) TestTargetGivesUserIsDeadReasonForFailure(checker *C) {
	canTarget, reasonForInvalidTarget := powercantarget.IsValidTarget(suite.teros.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos)
	checker.Assert(canTarget, Equals, true)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.TargetIsValid)

	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints())
	checker.Assert(suite.teros.Defense.IsDead(), Equals, true)

	canTarget, reasonForInvalidTarget = powercantarget.IsValidTarget(suite.teros.ID(), suite.axe.PowerID, suite.bandit.ID(), suite.repos)
	checker.Assert(canTarget, Equals, false)
	checker.Assert(reasonForInvalidTarget, Equals, powercantarget.UserIsDead)
}
