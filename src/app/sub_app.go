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

var (
	SystemErrorPackageDepth = 3
)

// SubApp has startup systems, repeating systems and cleanup systems.
//
// The repeating systems run at a fixed time. If running the systems takes longer then the tickRate, missed ticks
// will not be repeated.
// For example: if the tickRate is 1 second and a tick suddenly takes 4 seconds, the next tick will be run immediately
// after, and then after 1 second.
type SubApp struct {
	world                    *ecs.World
	schedules                map[scheduleType]*Scheduler
	resources                resourceStorage // resources that can be pulled by system params.
	logger                   Logger
	Name                     string
	tickRate                 *time.Duration // the rate at which the repeating systems run
	currentTick              uint
	lastDelta                *float64 // delta time of the last tick
	runner                   Runner
	outerWorlds              map[ecs.WorldId]*ecs.World
	EventStorage             *EventStorage
	scheduleSystemsIdCounter ScheduleSystemsId
	OnStartupSchedulesDone   func()
	features                 []IFeature // this slice will be process and emptied when starting this SubApp

	startupExecutor  Executor
	repeatedExecutor Executor
	cleanupExecutor  Executor
}

func New(logger Logger, worldConfigs ecs.WorldConfigs) (*SubApp, error) {
	if logger == nil {
		noOpLogger := NoOpLogger{}
		logger = &noOpLogger
	}

	world, err := ecs.NewWorld(worldConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to create world: %w", err)
	}

	resourceStorage := newResourceStorage()

	// The following resources are reserved by this app. If a user would add them, systems would
	// use the reserved resource instead if the inserted resource, which could cause confusion. So
	// we register them as blacklisted so that an error is logged when the user tries to add them.
	err = registerBlacklistedResource[*ecs.World](&resourceStorage)
	if err != nil {
		logger.Warn("App - failed to register blacklisted resource: %v", err)
	}

	subApp := SubApp{
		world: &world,
		schedules: map[scheduleType]*Scheduler{
			ScheduleTypeStartup:   new(NewScheduler()),
			ScheduleTypeRepeating: new(NewScheduler()),
			ScheduleTypeCleanup:   new(NewScheduler()),
		},
		resources:        resourceStorage,
		logger:           logger,
		Name:             "App",
		tickRate:         new(time.Second / 60.0),
		lastDelta:        new(0.0),
		outerWorlds:      map[ecs.WorldId]*ecs.World{},
		EventStorage:     new(newEventStorage()),
		startupExecutor:  &ConsecutiveExecutor{},
		repeatedExecutor: &ConsecutiveExecutor{},
		cleanupExecutor:  &ConsecutiveExecutor{},
	}
	subApp.UseFixedRunner()

	return &subApp, nil
}

// AddSystem adds a system that will be run when the schedule is run. Systems must be a function.
func (app *SubApp) AddSystem(schedule Schedule, system System) *SubApp {
	return app.addSystemWithSource(schedule, system, utils.Caller(2, SystemErrorPackageDepth))
}

// AddSystem adds a system that will be run when the schedule is run. Systems must be a function.
func (app *SubApp) addSystemWithSource(schedule Schedule, system System, source string) *SubApp {
	scheduler := app.getScheduler(schedule)
	if scheduler == nil {
		app.logger.Error("%s - failed to add system: schedule %s not found",
			app.Name,
			schedule,
		)
		return app
	}

	err := scheduler.AddSystem(schedule, system, source, app.world, &app.outerWorlds, app.logger, &app.resources, app.EventStorage)
	if err != nil {
		app.logger.Error("%s - failed to add system: %v",
			app.Name,
			err,
		)
	}

	return app
}

func (app *SubApp) getScheduler(schedule Schedule) *Scheduler {
	for _, scheduler := range app.schedules {
		if slices.Contains(scheduler.order, schedule) {
			return scheduler
		}
	}

	return nil
}

