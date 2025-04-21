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
	LogCaller       bool  // if true, log line include the file and line number of the method that called the log method
	CallerPathDepth int   // the number of caller directories to include. For example, use 3 for path/to/file.go
	Level           Level // skip logs that are lower than this level
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
		TimestampFormat: "15:04:05.0000",
		LogCaller:       true,
		CallerPathDepth: 3,
		Level:           LevelDebug,
	}
}

func (logger ConsoleLogger) Debug(message string) {
	logger.log(LevelDebug, logger.DebugColor, message)
}

func (logger ConsoleLogger) Info(message string) {
	logger.log(LevelInfo, logger.InfoColor, message)
}

func (logger ConsoleLogger) Warn(message string) {
	logger.log(LevelWarn, logger.WarnColor, message)
}

func (logger ConsoleLogger) Error(message string) {
	logger.log(LevelError, logger.ErrorColor, message)
}

func (logger ConsoleLogger) log(level Level, messageColor ansi.Color, message string) {
	if !logger.Level.Allows(level) {
		return
	}

	caller, ok := logger.getCaller()
	if ok {
		caller = fmt.Sprintf("| %s", ansi.Colorize(logger.CallerColor, caller))
	}

	fmt.Printf("%s | %-14s %s - %s\n",
		ansi.Colorize(logger.TimestampColor, time.Now().Format(logger.TimestampFormat)),
		ansi.Colorize(messageColor, strings.ToUpper(level.String())),
		caller,
		ansi.Colorize(messageColor, message),
	)
}

func (logger ConsoleLogger) getCaller() (caller string, ok bool) {
	if !logger.LogCaller {
		return "", false
	}

	_, fullPath, line, ok := runtime.Caller(3)
	if !ok {
		return "", false
	}

	var path string

	pathSplit := strings.Split(fullPath, "/")
	if len(pathSplit) <= logger.CallerPathDepth {
		path = strings.Join(pathSplit, "/")
	} else {
		path = strings.Join(pathSplit[len(pathSplit)-logger.CallerPathDepth:], "/")
	}

	return fmt.Sprintf("%s:%d", path, line), true
}
