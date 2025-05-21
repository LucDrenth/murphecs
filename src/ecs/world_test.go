package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorld(t *testing.T) {
	t.Run("default world does not panic", func(t *testing.T) {
		NewDefaultWorld()
	})

	t.Run("returns an error when using nil for ComponentCapacityStrategy", func(t *testing.T) {
		assert := assert.New(t)

		_, err := NewWorld(WorldConfigs{
			InitialComponentCapacityStrategy: nil,
		})
		assert.Error(err)
	})

	t.Run("returns an error when using nil for GrowComponentCapacityStrategy", func(t *testing.T) {
		assert := assert.New(t)

		_, err := NewWorld(WorldConfigs{
			InitialComponentCapacityStrategy: &StaticDefaultComponentCapacity{Capacity: 1024},
			ComponentCapacityGrowthStrategy:  nil,
		})
		assert.Error(err)
	})

	t.Run("succeeds when passing valid configs", func(t *testing.T) {
		assert := assert.New(t)

		_, err := NewWorld(WorldConfigs{
			InitialComponentCapacityStrategy: &StaticDefaultComponentCapacity{Capacity: 1024},
			ComponentCapacityGrowthStrategy:  &ComponentCapacityGrowthDouble{},
		})
		assert.NoError(err)
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
