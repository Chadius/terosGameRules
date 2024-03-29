package squaddieclass_test

import (
	"github.com/chadius/terosgamerules/entity/squaddieclass"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ClassProgressSuite struct{}

var _ = Suite(&ClassProgressSuite{})

func (suite *ClassProgressSuite) TestIfLevelIsConsumed(checker *C) {
	progress := squaddieclass.NewClassLevelsConsumed(
		"007",
		"Superspy",
		[]string{"1", "2", "3", "4", "5"},
	)
	checker.Assert(progress.IsLevelAlreadyConsumed("1"), Equals, true)
	checker.Assert(progress.IsLevelAlreadyConsumed("10"), Equals, false)
}
