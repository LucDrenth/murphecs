package ecs

import (
	"errors"
	"reflect"
)

// World contains all of the entities and their components.
type World struct {
	entityIdCounter          uint
	entities                 map[EntityId]*EntityData // TODO EntityData probably does not have to be pointer once archeType is implemented
	initialComponentCapacity initialComponentCapacityStrategy
	componentRegistry        componentRegistry
	archetypeStorage         archetypeStorage
}

// DefaultWorld returns a World with default configs.
func DefaultWorld() World {
	world, err := NewWorld(DefaultWorldConfigs())
	if err != nil {
		// Creating a world with default configs should never result in an error.
		// This is confirmed by a unit test, so we can 'safely' panic just in case
		// it happens.
		panic(err)
	}

	return world
}

// NewWorld returns a world that can contain entities and components.
func NewWorld(configs WorldConfigs) (World, error) {
	if configs.ComponentCapacityStrategy == nil {
		return World{}, errors.New("component capacity strategy can not be nil")
	}

	return World{
		entities:                 map[EntityId]*EntityData{},
		initialComponentCapacity: configs.ComponentCapacityStrategy,
		componentRegistry: componentRegistry{
			components: map[reflect.Type]uint{},
		},
		archetypeStorage: newArchetypeStorage(),
	}, nil
}

func (world *World) CountEntities() int {
	return len(world.entities)
}

func (world *World) CountComponents() int {
	return int(world.archetypeStorage.countComponents())
}

func (world *World) generateEntityId() EntityId {
	world.entityIdCounter++
	return EntityId(world.entityIdCounter)
}
