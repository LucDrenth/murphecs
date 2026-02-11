// Demonstrate how to use a resource from another world
package main

import (
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const startup app.Schedule = "Startup"
const update app.Schedule = "Update"

// 1. First we define a worldId
var appFooId = ecs.WorldId(10)

type targetWorldAppFoo struct{}

func (c targetWorldAppFoo) GetWorldId() *ecs.WorldId {
	return &appFooId
}

// 2. Now we implement QueryOption for it so that we can use it as target world
func (targetWorldAppFoo) GetCombinedQueryOptions(world *ecs.World) (ecs.CombinedQueryOptions, error) {
	return ecs.CombinedQueryOptions{TargetWorld: &appFooId}, nil
}

type myResource struct {
	counter int
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}

	appFoo, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	appFoo.
		AddSchedule(startup, app.ScheduleOptions{ScheduleType: app.ScheduleTypeStartup}).
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddResource(&logger)

	appBar, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	appBar.SetTickRate(time.Second)
	appFoo.SetTickRate(time.Second)

	appBar.
		AddSchedule(startup, app.ScheduleOptions{ScheduleType: app.ScheduleTypeStartup}).
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddResource(&logger)

	// 3. Register appFoo to appBar so that we can target appFoo from an appBar system
	appBar.RegisterOuterWorld(appFooId, appFoo.World())

	// 4. Add a resource that appBar will access
	appFoo.AddResource(&myResource{})

	// 5. Add a system for appBar that uses a resource from appFoo
	appBar.AddSystem(update, func(res ecs.OuterResource[*myResource, targetWorldAppFoo]) error {
		res.Value.counter += 1
		return nil
	})

	appFoo.AddSystem(update, func(res myResource, log app.Logger) {
		log.Info("counter: %d", res.counter)
	})

	run.RunApps(appFoo, appBar)
}
