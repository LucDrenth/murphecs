package app

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/lucdrenth/murph_engine/src/ecs"
	"github.com/lucdrenth/murph_engine/src/log"
	"github.com/lucdrenth/murph_engine/src/utils"
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
	features  map[reflect.Type]IFeature
	logger    log.Logger
	debugType string
}

func NewBasicSubApp(logger log.Logger) BasicSubApp {
	resourceStorage := newResourceStorage()

	// The following resources are reserved by this app. Even if a user would add them, it
	// would not be possible to fetch them because the reserved resource would be returned
	// instead. Thus we register them as blacklisted so that an error is logged when the user
	// tries to add them.
	registerBlacklistedResource[*ecs.World](&resourceStorage)
	registerBlacklistedResourceType(reflect.TypeOf(logger), &resourceStorage)

	return BasicSubApp{
		world: ecs.NewWorld(),
		schedules: map[ScheduleType]*Scheduler{
			ScheduleTypeStartup:   utils.PointerTo(NewScheduler()),
			ScheduleTypeRepeating: utils.PointerTo(NewScheduler()),
			ScheduleTypeCleanup:   utils.PointerTo(NewScheduler()),
		},
		resources: resourceStorage,
		features:  map[reflect.Type]IFeature{},
		logger:    logger,
		debugType: "App",
	}
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
	featureType := reflect.TypeOf(feature)
	if _, exists := app.features[featureType]; exists {
		app.logger.Error(fmt.Sprintf("%s - failed to add feature %s: already added", app.debugType, featureType.String()))
		return app
	}

	initHasPointerReceiver, err := utils.MethodHasPointerReceiver(feature, "Init")
	if err != nil {
		app.logger.Error(fmt.Sprintf("%s - failed to add feature %s: failed to validate: %v", app.debugType, featureType.String(), err))
		return app
	}
	if !initHasPointerReceiver {
		app.logger.Error(fmt.Sprintf("%s - failed to add feature %s: Init must be pointer receiver", app.debugType, featureType.String()))
		return app
	}

	app.features[featureType] = feature
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

func (app *BasicSubApp) processFeatures() {
	for _, feature := range app.features {
		feature.Init()
	}

	for _, feature := range app.features {
		resources := feature.getResources()
		for i := range resources {
			app.AddResource(resources[i])
		}
	}

	for _, feature := range app.features {
		systems := feature.getSystems()
		for i := range systems {
			app.AddSystem(systems[i].schedule, systems[i].system)
		}
	}

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
