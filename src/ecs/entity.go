package ecs

type EntityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = EntityId(0)

type EntityData struct {
	components map[ComponentId]uint // ComponentId --> componentRegistry index
}

func (e *EntityData) hasComponent(c ComponentId) bool {
	_, ok := e.components[c]
	return ok
}
