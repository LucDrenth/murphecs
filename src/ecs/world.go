package ecs

// world contains all of the entities and their components.
type world struct {
	entityIdCounter uint
	entities        map[EntityId]*entityData
	components      map[componentType]*componentRegistry
}

// NewWorld returns a world that can contain entities and components.
func NewWorld() world {
	return world{
		entities:   map[EntityId]*entityData{},
		components: map[componentType]*componentRegistry{},
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

func (world *world) createEntity() EntityId {
	world.entityIdCounter++
	entity := EntityId(world.entityIdCounter)
	world.entities[entity] = &entityData{components: map[componentType]uint{}}
	return entity
}

// getComponentRegistry creates a new component registry if it doesn't exist yet.
func (world *world) getComponentRegistry(componentType componentType) *componentRegistry {
	componentRegistry, ok := world.components[componentType]
	if !ok {
		newComponentRegistry := createComponentRegistry(1024, componentType)
		world.components[componentType] = &newComponentRegistry
		return &newComponentRegistry
	}

	return componentRegistry
}
