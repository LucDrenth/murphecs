package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorld(t *testing.T) {
	t.Run("default world does not panic", func(t *testing.T) {
		DefaultWorld()
	})

	t.Run("returns an error when using nil for componentCapacityStrategy", func(t *testing.T) {
		assert := assert.New(t)

		_, err := NewWorld(WorldConfigs{
			ComponentCapacityStrategy: nil,
		})
		assert.Error(err)
	})
}

func TestGenerateEntityId(t *testing.T) {
	assert := assert.New(t)

	world, err := NewWorld(DefaultWorldConfigs())
	assert.NoError(err)
	entity1 := world.generateEntityId()
	entity2 := world.generateEntityId()

	assert.NotEqual(entity1, entity2)
}
