// Demonstrate how to use app features. A feature is a combination of resources and systems that get
// processed (initialized and added to the app) before running the app.
//
// Features are useful because:
//  1. They encapsulate resources and systems in to 1 pluggable Feature that can easily be replaced.
//  2. In contrast to adding systems directly to an app, a resource used as a system param does not need
//     to be added before adding a system that uses that resource.
package main

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine"
	"github.com/lucdrenth/murph_engine/src/engine/features"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

func main() {
	e := engine.Default()

	e.App(engine.AppIDCore).AddFeature(&DebugPrinterFeature{
		AppName:          "Core",
		RepeatedSchedule: schedule.Update,
	})
	e.App(engine.AppIDRenderer).AddFeature(&DebugPrinterFeature{
		AppName:          "Renderer",
		RepeatedSchedule: schedule.Render,
	})

	e.Run()
}

type appNameResource struct {
	name string
}

// A feature that prints something at startup, cleanup and on every tick.
// Its fields make it usable for different sub apps.
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

func tickPrinter(logger log.Logger, tickCounter *features.TickCounter, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Tick number %d", appName.name, tickCounter.Count))
}

func cleanupPrinter(logger log.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Cleaning up", appName.name))
}
