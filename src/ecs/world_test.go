package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct{}
type B struct{}

func TestWorldSpawn(t *testing.T) {
	t.Run("Successfully spawns", func(t *testing.T) {
		world := NewWorld()

		entity, err := world.Spawn()
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(1))
		entity, err = world.Spawn(A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(2))
		entity, err = world.Spawn(A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(3))
		entity, err = world.Spawn(B{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(4))
		entity, err = world.Spawn(A{}, B{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(5))
		entity, err = world.Spawn(B{}, A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(6))

		assert.Equal(t, uint(6), world.entityIdCounter)
	})

	t.Run("returns error if there are duplicate components", func(t *testing.T) {
		world := NewWorld()

		_, err := world.Spawn(A{}, A{})
		assert.Error(t, err)
		_, err = world.Spawn(A{}, A{}, A{})
		assert.Error(t, err)
		_, err = world.Spawn(A{}, A{}, B{})
		assert.Error(t, err)
		_, err = world.Spawn(A{}, B{}, A{})
		assert.Error(t, err)
		_, err = world.Spawn(B{}, A{}, A{})
		assert.Error(t, err)

		assert.Equal(t, uint(0), world.entityIdCounter)
	})
}
