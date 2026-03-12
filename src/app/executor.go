package app

import "github.com/lucdrenth/murphecs/src/ecs"

// Executor runs systems
type Executor interface {
	Load(systems []*ecs.ScheduleSystems, world *ecs.World, logger Logger, appName string)
	Run(currentTick uint)
	ProcessEvents(currentTick uint)
}

// ConsecutiveExecutor runs systems one after another, never running them in parallel.
type ConsecutiveExecutor struct {
	systems []*ecs.ScheduleSystems

	world        *ecs.World
	logger       Logger
	appName      string
	eventStorage *ecs.EventStorage
}

func (executor *ConsecutiveExecutor) Load(systems []*ecs.ScheduleSystems, world *ecs.World, logger Logger, appName string) {
	executor.systems = systems
	executor.world = world
	executor.logger = logger
	executor.appName = appName
	executor.eventStorage = world.Events()
}

func (executor *ConsecutiveExecutor) Run(currentTick uint) {
	for _, scheduleSystems := range executor.systems {
		errors := scheduleSystems.Exec(executor.world, executor.world.OuterWorlds(), executor.eventStorage, currentTick)
		for _, err := range errors {
			executor.logger.Error("%s - system returned error: %v", executor.appName, err)
		}
	}

	executor.world.Process()
}

func (executor *ConsecutiveExecutor) ProcessEvents(currentTick uint) {
	for _, startupSystem := range executor.systems {
		executor.eventStorage.ProcessEvents(startupSystem.Id(), currentTick)
	}
}
