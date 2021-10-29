package squaddiestats_test

import (
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/usecase/squaddiestats"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type squaddieOffense struct {
	teros *squaddie.Squaddie

	spear *power.Power
	blot  *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository

	repos *repositories.RepositoryCollection
}

var _ = Suite(&squaddieOffense{})

func (suite *squaddieOffense) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()

	suite.spear = powerBuilder.Builder().Spear().Build()
	suite.blot = powerBuilder.Builder().Blot().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.spear, suite.blot})

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
}

func (suite *squaddieOffense) TestSquaddieMeasuresAim(checker *C) {
	suite.teros.Offense.SquaddieAim = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackDamageBonus:                   1,
		AttackExtraBarrierBurn:              0,
		AttackCanBeEquipped:                 true,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: -2,
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 3,
			Damage:                    5,
		},
	}

	suite.blot.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:       2,
		AttackDamageBonus:      0,
		AttackExtraBarrierBurn: 2,
		AttackCanBeEquipped:    true,
		AttackCanCounterAttack: false,
	}

	spearAim, spearErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)

	blotAim, blotErr := squaddiestats.GetSquaddieAimWithPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotAim, Equals, 3)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfSquaddieDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower("does not exist", suite.spear.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "squaddie could not be found, SquaddieID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotExist(checker *C) {
	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.ID(), "does not exist", suite.repos)
	checker.Assert(err, ErrorMatches, "power could not be found, SquaddieID: does not exist")
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerHasNoAttackEffect(checker *C) {
	wait := powerBuilder.Builder().WithID("powerWait").Build()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{wait})

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.spear.GetReference(),
			suite.blot.GetReference(),
			wait.GetReference(),
		},
		suite.repos,
	)

	_, err := squaddiestats.GetSquaddieAimWithPower(suite.teros.ID(), wait.ID(), suite.repos)
	checker.Assert(err, ErrorMatches, "cannot attack with power, SquaddieID: powerWait")
}

func (suite *squaddieOffense) TestGetRawDamageOfPhysicalPower(checker *C) {
	suite.teros.Offense.SquaddieStrength = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackDamageBonus:                   1,
		AttackExtraBarrierBurn:              0,
		AttackCanBeEquipped:                 true,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: -2,
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 3,
			Damage:                    5,
		},
	}

	spearDamage, spearErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.ID(), suite.spear.ID(), suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearDamage, Equals, 2)
}

func (suite *squaddieOffense) TestGetRawDamageOfSpellPower(checker *C) {
	suite.teros.Offense.SquaddieMind = 3

	suite.blot.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:       2,
		AttackDamageBonus:      0,
		AttackExtraBarrierBurn: 2,
		AttackCanBeEquipped:    true,
		AttackCanCounterAttack: false,
	}

	blotDamage, blotErr := squaddiestats.GetSquaddieRawDamageWithPower(suite.teros.ID(), suite.blot.ID(), suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 3)
}

func (suite *squaddieOffense) TestGetCriticalThresholdOfPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackDamageBonus:                   1,
		AttackExtraBarrierBurn:              0,
		AttackCanBeEquipped:                 true,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: -2,
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 2,
			Damage:                    5,
		},
	}

	spearCritThreat, critErr := squaddiestats.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), suite.spear.PowerID, suite.repos)
	checker.Assert(critErr, IsNil)
	checker.Assert(spearCritThreat, Equals, 4)
}

func (suite *squaddieOffense) TestReturnsAnErrorIfPowerDoesNotCrit(checker *C) {
	suite.blot.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackDamageBonus:                   1,
		AttackExtraBarrierBurn:              0,
		AttackCanBeEquipped:                 true,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: -2,
	}

	_, critErr := squaddiestats.GetSquaddieCriticalThresholdWithPower(suite.teros.ID(), suite.blot.PowerID, suite.repos)
	checker.Assert(critErr, ErrorMatches, "cannot critical hit with power, SquaddieID: powerBlot")
}

func (suite *squaddieOffense) TestGetCriticalDamageOfPower(checker *C) {
	suite.teros.Offense.SquaddieStrength = 1

	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackDamageBonus:                   1,
		AttackExtraBarrierBurn:              0,
		AttackCanBeEquipped:                 true,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: -2,
		CriticalEffect: &power.CriticalEffect{
			CriticalHitThresholdBonus: 2,
			Damage:                    5,
		},
	}

	spearDamage, damageErr := squaddiestats.GetSquaddieCriticalRawDamageWithPower(suite.teros.ID(), suite.spear.PowerID, suite.repos)
	checker.Assert(damageErr, IsNil)
	checker.Assert(spearDamage, Equals, 7)
}

