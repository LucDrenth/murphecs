package app

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Runner interface {
	Run(exitChannel <-chan struct{}, systems []*SystemSet)
}

// fixedRunner runs systems at a fixed interval
type fixedRunner struct {
	tickRate    *time.Duration
	delta       *float64
	world       *ecs.World
	outerWorlds *map[ecs.WorldId]*ecs.World
	logger      Logger
	appName     string
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

			runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName)

			if currentTickRate != *runner.tickRate {
				runner.Run(exitChannel, systems)
				return
			}
		}
	}
}

// uncappedRunner runs systems repeatedly and as fast as it can
type uncappedRunner struct {
	delta       *float64
	world       *ecs.World
	outerWorlds *map[ecs.WorldId]*ecs.World
	logger      Logger
	appName     string
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

		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName)
	}
}

// nTimesRunner runs systems n amount of times and then returns
type nTimesRunner struct {
	numberOfRuns int
	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
}

func (runner *nTimesRunner) Run(exitChannel <-chan struct{}, systems []*SystemSet) {
	for range runner.numberOfRuns {
		runSystemSet(systems, runner.world, runner.outerWorlds, runner.logger, runner.appName)
	}
}

func runSystemSet(systems []*SystemSet, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string) {
	for _, systemSet := range systems {
		errors := systemSet.Exec(world, outerWorlds)
		for _, err := range errors {
			logger.Error(fmt.Sprintf("%s - system returned error: %v", appName, err))
		}
	}
}
