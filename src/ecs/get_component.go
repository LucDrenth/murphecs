// functions to get components for a given entity
package ecs

import (
	"fmt"
)

// Get returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if the component is not found.
//
// WARNING: Do not store the component pointer
func Get[A IComponent](world *world, entity EntityId) (a *A, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, err
	}

	return a, nil
}

// Get2 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get2[A, B IComponent](world *world, entity EntityId) (a *A, b *B, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

// Get3 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get3[A, B, C IComponent](world *world, entity EntityId) (a *A, b *B, c *C, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

// Get4 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get4[A, B, C, D IComponent](world *world, entity EntityId) (a *A, b *B, c *C, d *D, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d); err != nil {
		return nil, nil, nil, nil, err
	}

	return a, b, c, d, nil
}

// Get5 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get5[A, B, C, D, E IComponent](world *world, entity EntityId) (a *A, b *B, c *C, d *D, e *E, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, nil
}

// Get6 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get6[A, B, C, D, E, F IComponent](world *world, entity EntityId) (a *A, b *B, c *C, d *D, e *E, f *F, err error) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, nil
}

// Get7 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get7[A, B, C, D, E, F, G IComponent](world *world, entity EntityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, err error,
) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &g); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, nil
}

// Get8 returns the component that belongs to the given entity.
//
// Returns an ErrEntityNotFound error if the entity is not found.
// Returns an ErrComponentNotFound error if any of the components is not found.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get8[A, B, C, D, E, F, G, H IComponent](world *world, entity EntityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, err error,
) {
	entry, err := world.getEntry(entity)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	if err = setComponentFromEntry(entry, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &g); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(entry, &h); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, h, nil
}

// If a component of type T exists in entry, make target point to that component.
// If a component of type T does not exist in entry, return an error
func setComponentFromEntry[T IComponent](entry *entry, target **T) error {
	newTarget, _, err := getComponentFromEntry[T](entry)
	if err != nil {
		return fmt.Errorf("%w: entity does not have component %s", err, getComponentDebugType[T]())
	}

	*target = newTarget

	return nil
}
