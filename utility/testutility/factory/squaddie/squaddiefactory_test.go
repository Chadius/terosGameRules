package squaddie_test

import (
	squaddieEntity "github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/factory/power"
	"github.com/chadius/terosbattleserver/utility/testutility/factory/squaddie"
	"github.com/chadius/terosbattleserver/utility/testutility/factory/squaddieclass"
	. "gopkg.in/check.v1"
)

type SquaddieIdentificationBuilder struct {}

var _ = Suite(&SquaddieIdentificationBuilder{})

func (suite *SquaddieIdentificationBuilder) TestBuildSquaddieWithName(checker *C) {
	teros := squaddie.SquaddieFactory().WithName("Teros").Build()
	checker.Assert("Teros", Equals, teros.Identification.Name)
}

func (suite *SquaddieIdentificationBuilder) TestBuildSquaddieWithID(checker *C) {
	teros := squaddie.SquaddieFactory().WithID("squaddieTeros").Build()
	checker.Assert("squaddieTeros", Equals, teros.Identification.ID)
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationPlayer(checker *C) {
	teros := squaddie.SquaddieFactory().AsPlayer().Build()
	checker.Assert(squaddieEntity.Player, Equals, teros.Identification.Affiliation)
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationEnemy(checker *C) {
	bandit := squaddie.SquaddieFactory().AsEnemy().Build()
	checker.Assert(squaddieEntity.Enemy, Equals, bandit.Identification.Affiliation)
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationAlly(checker *C) {
	citizen := squaddie.SquaddieFactory().AsAlly().Build()
	checker.Assert(squaddieEntity.Ally, Equals, citizen.Identification.Affiliation)
}

func (suite *SquaddieIdentificationBuilder) TestBuildIdentificationAffiliationNeutral(checker *C) {
	bomb := squaddie.SquaddieFactory().AsNeutral().Build()
	checker.Assert(squaddieEntity.Neutral, Equals, bomb.Identification.Affiliation)
}

type SquaddieOffenseBuilder struct {}

var _ = Suite(&SquaddieOffenseBuilder{})

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithAim(checker *C) {
	teros := squaddie.SquaddieFactory().Aim(3).Build()
	checker.Assert(3, Equals, teros.Offense.Aim)
}

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithStrength(checker *C) {
	teros := squaddie.SquaddieFactory().Strength(2).Build()
	checker.Assert(2, Equals, teros.Offense.Strength)
}

func (suite *SquaddieOffenseBuilder) TestBuildSquaddieWithMind(checker *C) {
	teros := squaddie.SquaddieFactory().Mind(4).Build()
	checker.Assert(4, Equals, teros.Offense.Mind)
}


type SquaddieDefenseBuilder struct {}

var _ = Suite(&SquaddieDefenseBuilder{})

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithHitPoints(checker *C) {
	teros := squaddie.SquaddieFactory().HitPoints(9).Build()
	checker.Assert(9, Equals, teros.Defense.CurrentHitPoints)
	checker.Assert(9, Equals, teros.Defense.MaxHitPoints)
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithBarrier(checker *C) {
	teros := squaddie.SquaddieFactory().Barrier(3).Build()
	checker.Assert(3, Equals, teros.Defense.MaxBarrier)
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithArmor(checker *C) {
	teros := squaddie.SquaddieFactory().Armor(2).Build()
	checker.Assert(2, Equals, teros.Defense.Armor)
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithDodge(checker *C) {
	teros := squaddie.SquaddieFactory().Dodge(1).Build()
	checker.Assert(1, Equals, teros.Defense.Dodge)
}

func (suite *SquaddieDefenseBuilder) TestBuildSquaddieWithDeflect(checker *C) {
	teros := squaddie.SquaddieFactory().Deflect(4).Build()
	checker.Assert(4, Equals, teros.Defense.Deflect)
}


type SquaddieMovementBuilder struct {}

var _ = Suite(&SquaddieMovementBuilder{})

func (suite *SquaddieMovementBuilder) TestBuildWithDistance(checker *C) {
	soldier := squaddie.SquaddieFactory().MoveDistance(3).Build()
	checker.Assert(3, Equals, soldier.Movement.Distance)
}

func (suite *SquaddieMovementBuilder) TestBuildMovementCanHitAndRun(checker *C) {
	runner := squaddie.SquaddieFactory().CanHitAndRun().Build()
	checker.Assert(true, Equals, runner.Movement.HitAndRun)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFoot(checker *C) {
	soldier := squaddie.SquaddieFactory().MovementFoot().Build()
	checker.Assert(squaddieEntity.Foot, Equals, soldier.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementLight(checker *C) {
	ninja := squaddie.SquaddieFactory().MovementLight().Build()
	checker.Assert(squaddieEntity.Light, Equals, ninja.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementFly(checker *C) {
	bird := squaddie.SquaddieFactory().MovementFly().Build()
	checker.Assert(squaddieEntity.Fly, Equals, bird.Movement.Type)
}

func (suite *SquaddieMovementBuilder) TestChangeMovementTeleport(checker *C) {
	wizard := squaddie.SquaddieFactory().MovementTeleport().Build()
	checker.Assert(squaddieEntity.Teleport, Equals, wizard.Movement.Type)
}


type SquaddiePowerBuilder struct {}

var _ = Suite(&SquaddiePowerBuilder{})

func (suite *SquaddiePowerBuilder) TestBuildAddPower(checker *C) {
	spear := power.PowerFactory().Spear().Build()
	teros := squaddie.SquaddieFactory().AddPower(spear).Build()
	checker.Assert(spear.ID, Equals, teros.PowerCollection.PowerReferences[0].ID)
}

type SquaddieClassBuilder struct {}

var _ = Suite(&SquaddieClassBuilder{})

func (suite *SquaddieClassBuilder) TestBuildAddClass(checker *C) {
	mageClass := squaddieclass.ClassFactory().WithID("A class ID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.SquaddieFactory().AddClass(mageClass).Build()
	checker.Assert(true, Equals, teros.ClassProgress.HasAddedClass(mageClass.ID))
}

func (suite *SquaddieClassBuilder) TestBuildSetClass(checker *C) {
	mageClass := squaddieclass.ClassFactory().WithID("A class ID").WithName("mage").WithInitialBigLevelID("level0").Build()
	teros := squaddie.SquaddieFactory().AddClass(mageClass).SetClass(mageClass).Build()
	checker.Assert(mageClass.ID, Equals, teros.ClassProgress.CurrentClass)
}

type SpecificSquaddieBuilder struct {}

var _ = Suite(&SpecificSquaddieBuilder{})

func (suite *SpecificSquaddieBuilder) TestBuildTeros(checker *C) {
	teros := squaddie.SquaddieFactory().Teros().Build()
	checker.Assert("Teros", Equals, teros.Identification.Name)
}

func (suite *SpecificSquaddieBuilder) TestBuildBandit(checker *C) {
	bandit := squaddie.SquaddieFactory().Bandit().Build()
	checker.Assert("Bandit", Equals, bandit.Identification.Name)
}

func (suite *SpecificSquaddieBuilder) TestBuildLini(checker *C) {
	lini := squaddie.SquaddieFactory().Lini().Build()
	checker.Assert("Lini", Equals, lini.Identification.Name)
}

func (suite *SpecificSquaddieBuilder) TestBuildMysticMage(checker *C) {
	mysticMage := squaddie.SquaddieFactory().MysticMage().Build()
	checker.Assert("Mystic Mage", Equals, mysticMage.Identification.Name)
}

