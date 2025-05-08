// Demonstrate how to use app resources
package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/lucdrenth/murph_engine/src/ecs"
	engine "github.com/lucdrenth/murph_engine/src/engine"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type position struct {
	ecs.Component
	x float64
	y float64
}

type velocity struct {
	ecs.Component
	x float64
	y float64
}

func main() {
	e, err := engine.Default()
	if err != nil {
		panic(err)
	}

	e.App(engine.AppIDCore).
		AddSystem(schedule.Startup, spawn).
		AddSystem(schedule.Update, updatePositions).
		AddSystem(schedule.Update, logPositions)

	e.Run()
}

func spawn(world *ecs.World) {
	for range 3 {
		ecs.Spawn(world,
			&position{
				x: rand.Float64() * 100,
				y: rand.Float64() * 100,
			},
			&velocity{
				x: 1.0,
				y: 0.5,
			},
		)
	}
}

func updatePositions(world *ecs.World, query *ecs.Query2[position, velocity, ecs.Default]) {
	query.Result().Iter(func(_ ecs.EntityId, position *position, velocity *velocity) error {
		position.x += velocity.x
		position.y += velocity.y
		return nil
	})
}

// Get the position as read-only so that this system can be ran parallel with other systems.
func logPositions(log log.Logger, world *ecs.World, query *ecs.Query1[position, ecs.QueryOptionsAllReadOnly]) {
	query.Result().Iter(func(entityId ecs.EntityId, position *position) error {
		log.Info(fmt.Sprintf("%d: %.2f, %.2f", entityId, position.x, position.y))
		return nil
	})
}
