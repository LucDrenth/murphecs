package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpawn(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("successfully spawns", func(t *testing.T) {
		world := NewDefaultWorld()

		entity, err := Spawn(&world)
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(1))
		entity, err = Spawn(&world, &componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(2))
		entity, err = Spawn(&world, &componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(3))
		entity, err = Spawn(&world, &componentB{})
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(4))
		entity, err = Spawn(&world, &componentA{}, &componentB{})
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(5))
		entity, err = Spawn(&world, &componentB{}, &componentA{})
		assert.Nil(t, err)
		assert.Equal(t, entity, EntityId(6))

		assert.Equal(t, 6, world.CountEntities())
		assert.Equal(t, 7, world.CountComponents())
	})

	t.Run("returns error if there are duplicate components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		_, err := Spawn(&world, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(&world, &componentA{}, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(&world, &componentA{}, &componentA{}, &componentB{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(&world, &componentA{}, &componentB{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(&world, &componentB{}, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)

		assert.Equal(0, world.CountEntities())
		assert.Equal(0, world.CountComponents())
	})
}

type requiredComponentA struct{ Component }
type requiredComponentB struct{ Component }

type withRequiredComponents struct{ Component }

func (a withRequiredComponents) RequiredComponents() []IComponent {
	return []IComponent{&requiredComponentA{}, &requiredComponentB{}}
}

func TestSpawnWithRequiredComponents(t *testing.T) {
	t.Run("successfully spawns required components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(&world, &withRequiredComponents{})

		assert.NoError(err)
		assert.Equal(1, world.CountEntities())
		assert.Equal(3, world.CountComponents())

		a, b, c, err := Get3[requiredComponentA, requiredComponentB, withRequiredComponents](&world, entity)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.NotNil(c)
		assert.NoError(err)
	})
}
