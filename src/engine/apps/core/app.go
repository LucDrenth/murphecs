package core

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type CoreApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) CoreApp {
	app := app.NewBasicSubApp(logger)

	app.AddStartupSchedule(schedule.Startup)
	app.AddStartupSystem(schedule.Startup, startup)

	app.AddSchedule(schedule.Update)
	app.AddSystem(schedule.Update, printer)

	app.AddCleanupSchedule(schedule.Cleanup)
	app.AddCleanupSystem(schedule.Cleanup, cleanup)

	return CoreApp{
		BasicSubApp: app,
	}
}

// TODO can we pull in random system params here?
func startup(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Init core")
	return nil
}

// TODO can we pull in random system params here?
func printer(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Running ...")
	return nil
}

// TODO can we pull in random system params here?
func cleanup(a app.SubApp, _ ...app.SystemParam) error {
	a.Logger().Info("Cleaning up core")
	return nil
}