func (suite *squaddieOffense) TestSquaddieCanCounterAttackWithPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackCanCounterAttack: true,
	}

	suite.blot.AttackEffect = &power.AttackingEffect{
		AttackCanCounterAttack: false,
	}

	spearCanCounter, spearErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(suite.teros.ID(), suite.spear.PowerID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearCanCounter, Equals, true)

	blotCanCounter, blotErr := squaddiestats.GetSquaddieCanCounterAttackWithPower(suite.teros.ID(), suite.blot.PowerID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotCanCounter, Equals, false)
}

func (suite *squaddieOffense) TestSquaddieShowsCounterAttackToHit(checker *C) {
	suite.teros.Offense.SquaddieAim = 2

	suite.spear.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:                    1,
		AttackCanCounterAttack:              true,
		AttackCounterAttackPenaltyReduction: 1,
	}

	spearAim, spearErr := squaddiestats.GetSquaddieCounterAttackAimWithPower(suite.teros.ID(), suite.spear.PowerID, suite.repos)
	checker.Assert(spearErr, IsNil)
	checker.Assert(spearAim, Equals, 2)
}

func (suite *squaddieOffense) TestGetTotalBarrierBurnOfAttacks(checker *C) {
	suite.teros.Offense.SquaddieMind = 3

	suite.blot.AttackEffect = &power.AttackingEffect{
		AttackToHitBonus:       2,
		AttackDamageBonus:      0,
		AttackExtraBarrierBurn: 2,
		AttackCanBeEquipped:    true,
		AttackCanCounterAttack: false,
	}

	blotDamage, blotErr := squaddiestats.GetSquaddieExtraBarrierBurnWithPower(suite.teros.ID(), suite.blot.PowerID, suite.repos)
	checker.Assert(blotErr, IsNil)
	checker.Assert(blotDamage, Equals, 2)
}

func (suite *squaddieOffense) TestCanCriticallyHitWithPower(checker *C) {
	suite.spear.AttackEffect = &power.AttackingEffect{
		CriticalEffect: &power.CriticalEffect{
			Damage: 1,
		},
	}

	spearCanCrit, spearCanCritErr := squaddiestats.GetSquaddieCanCriticallyHitWithPower(suite.teros.ID(), suite.spear.PowerID, suite.repos)
	checker.Assert(spearCanCritErr, IsNil)
	checker.Assert(spearCanCrit, Equals, true)

	suite.blot.AttackEffect = &power.AttackingEffect{}

	blotCanCrit, blotCanCritErr := squaddiestats.GetSquaddieCanCriticallyHitWithPower(suite.teros.ID(), suite.blot.PowerID, suite.repos)
	checker.Assert(blotCanCritErr, IsNil)
	checker.Assert(blotCanCrit, Equals, false)
}

type healingPower struct {
	lini         *squaddie.Squaddie
	healingStaff *power.Power

	powerRepo    *power.Repository
	squaddieRepo *squaddie.Repository

	repos *repositories.RepositoryCollection
}

var _ = Suite(&healingPower{})

func (suite *healingPower) SetUpTest(checker *C) {
	suite.lini = squaddieBuilder.Builder().Lini().Build()

	suite.healingStaff = powerBuilder.Builder().HealingStaff().Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.lini})

	suite.powerRepo = power.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{suite.healingStaff})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.lini,
		[]*power.Reference{
			suite.healingStaff.GetReference(),
		},
		suite.repos,
	)
}

func (suite *healingPower) TestSquaddieKnowsHealingPotential(checker *C) {
	suite.lini.Defense.SquaddieCurrentHitPoints = 1
	suite.lini.Offense.SquaddieMind = 1
	suite.healingStaff.HealingEffect = &power.HealingEffect{
		HealingHitPointsHealed:                  3,
		HealingHealingAdjustmentBasedOnUserMind: power.Full,
	}

	staffHeal, staffErr := squaddiestats.GetHitPointsHealedWithPower(suite.lini.ID(), suite.healingStaff.PowerID, suite.lini.ID(), suite.repos)
	checker.Assert(staffErr, IsNil)
	checker.Assert(staffHeal, Equals, 4)
}
