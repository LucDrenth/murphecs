package main

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/log"
)

func main() {
	logger := log.Console()

	// Log messages at different log levels
	logger.Debug("a debug message")
	logger.Info("some very useful info")
	logger.Warn("a warning about something")
	logger.Error("an error that has happened")

	println()

	// Don't log the same message multiple times if it comes from the same location in the code
	for i := range 5 {
		logger.InfoOnce("this is logged only once")
		logger.InfoOnce(fmt.Sprintf("this is logged %d/5 times because this message is not unique", i+1))
	}

	println()

	logger.Trace("lets find out where this is called from")
}
