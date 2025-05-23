package ecs

import "fmt"

// Delete removes an entity from the world.
//
// Returns an ErrEntityNotFound error if the entity did not exist in the world.
func Delete(world *World, entity EntityId) error {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	err := entityData.archetype.removeEntity(entity)
	if err != nil {
		return fmt.Errorf("failed to remove entity from archetype: %w", err)
	}

	delete(world.archetypeStorage.entityIdToArchetype, entity)
	delete(world.entities, entity)

	return nil
}
