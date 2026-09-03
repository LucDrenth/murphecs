package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasComponent(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns error when entity is not found", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := world.HasComponent[componentA](nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn()
		assert.NoError(err)

		result, err := world.HasComponent[componentA](entity)
		assert.NoError(err)
		assert.False(result)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{})
		assert.NoError(err)

		result, err := world.HasComponent[componentA](entity)
		assert.NoError(err)
		assert.True(result)
	})
}

func TestHasComponentId(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns error when entity is not found", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentId := ComponentIdFor[componentA](world)

		_, err := world.HasComponentId(nonExistingEntity, componentId)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn()
		assert.NoError(err)
		componentId := ComponentIdFor[componentA](world)

		result, err := world.HasComponentId(entity, componentId)
		assert.NoError(err)
		assert.False(result)
	})

	t.Run("returns false when entity does not have the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{})
		assert.NoError(err)
		componentId := ComponentIdFor[componentA](world)

		result, err := world.HasComponentId(entity, componentId)
		assert.NoError(err)
		assert.True(result)
	})
}
