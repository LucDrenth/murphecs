package ecs

import (
	"errors"
	"fmt"
	"sync"
)

type WorldId int

// World contains all of the entities and their components.
type World struct {
	id *WorldId // setting an id is optional

	entityIdCounter   uint
	entities          map[EntityId]*EntityData
	componentRegistry componentRegistry
	archetypeStorage  archetypeStorage

	resources resourceStorage
	observers observerRegistry
	events    EventStorage

	initialComponentCapacityStrategy initialComponentCapacityStrategy
	componentCapacityGrowthStrategy  componentCapacityGrowthStrategy

	scheduler                Scheduler
	outerWorlds              map[WorldId]*World
	logger                   Logger
	scheduleSystemsIdCounter ScheduleSystemsId

	Mutex sync.RWMutex
}

// NewDefaultWorld returns a World with default configs.
func NewDefaultWorld() *World {
	world, err := NewWorld(DefaultWorldConfigs())
	if err != nil {
		// Creating a world with default configs should never result in an error.
		// This is confirmed by a unit test, so we can 'safely' panic just in case
		// it happens.
		panic(err)
	}

	return &world
}

// NewWorld returns a world that can contain entities and components.
func NewWorld(configs WorldConfigs) (World, error) {
	if configs.InitialComponentCapacityStrategy == nil {
		return World{}, errors.New("config InitialComponentCapacityStrategy can not be nil")
	}

	if configs.ComponentCapacityGrowthStrategy == nil {
		return World{}, errors.New("config ComponentCapacityGrowthStrategy can not be nil")
	}

	logger := configs.Logger
	if logger == nil {
		logger = &NoOpLogger{}
	}

	return World{
		entities:                         map[EntityId]*EntityData{},
		id:                               configs.Id,
		initialComponentCapacityStrategy: configs.InitialComponentCapacityStrategy,
		componentCapacityGrowthStrategy:  configs.ComponentCapacityGrowthStrategy,
		componentRegistry:                newComponentRegistry(),
		archetypeStorage:                 newArchetypeStorage(),
		resources:                        newResourceStorage(),
		observers:                        newObserverRegistry(),
		events:                           NewEventStorage(),
		scheduler:                        newScheduler(),
		outerWorlds:                      map[WorldId]*World{},
		logger:                           logger,
	}, nil
}

// Process should be called on a regular basis (such as every tick).
//
// ! This call is NOT concurrency safe !
func (world *World) Process() {
	world.componentRegistry.processComponentIdRegistries()
}

func (world *World) CountEntities() int {
	return len(world.entities)
}

func (world *World) CountComponents() int {
	return int(world.archetypeStorage.countComponents())
}

func (world *World) CountArchetypes() int {
	return len(world.archetypeStorage.componentsHashToArchetype)
}

func (world *World) generateEntityId() EntityId {
	world.entityIdCounter++
	return EntityId(world.entityIdCounter)
}

func (world *World) Id() *WorldId {
	return world.id
}

func (world *World) Resources() *resourceStorage {
	return &world.resources
}

func (world *World) Events() *EventStorage {
	return &world.events
}

// AddSchedule adds a schedule that systems can be added to.
func (world *World) AddSchedule(schedule Schedule, order ScheduleOrder, isPaused bool) error {
	if order == nil {
		order = ScheduleLast{}
	}

	world.scheduleSystemsIdCounter++
	return world.scheduler.addSchedule(schedule, world.scheduleSystemsIdCounter, order, isPaused)
}

// AddSystem adds a system to the given schedule. Systems must be functions.
func (world *World) AddSystem(schedule Schedule, system System) error {
	return world.AddSystemWithSource(schedule, system, callerSource(1))
}

// AddSystemWithSource adds a system to the given schedule with an explicit source path for error messages.
func (world *World) AddSystemWithSource(schedule Schedule, system System, source string) error {
	return world.scheduler.addSystem(schedule, system, source, world, &world.outerWorlds, world.logger, world.Events())
}

// SetSchedulePaused pauses or unpauses a schedule.
func (world *World) SetSchedulePaused(schedule Schedule, isPaused bool) error {
	systems, exists := world.scheduler.systems[schedule]
	if !exists {
		return fmt.Errorf("%w: %s", ErrScheduleNotFound, schedule)
	}

	currentlyPaused := systems.isPaused.Load()
	if currentlyPaused == isPaused {
		return nil
	}

	if !currentlyPaused && isPaused {
		systems.isFirstExecSincePaused = true
	}
	systems.isPaused.Store(isPaused)

	return nil
}

// RegisterOuterWorld lets systems query components and resources from another world.
func (world *World) RegisterOuterWorld(id WorldId, other *World) error {
	if _, exists := world.outerWorlds[id]; exists {
		return fmt.Errorf("id %d is already registered", id)
	}

	world.outerWorlds[id] = other
	return nil
}

// OuterWorlds returns the map of registered outer worlds.
func (world *World) OuterWorlds() *map[WorldId]*World {
	return &world.outerWorlds
}

// GetScheduleSystems returns all [ScheduleSystems] in their execution order.
func (world *World) GetScheduleSystems() ([]*ScheduleSystems, error) {
	return world.scheduler.getScheduleSystems()
}

// GetScheduleSystemsBySchedules returns the [ScheduleSystems] for the given schedule names, in order.
func (world *World) GetScheduleSystemsBySchedules(schedules []Schedule) ([]*ScheduleSystems, error) {
	return world.scheduler.getScheduleSystemsBySchedules(schedules)
}

// PrepareSystems resolves outer-resource system params for all schedules.
func (world *World) PrepareSystems() error {
	scheduleSystems, err := world.scheduler.getScheduleSystems()
	if err != nil {
		return fmt.Errorf("failed to get schedule systems: %w", err)
	}

	for _, systems := range scheduleSystems {
		if err := systems.prepare(&world.outerWorlds); err != nil {
			return fmt.Errorf("failed to prepare systems: %w", err)
		}
	}

	return nil
}

// NumberOfSystems returns the total number of systems across all schedules.
func (world *World) NumberOfSystems() uint {
	return world.scheduler.numberOfSystems()
}

// NumberOfSchedules returns the total number of registered schedules.
func (world *World) NumberOfSchedules() uint {
	return world.scheduler.numberOfSchedules()
}

type WorldStats struct {
	NumberOfEntities   int
	NumberOfComponents int
	NumberOfArchetypes int
}

func (world *World) Stats() WorldStats {
	return WorldStats{
		NumberOfEntities:   world.CountEntities(),
		NumberOfComponents: world.CountEntities(),
		NumberOfArchetypes: world.CountArchetypes(),
	}
}
