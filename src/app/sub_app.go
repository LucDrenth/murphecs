package app

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
	"github.com/lucdrenth/murphecs/src/utils"
)

var (
	SystemErrorPackageDepth = ecs.SystemErrorPackageDepth
)

// SubApp has startup systems, repeating systems and cleanup systems.
//
// The repeating systems run at a fixed time. If running the systems takes longer then the tickRate, missed ticks
// will not be repeated.
// For example: if the tickRate is 1 second and a tick suddenly takes 4 seconds, the next tick will be run immediately
// after, and then after 1 second.
type SubApp struct {
	world         *ecs.World
	scheduleTypes map[scheduleType][]ecs.Schedule
	logger        Logger
	Name          string
	tickRate      *time.Duration // the rate at which the repeating systems run
	currentTick   uint
	lastDelta     *float64 // delta time of the last tick
	runner        Runner
	features      []IFeature // this slice will be processed and emptied when starting this SubApp

	OnStartupSchedulesDone func()

	startupExecutor  Executor
	repeatedExecutor Executor
	cleanupExecutor  Executor
}

func New(logger Logger, worldConfigs ecs.WorldConfigs) (*SubApp, error) {
	if logger == nil {
		noOpLogger := ecs.NoOpLogger{}
		logger = &noOpLogger
	}

	worldConfigs.Logger = logger

	world, err := ecs.NewWorld(worldConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to create world: %w", err)
	}

	// The following resources are reserved by this app. If a user would add them, systems would
	// use the reserved resource instead of the inserted resource, which could cause confusion. So
	// we register them as blacklisted so that an error is logged when the user tries to add them.
	err = ecs.RegisterBlacklistedResource[*ecs.World](world.Resources())
	if err != nil {
		logger.Warn("App - failed to register blacklisted resource: %v", err)
	}

	subApp := SubApp{
		world: &world,
		scheduleTypes: map[scheduleType][]ecs.Schedule{
			ScheduleTypeStartup:   {},
			ScheduleTypeRepeating: {},
			ScheduleTypeCleanup:   {},
		},
		logger:           logger,
		Name:             "App",
		tickRate:         new(time.Second / 60.0),
		lastDelta:        new(0.0),
		startupExecutor:  &ConsecutiveExecutor{},
		repeatedExecutor: &ConsecutiveExecutor{},
		cleanupExecutor:  &ConsecutiveExecutor{},
	}
	subApp.UseFixedRunner()

	return &subApp, nil
}

// AddSystem adds a system that will be run when the schedule is run. Systems must be a function.
func (app *SubApp) AddSystem(schedule ecs.Schedule, system ecs.System) *SubApp {
	return app.addSystemWithSource(schedule, system, utils.Caller(2, SystemErrorPackageDepth))
}

// addSystemWithSource adds a system with an explicit source path for error messages.
func (app *SubApp) addSystemWithSource(schedule ecs.Schedule, system ecs.System, source string) *SubApp {
	err := app.world.AddSystemWithSource(schedule, system, source)
	if err != nil {
		app.logger.Error("%s - failed to add system: %v", app.Name, err)
	}

	return app
}

// AddSchedule adds a schedule that systems can be added to.
func (app *SubApp) AddSchedule(schedule ecs.Schedule, options ScheduleOptions) *SubApp {
	if _, ok := app.scheduleTypes[options.ScheduleType]; !ok {
		app.logger.Error("%s - failed to add schedule %s: %v: %d", app.Name, schedule, ErrScheduleTypeNotFound, options.ScheduleType)
		return app
	}

	if options.Order == nil {
		options.Order = ecs.ScheduleLast{}
	}

	err := app.world.AddSchedule(schedule, options.Order, options.IsPaused)
	if err != nil {
		app.logger.Error("%s - failed to add schedule %s: %v", app.Name, schedule, err)
		return app
	}

	app.scheduleTypes[options.ScheduleType] = append(app.scheduleTypes[options.ScheduleType], schedule)
	return app
}

