package log

type Logger interface {
	// Log a debug message
	Debug(string)
	// Log a debug message. Skip if the same message is logged from the same location in the code.
	DebugOnce(string)
	// Log an info message
	Info(string)
	// Log an info message. Skip if the same message is logged from the same location in the code.
	InfoOnce(string)
	// Log a warning message
	Warn(string)
	// Log a warning message. Skip if the same message is logged from the same location in the code.
	WarnOnce(string)
	// Log an error message
	Error(string)
	// Log an error message. Skip if the same message is logged from the same location in the code.
	ErrorOnce(string)
}
