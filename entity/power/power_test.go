package power_test

import (
	powerBuilder "github.com/chadius/terosgamerules/entity/power"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type PowerChanceCheckSuite struct{}

var _ = Suite(&PowerChanceCheckSuite{})

func (suite *PowerChanceCheckSuite) TestPowerCanCrit(checker *C) {
	staticPower := powerBuilder.NewPowerBuilder().DealsDamage(1).Build()
	checker.Assert(staticPower.CanCriticallyHit(), Equals, false)

	criticalPower := powerBuilder.NewPowerBuilder().DealsDamage(1).CriticalDealsDamage(1).Build()
	checker.Assert(criticalPower.CanCriticallyHit(), Equals, true)
}
