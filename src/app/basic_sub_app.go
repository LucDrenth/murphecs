package app

import (
	"fmt"
	"slices"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/lucdrenth/murphecs/src/utils"
)

type ScheduleType int

const (
	ScheduleTypeStartup   ScheduleType = iota // run only once, on startup
	ScheduleTypeRepeating                     // runs repeatedly, in the main loop
	ScheduleTypeCleanup                       // runs only once, before quitting
)

type BasicSubApp struct {
	world     ecs.World
	schedules map[ScheduleType]*Scheduler
	resources resourceStorage // resources that can be pulled by system params.
	features  *Feature        // this 'master' feature is empty except for its nested features
	logger    Logger
	debugType string
	tickRate  time.Duration
	lastDelta float64 // delta time of the last tick
}

func NewBasicSubApp(logger Logger, worldConfigs ecs.WorldConfigs) (BasicSubApp, error) {
	if logger == nil {
		noOpLogger := NoOpLogger{}
		logger = &noOpLogger
	}

	world, err := ecs.NewWorld(worldConfigs)
	if err != nil {
		return BasicSubApp{}, fmt.Errorf("failed to create world: %w", err)
	}

	resourceStorage := newResourceStorage()

	// The following resources are reserved by this app. Even if a user would add them, it
	// would not be possible to fetch them because the reserved resource would be returned
	// instead. Thus we register them as blacklisted so that an error is logged when the user
	// tries to add them.
	registerBlacklistedResource[*ecs.World](&resourceStorage)

	return BasicSubApp{
		world: world,
		schedules: map[ScheduleType]*Scheduler{
			ScheduleTypeStartup:   utils.PointerTo(NewScheduler()),
			ScheduleTypeRepeating: utils.PointerTo(NewScheduler()),
			ScheduleTypeCleanup:   utils.PointerTo(NewScheduler()),
		},
		resources: resourceStorage,
		features:  &Feature{},
		logger:    logger,
		debugType: "App",
		tickRate:  time.Second / 60.0,
	}, nil
}

func (app *BasicSubApp) AddSystem(schedule Schedule, system System) SubApp {
	for _, scheduler := range app.schedules {
		if slices.Contains(scheduler.order, schedule) {
			err := scheduler.AddSystem(schedule, system, &app.world, app.logger, &app.resources)
			if err != nil {
				app.logger.Error(fmt.Sprintf("%s - failed to add system %s: %v",
					app.debugType,
					systemToDebugString(system),
					err,
				))
			}

			return app
		}
	}

	app.logger.Error(fmt.Sprintf("%s - failed to add system %s: schedule %s not found",
		app.debugType,
		systemToDebugString(system),
		schedule,
	))
	return app
}

func (app *BasicSubApp) AddSchedule(schedule Schedule, scheduleType ScheduleType) SubApp {
	scheduler, ok := app.schedules[scheduleType]
	if !ok {
		app.logger.Error(fmt.Sprintf("%s - failed to add schedule %s: invalid schedule type", app.debugType, schedule))
		return app
	}

	err := scheduler.AddSchedule(schedule)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to add schedule %s: %v", app.debugType, schedule, err))
	}

	return app
}

func (app *BasicSubApp) AddResource(resource Resource) SubApp {
	err := app.resources.add(resource)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to add resource %s: %v", app.debugType, getResourceDebugType(resource), err))
	}

	return app
}

func (app *BasicSubApp) AddFeature(feature IFeature) SubApp {
	app.features.AddFeature(feature)
	return app
}

func (app *BasicSubApp) Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool) {
	app.processFeatures()

	startupSystems, err := app.schedules[ScheduleTypeStartup].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get startup systems: %v", app.debugType, err))
		return
	}

	repeatedSystems, err := app.schedules[ScheduleTypeRepeating].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get repeated systems: %v", app.debugType, err))
		return
	}

	cleanupSystems, err := app.schedules[ScheduleTypeCleanup].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get cleanup systems: %v", app.debugType, err))
		return
	}

	app.runSystemSet(startupSystems)
	app.runRepeatedUntilExit(exitChannel, repeatedSystems)
	app.runSystemSet(cleanupSystems)
	isDoneChannel <- true
}

// processFeatures adds the resources and systems from all features
func (app *BasicSubApp) processFeatures() {
	features := app.features.GetFeatures()
	validatedFeatures := make([]IFeature, 0, len(features))
	for _, feature := range features {
		err := validateFeature(feature)
		if err != nil {
			app.logger.Error(fmt.Sprintf("%s - %v", app.debugType, err))
			continue
		}

		validatedFeatures = append(validatedFeatures, feature)
	}

	for _, feature := range validatedFeatures {
		resources := feature.GetResources()
		for i := range resources {
			app.AddResource(resources[i])
		}
	}

	for _, feature := range validatedFeatures {
		systems := feature.GetSystems()
		for i := range systems {
			app.AddSystem(systems[i].schedule, systems[i].system)
		}
	}

	// free up resources
	app.features = nil
}

func (app *BasicSubApp) runSystemSet(systemSets []*SystemSet) {
	for _, systemSet := range systemSets {
		errors := systemSet.exec(&app.world)
		for _, err := range errors {
			app.logger.Error(fmt.Sprintf("%s - system returned error: %v", app.debugType, err))
		}
	}
}

func (app *BasicSubApp) SetDebugType(debugType string) {
	app.debugType = debugType
}

func (app *BasicSubApp) SetTickRate(tickRate time.Duration) {
	if tickRate == 0 {
		app.logger.Error(fmt.Sprintf("%s - failed to set tickRate: can not be zero", app.debugType))
		return
	}
	app.tickRate = tickRate
}

func (app *BasicSubApp) Delta() float64 {
	return app.lastDelta
}

func (app *BasicSubApp) NumberOfResources() uint {
	return uint(len(app.resources.resources))
}

func (app *BasicSubApp) NumberOfSystems() uint {
	result := uint(0)

	for _, schedules := range app.schedules {
		result += schedules.NumberOfSystems()
	}
	return result
}

func (app *BasicSubApp) runRepeatedUntilExit(exitChannel <-chan struct{}, systems []*SystemSet) {
	ticker := time.NewTicker(app.tickRate)
	currentTickRate := app.tickRate
	var now int64
	var delta float64
	start := time.Now().UnixNano()

	for {
		select {
		case <-exitChannel:
			return

		case <-ticker.C:
			now = time.Now().UnixNano()
			delta = float64(now-start) / 1_000_000_000
			start = now

			app.lastDelta = delta
			app.runSystemSet(systems)

			if currentTickRate != app.tickRate {
				app.runRepeatedUntilExit(exitChannel, systems)
				return
			}
		}
	}
}
