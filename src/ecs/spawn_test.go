package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpawn(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("error when passing nil for component", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, nil)
		assert.ErrorIs(err, ErrComponentIsNil)
		assert.Equal(nonExistingEntity, entity)

		assert.Equal(0, world.CountEntities())
		assert.Equal(0, world.CountComponents())
	})

	t.Run("succeeds when passing component by value", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, componentA{})
		assert.NoError(err)
		assert.NotEqual(nonExistingEntity, entity)

		assert.Equal(1, world.CountEntities())
		assert.Equal(1, world.CountComponents())
	})

	t.Run("success when passing 1 component by pointer and passing 1 component by value", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, componentA{}, &componentB{})
		assert.NoError(err)
		assert.NotEqual(nonExistingEntity, entity)

		// retry with different component order
		entity, err = Spawn(world, &componentA{}, componentB{})
		assert.NoError(err)
		assert.NotEqual(nonExistingEntity, entity)

		assert.Equal(2, world.CountEntities())
		assert.Equal(4, world.CountComponents())
	})

	t.Run("returns error if there are duplicate components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		_, err := Spawn(world, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(world, &componentA{}, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(world, &componentA{}, &componentA{}, &componentB{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(world, &componentA{}, &componentB{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)
		_, err = Spawn(world, &componentB{}, &componentA{}, &componentA{})
		assert.ErrorIs(err, ErrDuplicateComponent)

		assert.Equal(0, world.CountEntities())
		assert.Equal(0, world.CountComponents())
	})

	t.Run("successfully spawns", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world)
		assert.NoError(err)
		assert.Equal(entity, EntityId(1))
		entity, err = Spawn(world, &componentA{})
		assert.NoError(err)
		assert.Equal(entity, EntityId(2))
		entity, err = Spawn(world, &componentA{})
		assert.NoError(err)
		assert.Equal(entity, EntityId(3))
		entity, err = Spawn(world, &componentB{})
		assert.NoError(err)
		assert.Equal(entity, EntityId(4))
		entity, err = Spawn(world, &componentA{}, &componentB{})
		assert.NoError(err)
		assert.Equal(entity, EntityId(5))
		entity, err = Spawn(world, &componentB{}, &componentA{})
		assert.NoError(err)
		assert.Equal(entity, EntityId(6))

		assert.Equal(6, world.CountEntities())
		assert.Equal(7, world.CountComponents())
	})
}

type requiredComponentA struct{ Component }
type requiredComponentB struct{ Component }

type withRequiredComponents struct{ Component }

func (a withRequiredComponents) RequiredComponents() []IComponent {
	return []IComponent{
		requiredComponentA{},
		&requiredComponentB{},
	}
}

func TestSpawnWithRequiredComponents(t *testing.T) {
	t.Run("successfully spawns required components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := Spawn(world, &withRequiredComponents{})

		assert.NoError(err)
		assert.Equal(1, world.CountEntities())
		assert.Equal(3, world.CountComponents())

		a, b, c, err := Get3[requiredComponentA, requiredComponentB, withRequiredComponents](world, entity)
		assert.NotNil(a)
		assert.NotNil(b)
		assert.NotNil(c)
		assert.NoError(err)
	})
}
