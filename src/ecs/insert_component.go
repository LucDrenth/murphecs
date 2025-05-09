package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Insert adds the given components and all their required components (that the entity does not yet have) to the given entity.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error when the given entity does not exist
//   - Returns an ErrDuplicateComponent error when any of the given components are of the same type.
//   - Returns an ErrComponentAlreadyPresent error if any of the components is already present while still inserting
//     the components that are not yet present.
//   - Returns an ErrComponentIsNotAPointer error if any of the given components, or their required components, are not
//     passed as a reference (e.g. componentA{} instead of &componentA{})
//   - Returns an ErrInvalidComponentStorageCapacity if the component storage capacity, that is decided through World
//     configs, is not valid
func Insert(world *World, entity EntityId, components ...IComponent) (resultErr error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	componentIds := toComponentIds(components)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		debugType := ComponentDebugStringOf(components[duplicateIndexA])
		return fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, debugType, duplicateIndexA, duplicateIndexB)
	}

	for i, component := range components {
		if _, componentExists := entityData.components[componentIds[i]]; componentExists {
			resultErr = fmt.Errorf("%w: %s", ErrComponentAlreadyPresent, ComponentDebugStringOf(component))
			continue
		}

		componentRegistry, err := world.getComponentRegistry(componentIds[i])
		if err != nil {
			resultErr = fmt.Errorf("failed to get component registry: %w", err)
			continue
		}

		componentIndex, err := componentRegistry.insert(component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert component: %w", err)
			continue
		}

		world.entities[entity].components[componentIds[i]] = componentIndex
	}

	requiredComponents := getAllRequiredComponents(&componentIds, components)
	componentIds = toComponentIds(requiredComponents)

	for i, component := range requiredComponents {
		if _, componentExists := entityData.components[componentIds[i]]; componentExists {
			continue
		}

		componentRegistry, err := world.getComponentRegistry(componentIds[i])
		if err != nil {
			resultErr = fmt.Errorf("failed to get component registry: %w", err)
			continue
		}

		componentIndex, err := componentRegistry.insert(component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert a required component: %w", err)
			continue
		}

		world.entities[entity].components[componentIds[i]] = componentIndex
	}

	return resultErr
}
