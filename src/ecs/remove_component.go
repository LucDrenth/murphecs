// functions to remove components from a specific entity
package ecs

import (
	"fmt"
)

// Remove removes the given component from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove[A IComponent](world *World, entity EntityId) error {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	return removeComponentFromEntityData[A](entityData, world)
}

// Remove2 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove2[A, B IComponent](world *World, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData, world); err != nil {
		result = err
	}

	return result
}

// Remove3 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove3[A, B, C IComponent](world *World, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[C](entityData, world); err != nil {
		result = err
	}

	return result
}

// Remove4 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove4[A, B, C, D IComponent](world *World, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[C](entityData, world); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[D](entityData, world); err != nil {
		result = err
	}

	return result
}

func removeComponentFromEntityData[T IComponent](entry *EntityData, world *World) error {
	componentId := ComponentIdFor[T](world)
	if _, ok := entry.components[componentId]; !ok {
		return fmt.Errorf("%w: %s", ErrComponentNotFound, ComponentDebugStringFor[T]())
	}

	delete(entry.components, componentId)
	return nil
}
