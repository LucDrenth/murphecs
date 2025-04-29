package renderer

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
	"github.com/lucdrenth/murph_engine/src/tick"
)

type RendererApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) RendererApp {
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

	return RendererApp{
		BasicSubApp: app,
	}
}

func startup(logger log.Logger) {
	logger.Info("Init renderer")
}

func printer(logger log.Logger, tickCounter *tick.Counter) {
	logger.Info(fmt.Sprintf("Renderer - tick number %d", tickCounter.Count))
}

func cleanup(logger log.Logger) {
	logger.Info("Cleaning up renderer")
}
