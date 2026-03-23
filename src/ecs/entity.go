package ecs

type EntityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = EntityId(0)

type EntityData struct {
	archetype *Archetype
	row       uint              // index of archetype its component storages
	observers *observerRegistry // lazily initialized to preserve space
}

func (e *EntityData) hasComponent(c ComponentId) bool {
	return e.archetype.HasComponent(c)
}

// EntityExists returns whether the entity is currently in the world.
// It returns false if the entity did exist but got despawned.
func EntityExists(world *World, entity EntityId) bool {
	_, exists := world.entities[entity]
	return exists
}
