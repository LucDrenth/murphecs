package core

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
	"github.com/lucdrenth/murph_engine/src/tick"
)

type CoreApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) CoreApp {
	app := app.NewBasicSubApp(logger)

	app.AddStartupSchedule(schedule.Startup)
	app.AddStartupSystem(schedule.Startup, startup)

	app.AddSchedule(schedule.PreUpdate)
	app.AddSchedule(schedule.Update)
	app.AddSchedule(schedule.PostUpdate)
	app.AddSystem(schedule.Update, printer)

	app.AddCleanupSchedule(schedule.Cleanup)
	app.AddCleanupSystem(schedule.Cleanup, cleanup)

	tick.Init(&app)

	return CoreApp{
		BasicSubApp: app,
	}
}

func startup(logger log.Logger) {
	logger.Info("Init core")
}

func printer(logger log.Logger, world *ecs.World, tickCounterQuery *ecs.Query1[tick.TickCounter, ecs.NoFilter, ecs.AllRequired]) {
	tickCounterQuery.Exec(world)

	count := uint(0)
	tickCounterQuery.Result().Iter(func(entityId ecs.EntityId, a *tick.TickCounter) error {
		count = a.Count
		return nil
	})

	logger.Info(fmt.Sprintf("Core - tick number %d", count))
}

func cleanup(logger log.Logger) {
	logger.Info("Cleaning up core")
}
