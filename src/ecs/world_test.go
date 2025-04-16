package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldGetEntry(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error if world is empty", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		entry, err := world.getEntry(nonExistingEntity)
		assert.Nil(entry)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if the given entity is not present", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()
		_, err := Spawn(&world)
		assert.NoError(err)

		entry, err := world.getEntry(nonExistingEntity)
		assert.Nil(entry)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully gets the right entry", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		_, err := Spawn(&world)
		assert.NoError(err)
		entity, err := Spawn(&world, componentA{})
		assert.NoError(err)
		_, err = Spawn(&world)
		assert.NoError(err)

		entry, err := world.getEntry(entity)
		assert.NoError(err)
		assert.Equal(1, entry.countComponents())
	})
}
