package affiliation_test

import (
	"github.com/chadius/terosgamerules/entity/affiliation"
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type FactorySuite struct{}

var _ = Suite(&FactorySuite{})

func (suite *FactorySuite) TestFactoryReturnsPlayerAffiliation(checker *C) {
	affiliationLogic := affiliation.NewAffiliationLogic("player")
	checker.Assert(reflect.TypeOf(affiliationLogic).String(), Equals, "*affiliation.Player")
}

func (suite *FactorySuite) TestFactoryReturnsAllyAffiliation(checker *C) {
	affiliationLogic := affiliation.NewAffiliationLogic("ally")
	checker.Assert(reflect.TypeOf(affiliationLogic).String(), Equals, "*affiliation.Ally")
}

func (suite *FactorySuite) TestFactoryReturnsEnemyAffiliation(checker *C) {
	affiliationLogic := affiliation.NewAffiliationLogic("enemy")
	checker.Assert(reflect.TypeOf(affiliationLogic).String(), Equals, "*affiliation.Enemy")
}

func (suite *FactorySuite) TestWhenUnknownKeyword_ThenFactoryReturnsNeutral(checker *C) {
	affiliationLogic := affiliation.NewAffiliationLogic("kwyjibo")
	checker.Assert(reflect.TypeOf(affiliationLogic).String(), Equals, "*affiliation.Neutral")
}
