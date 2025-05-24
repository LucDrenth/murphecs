package app

import (
	"fmt"
	"slices"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/lucdrenth/murphecs/src/utils"
)

type scheduleType int

const (
	ScheduleTypeStartup   scheduleType = iota // run only once, on startup
	ScheduleTypeRepeating                     // runs repeatedly, in the main loop
	ScheduleTypeCleanup                       // runs only once, before quitting
)

// SubApp has startup systems, repeating systems and cleanup systems.
//
// The repeating systems run at a fixed time. If running the systems takes longer then the tickRate, missed ticks
// will not be repeated.
// For example: if the tickRate is 1 second and a tick suddenly takes 4 seconds, the next tick will be run immediately
// after, and then after 1 second.
type SubApp struct {
	world       ecs.World
	schedules   map[scheduleType]*Scheduler
	resources   resourceStorage // resources that can be pulled by system params.
	logger      Logger
	name        string
	tickRate    *time.Duration // the rate at which the repeating systems run
	lastDelta   *float64       // delta time of the last tick
	runner      Runner
	outerWorlds map[ecs.WorldId]*ecs.World
}

func New(logger Logger, worldConfigs ecs.WorldConfigs) (SubApp, error) {
	if logger == nil {
		noOpLogger := NoOpLogger{}
		logger = &noOpLogger
	}

	world, err := ecs.NewWorld(worldConfigs)
	if err != nil {
		return SubApp{}, fmt.Errorf("failed to create world: %w", err)
	}

	resourceStorage := newResourceStorage()

	// The following resources are reserved by this app. Even if a user would add them, it
	// would not be possible to fetch them because the reserved resource would be returned
	// instead. Thus we register them as blacklisted so that an error is logged when the user
	// tries to add them.
	registerBlacklistedResource[*ecs.World](&resourceStorage)

	subApp := SubApp{
		world: world,
		schedules: map[scheduleType]*Scheduler{
			ScheduleTypeStartup:   utils.PointerTo(NewScheduler()),
			ScheduleTypeRepeating: utils.PointerTo(NewScheduler()),
			ScheduleTypeCleanup:   utils.PointerTo(NewScheduler()),
		},
		resources:   resourceStorage,
		logger:      logger,
		name:        "App",
		tickRate:    utils.PointerTo(time.Second / 60.0),
		lastDelta:   utils.PointerTo(0.0),
		outerWorlds: map[ecs.WorldId]*ecs.World{},
	}
	subApp.SetFixedRunner()

	return subApp, nil
}

func (app *SubApp) AddSystem(schedule Schedule, system System) *SubApp {
	for _, scheduler := range app.schedules {
		if slices.Contains(scheduler.order, schedule) {
			err := scheduler.AddSystem(schedule, system, &app.world, &app.outerWorlds, app.logger, &app.resources)
			if err != nil {
				app.logger.Error(fmt.Sprintf("%s - failed to add system %s: %v",
					app.name,
					systemToDebugString(system),
					err,
				))
			}

			return app
		}
	}

	app.logger.Error(fmt.Sprintf("%s - failed to add system %s: schedule %s not found",
		app.name,
		systemToDebugString(system),
		schedule,
	))
	return app
}

func (app *SubApp) AddSchedule(schedule Schedule, scheduleType scheduleType) *SubApp {
	scheduler, ok := app.schedules[scheduleType]
	if !ok {
		app.logger.Error(fmt.Sprintf("%s - failed to add schedule %s: invalid schedule type", app.name, schedule))
		return app
	}

	err := scheduler.AddSchedule(schedule)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to add schedule %s: %v", app.name, schedule, err))
	}

	return app
}

func (app *SubApp) AddResource(resource Resource) *SubApp {
	err := app.resources.add(resource)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to add resource %s: %v", app.name, getResourceDebugType(resource), err))
	}

	return app
}

func (app *SubApp) AddFeature(feature IFeature) *SubApp {
	feature.Init()
	features := feature.GetAndInitNestedFeatures()
	features = append(features, feature)

	validatedFeatures := make([]IFeature, 0, len(features))
	for _, feature := range features {
		err := validateFeature(feature)
		if err != nil {
			app.logger.Error(fmt.Sprintf("%s - %v", app.name, err))
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

	return app
}

func (app *SubApp) Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool) {
	startupSystems, err := app.schedules[ScheduleTypeStartup].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get startup systems: %v", app.name, err))
		return
	}

	repeatedSystems, err := app.schedules[ScheduleTypeRepeating].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get repeated systems: %v", app.name, err))
		return
	}

	cleanupSystems, err := app.schedules[ScheduleTypeCleanup].GetSystemSets()
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to get cleanup systems: %v", app.name, err))
		return
	}

	onceRunner := onceRunner{
		world:       &app.world,
		outerWorlds: &app.outerWorlds,
		logger:      app.logger,
		appName:     app.name,
	}

	onceRunner.Run(exitChannel, startupSystems)
	app.runner.Run(exitChannel, repeatedSystems)
	onceRunner.Run(exitChannel, cleanupSystems)
	isDoneChannel <- true
}

func (app *SubApp) SetName(name string) {
	app.name = name
}

// SetTickRate sets the interval at which the repeated systems are run. This can be safely changed while
// the app is already running, in which case it will be picked up after the next run.
func (app *SubApp) SetTickRate(tickRate time.Duration) {
	if tickRate == 0 {
		app.logger.Error(fmt.Sprintf("%s - failed to set tickRate: can not be zero", app.name))
		return
	}
	*app.tickRate = tickRate
}

func (app *SubApp) Delta() float64 {
	return *app.lastDelta
}

func (app *SubApp) NumberOfResources() uint {
	return uint(len(app.resources.resources))
}

func (app *SubApp) NumberOfSystems() uint {
	result := uint(0)

	for _, schedules := range app.schedules {
		result += schedules.NumberOfSystems()
	}
	return result
}

func (app *SubApp) World() *ecs.World {
	return &app.world
}

func (app *SubApp) OuterWorlds() *map[ecs.WorldId]*ecs.World {
	return &app.outerWorlds
}

func (app *SubApp) RegisterOuterWorld(id ecs.WorldId, world *ecs.World) error {
	if _, exists := app.outerWorlds[id]; exists {
		return fmt.Errorf("id %d is already registered", id)
	}

	app.outerWorlds[id] = world
	return nil
}

// SetRunner sets the runner for the repeated systems
func (app *SubApp) SetRunner(runner Runner) {
	if runner == nil {
		app.logger.Error(fmt.Sprintf("%s - failed to set runner: can not be nil", app.name))
		return
	}

	app.runner = runner
}

// SetFixedRunner sets the default fixedRunner, which runs systems at a fixed interval. To control the interval time,
// use `app.SetTickRate`.
func (app *SubApp) SetFixedRunner() {
	app.runner = &fixedRunner{
		tickRate:    app.tickRate,
		delta:       app.lastDelta,
		world:       &app.world,
		outerWorlds: &app.outerWorlds,
		logger:      app.logger,
		appName:     app.name,
	}
}
