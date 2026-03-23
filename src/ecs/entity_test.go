package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityExists(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns false for non-existing entity", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		assert.False(EntityExists(world, nonExistingEntity))
	})

	t.Run("returns true for spawned entity", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, &componentA{})
		assert.NoError(err)

		assert.True(EntityExists(world, entity))
	})

	t.Run("returns false after entity is despawned", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, &componentA{})
		assert.NoError(err)

		err = Despawn(world, entity)
		assert.NoError(err)

		assert.False(EntityExists(world, entity))
	})
}
