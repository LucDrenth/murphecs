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
	coreApp := app.NewBasicSubApp(logger)
	coreApp.SetDebugType("Core")

	coreApp.AddSchedule(schedule.Startup, app.ScheduleTypeStartup)
	coreApp.AddSchedule(schedule.PreUpdate, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.Update, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.PostUpdate, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.Cleanup, app.ScheduleTypeCleanup)

	tick.Init(&coreApp, schedule.PreUpdate)

	coreApp.AddSystem(schedule.Startup, startup)
	coreApp.AddSystem(schedule.Update, printer)
	coreApp.AddSystem(schedule.Cleanup, cleanup)

	return CoreApp{
		BasicSubApp: coreApp,
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
