package ecs

type entityId = uint

// This entityId can never exist in `world` because the inserted entityId's starts at 1.
// Useful for tests.
const nonExistingEntity = entityId(0)
