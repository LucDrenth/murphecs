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

func TestCreateEntity(t *testing.T) {
	assert := assert.New(t)

	world, err := NewWorld(DefaultWorldConfigs())
	assert.NoError(err)
	entity1 := world.createEntity()
	entity2 := world.createEntity()

	assert.NotEqual(entity1, entity2)
}

func TestGetComponentRegistry(t *testing.T) {
	type componentA struct{ Component }

	assert := assert.New(t)

	// create component registry if it is not present yet
	world, err := NewWorld(DefaultWorldConfigs())
	assert.NoError(err)
	componentRegistry, err := world.getComponentRegistry(ComponentIdFor[componentA](&world))
	assert.NoError(err)
	assert.NotNil(componentRegistry)

	// get the same component registry if it is already present
	componentRegistry2, err := world.getComponentRegistry(ComponentIdFor[componentA](&world))
	assert.NotNil(componentRegistry)
	assert.NoError(err)
	assert.Equal(componentRegistry, componentRegistry2)
}
