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

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const (
	startup app.Schedule = "Startup"
	update  app.Schedule = "Update"
	cleanup app.Schedule = "Cleanup"
)

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.NewBasicSubApp(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.AddResource(&logger)

	myApp.AddSchedule(startup, app.ScheduleTypeStartup).
		AddSchedule(update, app.ScheduleTypeRepeating).
		AddSchedule(cleanup, app.ScheduleTypeCleanup)

	myApp.AddFeature(&debugPrinterFeature{
		AppName: "MyApp",
	})

	run.RunSubApp(&myApp)
}

type appNameResource struct {
	name string
}

// A feature that prints something at startup, cleanup and on every update.
type debugPrinterFeature struct {
	app.Feature
	AppName string
}

type tickCounter struct {
	count int
}

func (f *debugPrinterFeature) Init() {
	f.AddResource(&appNameResource{name: f.AppName})
	f.AddResource(&tickCounter{})
	f.AddSystem(startup, startupPrinter)
	f.AddSystem(update, tickPrinter)
	f.AddSystem(cleanup, cleanupPrinter)
}

func startupPrinter(logger app.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Starting up", appName.name))
}

func tickPrinter(logger app.Logger, tickCounter *tickCounter, appName appNameResource) {
	tickCounter.count += 1
	logger.Info(fmt.Sprintf("%s - Tick number %d", appName.name, tickCounter.count))
}

func cleanupPrinter(logger app.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Cleaning up", appName.name))
}
