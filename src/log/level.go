package log

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError

	// The highest level so that it is always allowed to be logged.
	levelStackTrace
)

// String returns a lowercase representation of the log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case levelStackTrace:
		return "trace"
	default:
		return "unknown"
	}
}

func (l Level) Allows(logLevel Level) bool {
	return l <= logLevel
}
