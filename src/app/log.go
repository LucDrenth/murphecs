package app

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
)

// Logger is an alias for [ecs.Logger].
type Logger = ecs.Logger

// NoOpLogger is an alias for [ecs.NoOpLogger].
type NoOpLogger = ecs.NoOpLogger

var _ Logger = (*SimpleConsoleLogger)(nil)

// SimpleConsoleLogger is a logger that prints messages to the console with a level prefix
type SimpleConsoleLogger struct {
}

func (l *SimpleConsoleLogger) Debug(message string, arguments ...any) {
	fmt.Println("DEBUG: " + fmt.Sprintf(message, arguments...))
}
func (l *SimpleConsoleLogger) Info(message string, arguments ...any) {
	fmt.Println("INFO: " + fmt.Sprintf(message, arguments...))
}
func (l *SimpleConsoleLogger) Warn(message string, arguments ...any) {
	fmt.Println("WARN: " + fmt.Sprintf(message, arguments...))
}
func (l *SimpleConsoleLogger) Error(message string, arguments ...any) {
	fmt.Println("ERROR: " + fmt.Sprintf(message, arguments...))
}

// TestLogger counts the number of logs per level. It's main use is for tests.
type TestLogger struct {
	NumberOfDebugLogs uint
	NumberOfInfoLogs  uint
	NumberOfWarnLogs  uint
	NumberOfErrorLogs uint
}

func (l *TestLogger) Debug(message string, arguments ...any) {
	l.NumberOfDebugLogs++
}
func (l *TestLogger) Info(message string, arguments ...any) {
	l.NumberOfInfoLogs++
}
func (l *TestLogger) Warn(message string, arguments ...any) {
	l.NumberOfWarnLogs++
}
func (l *TestLogger) Error(message string, arguments ...any) {
	l.NumberOfErrorLogs++
}
