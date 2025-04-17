// functions to remove components from a specific entity
package ecs

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// Remove removes the given component from entity.
//
// Returns an ErrEntityNotFound error if the entity does not exist in world.
// Returns an ErrComponentNotFound error if the component is not present in the entity.
func Remove[A IComponent](world *world, entity EntityId) error {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	return removeComponentFromEntry[A](entry)
}

// Remove2 removes the given components from entity.
//
// Returns an ErrEntityNotFound error if the entity does not exist in world.
// Returns an ErrComponentNotFound error if any of the components are not present in entity, while still removing the ones that are present.
func Remove2[A, B IComponent](world *world, entity EntityId) (result error) {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntry[A](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[B](entry); err != nil {
		result = err
	}

	return result
}

// Remove3 removes the given components from entity.
//
// Returns an ErrEntityNotFound error if the entity does not exist in world.
// Returns an ErrComponentNotFound error if any of the components are not present in entity, while still removing the ones that are present.
func Remove3[A, B, C IComponent](world *world, entity EntityId) (result error) {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntry[A](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[B](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[C](entry); err != nil {
		result = err
	}

	return result
}

// Remove4 removes the given components from entity.
//
// Returns an ErrEntityNotFound error if the entity does not exist in world.
// Returns an ErrComponentNotFound error if any of the components are not present in entity, while still removing the ones that are present.
func Remove4[A, B, C, D IComponent](world *world, entity EntityId) (result error) {
	entry, ok := world.entities[entity]
	if !ok {
		return ErrEntityNotFound
	}

	if err := removeComponentFromEntry[A](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[B](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[C](entry); err != nil {
		result = err
	}
	if err := removeComponentFromEntry[D](entry); err != nil {
		result = err
	}

	return result
}

func removeComponentFromEntry[T IComponent](entry *entry) error {
	_, componentIndex, err := getComponentFromEntry[T](entry)

	if err != nil {
		return fmt.Errorf("%w: %s", err, getComponentDebugType[T]())
	}

	utils.RemoveFromSlice(&entry.components, componentIndex)
	return nil
}
