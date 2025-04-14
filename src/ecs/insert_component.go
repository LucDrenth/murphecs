package ecs

import (
	"fmt"
)

// Insert adds components to the given entity.
//
// Returns an ErrComponentAlreadyPresent error if any of the components is already present
// while still inserting the components that are not yet present.
func Insert(world *world, entity entityId, components ...IComponent) (err error) {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	for _, component := range components {
		if entry.containsComponentType(toComponentType(component)) {
			err = fmt.Errorf("%w: %s", ErrComponentAlreadyPresent, toComponentDebugType(component))
		} else {
			entry.components = append(entry.components, component)
		}
	}

	return err
}
