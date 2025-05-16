package app

import "fmt"

type Logger interface {
	// Log a debug message
	Debug(string)
	// Log an info message
	Info(string)
	// Log a warning message
	Warn(string)
	// Log an error message
	Error(string)
}

var _ Logger = (*NoOpLogger)(nil)

// NoOpLogger is a no-op logger that does nothing.
type NoOpLogger struct {
}

func (n *NoOpLogger) Debug(string) {}
func (n *NoOpLogger) Info(string)  {}
func (n *NoOpLogger) Warn(string)  {}
func (n *NoOpLogger) Error(string) {}

// SimpleConsoleLogger is a logger that prints messages to the console with a level prefix
type SimpleConsoleLogger struct {
}

func (l *SimpleConsoleLogger) Debug(message string) {
	fmt.Println("DEBUG: " + message)
}
func (l *SimpleConsoleLogger) Info(message string) {
	fmt.Println("INFO: " + message)

}
func (l *SimpleConsoleLogger) Warn(message string) {
	fmt.Println("WARN: " + message)

}
func (l *SimpleConsoleLogger) Error(message string) {
	fmt.Println("ERROR: " + message)
}
