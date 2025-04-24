package app

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/log"
)

type SubApp interface {
	Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool)
	Logger() log.Logger
}

type BasicSubApp struct {
	world             ecs.World
	startupSchedules  Scheduler // these systems only run once, on startup
	repeatedSchedules Scheduler // these systems run in the main loop
	cleanupSchedules  Scheduler // these systems only run once, before quitting
	logger            log.Logger
}

func NewBasicSubApp(logger log.Logger) BasicSubApp {
	return BasicSubApp{
		world:             ecs.NewWorld(),
		startupSchedules:  NewScheduler(),
		repeatedSchedules: NewScheduler(),
		cleanupSchedules:  NewScheduler(),
		logger:            logger,
	}
}

func (app *BasicSubApp) AddStartupSystem(schedule Schedule, system System) {
	err := app.startupSchedules.AddSystem(schedule, system, &app.world, app.logger)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add startup system: %v", err))
	}
}

func (app *BasicSubApp) AddStartupSchedule(schedule Schedule) {
	err := app.startupSchedules.AddSchedule(schedule)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add startup schedule %s", schedule))
	}
}

func (app *BasicSubApp) AddSystem(schedule Schedule, system System) {
	err := app.repeatedSchedules.AddSystem(schedule, system, &app.world, app.logger)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add system: %v", err))
	}
}

func (app *BasicSubApp) AddSchedule(schedule Schedule) {
	err := app.repeatedSchedules.AddSchedule(schedule)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add schedule %s", schedule))
	}
}

func (app *BasicSubApp) AddCleanupSystem(schedule Schedule, system System) {
	err := app.cleanupSchedules.AddSystem(schedule, system, &app.world, app.logger)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add cleanup system: %v", err))
	}
}

func (app *BasicSubApp) AddCleanupSchedule(schedule Schedule) {
	err := app.cleanupSchedules.AddSchedule(schedule)
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to add cleanup schedule %s", schedule))
	}
}

func (app *BasicSubApp) Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool) {
	startupSystems, err := app.startupSchedules.GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get startup systems: %v", err))
		return
	}

	repeatedSystems, err := app.repeatedSchedules.GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get repeated systems: %v", err))
		return
	}

	cleanupSystems, err := app.cleanupSchedules.GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("failed to get cleanup systems: %v", err))
		return
	}

	app.runSystemSet(startupSystems)

	exit := false

	// TODO a timed loop that runs x times a second
	for {
		select {
		case <-exitChannel:
			exit = true
		default:
			exit = false
		}

		if exit {
			break
		}

		app.runSystemSet(repeatedSystems)

		time.Sleep(time.Second * 1)
	}

	app.runSystemSet(cleanupSystems)
	isDoneChannel <- true
}

func (app *BasicSubApp) Logger() log.Logger {
	return app.logger
}

func (app *BasicSubApp) runSystemSet(systemSets []*SystemSet) {
	for _, systemSet := range systemSets {
		for i := range systemSet.systems {
			systemSet.systems[i].exec(app.logger)
		}
	}
}
