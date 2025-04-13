package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("return an error if the entity does not exist", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()

		err := Remove[componentB](&world, nonExistingEntity)
		assert.Error(err)
		assert.True(errors.Is(err, ErrEntityNotFound))
	})

	t.Run("return an error if the entity does not contain the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, componentA{})
		assert.NoError(err)

		err = Remove[componentB](&world, entity)
		assert.Error(err)
		assert.True(errors.Is(err, ErrComponentNotFound))
	})

	t.Run("successfully removes a component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, componentA{}, componentB{})
		assert.NoError(err)

		err = Remove[componentA](&world, entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := Get[componentA](&world, entity)
		assert.Error(err)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := Get[componentB](&world, entity)
		assert.NoError(err)
		assert.NotNil(b)
	})

	t.Run("successfully removes a component, no matter which one", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, componentA{}, componentB{})
		assert.NoError(err)

		err = Remove[componentB](&world, entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := Get[componentB](&world, entity)
		assert.Error(err)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := Get[componentA](&world, entity)
		assert.NoError(err)
		assert.NotNil(b)
	})
}
