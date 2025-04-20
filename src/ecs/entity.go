package ecs

type EntityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = EntityId(0)

type entityData struct {
	components map[componentType]uint // componentType --> componentRegistry index
}

func (e *entityData) hasComponent(c componentType) bool {
	_, ok := e.components[c]
	return ok
}
