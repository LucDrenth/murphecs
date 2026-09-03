package ecs

import (
	"errors"
	"fmt"
	"reflect"
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
	currentScheduleSystemsId ScheduleSystemsId // set to the running schedule's id during Exec, 0 otherwise

	Mutex sync.RWMutex

	isQuerying bool
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
		NumberOfComponents: world.CountComponents(),
		NumberOfArchetypes: world.CountArchetypes(),
	}
}

// GetComponentsForEntity returns all components that belong to the given entity, keyed by their ComponentId.
// Returned components are copies.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error if the entity is not found.
func (world *World) GetComponentsForEntity(entity EntityId) (map[ComponentId]any, error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, ErrEntityNotFound
	}

	components := make(map[ComponentId]any, len(entityData.archetype.componentIds))

	for componentId, storage := range entityData.archetype.components {
		componentPointer, err := storage.getComponentPointer(entityData.row)
		if err != nil {
			return nil, err
		}

		components[componentId] = reflect.NewAt(componentId.componentType, componentPointer).Elem().Interface()
	}

	return components, nil
}

func (world *World) Spawn(components ...AnyComponent) (EntityId, error) {
	return spawn(world, components...)
}
func (world *World) Insert(entity EntityId, components ...AnyComponent) error {
	return insert(world, entity, components...)
}
func (world *World) InsertOrOverwrite(entity EntityId, components ...AnyComponent) error {
	return insertOrOverwrite(world, entity, components...)
}
func (world *World) Remove1[A AnyComponent](entity EntityId) error {
	return remove1[A](world, entity)
}
func (world *World) Remove2[A, B AnyComponent](entity EntityId) error {
	return remove2[A, B](world, entity)
}
func (world *World) Remove3[A, B, C AnyComponent](entity EntityId) error {
	return remove3[A, B, C](world, entity)
}
func (world *World) Remove4[A, B, C, D AnyComponent](entity EntityId) error {
	return remove4[A, B, C, D](world, entity)
}
func (world *World) Despawn(entity EntityId) error {
	return despawn(world, entity)
}
func (world *World) HasComponent[C AnyComponent](entity EntityId) (bool, error) {
	return hasComponent[C](world, entity)
}
func (world *World) HasComponentId(entity EntityId, componentId ComponentId) (bool, error) {
	return hasComponentId(world, entity, componentId)
}
func (world *World) Get1[A AnyComponent](entity EntityId) (A, error) {
	return get1[A](world, entity)
}
func (world *World) Get2[A, B AnyComponent](entity EntityId) (A, B, error) {
	return get2[A, B](world, entity)
}
func (world *World) Get3[A, B, C AnyComponent](entity EntityId) (A, B, C, error) {
	return get3[A, B, C](world, entity)
}
func (world *World) Get4[A, B, C, D AnyComponent](entity EntityId) (A, B, C, D, error) {
	return get4[A, B, C, D](world, entity)
}
func (world *World) Get5[A, B, C, D, E AnyComponent](entity EntityId) (A, B, C, D, E, error) {
	return get5[A, B, C, D, E](world, entity)
}
func (world *World) Get6[A, B, C, D, E, F AnyComponent](entity EntityId) (A, B, C, D, E, F, error) {
	return get6[A, B, C, D, E, F](world, entity)
}
func (world *World) Get7[A, B, C, D, E, F, G AnyComponent](entity EntityId) (A, B, C, D, E, F, G, error) {
	return get7[A, B, C, D, E, F, G](world, entity)
}
func (world *World) Get8[A, B, C, D, E, F, G, H AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, error) {
	return get8[A, B, C, D, E, F, G, H](world, entity)
}
func (world *World) Get9[A, B, C, D, E, F, G, H, I AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, error) {
	return get9[A, B, C, D, E, F, G, H, I](world, entity)
}
func (world *World) Get10[A, B, C, D, E, F, G, H, I, J AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, error) {
	return get10[A, B, C, D, E, F, G, H, I, J](world, entity)
}
func (world *World) Get11[A, B, C, D, E, F, G, H, I, J, K AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, error) {
	return get11[A, B, C, D, E, F, G, H, I, J, K](world, entity)
}
func (world *World) Get12[A, B, C, D, E, F, G, H, I, J, K, L AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, L, error) {
	return get12[A, B, C, D, E, F, G, H, I, J, K, L](world, entity)
}
func (world *World) Get13[A, B, C, D, E, F, G, H, I, J, K, L, M AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, L, M, error) {
	return get13[A, B, C, D, E, F, G, H, I, J, K, L, M](world, entity)
}
func (world *World) Get14[A, B, C, D, E, F, G, H, I, J, K, L, M, N AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, error) {
	return get14[A, B, C, D, E, F, G, H, I, J, K, L, M, N](world, entity)
}
func (world *World) Get15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, error) {
	return get15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O](world, entity)
}
func (world *World) Get16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P AnyComponent](entity EntityId) (A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, error) {
	return get16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P](world, entity)
}
func (world *World) Observe[O AnyObserver](entity EntityId, action System) error {
	return observe[O](world, entity, action)
}
func (world *World) On[O AnyObserver](action System) error {
	return on[O](world, action)
}
func (world *World) Trigger[O AnyObserver](observed O) {
	trigger[O](world, observed)
}
func (world *World) TriggerEntity[O AnyObserver](entity EntityId, observed O) error {
	return triggerEntity[O](world, entity, observed)
}
func (world *World) GetResource[R Resource]() (result R, err error) {
	return getResource[R](world)
}
