// Demonstrate how to use app resources. Resources can be seen
// as singletons that can be used a system params.
package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const update app.Schedule = "Update"

type counter struct {
	value int
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	myApp.AddSchedule(update, app.ScheduleTypeRepeating)

	myApp.
		AddResource(&logger).
		AddResource(&counter{}).
		AddSystem(update, incrementCounter).
		AddSystem(update, logCounter)

	run.RunApp(&myApp)
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
func logCounter(counter counter, log app.Logger) {
	log.Info(fmt.Sprintf("counter value: %d", counter.value))
}