type ScheduleOptions struct {
	// ScheduleType can be one of:
	//   - [ScheduleTypeStartup] - systems in a schedule with this schedule type run once, when starting the app
	//   - [ScheduleTypeRepeating] - systems in a schedule with this schedule type run repeatedly, after startup
	//   - [ScheduleTypeCleanup] - systems in a schedule with this schedule type run once, when closing the app
	ScheduleType scheduleType

	// Order decides when the schedule systems should run, relative to schedules. It can be one of:
	//	- [ScheduleLast] - this is also the default if this field is left nil
	// 	- [ScheduleBefore]
	// 	- [ScheduleAfter]
	Order ScheduleOrder

	// IsPaused determines the initial pause state of the schedule
	IsPaused bool
}

// AddSchedule adds a schedule that systems can be added to.
func (app *SubApp) AddSchedule(schedule Schedule, options ScheduleOptions) *SubApp {
	scheduler, ok := app.schedules[options.ScheduleType]
	if !ok {
		app.logger.Error("%s - failed to add schedule %s: %w: %d", app.Name, schedule, ErrScheduleTypeNotFound, options.ScheduleType)
		return app
	}

	if options.Order == nil {
		options.Order = ScheduleLast{}
	}

	app.scheduleSystemsIdCounter++
	err := scheduler.AddSchedule(schedule, app.scheduleSystemsIdCounter, options.Order, options.IsPaused)
	if err != nil {
		app.logger.Error("%s - failed to add schedule %s: %v", app.Name, schedule, err)
	}

	return app
}

func (app *SubApp) SetSchedulePaused(schedule Schedule, scheduleType scheduleType, isPaused bool) error {
	scheduler, exists := app.schedules[scheduleType]
	if !exists {
		return fmt.Errorf("%w: %d", ErrScheduleTypeNotFound, scheduleType)
	}

	systems, exists := scheduler.systems[schedule]
	if !exists {
		return fmt.Errorf("%w: %s", ErrScheduleNotFound, schedule)
	}

	currentlyPaused := systems.isPaused.Load()
	if currentlyPaused == isPaused {
		return nil
	}

	if !currentlyPaused && isPaused {
		systems.isFirstExecSincePaused = true
	}
	systems.isPaused.Store(isPaused)

	return nil
}

// AddResource adds a resource that can then be used in system params. There can only be 1 one of each resource type.
//
// Struct resources must be passed by reference.
//   - You can use it in a system param as a pointer to get a reference
//   - You can use it without a pointer to get a copy of the current resource value.
//
// Interface resources can be passed by either reference or by value.
//   - If you pass it by reference, you must use it in a system param by the interface type, but not as a pointer.
//     You can then use it by reference, as is normally the case with interface values.
//   - If you pass it by value, you must use it in a system param by its struct implementation. You can then use that
//     as a regular struct resource (see above for explanation).
func (app *SubApp) AddResource(resource Resource) *SubApp {
	err := app.resources.add(resource)
	if err != nil {
		app.logger.Error("%s - failed to add resource %s: %v", app.Name, getResourceDebugType(resource), err)
	}

	return app
}

func (app *SubApp) AddFeature(feature IFeature) *SubApp {
	app.features = append(app.features, feature)
	return app
}

func (app *SubApp) ProcessFeatures() {
	validatedFeatures := []IFeature{}

	for _, feature := range app.features {
		feature.Init()
		features := []IFeature{feature}
		features = append(features, feature.GetAndInitNestedFeatures()...)

		for _, feature := range features {
			err := validateFeature(feature)
			if err != nil {
				app.logger.Error("%s - %v", app.Name, err)
				continue
			}

			validatedFeatures = append(validatedFeatures, feature)
		}
	}

	// add all resources before the systems so that when adding the systems, resources
	// as system params can be validated.
	for _, feature := range validatedFeatures {
		resources := feature.GetResources()
		for i := range resources {
			app.AddResource(resources[i])
		}
	}

	for _, feature := range validatedFeatures {
		systems := feature.GetSystems()
		for i := range systems {
			app.addSystemWithSource(systems[i].schedule, systems[i].system, systems[i].source)
		}
	}

	app.features = []IFeature{}
}

