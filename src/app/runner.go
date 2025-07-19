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

// RunnerBasis provides base functionality of a runner. This includes tracking delta time.
type RunnerBasis struct {
	onFirstRunDone func()
	onRunDone      func()
	currentTick    *uint
	isFirstRun     bool

	delta           *float64
	timeLoopStarted int64 // in nano seconds
}

func NewRunnerBasis(app *SubApp) RunnerBasis {
	return RunnerBasis{
		onFirstRunDone: func() {},
		onRunDone:      func() {},
		currentTick:    &app.currentTick,
		isFirstRun:     true,
		delta:          app.lastDelta,
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

func (runner *RunnerBasis) Start() {
	now := time.Now().UnixNano()

	if runner.isFirstRun {
		runner.timeLoopStarted = now
	}

	*runner.delta = float64(now-runner.timeLoopStarted) / 1_000_000_000.0
	runner.timeLoopStarted = now
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
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (runner *fixedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	ticker := time.NewTicker(*runner.tickRate)
	currentTickRate := *runner.tickRate

	for {
		select {
		case <-exitChannel:
			return

		case <-ticker.C:
			runner.Start()
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
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (runner *uncappedRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	for {
		select {
		case <-exitChannel:
			return
		default:
		}

		runner.Start()
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

		runner.Start()
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
