package ecs

import (
	"github.com/lucdrenth/murph_engine/src/log"
)

// World contains all of the entities and their components.
type World struct {
	entityIdCounter uint
	entities        map[EntityId]*entityData
	components      map[componentType]*componentRegistry
	logger          log.Logger
}

// NewWorld returns a world that can contain entities and components.
func NewWorld() World {
	logger := log.Console()

	return World{
		entities:   map[EntityId]*entityData{},
		components: map[componentType]*componentRegistry{},
		logger:     &logger,
	}
}

func (world *World) CountEntities() int {
	return len(world.entities)
}

func (world *World) CountComponents() int {
	result := 0

	for _, entry := range world.entities {
		result += len(entry.components)
	}

	return result
}

func (world *World) createEntity() EntityId {
	world.entityIdCounter++
	entity := EntityId(world.entityIdCounter)
	world.entities[entity] = &entityData{components: map[componentType]uint{}}
	return entity
}

// getComponentRegistry creates a new component registry if it doesn't exist yet.
func (world *World) getComponentRegistry(componentType componentType) *componentRegistry {
	componentRegistry, ok := world.components[componentType]
	if !ok {
		newComponentRegistry := createComponentRegistry(1024, componentType)
		world.components[componentType] = &newComponentRegistry
		return &newComponentRegistry
	}

	return componentRegistry
}
