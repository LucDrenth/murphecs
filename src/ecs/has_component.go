package ecs

// HasComponent returns wether entity has component C.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error if the entity is not found.
func HasComponent[C IComponent](world *World, entity EntityId) (bool, error) {
	entityData, exists := world.entities[entity]
	if !exists {
		return false, ErrEntityNotFound
	}

	return entityData.archetype.HasComponent(ComponentIdFor[C](world)), nil
}

// HasComponentId returns wether entity has a component with id componentId. This is more performant than
// HasComponent if you use this method multiple times with the same componentId because the componentId does
// not have to be calculated every call.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error if the entity is not found.
func HasComponentId(world *World, entity EntityId, componentId ComponentId) (bool, error) {
	entityData, exists := world.entities[entity]
	if !exists {
		return false, ErrEntityNotFound
	}

	return entityData.archetype.HasComponent(componentId), nil
}
