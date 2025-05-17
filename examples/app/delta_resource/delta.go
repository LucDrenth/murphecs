// Demonstrate how to set a delta resource that can be used in systems
package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const beforeUpdate app.Schedule = "BeforeUpdate" // this schedule is where the new delta resource is set.
const update app.Schedule = "Update"             // this schedule can use the updated delta resource

type delta struct {
	secs float64
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	myApp.AddSchedule(beforeUpdate, app.ScheduleTypeRepeating)
	myApp.AddSchedule(update, app.ScheduleTypeRepeating)

	myApp.AddResource(&logger)
	myApp.AddResource(&delta{})
	myApp.AddSystem(beforeUpdate, createSetDeltaResourceSystem(&myApp))
	myApp.AddSystem(update, logCurrentDelta)

	run.RunApp(&myApp)
}

func createSetDeltaResourceSystem(app *app.SubApp) func(delta *delta) {
	return func(delta *delta) {
		delta.secs = app.Delta()
	}
}

func logCurrentDelta(log app.Logger, delta *delta) {
	log.Info(fmt.Sprintf("delta: %f", delta.secs))
}
