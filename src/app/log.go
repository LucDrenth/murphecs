package app

import "fmt"

type Logger interface {
	// Log a debug message
	Debug(message string, arguments ...any)
	// Log an info message
	Info(message string, arguments ...any)
	// Log a warning message
	Warn(message string, arguments ...any)
	// Log an error message
	Error(message string, arguments ...any)
}

var _ Logger = (*NoOpLogger)(nil)

// NoOpLogger is a no-op logger that does nothing.
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(message string, arguments ...any) {}
func (n *NoOpLogger) Info(message string, arguments ...any)  {}
func (n *NoOpLogger) Warn(message string, arguments ...any)  {}
func (n *NoOpLogger) Error(message string, arguments ...any) {}

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
