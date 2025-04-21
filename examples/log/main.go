package main

import "github.com/lucdrenth/murph_engine/src/log"

func main() {
	logger := log.Console()
	logger.Debug("a debug message")
	logger.Info("some very useful info")
	logger.Warn("a warning about something")
	logger.Error("an error that has happened")
}
