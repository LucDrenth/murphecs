package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct{}
type B struct{}

func TestWorldInsert(t *testing.T) {
	t.Run("Successfully inserts components", func(t *testing.T) {
		world := NewWorld()

		entity, err := world.Insert()
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(1))
		entity, err = world.Insert(A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(2))
		entity, err = world.Insert(A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(3))
		entity, err = world.Insert(B{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(4))
		entity, err = world.Insert(A{}, B{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(5))
		entity, err = world.Insert(B{}, A{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(6))

		assert.Equal(t, uint(6), world.entityIdCounter)
	})

	t.Run("returns error if there are duplicate components", func(t *testing.T) {
		world := NewWorld()

		_, err := world.Insert(A{}, A{})
		assert.Error(t, err)
		_, err = world.Insert(A{}, A{}, A{})
		assert.Error(t, err)
		_, err = world.Insert(A{}, A{}, B{})
		assert.Error(t, err)
		_, err = world.Insert(A{}, B{}, A{})
		assert.Error(t, err)
		_, err = world.Insert(B{}, A{}, A{})
		assert.Error(t, err)

		assert.Equal(t, uint(0), world.entityIdCounter)
	})
}
