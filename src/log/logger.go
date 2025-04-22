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
	// Trace logs a stack trace. Gets logged no matter the log level.
	Trace(string)
	// Trace logs a stack trace. Gets logged no matter the log level. Skip if the same message is logged from the same location in the code.
	TraceOnce(string)

	// ClearStorage clears the storage that is used for the LogOnce methods.
	//
	// It might be a good idea to call this periodically if you make excessive use of those methods
	// to prevent memory from continuously increasing. But keep in mind that this will make the logs
	// from LogOnce methods log once again.
	ClearStorage()
}
