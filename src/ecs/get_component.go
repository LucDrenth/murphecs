// functions to get components for a given entity
package ecs

import (
	"fmt"
)

// Get returns the component that belongs to the given entity.
// Returns an error if either the entity or the component is not found.
//
// WARNING: Do not store the component pointer
func Get[A IComponent](world *world, entity entityId) (a *A, err error) {
	entry, err := getEntry(world, entity)
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
func Get2[A IComponent, B IComponent](world *world, entity entityId) (a *A, b *B, err error) {
	entry, err := getEntry(world, entity)
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
func Get3[A IComponent, B IComponent, C IComponent](world *world, entity entityId) (a *A, b *B, c *C, err error) {
	entry, err := getEntry(world, entity)
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
func Get4[A IComponent, B IComponent, C IComponent, D IComponent](world *world, entity entityId) (a *A, b *B, c *C, d *D, err error) {
	entry, err := getEntry(world, entity)
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

// Get5 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get5[A IComponent, B IComponent, C IComponent, D IComponent, E IComponent](world *world, entity entityId) (a *A, b *B, c *C, d *D, e *E, err error) {
	entry, err := getEntry(world, entity)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d, 4); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e, 5); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, nil
}

// Get6 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get6[A IComponent, B IComponent, C IComponent, D IComponent, E IComponent, F IComponent](world *world, entity entityId) (a *A, b *B, c *C, d *D, e *E, f *F, err error) {
	entry, err := getEntry(world, entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d, 4); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e, 5); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f, 6); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, nil
}

// Get7 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get7[A IComponent, B IComponent, C IComponent, D IComponent, E IComponent, F IComponent, G IComponent](world *world, entity entityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, err error,
) {
	entry, err := getEntry(world, entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d, 4); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e, 5); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f, 6); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &g, 7); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, nil
}

// Get8 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get8[A IComponent, B IComponent, C IComponent, D IComponent, E IComponent, F IComponent, G IComponent, H IComponent](world *world, entity entityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, err error,
) {
	entry, err := getEntry(world, entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a, 1); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b, 2); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c, 3); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d, 4); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e, 5); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f, 6); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &g, 7); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &h, 8); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, h, nil
}

// getEntry returns the entry that correspond to entity, or an error if it wasn't found.
func getEntry(world *world, entity entityId) (*entry, error) {
	entry, ok := world.entities[entity]
	if !ok {
		return nil, ErrEntityNotFound
	}

	return entry, nil
}

// If a component of type T exists in entry, make target point to that component.
// If a component of type T does not exist in entry, return an error
func setComponentFromEntry[T IComponent](entry *entry, target **T, genericPosition int) error {
	newTarget, _, err := getComponentFromEntry[T](entry)
	if err != nil {
		return fmt.Errorf("%w: entity does not have component at generic position %d", err, genericPosition)
	}

	*target = newTarget

	return nil

}
