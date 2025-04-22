package ecs

// Delete removes an entity from the world if the entity exists in the world.
//
// Returns an ErrEntityNotFound error if the entity did not exist in the world.
func Delete(world *World, entity EntityId) error {
	if _, ok := world.entities[entity]; !ok {
		return ErrEntityNotFound
	}

	delete(world.entities, entity)
	return nil
}
