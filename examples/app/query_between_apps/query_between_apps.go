// Demonstrate how to query the world of another SubApp.
package main

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const startup app.Schedule = "Startup"
const update app.Schedule = "Update"

type myComponent struct {
	ecs.Component
	value int
}

// 1. First we define a worldId
var appFooId = ecs.WorldId(10)

type targetWorldAppFoo struct{}

func (c targetWorldAppFoo) GetWorldId() *ecs.WorldId {
	return &appFooId
}

// 2. Now we implement QueryOption for it so that we use it as a query parameter
func (targetWorldAppFoo) GetCombinedQueryOptions(world *ecs.World) (ecs.CombinedQueryOptions, error) {
	return ecs.CombinedQueryOptions{TargetWorld: &appFooId}, nil
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}

	appFoo, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	appFoo.AddResource(&logger)
	appFoo.AddSchedule(startup, app.ScheduleTypeStartup)
	appFoo.AddSchedule(update, app.ScheduleTypeRepeating)

	appBar, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	appBar.SetTickRate(time.Second * 2)
	appBar.AddResource(&logger)
	appBar.AddSchedule(startup, app.ScheduleTypeStartup)
	appBar.AddSchedule(update, app.ScheduleTypeRepeating)

	// 3. Register appFoo to appBar so that we can query appFoo from a appBar system
	appBar.RegisterOuterWorld(appFooId, appFoo.World())

	// 4. We register a startup system for appFoo that spawns components that appBar will query
	appFoo.AddSystem(startup, func(world *ecs.World, log app.Logger) error {
		entity, err := ecs.Spawn(world, &myComponent{value: 100})
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("spawned entity %d with value 100", entity))

		entity, err = ecs.Spawn(world, &myComponent{value: 200})
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("spawned entity %d with value 200", entity))

		return nil
	})

	// 5. Now register a system for appBar that queries appFoo
	appBar.AddSystem(update, func(query *ecs.Query1[myComponent, targetWorldAppFoo], log app.Logger) error {
		query.Result().Iter(func(entityId ecs.EntityId, a *myComponent) error {
			log.Info(fmt.Sprintf("%d: %d", entityId, a.value))
			return nil
		})

		fmt.Println()
		return nil
	})

	run.RunApp(&appFoo, &appBar)
}
