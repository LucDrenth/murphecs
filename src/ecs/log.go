package ecs

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
