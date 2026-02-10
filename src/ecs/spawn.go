package ecs

import (
	"fmt"
	"reflect"

	"github.com/lucdrenth/murphecs/src/utils"
)

// Spawn spawns the given components and all their required components that are not declared in the component parameters.
// Return the associated entityId of the newly created entity on success.
//
// Can return the following errors:
//   - Returns an ErrComponentIsNil error when any of the given components is nil
//   - Returns an ErrDuplicateComponent error when any of the given components are of the same type.
func Spawn(world *World, components ...AnyComponent) (EntityId, error) {
	for i, component := range components {
		if component == nil {
			return nonExistingEntity, fmt.Errorf("%w: at position %d", ErrComponentIsNil, i+1)
		}
	}

	componentIds := toComponentIds(components, world)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		debugType := ComponentDebugStringOf(components[duplicateIndexA])
		return nonExistingEntity, fmt.Errorf("%w: %s at positions %d and %d", ErrComponentDuplicate, debugType, duplicateIndexA, duplicateIndexB)
	}

	// get required components
	requiredComponents := getAllRequiredComponents(&componentIds, components, world)
	components = append(components, requiredComponents...)

	// collect reflect.Value's. If any component is not a pointer, convert it to a pointer.
	componentValues := make([]reflect.Value, len(components))
	for i, component := range components {
		componentValue := reflect.ValueOf(component)
		if componentValue.Kind() != reflect.Pointer {
			ptrValue := reflect.New(componentValue.Type())
			ptrValue.Elem().Set(componentValue)
			components[i] = ptrValue.Interface().(AnyComponent)
			componentValue = ptrValue
		}

		componentValues[i] = componentValue
	}

	// spawn components
	entityId := world.generateEntityId()

	var returnedErr error = nil

	archetype, err := world.archetypeStorage.getArchetype(world, componentIds)
	if err != nil {
		return nonExistingEntity, err
	}

	var row uint
	for i, component := range components {
		// We can not reuse componentIds because it is not in the same order as components
		componentId := ComponentIdOf(component, world)

		storage := archetype.components[componentId]
		row, err = storage.insertValue(world, &componentValues[i])
		if err != nil {
			return nonExistingEntity, fmt.Errorf("failed to insert component %s in to component registry: %w", componentId.DebugString(), err)
		}
	}

	world.archetypeStorage.entityIdToArchetype[entityId] = archetype
	archetype.entities = append(archetype.entities, entityId)

	world.entities[entityId] = &EntityData{
		row:       row,
		archetype: archetype,
	}

	return entityId, returnedErr
}
