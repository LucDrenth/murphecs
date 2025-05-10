// functions to get components for a given entity
package ecs

// Get1 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - Returns an ErrEntityNotFound error if the entity is not found.
//   - Returns an ErrComponentNotFound error if the entity does not have the component.
//
// WARNING: Do not store the component pointer
func Get1[A IComponent](world *World, entity EntityId) (a *A, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, err
	}

	return a, nil
}

// Get2 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get2[A, B IComponent](world *World, entity EntityId) (a *A, b *B, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

// Get3 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get3[A, B, C IComponent](world *World, entity EntityId) (a *A, b *B, c *C, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

// Get4 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get4[A, B, C, D IComponent](world *World, entity EntityId) (a *A, b *B, c *C, d *D, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &d); err != nil {
		return nil, nil, nil, nil, err
	}

	return a, b, c, d, nil
}

// Get5 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get5[A, B, C, D, E IComponent](world *World, entity EntityId) (a *A, b *B, c *C, d *D, e *E, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &d); err != nil {
		return nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &e); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, nil
}

// Get6 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get6[A, B, C, D, E, F IComponent](world *World, entity EntityId) (a *A, b *B, c *C, d *D, e *E, f *F, err error) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, nil
}

// Get7 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get7[A, B, C, D, E, F, G IComponent](world *World, entity EntityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, err error,
) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &g); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, nil
}

// Get8 returns the component that belongs to the given entity.
//
// Can return the following errors:
//   - ErrEntityNotFound error if the entity is not found.
//   - ErrComponentNotFound error if the entity does not have any of the components.
//
// Returns the same component pointer multiple times if multiple component of the same type are given.
//
// WARNING: Do not store any of the component pointers
func Get8[A, B, C, D, E, F, G, H IComponent](world *World, entity EntityId) (
	a *A, b *B, c *C, d *D, e *E, f *F, g *G, h *H, err error,
) {
	entityData, ok := world.entities[entity]
	if !ok {
		return nil, nil, nil, nil, nil, nil, nil, nil, ErrEntityNotFound
	}

	if err = setComponentFromEntry(world, entityData, &a); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &b); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &c); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &d); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &e); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &f); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &g); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	} else if err = setComponentFromEntry(world, entityData, &h); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return a, b, c, d, e, f, g, h, nil
}

// If a component of type T exists in entry, make target point to that component.
//
// Can return the following errors:
//   - ErrComponentNotFound error when the entity does not have a component of type T.
func setComponentFromEntry[T IComponent](world *World, entityData *EntityData, target **T) error {
	componentId := ComponentIdFor[T](world)

	componentRegistryIndex, ok := entityData.components[componentId]
	if !ok {
		return ErrComponentNotFound
	}

	result, err := getComponentFromComponentRegistry[T](world.components[componentId], componentRegistryIndex)
	if err != nil {
		return err
	}

	*target = result
	return nil
}
