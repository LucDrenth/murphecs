package ecs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldSpawn(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("successfully spawns", func(t *testing.T) {
		world := NewWorld()

		entity, err := Spawn(&world)
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(1))
		entity, err = Spawn(&world, componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(2))
		entity, err = Spawn(&world, componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(3))
		entity, err = Spawn(&world, componentB{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(4))
		entity, err = Spawn(&world, componentA{}, componentB{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(5))
		entity, err = Spawn(&world, componentB{}, componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, entityId(6))

		assert.Equal(t, 6, world.CountEntities())
		assert.Equal(t, 7, world.CountComponents())
	})

	t.Run("returns error if there are duplicate components", func(t *testing.T) {
		world := NewWorld()

		_, err := Spawn(&world, componentA{}, componentA{})
		assert.Error(t, err)
		_, err = Spawn(&world, componentA{}, componentA{}, componentA{})
		assert.Error(t, err)
		_, err = Spawn(&world, componentA{}, componentA{}, componentB{})
		assert.Error(t, err)
		_, err = Spawn(&world, componentA{}, componentB{}, componentA{})
		assert.Error(t, err)
		_, err = Spawn(&world, componentB{}, componentA{}, componentA{})
		assert.Error(t, err)

		assert.Equal(t, 0, world.CountEntities())
		assert.Equal(t, 0, world.CountComponents())
	})
}

type withRequiredComponents struct{ Component }

func (a withRequiredComponents) RequiredComponents() []IComponent {
	return []IComponent{componentA{}, componentB{}}
}

func TestRequiredComponents(t *testing.T) {
	t.Run("successfully spawns required components", func(t *testing.T) {
		world := NewWorld()

		_, err := Spawn(&world, withRequiredComponents{})

		assert.NoError(t, err)
		assert.Equal(t, 1, world.CountEntities())
		assert.Equal(t, 3, world.CountComponents())
	})
}

func TestGetComponentFromEntry(t *testing.T) {
	type componentA struct {
		value int
		Component
	}
	const expectedValueA = 101
	type componentB struct {
		value int
		Component
	}
	const expectedValueB = 102

	t.Run("gets the component if its present in entry", func(t *testing.T) {
		assert := assert.New(t)
		entry := entry{components: []IComponent{
			componentA{value: expectedValueA},
			componentB{value: expectedValueB},
		}}

		componentA, _, err := getComponentFromEntry[componentA](&entry)
		assert.NoError(err)
		assert.Equal(expectedValueA, (*componentA).value)

		componentB, _, err := getComponentFromEntry[componentB](&entry)
		assert.NoError(err)
		assert.Equal(expectedValueB, (*componentB).value)
	})

	t.Run("return an error if the entry does not contain the component", func(t *testing.T) {
		assert := assert.New(t)

		entry := entry{components: []IComponent{
			componentA{},
		}}

		componentA, _, err := getComponentFromEntry[componentB](&entry)
		assert.Error(err)
		assert.True(errors.Is(err, ErrComponentNotFound))
		assert.Nil(componentA)
	})
}
