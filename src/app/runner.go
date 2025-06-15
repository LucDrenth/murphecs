package app

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Runner interface {
	Run(exitChannel <-chan struct{}, systems []*SystemSet)
	SetOnFirstRunDone(func())
	SetOnRunDone(func())
}

// fixedRunner runs systems at a fixed interval
type fixedRunner struct {
	tickRate       *time.Duration
	delta          *float64
	world          *ecs.World
	outerWorlds    *map[ecs.WorldId]*ecs.World
	logger         Logger
	appName        string
	eventStorage   *EventStorage
	onFirstRunDone func()
	onRunDone      func()
	currentTick    *uint
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

			runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage, *runner.currentTick)

			if isFirstRun {
				if runner.onFirstRunDone != nil {
					runner.onFirstRunDone()
				}
				isFirstRun = false
			}
			if runner.onRunDone != nil {
				runner.onRunDone()
			}

			if currentTickRate != *runner.tickRate {
				runner.Run(exitChannel, systems)
				return
			}
		}
	}
}

func (runner *fixedRunner) SetOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *fixedRunner) SetOnRunDone(handler func()) {
	runner.onRunDone = handler
}

// uncappedRunner runs systems repeatedly and as fast as it can
type uncappedRunner struct {
	delta          *float64
	world          *ecs.World
	outerWorlds    *map[ecs.WorldId]*ecs.World
	logger         Logger
	appName        string
	eventStorage   *EventStorage
	onFirstRunDone func()
	onRunDone      func()
	currentTick    *uint
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

		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage, *runner.currentTick)

		if isFirstRun {
			if runner.onFirstRunDone != nil {
				runner.onFirstRunDone()
			}
			isFirstRun = false
		}
		if runner.onRunDone != nil {
			runner.onRunDone()
		}
	}
}

func (runner *uncappedRunner) SetOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *uncappedRunner) SetOnRunDone(handler func()) {
	runner.onRunDone = handler
}

// nTimesRunner runs systems n amount of times and then returns
type nTimesRunner struct {
	numberOfRuns   int
	world          *ecs.World
	outerWorlds    *map[ecs.WorldId]*ecs.World
	logger         Logger
	appName        string
	eventStorage   *EventStorage
	onFirstRunDone func()
	onRunDone      func()
	currentTick    *uint
}

func (runner *nTimesRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	isFirstRun := true

	for range runner.numberOfRuns {
		select {
		case <-exitChannel:
			return
		default:
		}

		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage, *runner.currentTick)

		if isFirstRun {
			if runner.onFirstRunDone != nil {
				runner.onFirstRunDone()
			}
			isFirstRun = false
		}
		if runner.onRunDone != nil {
			runner.onRunDone()
		}
	}
}

func (runner *nTimesRunner) SetOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *nTimesRunner) SetOnRunDone(handler func()) {
	runner.onRunDone = handler
}

func runSystemSet(systems []*SystemSet, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string, eventStorage *EventStorage, currentTick uint) {
	for _, systemSet := range systems {
		errors := systemSet.Exec(world, outerWorlds, eventStorage, currentTick)
		for _, err := range errors {
			logger.Error(fmt.Sprintf("%s - system returned error: %v", appName, err))
		}
	}
}
