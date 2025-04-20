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
func Remove[A IComponent](world *world, entity EntityId) error {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	return removeComponentFromEntityData[A](entityData)
}

// Remove2 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove2[A, B IComponent](world *world, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData); err != nil {
		result = err
	}

	return result
}

// Remove3 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove3[A, B, C IComponent](world *world, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[C](entityData); err != nil {
		result = err
	}

	return result
}

// Remove4 removes the given components from entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity does not exist in world.
//   - ErrComponentNotFound error if the component is not present in the entity.
func Remove4[A, B, C, D IComponent](world *world, entity EntityId) (result error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntityData[A](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[B](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[C](entityData); err != nil {
		result = err
	}
	if err := removeComponentFromEntityData[D](entityData); err != nil {
		result = err
	}

	return result
}

func removeComponentFromEntityData[T IComponent](entry *entityData) error {
	componentType := getComponentType[T]()
	if _, ok := entry.components[componentType]; !ok {
		return fmt.Errorf("%w: %s", ErrComponentNotFound, getComponentDebugType[T]())
	}

	delete(entry.components, componentType)
	return nil
}
