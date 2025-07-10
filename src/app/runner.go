package app

import (
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Runner interface {
	Run(exitChannel <-chan struct{}, systems []*SystemSet)
	setOnFirstRunDone(func())
	setOnRunDone(func())
}

type RunnerBasis struct {
	onFirstRunDone func()
	onRunDone      func()
	currentTick    *uint
	isFirstRun     bool
}

func NewRunnerBasis(app *SubApp) RunnerBasis {
	return RunnerBasis{
		onFirstRunDone: func() {},
		onRunDone:      func() {},
		currentTick:    &app.currentTick,
		isFirstRun:     true,
	}
}

func (runner *RunnerBasis) setOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *RunnerBasis) setOnRunDone(handler func()) {
	runner.onRunDone = handler
}

func (runner *RunnerBasis) CurrentTick() uint {
	return *runner.currentTick
}

func (runner *RunnerBasis) Done() {
	if runner.isFirstRun {
		runner.onFirstRunDone()
		runner.isFirstRun = false
	}

	runner.onRunDone()
}

// fixedRunner runs systems at a fixed interval
type fixedRunner struct {
	RunnerBasis
	tickRate     *time.Duration
	delta        *float64
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (runner *fixedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	ticker := time.NewTicker(*runner.tickRate)
	currentTickRate := *runner.tickRate
	var now int64
	start := time.Now().UnixNano()

	for {
		select {
		case <-exitChannel:
			return

		case <-ticker.C:
			now = time.Now().UnixNano()
			*runner.delta = float64(now-start) / 1_000_000_000.0
			start = now

			runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage, *runner.currentTick)

			runner.Done()

			if currentTickRate != *runner.tickRate {
				runner.Run(exitChannel, systems)
				return
			}
		}
	}
}

// uncappedRunner runs systems repeatedly and as fast as it can
type uncappedRunner struct {
	RunnerBasis
	delta        *float64
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (runner *uncappedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	var now int64
	start := time.Now().UnixNano()

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

		runner.Done()
	}
}

func (runner *uncappedRunner) setOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *uncappedRunner) setOnRunDone(handler func()) {
	runner.onRunDone = handler
}

// nTimesRunner runs systems n amount of times and then returns
type nTimesRunner struct {
	RunnerBasis
	numberOfRuns int
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (runner *nTimesRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	for range runner.numberOfRuns {
		select {
		case <-exitChannel:
			return
		default:
		}

		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName, runner.eventStorage, *runner.currentTick)

		runner.Done()
	}
}

func (runner *nTimesRunner) setOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *nTimesRunner) setOnRunDone(handler func()) {
	runner.onRunDone = handler
}

func runSystemSet(systems []*SystemSet, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string, eventStorage *EventStorage, currentTick uint) {
	for _, systemSet := range systems {
		errors := systemSet.Exec(world, outerWorlds, eventStorage, currentTick)
		for _, err := range errors {
			logger.Error("%s - system returned error: %v", appName, err)
		}
	}
}
