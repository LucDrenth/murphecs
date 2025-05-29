package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	t.Run("return an error if the entity was not found", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		err := Delete(world, nonExistingEntity)

		assert.Error(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the entity", func(t *testing.T) {
		type structA struct{ Component }

		assert := assert.New(t)

		world := NewDefaultWorld()
		entity1, err := Spawn(world, &structA{})
		assert.NoError(err)
		entity2, err := Spawn(world, &structA{})
		assert.NoError(err)
		entity3, err := Spawn(world, &structA{})
		assert.NoError(err)

		err = Delete(world, entity2)
		assert.NoError(err)

		// check that we can still get entity1 and entity3
		_, err = Get1[structA](world, entity1)
		assert.NoError(err)
		_, err = Get1[structA](world, entity3)
		assert.NoError(err)
	})
}
