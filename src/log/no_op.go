package log

type noOpLogger struct{}

var _ Logger = noOpLogger{}

// NoOp returns a no-op logger, meaning it does nothing when its methods are called.
func NoOp() noOpLogger {
	return noOpLogger{}
}

func (noOpLogger) Debug(string) {}
func (noOpLogger) Info(string)  {}
func (noOpLogger) Warn(string)  {}
func (noOpLogger) Error(string) {}
