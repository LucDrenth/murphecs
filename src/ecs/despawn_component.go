package ecs

import "fmt"

// Despawn removes an entity from the world.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity did not exist in the world.
//   - ErrWorldIsLocked error while querying
func Despawn(world *World, entity EntityId) error {
	if world.isQuerying {
		// Prevent messing with query results
		return ErrWorldIsLocked
	}

	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	componentIds := entityData.archetype.componentIds

	err := entityData.archetype.removeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to remove entity from archetype: %w", err)
	}

	delete(world.archetypeStorage.entityIdToArchetype, entity)
	delete(world.entities, entity)

	world.observers.triggerDespawnObservers(world, componentIds, entity)
	if entityData.observers != nil {
		entityData.observers.triggerDespawnObservers(world, componentIds, entity)
	}

	return nil
}
