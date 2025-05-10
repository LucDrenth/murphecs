package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Spawn spawns the given components and all their required components that are not declared in the component parameters.
// Return the associated entityId of the newly created entity on success.
//
// Returns an ErrDuplicateComponent error when any of the given components are of the same type.
//
// Returns an ErrComponentIsNotAPointer error when any of the given component or their required components are not passed as
// a reference, while still inserting all other components that are valid.
//
// Calling Spawn without any components to generate an entityId is allowed.
func Spawn(world *World, components ...IComponent) (EntityId, error) {
	componentIds := toComponentIds(components, world)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		debugType := ComponentDebugStringOf(components[duplicateIndexA])
		return 0, fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, debugType, duplicateIndexA, duplicateIndexB)
	}

	// get required components
	requiredComponents := getAllRequiredComponents(&componentIds, components, world)
	components = append(components, requiredComponents...)

	// spawn components
	entity := world.createEntity()

	var returnedErr error = nil

	for _, component := range components {
		componentId := ComponentIdOf(component, world)
		componentRegistry, err := world.getComponentRegistry(componentId)
		if err != nil {
			returnedErr = fmt.Errorf("failed to get component registry: %w", err)
			continue
		}

		componentIndex, err := componentRegistry.insert(component)
		if err != nil {
			returnedErr = fmt.Errorf("failed to insert component in to component registry: %w", err)
			continue
		}

		world.entities[entity].components[componentId] = componentIndex
	}

	return entity, returnedErr
}
