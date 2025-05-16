// Demonstrate how to use app resources
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

const update app.Schedule = "Update"

type counter struct {
	value int
}

func main() {
	var logger log.Logger = &log.SimpleConsoleLogger{}
	myApp, err := app.NewBasicSubApp(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	myApp.AddSchedule(update, app.ScheduleTypeRepeating)

	myApp.AddResource(&logger)
	myApp.AddResource(&counter{})
	myApp.AddSystem(update, incrementCounter)
	myApp.AddSystem(update, logCounter)

	runApp(&myApp)
}

// Define the counter as a parameter of this system to use it.
// We can set any resource as a parameter as long as its added to the app.
// Note that we use a pointer to counter so that we can mutate it.
func incrementCounter(counter *counter) {
	counter.value++
}

// We define the counter param as a value (not a pointer to the counter), resulting in
// us getting a copy of the actual counter.
// This allows this system to be ran in parallel with other systems.
//
// The log.logger resource is a special resource that is added to every app by default.
// We do not need to add it ourselves.
func logCounter(counter counter, log log.Logger) {
	log.Info(fmt.Sprintf("counter value: %d", counter.value))
}

func runApp(subApp *app.BasicSubApp) {
	exitChannel := make(chan struct{})
	isDoneChannel := make(chan bool)

	subApp.Run(exitChannel, isDoneChannel)

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	<-cancelChan
}
