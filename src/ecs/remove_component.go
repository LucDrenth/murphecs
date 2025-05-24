// functions to remove components from a specific entity
package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murphecs/src/utils"
)

// Remove1 removes the given component from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove1[A IComponent](world *World, entity EntityId) error {
	return removeComponents(world, entity, []ComponentId{
		ComponentIdFor[A](world),
	})
}

// Remove2 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove2[A, B IComponent](world *World, entity EntityId) (result error) {
	return removeComponents(world, entity, []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
	})
}

// Remove3 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove3[A, B, C IComponent](world *World, entity EntityId) (result error) {
	return removeComponents(world, entity, []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
	})
}

// Remove4 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove4[A, B, C, D IComponent](world *World, entity EntityId) (result error) {
	return removeComponents(world, entity, []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
	})
}

func removeComponents(world *World, entityId EntityId, componentIds []ComponentId) (resultErr error) {
	entity, ok := world.entities[entityId]
	if !ok {
		return ErrEntityNotFound
	}

	duplicate, duplicateIndexA, duplicateIndexB := utils.GetFirstDuplicate(componentIds)
	if duplicate != nil {
		return fmt.Errorf("%w: %s at positions %d and %d", ErrDuplicateComponent, duplicate.DebugString(), duplicateIndexA, duplicateIndexB)
	}

	componentIdsToRemove := make([]ComponentId, 0, len(componentIds))
	for _, componentId := range componentIds {
		if !entity.archetype.HasComponent(componentId) {
			resultErr = fmt.Errorf("%w: %s", ErrComponentNotFound, componentId.DebugString())
		} else {
			componentIdsToRemove = append(componentIdsToRemove, componentId)
		}
	}

	if len(componentIdsToRemove) == 0 {
		return resultErr
	}

	oldArchetype := entity.archetype

	newComponentIds := make([]ComponentId, 0, len(oldArchetype.componentIds))
	for _, componentId := range oldArchetype.componentIds {
		if !slices.Contains(componentIdsToRemove, componentId) {
			newComponentIds = append(newComponentIds, componentId)
		}
	}

	newArchetype, err := world.archetypeStorage.getArchetype(world, newComponentIds)
	if err != nil {
		return err
	}

	var newRow uint
	var movedComponent *movedComponent

	for componentId, oldStorage := range oldArchetype.components {
		if newArchetype.HasComponent(componentId) {
			rawComponent, err := oldStorage.getComponentPointer(entity.row)
			if err != nil {
				resultErr = err
				continue
			}

			storage := newArchetype.components[componentId]
			newRow, err = storage.insertRaw(world, rawComponent)
			if err != nil {
				resultErr = err
				continue
			}
		}

		err, removeResult := oldStorage.remove(entity.row)
		if err != nil {
			resultErr = err
			continue
		}

		movedComponent = removeResult
	}

	handleComponentStorageIndexMove(world, movedComponent, oldArchetype)
	err = oldArchetype.removeEntity(entityId)
	if err != nil {
		resultErr = fmt.Errorf("failed to remove entity from old archetype: %w", err)
	}

	entity.archetype = newArchetype
	entity.row = newRow
	world.archetypeStorage.entityIdToArchetype[entityId] = newArchetype
	newArchetype.entities = append(newArchetype.entities, entityId)

	return resultErr
}

func handleComponentStorageIndexMove(world *World, movedComponent *movedComponent, archetype *Archetype) {
	if movedComponent == nil {
		return
	}

	for _, entityId := range archetype.entities {
		entityData := world.entities[entityId]

		if entityData.row == movedComponent.fromIndex {
			entityData.row = movedComponent.toIndex
		}
	}
}
