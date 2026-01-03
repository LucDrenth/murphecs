package app

import "github.com/lucdrenth/murphecs/src/ecs"

// Executor runs system
type Executor interface {
	Load(systems []*ScheduleSystems, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string, eventStorage *EventStorage)
	Run(currentTick uint)
	ProcessEvents(currentTick uint)
}

// ConsecutiveExecutor runs systems one after another, never running them in parallel.
type ConsecutiveExecutor struct {
	systems []*ScheduleSystems

	world        *ecs.World
	outerWorlds  *map[ecs.WorldId]*ecs.World
	logger       Logger
	appName      string
	eventStorage *EventStorage
}

func (executor *ConsecutiveExecutor) Load(systems []*ScheduleSystems, world *ecs.World, outerWorlds *map[ecs.WorldId]*ecs.World, logger Logger, appName string, eventStorage *EventStorage) {
	executor.systems = systems
	executor.world = world
	executor.outerWorlds = outerWorlds
	executor.logger = logger
	executor.appName = appName
	executor.eventStorage = eventStorage
}

func (executor *ConsecutiveExecutor) Run(currentTick uint) {
	for _, scheduleSystems := range executor.systems {
		errors := scheduleSystems.Exec(executor.world, executor.outerWorlds, executor.eventStorage, currentTick)
		for _, err := range errors {
			executor.logger.Error("%s - system returned error: %v", executor.appName, err)
		}
	}

	executor.world.Process()
}

func (executor *ConsecutiveExecutor) ProcessEvents(currentTick uint) {
	for _, startupSystem := range executor.systems {
		executor.eventStorage.ProcessEvents(startupSystem.id, currentTick)
	}
}