func (app *SubApp) SetSchedulePaused(schedule ecs.Schedule, isPaused bool) error {
	return app.world.SetSchedulePaused(schedule, isPaused)
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
func (app *SubApp) AddResource(resource ecs.Resource) *SubApp {
	err := app.world.Resources().Add(resource)
	if err != nil {
		app.logger.Error("%s - failed to add resource %s: %v", app.Name, ecs.GetResourceDebugType(resource), err)
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
	err := app.PrepareForRun()
	if err != nil {
		app.logger.Error("%s - prepare failed: %v", app.Name, err)
		return
	}

	app.RunStartupSchedules(exitChannel)
	app.RunRepeatedSchedules(exitChannel)
	app.RunCleanupSchedules(exitChannel)

	isDoneChannel <- true
}

func (app *SubApp) PrepareForRun() error {
	app.ProcessFeatures()

	err := app.world.PrepareSystems()
	if err != nil {
		return fmt.Errorf("prepare systems failed: %w", err)
	}

	err = app.prepareExecutors()
	if err != nil {
		app.logger.Error("%s - prepare executors: %v", app.Name, err)
		return fmt.Errorf("prepare executors failed: %w", err)
	}

	return nil
}

func (app *SubApp) RunStartupSchedules(exitChannel <-chan struct{}) {
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
}

func (app *SubApp) RunRepeatedSchedules(exitChannel <-chan struct{}) {
	app.runner.Run(exitChannel, app.repeatedExecutor)
}

func (app *SubApp) RunCleanupSchedules(exitChannel <-chan struct{}) {
	onceRunner := app.newNTimesRunner(1)
	onceRunner.Run(exitChannel, app.cleanupExecutor)
}

func (app *SubApp) prepareExecutors() error {
	startupSystems, err := app.world.GetScheduleSystemsBySchedules(app.scheduleTypes[ScheduleTypeStartup])
	if err != nil {
		return fmt.Errorf("failed to get startup systems: %v", err)
	}
	app.startupExecutor.Load(startupSystems, app.world, app.logger, app.Name)

	repeatedSystems, err := app.world.GetScheduleSystemsBySchedules(app.scheduleTypes[ScheduleTypeRepeating])
	if err != nil {
		return fmt.Errorf("failed to get repeated systems: %v", err)
	}
	app.repeatedExecutor.Load(repeatedSystems, app.world, app.logger, app.Name)

	cleanupSystems, err := app.world.GetScheduleSystemsBySchedules(app.scheduleTypes[ScheduleTypeCleanup])
	if err != nil {
		return fmt.Errorf("failed to get cleanup systems: %v", err)
	}
	app.cleanupExecutor.Load(cleanupSystems, app.world, app.logger, app.Name)

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
	return app.world.Resources().Count()
}

func (app *SubApp) NumberOfSystems() uint {
	return app.world.NumberOfSystems()
}

func (app *SubApp) NumberOfSchedules() uint {
	return app.world.NumberOfSchedules()
}

func (app *SubApp) World() *ecs.World {
	return app.world
}

func (app *SubApp) OuterWorlds() *map[ecs.WorldId]*ecs.World {
	return app.world.OuterWorlds()
}

// RegisterOuterWorld lets you use the outer world in system param queries.
func (app *SubApp) RegisterOuterWorld(id ecs.WorldId, world *ecs.World) error {
	return app.world.RegisterOuterWorld(id, world)
}

// SetRunner sets the runner for the repeated systems
func (app *SubApp) SetRunner(runner Runner) {
	if runner == nil {
		app.logger.Error("%s - failed to set runner: can not be nil", app.Name)
		return
	}

	app.runner = runner
}

// UseFixedRunner makes the systems run repeatedly, at a fixed interval. To control the interval time, use `app.SetTickRate`.
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

func (app *SubApp) UseOnceRunner() {
	app.runner = &onceRunner{RunnerBasis: NewRunnerBasis(app)}
}

// newNTimesRunner creates a runner that runs systems [numberOfRuns] amount of times
func (app *SubApp) newNTimesRunner(numberOfRuns int) nTimesRunner {
	return nTimesRunner{
		numberOfRuns: numberOfRuns,
		RunnerBasis:  NewRunnerBasis(app),
	}
}
