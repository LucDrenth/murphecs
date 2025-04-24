package tick

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
)

func Init(app *app.BasicSubApp) {
	app.AddStartupSystem(schedule.Startup, spawnCounter)
	app.AddSystem(schedule.PreUpdate, updateCounter)
}

type TickCounter struct {
	Count uint
	ecs.Component
}

func spawnCounter(world *ecs.World) error {
	_, err := ecs.Spawn(world, &TickCounter{Count: 0})
	return err
}

func updateCounter(world *ecs.World, query *ecs.Query1[TickCounter, ecs.NoFilter, ecs.AllRequired]) {
	query.Exec(world)

	query.Result().Iter(func(_ ecs.EntityId, counter *TickCounter) error {
		counter.Count++
		return nil
	})
}
