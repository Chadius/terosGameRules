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
	powerFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	squaddieFactory "github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
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
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
	suite.mysticMage = squaddieFactory.SquaddieFactory().MysticMage().Build()
	suite.bandit = squaddieFactory.SquaddieFactory().Bandit().Build()
	suite.bandit2 = squaddieFactory.SquaddieFactory().Bandit().WithID("bandit2ID").WithName("bandit2").Build()

	suite.spear = powerFactory.PowerFactory().Spear().Build()
	suite.blot = powerFactory.PowerFactory().Blot().Build()
	suite.axe = powerFactory.PowerFactory().Axe().Build()
	suite.fireball = powerFactory.PowerFactory().IsSpell().DealsDamage(3).Build()

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
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
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
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
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
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.mysticMage.Identification.ID},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.forecastFireballOnBandits = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.mysticMage.Identification.ID,
			PowerID:         suite.fireball.ID,
			Targets:         []string{suite.bandit.Identification.ID, suite.bandit2.Identification.ID},
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
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 0)
}

func (suite *resultOnAttack) TestAttackCanHitButNotCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Mind = 2

	suite.blot.AttackEffect.DamageBonus = 3

	suite.bandit.Defense.CurrentBarrier = 3
	suite.bandit.Defense.Armor = 1

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 2)

	checker.Assert(
		suite.bandit.Defense.CurrentHitPoints,
		Equals,
		suite.bandit.Defense.MaxHitPoints-suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestAttackCanHitCritically(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.blot.AttackEffect = &power.AttackingEffect{
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 9000,
			Damage:                    3,
		},
	}

	suite.teros.Offense.Mind = 2

	suite.blot.AttackEffect.DamageBonus = 3

	suite.bandit.Defense.CurrentBarrier = 3
	suite.bandit.Defense.Armor = 1
	suite.bandit.Defense.MaxHitPoints = 1
	suite.bandit.Defense.SetHPToMax()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].PowerID, Equals, suite.blot.ID)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.CriticallyHitTarget, Equals, true)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByBarrier, Equals, 3)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultBlotOnBandit.ResultPerTarget[0].Attack.Damage.RawDamageDealt, Equals, 5)

	checker.Assert(
		suite.bandit.Defense.CurrentHitPoints,
		Equals,
		0,
	)
}

func (suite *resultOnAttack) TestCounterAttacks(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Strength = 2
	suite.teros.Defense.Armor = 0
	suite.teros.Defense.CurrentBarrier = 0

	suite.spear.AttackEffect.DamageBonus = 3

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 3
	suite.bandit.Offense.Strength = 0
	suite.bandit.Defense.Armor = 1
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].PowerID, Equals, suite.axe.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].UserID, Equals, suite.bandit.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].TargetID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.HitTarget, Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.IsCounterAttack, Equals, true)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.CriticallyHitTarget, Equals, false)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByBarrier, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.DamageAbsorbedByArmor, Equals, 0)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.RawDamageDealt, Equals, 3)

	checker.Assert(
		suite.teros.Defense.CurrentHitPoints,
		Equals,
		suite.teros.Defense.MaxHitPoints-suite.resultSpearOnBandit.ResultPerTarget[1].Attack.Damage.RawDamageDealt,
	)
}

func (suite *resultOnAttack) TestCounterAttacksApplyLast(checker *C) {
	suite.resultFireballOnBandits.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 3
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)
	powerequip.SquaddieEquipPower(suite.bandit2, suite.axe.ID, suite.repos)

	suite.bandit.Defense.MaxHitPoints = suite.fireball.AttackEffect.DamageBonus + suite.mysticMage.Offense.Mind + 1
	suite.bandit.Defense.SetHPToMax()

	suite.bandit2.Defense.MaxHitPoints = suite.fireball.AttackEffect.DamageBonus + suite.mysticMage.Offense.Mind + 1
	suite.bandit2.Defense.SetHPToMax()

	suite.forecastFireballOnBandits.CalculateForecast()
	suite.resultFireballOnBandits.Commit()

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].PowerID, Equals, suite.fireball.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].UserID, Equals, suite.mysticMage.Identification.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].PowerID, Equals, suite.fireball.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].UserID, Equals, suite.mysticMage.Identification.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[1].TargetID, Equals, suite.bandit2.Identification.ID)

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].PowerID, Equals, suite.axe.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].UserID, Equals, suite.bandit.Identification.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[2].TargetID, Equals, suite.mysticMage.Identification.ID)

	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].PowerID, Equals, suite.axe.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].UserID, Equals, suite.bandit2.Identification.ID)
	checker.Assert(suite.resultFireballOnBandits.ResultPerTarget[3].TargetID, Equals, suite.mysticMage.Identification.ID)
}

