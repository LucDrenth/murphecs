package ecs

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Spawn spawns the given components and all their required components that are not declared in the component parameters.
// Return the associated entityId of the newly created entity on success.
//
// Returns an ErrDuplicateComponent error when any of the given components are of the same type.
//
// Calling Spawn without any components to generate an entityId is allowed.
func Spawn(world *world, components ...IComponent) (entityId, error) {
	componentTypes := toComponentTypes(components)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentTypes)
	if duplicate != nil {
		componentType := reflect.TypeOf(components[duplicateIndexA]).String()
		return 0, fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, componentType, duplicateIndexA, duplicateIndexB)
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
