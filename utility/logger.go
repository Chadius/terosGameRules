package utility

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Logger is a module-wide way to access an injected logger.
var Logger LoggerInterface

// LoggerInterface interface abstracts the Go log module.
type LoggerInterface interface {
	LogMessage(string, int, LogSeverity, int)
}

// LogMessage is a structured way of writing error messages.
type LogMessage struct {
	Message string
	InvocationDescription string
	Severity LogSeverity
}

// LogSeverity indicates how important the message is.
type LogSeverity string

const (
	// Error means the application will probably fail to work correctly, but it won't crash.
	Error LogSeverity = "ERROR"
)

// InMemoryLogger keeps all of the LogMessages in an in-memory list.
type InMemoryLogger struct {
	Messages []LogMessage
}

// LogMessage logs a message with full options
func (logger *InMemoryLogger) LogMessage(message string, indents int, severity LogSeverity, invocationLevels int) {
	spaceIndents := ""
	for indentIndex:=0; indentIndex<indents; indentIndex++ {
		spaceIndents = spaceIndents + "  "
	}

	newMessage := LogMessage{
		spaceIndents + message,
		getInvocationDescription(invocationLevels),
		severity,
	}
	logger.Messages = append(logger.Messages, newMessage)
}

func getInvocationDescription(invocationLevels int) string {
	pc, fileName, line, _ := runtime.Caller(invocationLevels)
	expectedErrorMessage := fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fileName, line)
	return expectedErrorMessage
}

// FileLogger Logs messages to a default file.
type FileLogger struct {
	logFile *os.File
}

// LogMessage logs a message with full options
func (logger *FileLogger) LogMessage(message string, indents int, severity LogSeverity, invocationLevels int) {
	spaceIndents := ""
	for indentIndex:=0; indentIndex<indents; indentIndex++ {
		spaceIndents = spaceIndents + "  "
	}
	logger.openFileIfNeeded()
	println(logger.logFile)
	log.Printf("%s %s %s", string(severity)[0:5], spaceIndents + message, getInvocationDescription(invocationLevels))
	defer logger.logFile.Close()
}

func (logger *FileLogger) openFileIfNeeded() {
	logger.logFile, _ = os.OpenFile("log.text", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(logger.logFile)
}

// Log uses the module's Logger and calls it (or prints an error if the log isn't defined)
func Log(message string, indents int, severity LogSeverity) {
	if Logger == nil {
		return
	}

	spaceIndents := ""
	for indentIndex:=0; indentIndex<indents; indentIndex++ {
		spaceIndents = spaceIndents + "  "
	}

	Logger.LogMessage(message, indents, severity, 3)
}