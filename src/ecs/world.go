package ecs

// world contains all of the entities and their components.
type world struct {
	entityIdCounter uint
	entities        map[entityId]*entry
}

// NewWorld returns a world that can contain entities and components.
func NewWorld() world {
	return world{
		entities: map[entityId]*entry{},
	}
}

func (world *world) CountEntities() int {
	return len(world.entities)
}

func (world *world) CountComponents() int {
	result := 0

	for _, entry := range world.entities {
		result += len(entry.components)
	}

	return result
}

// getEntry returns the entry that correspond to entity, or an ErrEntityNotFound error if it wasn't found.
func (world *world) getEntry(entity entityId) (*entry, error) {
	entry, ok := world.entities[entity]
	if !ok {
		return nil, ErrEntityNotFound
	}

	return entry, nil
}
