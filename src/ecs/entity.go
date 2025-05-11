package ecs

type EntityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = EntityId(0)

type EntityData struct {
	archetype *Archetype
	row       uint // index of archetype its component storages
}

func (e *EntityData) hasComponent(c ComponentId) bool {
	return e.archetype.HasComponent(c)
}
