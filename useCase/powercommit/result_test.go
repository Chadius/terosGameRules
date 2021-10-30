package powercommit_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powercommit"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility/testutility"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type resultOnAttack struct {
	teros      *squaddie.Squaddie
	bandit     *squaddie.Squaddie
	bandit2    *squaddie.Squaddie
	mysticMage *squaddie.Squaddie

	spear    *power.Power
	blot     *power.Power
	fireball *power.Power
	axe      *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit   *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit   *powercommit.Result

	forecastSpearOnMysticMage *powerattackforecast.Forecast

	forecastFireballOnBandits *powerattackforecast.Forecast
	resultFireballOnBandits   *powercommit.Result
}

var _ = Suite(&resultOnAttack{})

func (suite *resultOnAttack) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.mysticMage = squaddieBuilder.Builder().MysticMage().Build()
	suite.bandit = squaddieBuilder.Builder().Bandit().Build()
	suite.bandit2 = squaddieBuilder.Builder().Bandit().WithID("bandit2ID").WithName("bandit2").Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().Build()
	suite.axe = powerBuilder.Builder().Axe().Build()
	suite.fireball = powerBuilder.Builder().IsSpell().DealsDamage(3).Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.bandit2, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.axe, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.axe.GetReference(),
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
		suite.mysticMage,
		[]*power.Reference{
			suite.fireball.GetReference(),
		},
		suite.repos,
	)

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultSpearOnBandit = &powercommit.Result{
		Forecast: suite.forecastSpearOnBandit,
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastSpearOnMysticMage = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.ID(),
			Targets:         []string{suite.mysticMage.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.forecastFireballOnBandits = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.mysticMage.ID(),
			PowerID:         suite.fireball.ID(),
			Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultFireballOnBandits = &powercommit.Result{
		Forecast: suite.forecastFireballOnBandits,
	}
}

func (suite *resultOnAttack) TestAttackCanMiss(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget, HasLen, 1)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 0)
}

func (suite *resultOnAttack) TestAttackCanHitButNotCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.SquaddieMind = 2

	suite.blot.AttackEffect.AttackDamageBonus = 3

	suite.bandit.Defense.SquaddieCurrentBarrier = 3
	suite.bandit.Defense.SquaddieArmor = 1

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID())
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 2)

	checker.Assert(
		suite.bandit.CurrentHitPoints(),
		Equals,
		suite.bandit.MaxHitPoints()-suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestAttackCanHitCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.blot.AttackEffect = &power.AttackingEffect{
		CriticalEffect: powerBuilder.CriticalEffectBuilder().CriticalHitThresholdBonus(9000).DealsDamage(3).Build(),
	}

	suite.teros.Offense.SquaddieMind = 2

	suite.blot.AttackEffect.AttackDamageBonus = 3

	suite.bandit.Defense.SquaddieCurrentBarrier = 3
	suite.bandit.Defense.SquaddieArmor = 1
	suite.bandit.Defense.SquaddieMaxHitPoints = 1
	suite.bandit.Defense.SetHPToMax()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID())
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 5)

	checker.Assert(
		suite.bandit.CurrentHitPoints(),
		Equals,
		0,
	)
}

func (suite *resultOnAttack) TestCounterAttacks(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.SquaddieStrength = 2
	suite.teros.Defense.SquaddieArmor = 0
	suite.teros.Defense.SquaddieCurrentBarrier = 0

	suite.spear.AttackEffect.AttackDamageBonus = 3

	suite.axe.AttackEffect.AttackCanCounterAttack = true
	suite.axe.AttackEffect.AttackDamageBonus = 3
	suite.bandit.Offense.SquaddieStrength = 0
	suite.bandit.Defense.SquaddieArmor = 1
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID(), suite.repos)

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.PowerID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.ID())

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].PowerID, Equals, suite.axe.PowerID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].UserID, Equals, suite.bandit.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].TargetID, Equals, suite.teros.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.IsCounterAttack, Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.RawDamageDealt, Equals, 3)

	checker.Assert(
		suite.teros.CurrentHitPoints(),
		Equals,
		suite.teros.MaxHitPoints()-suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestCounterAttacksApplyLast(checker *C) {
	suite.resultFireballOnBandits.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.axe.AttackEffect.AttackCanCounterAttack = true
	suite.axe.AttackEffect.AttackDamageBonus = 3
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.PowerID, suite.repos)
	powerequip.SquaddieEquipPower(suite.bandit2, suite.axe.PowerID, suite.repos)

	suite.bandit.Defense.SquaddieMaxHitPoints = suite.fireball.DamageBonus() + suite.mysticMage.Mind() + 1
	suite.bandit.Defense.SetHPToMax()

	suite.bandit2.Defense.SquaddieMaxHitPoints = suite.fireball.DamageBonus() + suite.mysticMage.Mind() + 1
	suite.bandit2.Defense.SetHPToMax()

	suite.forecastFireballOnBandits.CalculateForecast()
	suite.resultFireballOnBandits.Commit()

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].PowerID, Equals, suite.fireball.PowerID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].UserID, Equals, suite.mysticMage.ID())
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].TargetID, Equals, suite.bandit.ID())

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].PowerID, Equals, suite.fireball.PowerID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].UserID, Equals, suite.mysticMage.ID())
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].TargetID, Equals, suite.bandit2.ID())

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].PowerID, Equals, suite.axe.PowerID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].UserID, Equals, suite.bandit.ID())
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].TargetID, Equals, suite.mysticMage.ID())

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].PowerID, Equals, suite.axe.PowerID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].UserID, Equals, suite.bandit2.ID())
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].TargetID, Equals, suite.mysticMage.ID())
}

