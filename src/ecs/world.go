package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph/engine/src/utils"
)

// world contains all of the entities and their components.
type world struct {
	entityIdCounter uint
	entities        map[entityId]*entry
}

func NewWorld() world {
	return world{
		entities: map[entityId]*entry{},
	}
}

// Spawn spawns the given components and all their required components that are not declared in the component parameters. Return the associated entityId on success.
//
// An error is returned when any of the given components are of the same type.
//
// Spawn without any components to generate an entityId is allowed.
func Spawn(world *world, components ...IComponent) (entityId, error) {
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
	world.entities[entityId] = &entry{
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

// getEntry returns the entry that correspond to entity, or an ErrEntityNotFound error if it wasn't found.
func (world *world) getEntry(entity entityId) (*entry, error) {
	entry, ok := world.entities[entity]
	if !ok {
		return nil, ErrEntityNotFound
	}

	return entry, nil
}
