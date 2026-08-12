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

func TestGetComponentsForEntity(t *testing.T) {
	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := world.GetComponentsForEntity(nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns all components that belong to the entity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := Spawn(world, &componentA{value: expectedValueA}, &componentB{value: expectedValueB})
		assert.NoError(err)

		components, err := world.GetComponentsForEntity(entity)
		assert.NoError(err)
		assert.Len(components, 2)

		idA := ComponentIdFor[componentA](world)
		idB := ComponentIdFor[componentB](world)

		a, ok := components[idA].(componentA)
		assert.True(ok)
		assert.Equal(expectedValueA, a.value)

		b, ok := components[idB].(componentB)
		assert.True(ok)
		assert.Equal(expectedValueB, b.value)
	})

	t.Run("returned components are not affected by later mutations", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := Spawn(world, &componentA{value: expectedValueA})
		assert.NoError(err)

		components, err := world.GetComponentsForEntity(entity)
		assert.NoError(err)

		a, err := Get1[*componentA](world, entity)
		assert.NoError(err)
		a.value += 1

		idA := ComponentIdFor[componentA](world)
		snapshot, ok := components[idA].(componentA)
		assert.True(ok)
		assert.Equal(expectedValueA, snapshot.value)
	})
}

func TestStats(t *testing.T) {
	t.Run("world returns the correct stats after inserting", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		_, err := Spawn(world, &emptyComponentA{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentB{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentC{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentD{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentA{}) // existing archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentA{}, &emptyComponentB{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentB{}, &emptyComponentA{}) // existing archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}) // new archetype
		assert.NoError(err)

		assert.Equal(8, world.CountEntities())
		assert.Equal(12, world.CountComponents())
		assert.Equal(6, world.CountArchetypes())
	})

	t.Run("stats do not change when inserting, removing, spawning and deleting with existing archetypes", func(t *testing.T) {
		assert := assert.New(t)

		// Spawn components to create the following archetypes:
		// 	- emptyComponentA
		// 	- emptyComponentB
		// 	- emptyComponentA + emptyComponentB
		world := NewDefaultWorld()
		_, err := Spawn(world, &emptyComponentA{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentB{}) // new archetype
		assert.NoError(err)
		_, err = Spawn(world, &emptyComponentA{}, &emptyComponentB{}) // new archetype
		assert.NoError(err)

		// Spawning and then deleting an entity does not alter stats
		{
			originalStats := world.Stats()

			entity, err := Spawn(world, &emptyComponentA{})
			assert.NoError(err)
			err = Despawn(world, entity)
			assert.NoError(err)

			assert.Equal(originalStats, world.Stats())
		}

		// Inserting and then removing a component does not alter stats
		{
			entity, err := Spawn(world, &emptyComponentA{})
			assert.NoError(err)
			originalStats := world.Stats()

			err = Insert(world, entity, &emptyComponentB{})
			assert.NoError(err)
			err = Remove1[emptyComponentB](world, entity)
			assert.NoError(err)

			assert.Equal(originalStats, world.Stats())
		}
	})
}
