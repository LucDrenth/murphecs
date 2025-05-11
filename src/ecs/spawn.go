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
		return nonExistingEntity, fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, debugType, duplicateIndexA, duplicateIndexB)
	}

	// get required components
	requiredComponents := getAllRequiredComponents(&componentIds, components, world)
	components = append(components, requiredComponents...)

	// spawn components
	entityId := world.generateEntityId()

	var returnedErr error = nil

	archetype, err := world.archetypeStorage.getArchetype(world, componentIds)
	if err != nil {
		return nonExistingEntity, err
	}
	world.archetypeStorage.entityIdToArchetype[entityId] = archetype

	var row uint
	for _, component := range components {
		// We can not reuse componentIds because it is not in the same order as components
		componentId := ComponentIdOf(component, world)

		storage := archetype.components[componentId]
		row, err = storage.insert(component)
		if err != nil {
			returnedErr = fmt.Errorf("failed to insert component %s in to component registry: %w", componentId.DebugString(), err)
			continue
		}
	}

	world.entities[entityId] = &EntityData{
		row:       row,
		archetype: archetype,
	}

	return entityId, returnedErr
}
