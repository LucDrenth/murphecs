package ecs

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/utils"
)

// Insert adds the given components and all their required components (that the entity does not yet have) to the given entity.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error when the given entity does not exist
//   - Returns an ErrComponentIsNil error when any of the given components is nil
//   - Returns an ErrDuplicateComponent error when any of the given components are of the same type.
//   - Returns an ErrComponentAlreadyPresent error if any of the components is already present while still inserting
//     the components that are not yet present.
//   - Returns an ErrInvalidComponentStorageCapacity if the component storage capacity, that is decided through World
//     configs, is not valid
func Insert(world *World, entity EntityId, components ...AnyComponent) (resultErr error) {
	if len(components) == 0 {
		return nil
	}

	for i, component := range components {
		if component == nil {
			return fmt.Errorf("%w: at position %d", ErrComponentIsNil, i+1)
		}
	}

	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	componentIds := toComponentIds(components, world)

	// check for duplicates
	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		debugType := ComponentDebugStringOf(components[duplicateIndexA])
		return fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, debugType, duplicateIndexA, duplicateIndexB)
	}

	oldArchetype := entityData.archetype

	componentIdsToAdd := make([]ComponentId, 0, len(componentIds))
	componentsToAdd := make([]AnyComponent, 0, len(components))
	for i, componentId := range componentIds {
		if oldArchetype.HasComponent(componentId) {
			resultErr = fmt.Errorf("%w: %s", ErrComponentAlreadyPresent, componentId.DebugString())
		} else {
			componentIdsToAdd = append(componentIdsToAdd, componentId)
			componentsToAdd = append(componentsToAdd, components[i])
		}
	}

	if len(componentIdsToAdd) == 0 {
		return resultErr
	}

	// move archetype
	newComponentIds := append(componentIdsToAdd, oldArchetype.componentIds...)
	requiredComponents := getAllRequiredComponents(&newComponentIds, componentsToAdd, world)

	newArchetype, err := world.archetypeStorage.getArchetype(world, newComponentIds)
	if err != nil {
		return err
	}

	var newRow uint
	var movedComponent *movedComponent

	for componentId, oldStorage := range oldArchetype.components {
		rawComponent, err := oldStorage.getComponentPointer(entityData.row)
		if err != nil {
			return err
		}

		newRow, err = newArchetype.components[componentId].insertRaw(world, rawComponent)
		if err != nil {
			return err
		}

		removeResult, err := oldStorage.remove(entityData.row)
		if err != nil {
			return err
		}

		movedComponent = removeResult
	}

	handleComponentStorageIndexMove(world, movedComponent, oldArchetype)
	err = oldArchetype.removeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to remove entity from old archetype: %w", err)
	}

	// insert new component
	for i, component := range componentsToAdd {
		storage := newArchetype.components[componentIdsToAdd[i]]
		newRow, err = storage.insert(world, component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert component %s in to component registry: %w", componentIdsToAdd[i].DebugString(), err)
			continue
		}
	}

	for _, component := range requiredComponents {
		componentId := ComponentIdOf(component, world)
		storage := newArchetype.components[componentId]
		newRow, err = storage.insert(world, component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert required component %s in to component registry: %w", componentId.DebugString(), err)
			continue
		}
	}

	entityData.archetype = newArchetype
	entityData.row = newRow
	world.archetypeStorage.entityIdToArchetype[entity] = newArchetype
	newArchetype.entities = append(newArchetype.entities, entity)

	return resultErr
}

// Insert adds the given components and all their required components (that the entity does not yet have) to the given entity.
//
// If the entity already has any of the given components, those components will overwrite the existing component. Its required
// components will be ignored.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error when the given entity does not exist
//   - Returns an ErrDuplicateComponent error when any of the given components are of the same type.
//   - Returns an ErrInvalidComponentStorageCapacity if the component storage capacity, that is decided through World
//     configs, is not valid
func InsertOrOverwrite(world *World, entity EntityId, components ...AnyComponent) (resultErr error) {
	if len(components) == 0 {
		return nil
	}

	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	componentIds := toComponentIds(components, world)

	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		debugType := ComponentDebugStringOf(components[duplicateIndexA])
		return fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, debugType, duplicateIndexA, duplicateIndexB)
	}

	oldArchetype := entityData.archetype

	componentIdsToAdd := make([]ComponentId, 0, len(componentIds))
	componentsToAdd := make([]AnyComponent, 0, len(components))
	for i, componentId := range componentIds {
		if oldArchetype.HasComponent(componentId) {
			err := oldArchetype.components[componentId].set(components[i], entityData.row)
			if err != nil {
				resultErr = err
			}
		} else {
			componentIdsToAdd = append(componentIdsToAdd, componentId)
			componentsToAdd = append(componentsToAdd, components[i])
		}
	}

	if len(componentIdsToAdd) == 0 {
		return resultErr
	}

	// move archetype
	newComponentIds := append(componentIdsToAdd, oldArchetype.componentIds...)
	requiredComponents := getAllRequiredComponents(&newComponentIds, componentsToAdd, world)

	newArchetype, err := world.archetypeStorage.getArchetype(world, newComponentIds)
	if err != nil {
		return err
	}

	var newRow uint
	var movedComponent *movedComponent

	for componentId, oldStorage := range oldArchetype.components {
		rawComponent, err := oldStorage.getComponentPointer(entityData.row)
		if err != nil {
			return err
		}

		newRow, err = newArchetype.components[componentId].insertRaw(world, rawComponent)
		if err != nil {
			return err
		}

		removeResult, err := oldStorage.remove(entityData.row)
		if err != nil {
			return err
		}

		movedComponent = removeResult
	}

	handleComponentStorageIndexMove(world, movedComponent, oldArchetype)
	err = oldArchetype.removeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to remove entity from old archetype: %w", err)
	}

	// insert new component
	for i, component := range componentsToAdd {
		storage := newArchetype.components[componentIdsToAdd[i]]
		newRow, err = storage.insert(world, component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert component %s in to component registry: %w", componentIdsToAdd[i].DebugString(), err)
			continue
		}
	}

	for _, component := range requiredComponents {
		componentId := ComponentIdOf(component, world)
		storage := newArchetype.components[componentId]
		newRow, err = storage.insert(world, component)
		if err != nil {
			resultErr = fmt.Errorf("failed to insert required component %s in to component registry: %w", componentId.DebugString(), err)
			continue
		}
	}

	entityData.archetype = newArchetype
	entityData.row = newRow
	world.archetypeStorage.entityIdToArchetype[entity] = newArchetype
	newArchetype.entities = append(newArchetype.entities, entity)

	return resultErr
}
