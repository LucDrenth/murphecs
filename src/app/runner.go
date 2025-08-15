package app

import (
	"time"
)

type Runner interface {
	Run(exitChannel <-chan struct{}, executor Executor)
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
	tickRate *time.Duration
}

func (runner *fixedRunner) Run(exitChannel <-chan struct{}, executor Executor) {
	ticker := time.NewTicker(*runner.tickRate)
	currentTickRate := *runner.tickRate

	for {
		select {
		case <-exitChannel:
			return

		case <-ticker.C:
			runner.Start()
			executor.Run(*runner.currentTick)
			runner.Done()

			if currentTickRate != *runner.tickRate {
				runner.Run(exitChannel, executor)
				return
			}
		}
	}
}

// uncappedRunner runs systems repeatedly and as fast as it can
type uncappedRunner struct {
	RunnerBasis
}

func (runner *uncappedRunner) Run(exitChannel <-chan struct{}, executor Executor) {
	for {
		select {
		case <-exitChannel:
			return
		default:
		}

		runner.Start()
		executor.Run(*runner.currentTick)
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
}

func (runner *nTimesRunner) Run(exitChannel <-chan struct{}, executor Executor) {
	for range runner.numberOfRuns {
		select {
		case <-exitChannel:
			return
		default:
		}

		runner.Start()
		executor.Run(*runner.currentTick)
		runner.Done()
	}
}

func (runner *nTimesRunner) setOnFirstRunDone(handler func()) {
	runner.onFirstRunDone = handler
}

func (runner *nTimesRunner) setOnRunDone(handler func()) {
	runner.onRunDone = handler
}
