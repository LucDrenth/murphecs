package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasComponent(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns error when entity is not found", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		_, err := HasComponent[componentA](&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world)
		assert.NoError(err)

		result, err := HasComponent[componentA](&world, entity)
		assert.NoError(err)
		assert.False(result)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)

		result, err := HasComponent[componentA](&world, entity)
		assert.NoError(err)
		assert.True(result)
	})
}

func TestHasComponentId(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns error when entity is not found", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentId := ComponentIdFor[componentA](&world)

		_, err := HasComponentId(&world, nonExistingEntity, componentId)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world)
		assert.NoError(err)
		componentId := ComponentIdFor[componentA](&world)

		result, err := HasComponentId(&world, entity, componentId)
		assert.NoError(err)
		assert.False(result)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)
		componentId := ComponentIdFor[componentA](&world)

		result, err := HasComponentId(&world, entity, componentId)
		assert.NoError(err)
		assert.True(result)
	})
}
