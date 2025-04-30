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
	rendererApp := app.NewBasicSubApp(logger)
	rendererApp.SetDebugType("Renderer")

	rendererApp.AddSchedule(schedule.Startup, app.ScheduleTypeStartup)
	rendererApp.AddSchedule(schedule.PreRender, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.Render, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.PostRender, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.Cleanup, app.ScheduleTypeCleanup)

	tick.Init(&rendererApp, schedule.PreRender)

	rendererApp.AddSystem(schedule.Startup, startup)
	rendererApp.AddSystem(schedule.Render, printer)
	rendererApp.AddSystem(schedule.Cleanup, cleanup)

	return RendererApp{
		BasicSubApp: rendererApp,
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
