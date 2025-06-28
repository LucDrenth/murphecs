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

// testLogger counts the number times a certain level is logged. Useful for tests
type testLogger struct {
	debug uint
	info  uint
	warn  uint
	err   uint
}

func (l *testLogger) Debug(message string, arguments ...any) {
	l.debug++
}
func (l *testLogger) Info(message string, arguments ...any) {
	l.info++
}
func (l *testLogger) Warn(message string, arguments ...any) {
	l.warn++
}
func (l *testLogger) Error(message string, arguments ...any) {
	l.err++
}
