package utility_test

import (
	"fmt"
	"github.com/chadius/terosgamerules/utility"
	. "gopkg.in/check.v1"
	"runtime"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type LoggerInterfaceUsingMemory struct {
	logger *utility.InMemoryLogger
}

var _ = Suite(&LoggerInterfaceUsingMemory{})

func (suite *LoggerInterfaceUsingMemory) SetUpTest(checker *C) {
	utility.Logger = nil
	suite.logger = &utility.InMemoryLogger{}
}

func getInvocationFunctionAndFilename(stackTraceLevels int) string {
	pc, fileName, _, _ := runtime.Caller(stackTraceLevels)
	return fmt.Sprintf("%s[%s:", runtime.FuncForPC(pc).Name(), fileName)
}

func (suite *LoggerInterfaceUsingMemory) TestLoggerClassCanLogAllOptions(checker *C) {
	suite.logger.LogMessage("Oh no, an error!", 1, utility.Error, 2)
	invocationFunctionAndFilename := getInvocationFunctionAndFilename(1)
	checker.Assert(suite.logger.Messages, HasLen, 1)
	checker.Assert(strings.HasPrefix(suite.logger.Messages[0].InvocationDescription, invocationFunctionAndFilename), Equals, true)
	checker.Assert(suite.logger.Messages[0].Message, Equals, "  Oh no, an error!")
	checker.Assert(suite.logger.Messages[0].Severity, Equals, utility.Error)
}

func (suite *LoggerInterfaceUsingMemory) TestInjectingLogger(checker *C) {
	utility.Logger = suite.logger

	utility.Log("Oh no, an error!", 1, utility.Error)
	invocationFunctionAndFilename := getInvocationFunctionAndFilename(1)
	checker.Assert(suite.logger.Messages, HasLen, 1)
	checker.Assert(strings.HasPrefix(suite.logger.Messages[0].InvocationDescription, invocationFunctionAndFilename), Equals, true)
	checker.Assert(suite.logger.Messages[0].Message, Equals, "  Oh no, an error!")
	checker.Assert(suite.logger.Messages[0].Severity, Equals, utility.Error)
}
