package ecs

import (
	"errors"
	"reflect"
)

type WorldId int

// World contains all of the entities and their components.
type World struct {
	id *WorldId // setting an id is optional

	entityIdCounter   uint
	entities          map[EntityId]*EntityData
	componentRegistry componentRegistry
	archetypeStorage  archetypeStorage

	initialComponentCapacityStrategy initialComponentCapacityStrategy
	componentCapacityGrowthStrategy  componentCapacityGrowthStrategy
}

// NewDefaultWorld returns a World with default configs.
func NewDefaultWorld() World {
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
	if configs.InitialComponentCapacityStrategy == nil {
		return World{}, errors.New("config InitialComponentCapacityStrategy can not be nil")
	}

	if configs.ComponentCapacityGrowthStrategy == nil {
		return World{}, errors.New("config ComponentCapacityGrowthStrategy can not be nil")
	}

	return World{
		entities:                         map[EntityId]*EntityData{},
		id:                               configs.Id,
		initialComponentCapacityStrategy: configs.InitialComponentCapacityStrategy,
		componentCapacityGrowthStrategy:  configs.ComponentCapacityGrowthStrategy,
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
