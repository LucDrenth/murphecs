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
	for entityToMatch, entry := range world.entities {
		if entity != entityToMatch {
			continue
		}

		for _, component := range entry.components {
			maybeA, ok := component.(A)
			if ok {
				a = &maybeA
				break
			}

			// its not the component we're after
		}

		if a == nil {
			return nil, fmt.Errorf("entity does not have component A (at generic position 1)")
		}

		return a, err
	}

	return nil, fmt.Errorf("entity not found")
}

// Get2 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get2[A IComponent, B IComponent](entity entityId, world *world) (a *A, b *B, err error) {
	for entityToMatch, entry := range world.entities {
		if entity != entityToMatch {
			continue
		}

		for _, component := range entry.components {
			if a == nil {
				if maybeA, ok := component.(A); ok {
					a = &maybeA
					continue
				}
			}
			if b == nil {
				if maybeB, ok := component.(B); ok {
					b = &maybeB
					continue
				}
			}

			// its not the component we're after
		}

		if a == nil {
			return nil, nil, fmt.Errorf("entity does not have component A (at generic position 1)")
		} else if b == nil {
			return nil, nil, fmt.Errorf("entity does not have component B (at generic position 2)")
		}

		return a, b, err
	}

	return nil, nil, fmt.Errorf("entity not found")
}

// Get3 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get3[A IComponent, B IComponent, C IComponent](entity entityId, world *world) (a *A, b *B, c *C, err error) {
	for entityToMatch, entry := range world.entities {
		if entity != entityToMatch {
			continue
		}

		for _, component := range entry.components {
			if a == nil {
				if maybeA, ok := component.(A); ok {
					a = &maybeA
					continue
				}
			}
			if b == nil {
				if maybeB, ok := component.(B); ok {
					b = &maybeB
					continue
				}
			}
			if c == nil {
				if maybeC, ok := component.(C); ok {
					c = &maybeC
					continue
				}
			}

			// its not the component we're after
		}

		if a == nil {
			return nil, nil, nil,
				fmt.Errorf("entity does not have component A (at generic position 1)")
		} else if b == nil {
			return nil, nil, nil,
				fmt.Errorf("entity does not have component B (at generic position 2)")
		} else if c == nil {
			return nil, nil, nil,
				fmt.Errorf("entity does not have component C (at generic position 3)")
		}

		return a, b, c, err
	}

	return nil, nil, nil,
		fmt.Errorf("entity not found")
}

// Get4 returns the component that belongs to the given entity.
// Returns an error if either the entity or any of the components is not found.
//
// Returns an "entity does not have component ..." error if duplicate components are given.
//
// WARNING: Do not store any of the component pointers
func Get4[A IComponent, B IComponent, C IComponent, D IComponent](entity entityId, world *world) (a *A, b *B, c *C, d *D, err error) {
	for entityToMatch, entry := range world.entities {
		if entity != entityToMatch {
			continue
		}

		for _, component := range entry.components {
			if a == nil {
				if maybeA, ok := component.(A); ok {
					a = &maybeA
					continue
				}
			}
			if b == nil {
				if maybeB, ok := component.(B); ok {
					b = &maybeB
					continue
				}
			}
			if c == nil {
				if maybeC, ok := component.(C); ok {
					c = &maybeC
					continue
				}
			}
			if d == nil {
				if maybeD, ok := component.(D); ok {
					d = &maybeD
					continue
				}
			}

			// its not the component we're after
		}

		if a == nil {
			return nil, nil, nil, nil,
				fmt.Errorf("entity does not have component A (at generic position 1)")
		} else if b == nil {
			return nil, nil, nil, nil,
				fmt.Errorf("entity does not have component B (at generic position 2)")
		} else if c == nil {
			return nil, nil, nil, nil,
				fmt.Errorf("entity does not have component C (at generic position 3)")
		} else if d == nil {
			return nil, nil, nil, nil,
				fmt.Errorf("entity does not have component D (at generic position 4)")
		}

		return a, b, c, d, err
	}

	return nil, nil, nil, nil,
		fmt.Errorf("entity not found")
}
