// Demonstrate how to use to Query the ECS with systems.
package main

import (
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const (
	startup app.Schedule = "Startup"
	update  app.Schedule = "Update"
)

type position struct {
	ecs.Component

	x int
	y int
}

type velocity struct {
	ecs.Component

	x int
	y int
}

type player struct{ ecs.Component }

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.SetTickRate(time.Second)

	myApp.
		AddSchedule(startup, app.ScheduleOptions{ScheduleType: app.ScheduleTypeStartup}).
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddResource(&logger)

	myApp.
		AddSystem(startup, spawn).
		AddSystem(update, applyVelocity)

	run.RunApps(myApp)
}

func spawn(world *ecs.World) error {
	// this has the player component so this entity will be queried
	_, err := ecs.Spawn(world, position{}, velocity{x: 1, y: 2}, player{})
	if err != nil {
		return err
	}

	// this entity does not have the [player] component so it will not be queried
	_, err = ecs.Spawn(world, position{x: 1_000, y: 1_000}, velocity{x: 10, y: 10})
	if err != nil {
		return err
	}

	return nil
}

func applyVelocity(query *ecs.Query2[*position, velocity, ecs.With[player]], log app.Logger) {
	query.Iter(func(entityId ecs.EntityId, position *position, velocity velocity) {
		position.x += velocity.x
		position.y += velocity.y

		log.Info("%d, %d", position.x, position.y)
	})
}
