package core

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/app"
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
	app.AddSchedule(schedule.PreUpdate)
	app.AddSchedule(schedule.Update)
	app.AddSchedule(schedule.PostUpdate)
	app.AddCleanupSchedule(schedule.Cleanup)

	tick.Init(&app)

	app.AddStartupSystem(schedule.Startup, startup)
	app.AddSystem(schedule.Update, printer)
	app.AddCleanupSystem(schedule.Cleanup, cleanup)

	return CoreApp{
		BasicSubApp: app,
	}
}

func startup(logger log.Logger) {
	logger.Info("Init core")
}

func printer(logger log.Logger, tickCounter *tick.Counter) {
	logger.Info(fmt.Sprintf("Core - tick number %d", tickCounter.Count))
}

func cleanup(logger log.Logger) {
	logger.Info("Cleaning up core")
}
