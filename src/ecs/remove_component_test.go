package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }

	t.Run("return an error if the entity does not exist", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.Remove1[componentB](nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("return an error if the entity does not contain the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{})
		assert.NoError(err)

		err = world.Remove1[componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})

	t.Run("successfully removes a component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)

		err = world.Remove1[componentA](entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := world.Get1[*componentB](entity)
		assert.NoError(err)
		assert.NotNil(b)
	})

	t.Run("successfully removes a component, no matter which one", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)

		err = world.Remove1[componentB](entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := world.Get1[*componentA](entity)
		assert.NoError(err)
		assert.NotNil(b)
	})
}

func TestRemove2(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("return an error if the entity does not exist", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.Remove2[componentA, componentB](nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if any of the components is not present in the entity, but still removes the other one", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)
		err = world.Remove2[componentB, componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentA](entity)
		assert.NoError(err)

		entity, err = world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)
		err = world.Remove2[componentC, componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentA](entity)
		assert.NoError(err)

		entity, err = world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)
		err = world.Remove2[componentA, componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.NoError(err)

		entity, err = world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)
		err = world.Remove2[componentC, componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.NoError(err)
	})

	t.Run("successfully removes the right components, no matter the order of the given components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentA, componentB](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.NoError(err)

		entity, err = world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentB, componentC](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentA, componentC](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentB, componentA](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.NoError(err)

		entity, err = world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentC, componentB](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = world.Spawn(&componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = world.Remove2[componentC, componentA](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentC](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})
}

func TestRemove3(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }

	t.Run("return an error if the entity does not exist", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.Remove3[componentA, componentB, componentC](nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the right components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := world.Spawn(&componentA{}, &componentB{}, &componentC{}, &componentD{})
		assert.NoError(err)
		err = world.Remove3[componentB, componentD, componentA](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentD](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})
}

func TestRemove4(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }
	type componentE struct{ Component }

	t.Run("return an error if the entity does not exist", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.Remove4[componentA, componentB, componentC, componentD](nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the right components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		entity, err := world.Spawn(&componentA{}, &componentB{}, &componentC{}, &componentD{}, &componentE{})
		assert.NoError(err)
		err = world.Remove4[componentB, componentD, componentA, componentE](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentA](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentB](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentC](entity)
		assert.NoError(err)
		_, err = world.Get1[*componentD](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = world.Get1[*componentE](entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})
}
