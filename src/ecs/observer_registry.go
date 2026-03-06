package ecs

import "reflect"

type observerRegistry struct {
	observers        map[observerId][]any
	spawnObservers   map[ComponentId][]func(*World, EntityId)
	despawnObservers map[ComponentId][]func(*World, EntityId)
}

func newObserverRegistry() observerRegistry {
	return observerRegistry{
		observers:        map[observerId][]any{},
		spawnObservers:   map[ComponentId][]func(*World, EntityId){},
		despawnObservers: map[ComponentId][]func(*World, EntityId){},
	}
}

func (registry observerRegistry) triggerDespawnObservers(world *World, componentIds []ComponentId, entity EntityId) {
	for _, componentId := range componentIds {
		observers, exists := registry.despawnObservers[componentId]
		if !exists {
			continue
		}

		for _, observer := range observers {
			observer(world, entity)
		}
	}
}

func (registry observerRegistry) triggerSpawnObservers(world *World, componentIds []ComponentId, entity EntityId) {
	for _, componentId := range componentIds {
		observers, exists := registry.spawnObservers[componentId]
		if !exists {
			continue
		}

		for _, observer := range observers {
			observer(world, entity)
		}
	}
}

func registerObserver[O AnyObserver](registry *observerRegistry, world *World, action func(world *World, observed O)) {
	if action == nil {
		return
	}

	var observer O

	reflectedObserver := reflect.TypeFor[O]()
	if reflectedObserver.Kind() == reflect.Pointer {
		panic("can not register observer pointer")
	}

	switch observer.getObserverType() {
	case customObserver:
		{
			registry.observers[reflectedObserver] = append(registry.observers[reflectedObserver], action)
		}

	case spawnObserver:
		{
			componentId := observer.componentId(world)

			observerType := reflect.TypeOf(observer)
			o := reflect.New(observerType).Elem()
			entityField, _ := observerType.FieldByName("Entity")
			entityFieldIndex := entityField.Index[0]

			registry.spawnObservers[componentId] = append(registry.spawnObservers[componentId], func(w *World, entityId EntityId) {
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

			registry.despawnObservers[componentId] = append(registry.despawnObservers[componentId], func(w *World, entityId EntityId) {
				o.Field(entityFieldIndex).Set(reflect.ValueOf(entityId))
				action(w, o.Interface().(O))
			})
		}

	default:
		panic("unhandled observer type")
	}
}

func triggerObserver[O AnyObserver](world *World, registry *observerRegistry, observed O) {
	observerId := reflect.TypeFor[O]()
	observers, exists := registry.observers[observerId]
	if !exists {
		return
	}

	for _, observer := range observers {
		observer.(func(*World, O))(world, observed)
	}
}
