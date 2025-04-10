package ecs

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/utils"
)

type entityId = uint

type entry struct {
	components []IComponent
}

type world struct {
	entityIdCounter uint
	entities        map[entityId]entry
}

func NewWorld() world {
	return world{
		entities: map[entityId]entry{},
	}
}

func (world *world) Spawn(components ...IComponent) (entityId, error) {
	componentTypes := toComponentTypes(components)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentTypes)
	if duplicate != nil {
		return 0, fmt.Errorf("can not spawn duplicate component: %s at positions %d and %d", *duplicate, duplicateIndexA, duplicateIndexB)
	}

	// get required components
	requiredComponents := getAllRequiredComponents(&componentTypes, components)
	components = append(components, requiredComponents...)

	// spawn components
	world.entityIdCounter++
	entityId := world.entityIdCounter
	world.entities[entityId] = entry{
		components: components,
	}

	return entityId, nil
}

func (world *world) CountEntities() int {
	return len(world.entities)
}

func (world *world) CountComponents() int {
	result := 0

	for _, entry := range world.entities {
		result += len(entry.components)
	}

	return result
}
