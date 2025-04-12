package ecs

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/utils"
)

type entityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = entityId(0)

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

// Spawn spawns the given components and all their required components that are not declared in the component parameters. Return the associated entityId on success.
//
// An error is returned when any of the given components are of the same type.
//
// Spawn without any components to generate an entityId is allowed.
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
