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
	"os"
	"os/signal"
	"syscall"

	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/lucdrenth/murphecs/src/log"
)

const (
	Startup app.Schedule = "Startup"
	Update  app.Schedule = "Update"
	Cleanup app.Schedule = "Cleanup"
)

func main() {
	var logger log.Logger = &log.SimpleConsoleLogger{}
	myApp, err := app.NewBasicSubApp(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.AddResource(&logger)

	myApp.AddSchedule(Startup, app.ScheduleTypeStartup).
		AddSchedule(Update, app.ScheduleTypeRepeating).
		AddSchedule(Cleanup, app.ScheduleTypeCleanup)

	myApp.AddFeature(&debugPrinterFeature{
		AppName: "MyApp",
	})

	runApp(&myApp)
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
	f.AddSystem(Startup, startupPrinter)
	f.AddSystem(Update, tickPrinter)
	f.AddSystem(Cleanup, cleanupPrinter)
}

func startupPrinter(logger log.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Starting up", appName.name))
}

func tickPrinter(logger log.Logger, tickCounter *tickCounter, appName appNameResource) {
	tickCounter.count += 1
	logger.Info(fmt.Sprintf("%s - Tick number %d", appName.name, tickCounter.count))
}

func cleanupPrinter(logger log.Logger, appName appNameResource) {
	logger.Info(fmt.Sprintf("%s - Cleaning up", appName.name))
}

func runApp(subApp *app.BasicSubApp) {
	exitChannel := make(chan struct{})
	isDoneChannel := make(chan bool)

	subApp.Run(exitChannel, isDoneChannel)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan
}
