package ecs

import (
	"reflect"
)

type observerType uint

const (
	// customObserver is any user defined observer
	customObserver observerType = iota
	spawnObserver
	despawnObserver
)

type AnyObserver interface {
	getObserverType() observerType
	componentId(world *World) ComponentId
}

type observerId reflect.Type

// Observer can be embedded in to a struct to make a custom observer.
//
// It can be used in two ways:
//   - globally: register it with [On] and trigger it with [Trigger]
//   - for an entity: register it with [Observe] and trigger it with [TriggerEntity]
type Observer struct{}

func (Observer) getObserverType() observerType {
	return customObserver
}

func (Observer) componentId(_ *World) ComponentId {
	panic("unexpected call to componentId")
}

// OnSpawn is triggered when:
//   - an entity with component [C] is spawned using [Spawn]
//   - component [C] is added to an entity using [Insert] or [InsertOrOverwrite]
//
// Global OnSpawn observers get triggered before entity-specific observers.
type OnSpawn[C AnyComponent] struct {
	Observer
	Entity EntityId
}

func (OnSpawn[C]) getObserverType() observerType {
	return spawnObserver
}

func (OnSpawn[C]) componentId(world *World) ComponentId {
	return ComponentIdFor[C](world)
}

// OnSpawn is triggered when:
//   - component [C] gets removed from an entity using [Remove1], [Remove2] and so on
//   - an entity with component [C] gets despawned using [Despawn]
//
// Global OnDespawn observers get triggered before entity-specific observers.
type OnDespawn[C AnyComponent] struct {
	Observer
	Entity EntityId
}

func (OnDespawn[C]) getObserverType() observerType {
	return despawnObserver
}

func (OnDespawn[C]) componentId(world *World) ComponentId {
	return ComponentIdFor[C](world)
}

// On registers a global observer
func On[O AnyObserver](world *World, action func(world *World, observed O)) {
	registerObserver(&world.observers, world, action)
}

// Trigger triggers all registered observers for the given observer
func Trigger[O AnyObserver](world *World, observed O) {
	triggerObserver(world, &world.observers, observed)
}

// Trigger triggers all registered observers for the given observer
func TriggerEntity[O AnyObserver](world *World, entity EntityId, observed O) error {
	entityData, exists := world.entities[entity]
	if !exists {
		return ErrEntityNotFound
	}
	if entityData.observers == nil {
		return nil
	}

	triggerObserver(world, entityData.observers, observed)
	return nil
}

// Observe registers an entity-specific observer
func Observe[O AnyObserver](world *World, entity EntityId, action func(world *World, observed O)) error {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if entityData.observers == nil {
		entityData.observers = new(newObserverRegistry())
	}
	registerObserver(entityData.observers, world, action)

	return nil
}