func (app *SubApp) Run(exitChannel <-chan struct{}, isDoneChannel chan<- bool) {
	app.ProcessFeatures()

	err := app.prepareExecutors()
	if err != nil {
		app.logger.Error("%s - %v", err)
		return
	}

	onceRunner := app.newNTimesRunner(1)
	onceRunner.Run(exitChannel, app.startupExecutor)
	if app.OnStartupSchedulesDone != nil {
		app.OnStartupSchedulesDone()
	}

	app.runner.setOnFirstRunDone(func() {
		// Events written by EventWriter's in startup systems do not get cleared by default so
		// that they can be read by the repeated schedules.
		app.startupExecutor.ProcessEvents(app.currentTick)
	})
	app.runner.setOnRunDone(func() {
		app.currentTick++
	})

	app.runner.Run(exitChannel, app.repeatedExecutor)
	onceRunner.Run(exitChannel, app.cleanupExecutor)
	isDoneChannel <- true
}

func (app *SubApp) prepareExecutors() error {
	startupSystems, err := app.schedules[ScheduleTypeStartup].GetScheduleSystems()
	if err != nil {
		return fmt.Errorf("failed to get startup systems: %v", err)
	}
	app.startupExecutor.Load(startupSystems, app.world, &app.outerWorlds, app.logger, app.Name, app.EventStorage)

	repeatedSystems, err := app.schedules[ScheduleTypeRepeating].GetScheduleSystems()
	if err != nil {
		return fmt.Errorf("failed to get repeated systems: %v", err)
	}
	app.repeatedExecutor.Load(repeatedSystems, app.world, &app.outerWorlds, app.logger, app.Name, app.EventStorage)

	cleanupSystems, err := app.schedules[ScheduleTypeCleanup].GetScheduleSystems()
	if err != nil {
		return fmt.Errorf("failed to get cleanup systems: %v", err)
	}
	app.cleanupExecutor.Load(cleanupSystems, app.world, &app.outerWorlds, app.logger, app.Name, app.EventStorage)

	return nil
}

// SetTickRate sets the interval at which the repeated systems are run. This can be safely changed while
// the app is already running, in which case it will be picked up after the next run.
func (app *SubApp) SetTickRate(tickRate time.Duration) {
	*app.tickRate = tickRate
}

func (app *SubApp) GetCurrentTick() *uint {
	return &app.currentTick
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

func (app *SubApp) NumberOfSchedules() uint {
	result := uint(0)

	for _, schedules := range app.schedules {
		result += uint(len(schedules.systems))
	}

	return result
}

func (app *SubApp) World() *ecs.World {
	return app.world
}

func (app *SubApp) OuterWorlds() *map[ecs.WorldId]*ecs.World {
	return &app.outerWorlds
}

// RegisterOuterWorld lets you use the outer world in system param queries.
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
		app.logger.Error("%s - failed to set runner: can not be nil", app.Name)
		return
	}

	app.runner = runner
}

// SetFixedRunner makes the systems run repeatedly, at a fixed interval. To control the interval time, use `app.SetTickRate`.
func (app *SubApp) UseFixedRunner() {
	app.runner = &fixedRunner{
		tickRate:    app.tickRate,
		RunnerBasis: NewRunnerBasis(app),
	}
}

// UseUncappedRunner makes the systems run repeatedly, as frequent possible
func (app *SubApp) UseUncappedRunner() {
	app.runner = &uncappedRunner{
		RunnerBasis: NewRunnerBasis(app),
	}
}

func (app *SubApp) UseNTimesRunner(numberOfRuns int) {
	runner := app.newNTimesRunner(numberOfRuns)
	app.runner = &runner
}

// NewNTimesRunner creates a runner that runs systems [numberOfRuns] amount of  times
func (app *SubApp) newNTimesRunner(numberOfRuns int) nTimesRunner {
	return nTimesRunner{
		numberOfRuns: numberOfRuns,
		RunnerBasis:  NewRunnerBasis(app),
	}
}
