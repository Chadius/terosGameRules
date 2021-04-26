package powercounter_test

import (
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/powerusagecontext"
	"github.com/cserrant/terosBattleServer/entity/squaddie"
	"github.com/cserrant/terosBattleServer/usecase/powercommit"
	"github.com/cserrant/terosBattleServer/usecase/powerequip"
	"github.com/cserrant/terosBattleServer/usecase/powerforecast"
	"github.com/cserrant/terosBattleServer/utility/testutility"
	. "gopkg.in/check.v1"
"testing"
)

func Test(t *testing.T) { TestingT(t) }


type TargetAttemptsCounterSuite struct {
	teros *squaddie.Squaddie
	spear *power.Power
	shield *power.Power
	blot *power.Power

	bandit *squaddie.Squaddie
	axe *power.Power

	powerRepo *power.Repository
	squaddieRepo *squaddie.Repository
}

var _ = Suite(&TargetAttemptsCounterSuite{})

func (suite *TargetAttemptsCounterSuite) SetUpTest(checker *C) {
	suite.teros = squaddie.NewSquaddie("teros")
	suite.spear = power.NewPower("spear")
	suite.spear.AttackEffect.CanBeEquipped = true
	suite.spear.AttackEffect.CanCounterAttack = true
	suite.spear.AttackEffect.CounterAttackToHitPenalty = -2

	suite.shield = power.NewPower("shield")
	suite.shield.AttackEffect.CanBeEquipped = true
	suite.shield.AttackEffect.CanCounterAttack = false

	suite.axe = power.NewPower("axe the second")
	suite.axe.AttackEffect.CanBeEquipped = true
	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.CounterAttackToHitPenalty = -2

	suite.blot = power.NewPower("blot")
	suite.blot.PowerType = power.Spell

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.spear,
		suite.shield,
		suite.blot,
		suite.axe,
	})

	terosPowerReferences := []*power.Reference{
		suite.spear.GetReference(),
		suite.shield.GetReference(),
		suite.blot.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.teros, terosPowerReferences, suite.powerRepo)

	suite.bandit = squaddie.NewSquaddie("bandit")
	suite.bandit.Identification.Name = "bandit"
	banditPowerReferences := []*power.Reference{
		suite.axe.GetReference(),
	}
	powerequip.LoadAllOfSquaddieInnatePowers(suite.bandit, banditPowerReferences, suite.powerRepo)

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
	})
}

func (suite *TargetAttemptsCounterSuite) TestTargetWillCounterAttackWithEquippedCounterablePower(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.spear.ID, suite.powerRepo)
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)

	expectedTerosCounterAttackForecast := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.teros.Identification.ID,
			TargetSquaddieIDs: []string{suite.bandit.Identification.ID},
			PowerID:           suite.spear.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:			suite.spear.ID,
			AttackerID:			suite.teros.Identification.ID,
			TargetID:			suite.bandit.Identification.ID,
			IsCounterAttack:	false,
		},
	)
	terosHitRate := expectedTerosCounterAttackForecast.HitRate

	banditAttackForecast := powerforecast.GetExpectedDamage(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.bandit.Identification.ID,
			TargetSquaddieIDs: []string{suite.teros.Identification.ID},
			PowerID:           suite.axe.ID,
			PowerRepo:         suite.powerRepo,
		},
		&powerusagecontext.AttackContext{
			PowerID:			suite.axe.ID,
			AttackerID:			suite.bandit.Identification.ID,
			TargetID:			suite.teros.Identification.ID,
			IsCounterAttack:	false,
		},
	)
	checker.Assert(banditAttackForecast.CounterAttack, NotNil)
	checker.Assert(banditAttackForecast.CounterAttack.IsACounterAttack, Equals, true)
	checker.Assert(banditAttackForecast.CounterAttack.AttackingSquaddieID, Equals, suite.teros.Identification.ID)
	checker.Assert(banditAttackForecast.CounterAttack.PowerID, Equals, suite.spear.ID)
	checker.Assert(banditAttackForecast.CounterAttack.TargetSquaddieID, Equals, suite.bandit.Identification.ID)
	checker.Assert(banditAttackForecast.CounterAttack.HitRate, Equals, terosHitRate - 2)
}

func (suite *TargetAttemptsCounterSuite) TestTargetWillCommitToCounterAttack(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.spear.ID, suite.powerRepo)
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)

	banditAttackTerosCounterattackReport := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.bandit.Identification.ID,
			TargetSquaddieIDs: []string{suite.teros.Identification.ID},
			PowerID:           suite.axe.ID,
			PowerRepo:         suite.powerRepo,
		},
		testutility.AlwaysHitDieRoller{},
	)

	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports, HasLen, 2)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].AttackerID, Equals, suite.bandit.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].TargetID, Equals, suite.teros.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].PowerID, Equals, suite.axe.ID)

	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[1].AttackerID, Equals, suite.teros.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[1].TargetID, Equals, suite.bandit.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[1].PowerID, Equals, suite.spear.ID)
}

func (suite *TargetAttemptsCounterSuite) TestTargetCannotCommitToCounterAttackIfCounterattackIsNotEquipped(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.shield.ID, suite.powerRepo)
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.powerRepo)

	banditAttackTerosCounterattackReport := powercommit.UsePowerAgainstSquaddiesAndGetReport(
		&powerusagecontext.PowerUsageContext{
			SquaddieRepo:      suite.squaddieRepo,
			ActingSquaddieID:  suite.bandit.Identification.ID,
			TargetSquaddieIDs: []string{suite.teros.Identification.ID},
			PowerID:           suite.axe.ID,
			PowerRepo:         suite.powerRepo,
		},
		testutility.AlwaysHitDieRoller{},
	)

	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports, HasLen, 1)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].AttackerID, Equals, suite.bandit.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].TargetID, Equals, suite.teros.Identification.ID)
	checker.Assert(banditAttackTerosCounterattackReport.AttackingPowerReports[0].PowerID, Equals, suite.axe.ID)
}