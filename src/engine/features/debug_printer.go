package features

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type appNameResource struct {
	name string
}

type DebugPrinterFeature struct {
	app.Feature
	RepeatedSchedule app.Schedule
	AppName          string
}

func (f *DebugPrinterFeature) Init() {
	f.AddResource(&appNameResource{name: f.AppName})
	f.AddSystem(schedule.Startup, startupPrinter)
	f.AddSystem(f.RepeatedSchedule, tickPrinter)
	f.AddSystem(schedule.Cleanup, cleanupPrinter)
}

func startupPrinter(logger log.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Starting up", appName.name))
}

func tickPrinter(logger log.Logger, tickCounter *TickCounter, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Tick number %d", appName.name, tickCounter.Count))
}

func cleanupPrinter(logger log.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Cleaning up", appName.name))
}
