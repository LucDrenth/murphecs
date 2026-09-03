package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDespawn(t *testing.T) {
	t.Run("return an error if the entity was not found", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		err := world.Despawn(nonExistingEntity)

		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the entity", func(t *testing.T) {
		type structA struct{ Component }

		assert := assert.New(t)

		world := NewDefaultWorld()
		entity1, err := world.Spawn(&structA{})
		assert.NoError(err)
		entity2, err := world.Spawn(&structA{})
		assert.NoError(err)
		entity3, err := world.Spawn(&structA{})
		assert.NoError(err)

		err = world.Despawn(entity2)
		assert.NoError(err)

		// check that we can still get entity1 and entity3
		_, err = world.Get1[structA](entity1)
		assert.NoError(err)
		_, err = world.Get1[structA](entity3)
		assert.NoError(err)
	})
}
