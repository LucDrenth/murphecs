package core

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/engine/features"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
	"github.com/lucdrenth/murph_engine/src/log"
)

type CoreApp struct {
	app.BasicSubApp
}

func New(logger log.Logger) (CoreApp, error) {
	coreApp, err := app.NewBasicSubApp(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		return CoreApp{}, err
	}

	coreApp.SetDebugType("Core")

	coreApp.AddSchedule(schedule.Startup, app.ScheduleTypeStartup)
	coreApp.AddSchedule(schedule.PreUpdate, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.Update, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.PostUpdate, app.ScheduleTypeRepeating)
	coreApp.AddSchedule(schedule.Cleanup, app.ScheduleTypeCleanup)

	coreApp.AddFeature(&features.TickCounterFeature{Schedule: schedule.PreUpdate})
	// coreApp.AddFeature(&features.DebugPrinterFeature{
	// 	AppName:          "Core",
	// 	RepeatedSchedule: schedule.Update,
	// })

	return CoreApp{
		BasicSubApp: coreApp,
	}, nil
}