func (suite *resultOnAttack) TestDeadSquaddiesCannotCounterAttack(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.SquaddieStrength = suite.bandit.MaxHitPoints()
	suite.teros.Defense.SquaddieArmor = 0
	suite.teros.Defense.SquaddieCurrentBarrier = 0

	suite.spear.AttackEffect.AttackDamageBonus = 3

	suite.axe.AttackEffect.AttackCanCounterAttack = true
	suite.axe.AttackEffect.AttackDamageBonus = 3
	suite.bandit.Offense.SquaddieStrength = 0
	suite.bandit.Defense.SquaddieArmor = 0
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.PowerID, suite.repos)

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.PowerID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.ID())
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.ID())

	checker.Assert(suite.bandit.Defense.IsDead(), Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget, HasLen, 1)
}

type EquipPowerWhenCommitting struct {
	teros      *squaddie.Squaddie
	bandit     *squaddie.Squaddie
	mysticMage *squaddie.Squaddie

	spear    *power.Power
	blot     *power.Power
	fireball *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastSpearOnBandit *powerattackforecast.Forecast
	resultSpearOnBandit   *powercommit.Result

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit   *powercommit.Result

	forecastFireballOnBandit *powerattackforecast.Forecast
	resultFireballOnBandit   *powercommit.Result
}

var _ = Suite(&EquipPowerWhenCommitting{})

func (suite *EquipPowerWhenCommitting) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.mysticMage = squaddieBuilder.Builder().MysticMage().Build()
	suite.bandit = squaddieBuilder.Builder().Bandit().Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().CannotBeEquipped().Build()
	suite.fireball = powerBuilder.Builder().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.bandit, suite.mysticMage})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot, suite.fireball})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.mysticMage,
		[]*power.Reference{
			suite.fireball.GetReference(),
		},
		suite.repos,
	)

	suite.forecastSpearOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.spear.PowerID,
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultSpearOnBandit = &powercommit.Result{
		Forecast: suite.forecastSpearOnBandit,
	}

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.PowerID,
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastFireballOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.mysticMage.ID(),
			PowerID:         suite.fireball.PowerID,
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultFireballOnBandit = &powercommit.Result{
		Forecast: suite.forecastFireballOnBandit,
	}
}

func (suite *EquipPowerWhenCommitting) TestCommitWillTryToEquipPower(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.teros.PowerCollection.CurrentlyEquippedPowerID, Equals, suite.spear.PowerID)
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillKeepPreviousPowerIfCommitPowerCannotBeEquipped(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.spear.PowerID, suite.repos)
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.PowerID)
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillNotEquipPowerIfNoneExistAfterCommitting(checker *C) {
	suite.resultFireballOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastFireballOnBandit.CalculateForecast()
	suite.resultFireballOnBandit.Commit()
	checker.Assert(suite.mysticMage.PowerCollection.HasEquippedPower(), Equals, false)
}

type ResultOnHealing struct {
	lini  *squaddie.Squaddie
	teros *squaddie.Squaddie
	vale  *squaddie.Squaddie

	healingStaff *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros   *powercommit.Result

	forecastHealingStaffOnTerosAndVale *powerattackforecast.Forecast
	resultHealingStaffOnTerosAndVale   *powercommit.Result
}

var _ = Suite(&ResultOnHealing{})

func (suite *ResultOnHealing) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.lini = squaddieBuilder.Builder().Lini().Build()
	suite.vale = squaddieBuilder.Builder().WithName("Vale").AsPlayer().Build()

	suite.healingStaff = powerBuilder.Builder().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.lini, suite.vale})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.ID(),
			PowerID:         suite.healingStaff.PowerID,
			Targets:         []string{suite.teros.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.resultHealingStaffOnTeros = &powercommit.Result{
		Forecast: suite.forecastHealingStaffOnTeros,
	}

	suite.forecastHealingStaffOnTerosAndVale = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.ID(),
			PowerID:         suite.healingStaff.PowerID,
			Targets:         []string{suite.teros.ID(), suite.vale.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.resultHealingStaffOnTerosAndVale = &powercommit.Result{
		Forecast: suite.forecastHealingStaffOnTerosAndVale,
	}
}

func (suite *ResultOnHealing) TestHealResultShowsHitPointsRestored(checker *C) {
	suite.teros.Defense.SquaddieCurrentHitPoints = 1
	suite.lini.Offense.SquaddieMind = 1

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()

	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget[0].PowerID, Equals, suite.healingStaff.PowerID)
	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget[0].Healing.HitPointsRestored, Equals, 4)

	checker.Assert(
		suite.teros.Defense.SquaddieCurrentHitPoints,
		Equals,
		1+suite.resultHealingStaffOnTeros.ResultPerTarget[0].Healing.HitPointsRestored,
	)
}

func (suite *ResultOnHealing) TestHealResultShowsForEachTarget(checker *C) {
	suite.teros.Defense.SquaddieCurrentHitPoints = 1
	suite.vale.Defense.SquaddieCurrentHitPoints = 2
	suite.lini.Offense.SquaddieMind = 1

	suite.forecastHealingStaffOnTerosAndVale.CalculateForecast()
	suite.resultHealingStaffOnTerosAndVale.Commit()

	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget, HasLen, 2)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].PowerID, Equals, suite.healingStaff.PowerID)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].Healing.HitPointsRestored, Equals, 4)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].PowerID, Equals, suite.healingStaff.PowerID)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].Healing.HitPointsRestored, Equals, 3)

	checker.Assert(
		suite.teros.Defense.SquaddieCurrentHitPoints,
		Equals,
		1+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].Healing.HitPointsRestored,
	)

	checker.Assert(
		suite.vale.Defense.SquaddieCurrentHitPoints,
		Equals,
		2+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].Healing.HitPointsRestored,
	)
}
