// functions to get components for a given entity
package ecs

import (
	"fmt"
)

type Luc struct{ Component }

// Get returns the component that belongs to the given entity.
// Returns an error if either the entity or the component is not found.
//
// WARNING: Do not store the component pointer
func Get[A IComponent](entity entityId, world *world) (a *A, err error) {
	entry, err := getEntry(entity, world)
	if err != nil {
		return nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, err
	}

	return a, nil
}

// Get2 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get2[A IComponent, B IComponent](entity entityId, world *world) (a *A, b *B, err error) {
	entry, err := getEntry(entity, world)
	if err != nil {
		return nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

// Get3 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get3[A IComponent, B IComponent, C IComponent](entity entityId, world *world) (a *A, b *B, c *C, err error) {
	entry, err := getEntry(entity, world)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

// Get4 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get4[A IComponent, B IComponent, C IComponent, D IComponent](entity entityId, world *world) (a *A, b *B, c *C, d *D, err error) {
	entry, err := getEntry(entity, world)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d, 4); err != nil {
		return nil, nil, nil, nil, err
	}

	return a, b, c, d, nil
}

func getEntry(entity entityId, world *world) (*entry, error) {
	entry, ok := world.entities[entity]
	if !ok {
		return nil, fmt.Errorf("entity not found")
	}

	return &entry, nil
}

// If a component of type T exists in entry, make target point to that component.
// If a component of type T does not exist in entry, return an error
func setComponentFromEntry[T IComponent](entry *entry, target **T, genericPosition int) error {
	for _, component := range entry.components {
		if maybeTarget, ok := component.(T); ok {
			*target = &maybeTarget
			return nil
		}
	}

	return fmt.Errorf("entity does not have component at generic position %d", genericPosition)
}
