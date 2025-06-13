package app

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Runner interface {
	Run(exitChannel <-chan struct{}, systems []*SystemSet)
	SetStartupSystemSetIds([]SystemSetId)
}

// fixedRunner runs systems at a fixed interval
type fixedRunner struct {
	tickRate            *time.Duration
	delta               *float64
	world               *ecs.World
	outerWorlds         *map[ecs.WorldId]*ecs.World
	logger              Logger
	appName             string
	eventStorage        *EventStorage
	startupSystemSetIds []SystemSetId
}

func (runner *fixedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	ticker := time.NewTicker(*runner.tickRate)
	currentTickRate := *runner.tickRate
	var now int64
	start := time.Now().UnixNano()
	isFirstRun := true

	for {
		select {
		case <-exitChannel:
			return

		case <-ticker.C:
			now = time.Now().UnixNano()
			*runner.delta = float64(now-start) / 1_000_000_000.0
			start = now

			runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage)

			if isFirstRun {
				// clear events that got written during the startup schedules
				for _, id := range runner.startupSystemSetIds {
					runner.eventStorage.ProcessEvents(id)
				}
			}
			isFirstRun = false

			if currentTickRate != *runner.tickRate {
				runner.Run(exitChannel, systems)
				return
			}
		}
	}
}

func (runner *fixedRunner) SetStartupSystemSetIds(ids []SystemSetId) {
	runner.startupSystemSetIds = ids
}

// uncappedRunner runs systems repeatedly and as fast as it can
type uncappedRunner struct {
	delta               *float64
	world               *ecs.World
	outerWorlds         *map[ecs.WorldId]*ecs.World
	logger              Logger
	appName             string
	eventStorage        *EventStorage
	startupSystemSetIds []SystemSetId
}

func (runner *uncappedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	var now int64
	start := time.Now().UnixNano()
	isFirstRun := true

	for {
		select {
		case <-exitChannel:
			return
		default:
		}

		now = time.Now().UnixNano()
		*runner.delta = float64(now-start) / 1_000_000_000
		start = now

		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage)

		if isFirstRun {
			// clear events that got written during the startup schedules
			for _, id := range runner.startupSystemSetIds {
				runner.eventStorage.ProcessEvents(id)
			}
		}
		isFirstRun = false
	}
}

func (runner *uncappedRunner) SetStartupSystemSetIds(ids []SystemSetId) {
	runner.startupSystemSetIds = ids
}

// nTimesRunner runs systems n amount of times and then returns
type nTimesRunner struct {
	numberOfRuns        int
	world               *ecs.World
	outerWorlds         *map[ecs.WorldId]*ecs.World
	logger              Logger
	appName             string
	eventStorage        *EventStorage
	startupSystemSetIds []SystemSetId
}

func (runner *nTimesRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	isFirstRun := true

	for range runner.numberOfRuns {
		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage)

		if isFirstRun {
			// clear events that got written during the startup schedules
			for _, id := range runner.startupSystemSetIds {
				runner.eventStorage.ProcessEvents(id)
			}
		}
		isFirstRun = false
	}
}

func (runner *nTimesRunner) SetStartupSystemSetIds(ids []SystemSetId) {
	runner.startupSystemSetIds = ids
}

func runSystemSet(systems []*SystemSet, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string, eventStorage *EventStorage) {
	for _, systemSet := range systems {
		errors := systemSet.Exec(world, outerWorlds, eventStorage)
		for _, err := range errors {
			logger.Error(fmt.Sprintf("%s - system returned error: %v", appName, err))
		}
	}
}
