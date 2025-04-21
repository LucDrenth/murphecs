// a wrapper around "log/slog"
package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/lucdrenth/murph_engine/src/log/ansi"
)

type ConsoleLogger struct {
	DebugColor ansi.Color
	InfoColor  ansi.Color
	WarnColor  ansi.Color
	ErrorColor ansi.Color

	TimestampColor  ansi.Color
	TimestampFormat string

	CallerColor     ansi.Color
	LogCaller       bool // if true, log line include the file and line number of the method that called the log method
	CallerPathDepth int  // the number of caller directories to include. For example, use 3 for path/to/file.go
}

var _ Logger = ConsoleLogger{}

// Console returns a logger that prints colored output to the console. Output format looks like this:
//
// 12:30:01 | INFO | main.go:15 - a message with a nice color
func Console() ConsoleLogger {
	return ConsoleLogger{
		DebugColor:      ansi.ColorGrey,
		InfoColor:       ansi.ColorWhite,
		WarnColor:       ansi.ColorYellow,
		ErrorColor:      ansi.ColorBrightRed,
		CallerColor:     ansi.ColorCyan,
		TimestampColor:  ansi.ColorGreen,
		TimestampFormat: "15:04:05.000",
		LogCaller:       true,
		CallerPathDepth: 3,
	}
}

func (s ConsoleLogger) Debug(message string) {
	s.log("DEBUG", s.DebugColor, message)
}

func (s ConsoleLogger) Info(message string) {
	s.log("INFO ", s.InfoColor, message)
}

func (s ConsoleLogger) Warn(message string) {
	s.log("WARN ", s.WarnColor, message)
}

func (s ConsoleLogger) Error(message string) {
	s.log("ERROR", s.ErrorColor, message)
}

func (s ConsoleLogger) log(level string, messageColor ansi.Color, message string) {
	caller, ok := s.getCaller()
	if ok {
		caller = fmt.Sprintf("| %s", ansi.Colorize(s.CallerColor, caller))
	}

	fmt.Printf("%s | %s %s - %s\n",
		ansi.Colorize(s.TimestampColor, time.Now().Format(s.TimestampFormat)),
		ansi.Colorize(messageColor, level),
		caller,
		ansi.Colorize(messageColor, message),
	)
}

func (s ConsoleLogger) getCaller() (caller string, ok bool) {
	if !s.LogCaller {
		return "", false
	}

	_, fullPath, line, ok := runtime.Caller(3)
	if !ok {
		return "", false
	}

	var path string

	pathSplit := strings.Split(fullPath, "/")
	if len(pathSplit) <= s.CallerPathDepth {
		path = strings.Join(pathSplit, "/")
	} else {
		path = strings.Join(pathSplit[len(pathSplit)-s.CallerPathDepth:], "/")
	}

	return fmt.Sprintf("%s:%d", path, line), true
}
