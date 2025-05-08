package renderer

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/engine/features"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type RendererApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) (RendererApp, error) {
	rendererApp, err := app.NewBasicSubApp(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		return RendererApp{}, err
	}

	rendererApp.SetDebugType("Renderer")

	rendererApp.AddSchedule(schedule.Startup, app.ScheduleTypeStartup)
	rendererApp.AddSchedule(schedule.PreRender, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.Render, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.PostRender, app.ScheduleTypeRepeating)
	rendererApp.AddSchedule(schedule.Cleanup, app.ScheduleTypeCleanup)

	rendererApp.AddFeature(&features.TickCounterFeature{Schedule: schedule.PreRender})
	// rendererApp.AddFeature(&features.DebugPrinterFeature{
	// 	AppName:          "Renderer",
	// 	RepeatedSchedule: schedule.Render,
	// })

	return RendererApp{
		BasicSubApp: rendererApp,
	}, nil
}
