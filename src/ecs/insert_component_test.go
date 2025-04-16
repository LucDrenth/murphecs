package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("no error when passing an empty list of components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world)
		assert.NoError(err)

		err = Insert(&world, entity)
		assert.NoError(err)
	})

	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		err := Insert(&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if any of the components are already present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, componentA{}, componentB{})
		assert.NoError(err)

		// one component that is already present
		err = Insert(&world, entity, componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// another component that is already present
		err = Insert(&world, entity, componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present
		err = Insert(&world, entity, componentA{}, componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present - different order
		err = Insert(&world, entity, componentB{}, componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present
		err = Insert(&world, entity, componentB{}, componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present - different order
		err = Insert(&world, entity, componentC{}, componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)
	})

	t.Run("if any component is already present, still inserts the other components that are not present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, componentB{})
		assert.NoError(err)

		err = Insert(&world, entity, componentA{}, componentB{}, componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts the components, and only to the given entity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entityA, err := Spawn(&world)
		assert.NoError(err)
		entityB, err := Spawn(&world, componentB{})
		assert.NoError(err)

		Insert(&world, entityA, componentA{}, componentC{})

		a, err := Get[componentA](&world, entityA)
		assert.NoError(err)
		assert.NotNil(a)
		a, err = Get[componentA](&world, entityB)
		assert.Error(err)
		assert.Nil(a)

		assert.Equal(3, world.CountComponents())
	})
}
