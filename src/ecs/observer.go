package ecs

import "reflect"

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

type Observer struct{}

func (Observer) getObserverType() observerType {
	return customObserver
}

func (Observer) componentId(_ *World) ComponentId {
	panic("unexpected call to componentId")
}

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

// On registers an observer
func On[O AnyObserver](world *World, action func(world *World, observed O)) {
	var observer O

	switch observer.getObserverType() {
	case customObserver:
		{
			observerId := reflect.TypeFor[O]()
			world.observers[observerId] = append(world.observers[observerId], action)
		}

	case spawnObserver:
		{
			componentId := observer.componentId(world)

			observerType := reflect.TypeOf(observer)
			o := reflect.New(observerType).Elem()
			entityField, _ := observerType.FieldByName("Entity")
			entityFieldIndex := entityField.Index[0]

			world.spawnObservers[componentId] = append(world.spawnObservers[componentId], func(w *World, entityId EntityId) {
				o.Field(entityFieldIndex).Set(reflect.ValueOf(entityId))
				action(w, o.Interface().(O))
			})
		}

	case despawnObserver:
		{
			componentId := observer.componentId(world)

			observerType := reflect.TypeOf(observer)
			o := reflect.New(observerType).Elem()
			entityField, _ := observerType.FieldByName("Entity")
			entityFieldIndex := entityField.Index[0]

			world.despawnObservers[componentId] = append(world.despawnObservers[componentId], func(w *World, entityId EntityId) {
				o.Field(entityFieldIndex).Set(reflect.ValueOf(entityId))
				action(w, o.Interface().(O))
			})
		}

	default:
		panic("unhandled observer type")
	}
}

// Trigger triggers all registered observers for the given observer
func Trigger[O AnyObserver](world *World, observed O) {
	observerId := reflect.TypeFor[O]()
	observers, exists := world.observers[observerId]
	if !exists {
		return
	}

	for _, observer := range observers {
		observer.(func(*World, O))(world, observed)
	}
}
