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

	t.Run("uses ID config", func(t *testing.T) {
		assert := assert.New(t)

		worldId := WorldId(3)
		worldConfigs := DefaultWorldConfigs()
		worldConfigs.Id = &worldId

		world, err := NewWorld(worldConfigs)
		assert.NoError(err)
		assert.Equal(worldId, *world.Id())
	})
}

func TestGenerateEntityId(t *testing.T) {
	assert := assert.New(t)

	world := NewDefaultWorld()
	entity1 := world.generateEntityId()
	entity2 := world.generateEntityId()

	assert.NotEqual(entity1, entity2)
}

func TestStats(t *testing.T) {
	assert := assert.New(t)

	world := NewDefaultWorld()
	_, err := Spawn(&world, &emptyComponentA{}) // new archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentB{}) // new archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentC{}) // new archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentD{}) // new archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentA{}) // existing archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentA{}, &emptyComponentB{}) // new archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentB{}, &emptyComponentA{}) // existing archetype
	assert.NoError(err)
	_, err = Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}) // new archetype
	assert.NoError(err)

	assert.Equal(8, world.CountEntities())
	assert.Equal(12, world.CountComponents())
	assert.Equal(6, world.CountArchetypes())
}
