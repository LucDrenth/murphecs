package ecs

import "fmt"

// Despawn removes an entity from the world.
//
// Returns an ErrEntityNotFound error if the entity did not exist in the world.
func Despawn(world *World, entity EntityId) error {
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

	handleDespawnObservers(world, componentIds, entity)

	return nil
}

func handleDespawnObservers(world *World, componentIds []ComponentId, entity EntityId) {
	for _, componentId := range componentIds {
		observers, exists := world.despawnObservers[componentId]
		if !exists {
			continue
		}

		for _, observer := range observers {
			observer(world, entity)
		}
	}
}