func (suite *resultOnAttack) TestDeadSquaddiesCannotCounterAttack(checker *C) {
	suite.resultSpearOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense.Strength = suite.bandit.Defense.MaxHitPoints
	suite.teros.Defense.Armor = 0
	suite.teros.Defense.CurrentBarrier = 0

	suite.spear.AttackEffect.DamageBonus = 3

	suite.axe.AttackEffect.CanCounterAttack = true
	suite.axe.AttackEffect.DamageBonus = 3
	suite.bandit.Offense.Strength = 0
	suite.bandit.Defense.Armor = 0
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.ID, suite.repos)

	suite.forecastSpearOnBandit.CalculateForecast()
	suite.resultSpearOnBandit.Commit()

	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].PowerID, Equals, suite.spear.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].UserID, Equals, suite.teros.Identification.ID)
	checker.Assert(suite.resultSpearOnBandit.ResultPerTarget[0].TargetID, Equals, suite.bandit.Identification.ID)

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
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
	suite.mysticMage = squaddieFactory.SquaddieFactory().MysticMage().Build()
	suite.bandit = squaddieFactory.SquaddieFactory().Bandit().Build()

	suite.spear = powerFactory.PowerFactory().Spear().Build()
	suite.blot = powerFactory.PowerFactory().Blot().CannotBeEquipped().Build()
	suite.fireball = powerFactory.PowerFactory().Build()

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
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.spear.ID,
			Targets:         []string{suite.bandit.Identification.ID},
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
			UserID:          suite.teros.Identification.ID,
			PowerID:         suite.blot.ID,
			Targets:         []string{suite.bandit.Identification.ID},
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
			UserID:          suite.mysticMage.Identification.ID,
			PowerID:         suite.fireball.ID,
			Targets:         []string{suite.bandit.Identification.ID},
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

	checker.Assert(suite.teros.PowerCollection.CurrentlyEquippedPowerID, Equals, suite.spear.ID)
}

func (suite *EquipPowerWhenCommitting) TestSquaddieWillKeepPreviousPowerIfCommitPowerCannotBeEquipped(checker *C) {
	powerequip.SquaddieEquipPower(suite.teros, suite.spear.ID, suite.repos)
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	checker.Assert(suite.teros.PowerCollection.HasEquippedPower(), Equals, true)
	checker.Assert(suite.teros.PowerCollection.GetEquippedPowerID(), Equals, suite.spear.ID)
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
	suite.teros = squaddieFactory.SquaddieFactory().Teros().Build()
	suite.lini = squaddieFactory.SquaddieFactory().Lini().Build()
	suite.vale = squaddieFactory.SquaddieFactory().WithName("Vale").AsPlayer().Build()

	suite.healingStaff = powerFactory.PowerFactory().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros, suite.lini, suite.vale})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{PowerRepo: suite.powerRepo, SquaddieRepo: suite.squaddieRepo}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID},
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
			UserID:          suite.lini.Identification.ID,
			PowerID:         suite.healingStaff.ID,
			Targets:         []string{suite.teros.Identification.ID, suite.vale.Identification.ID},
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
	suite.teros.Defense.CurrentHitPoints = 1
	suite.lini.Offense.Mind = 1

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()

	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget[0].PowerID, Equals, suite.healingStaff.ID)
	checker.Assert(suite.resultHealingStaffOnTeros.ResultPerTarget[0].Healing.HitPointsRestored, Equals, 4)

	checker.Assert(
		suite.teros.Defense.CurrentHitPoints,
		Equals,
		1+suite.resultHealingStaffOnTeros.ResultPerTarget[0].Healing.HitPointsRestored,
	)
}

func (suite *ResultOnHealing) TestHealResultShowsForEachTarget(checker *C) {
	suite.teros.Defense.CurrentHitPoints = 1
	suite.vale.Defense.CurrentHitPoints = 2
	suite.lini.Offense.Mind = 1

	suite.forecastHealingStaffOnTerosAndVale.CalculateForecast()
	suite.resultHealingStaffOnTerosAndVale.Commit()

	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget, HasLen, 2)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].PowerID, Equals, suite.healingStaff.ID)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].Healing.HitPointsRestored, Equals, 4)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].PowerID, Equals, suite.healingStaff.ID)
	checker.Assert(suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].Healing.HitPointsRestored, Equals, 3)

	checker.Assert(
		suite.teros.Defense.CurrentHitPoints,
		Equals,
		1+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[0].Healing.HitPointsRestored,
	)

	checker.Assert(
		suite.vale.Defense.CurrentHitPoints,
		Equals,
		2+suite.resultHealingStaffOnTerosAndVale.ResultPerTarget[1].Healing.HitPointsRestored,
	)
}
