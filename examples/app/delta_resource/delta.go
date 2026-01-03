// Demonstrate how to set a delta resource that can be used in systems
package main

import (
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

	myApp.
		AddSchedule(beforeUpdate, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating})

	myApp.
		AddResource(&logger).
		AddResource(&delta{}).
		AddSystem(beforeUpdate, createSetDeltaResourceSystem(myApp)).
		AddSystem(update, logCurrentDelta)

	run.RunApps(myApp)
}

func createSetDeltaResourceSystem(app *app.SubApp) func(delta *delta) {
	return func(delta *delta) {
		delta.secs = app.Delta()
	}
}

func logCurrentDelta(log app.Logger, delta *delta) {
	log.Info("delta: %f", delta.secs)
}
