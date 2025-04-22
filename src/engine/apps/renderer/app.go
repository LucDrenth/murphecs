package renderer

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type RendererApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) RendererApp {
	app := app.NewBasicSubApp(logger)

	app.AddStartupSchedule(schedule.Startup)
	app.AddStartupSystem(schedule.Startup, startup)

	app.AddSchedule(schedule.Render)
	app.AddSystem(schedule.Render, printer)

	app.AddCleanupSchedule(schedule.Cleanup)
	app.AddCleanupSystem(schedule.Cleanup, cleanup)

	return RendererApp{
		BasicSubApp: app,
	}
}

// TODO can we pull in random system params here?
func startup(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Init renderer")
	return nil
}

// TODO can we pull in random system params here?
func printer(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Render ...")
	return nil
}

// TODO can we pull in random system params here?
func cleanup(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Cleaning up renderer")
	return nil
}
