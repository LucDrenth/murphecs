// functions to remove components from a specific entity
package ecs

import "github.com/lucdrenth/murph/engine/src/utils"

// Remove removes the given component from entity.
//
// Returns an ErrEntityNotFound error if the entity does not exist in world.
// Returns an ErrComponentNotFound error if the entity does not have the component.
func Remove[T IComponent](world *world, entity entityId) error {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	_, componentIndex, err := getComponentFromEntry[T](entry)
	if err != nil {
		return err
	}

	utils.RemoveFromSlice(&entry.components, componentIndex)

	return nil
}
