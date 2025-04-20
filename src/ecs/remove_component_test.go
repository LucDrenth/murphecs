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
		world := NewWorld()

		err := Remove[componentB](&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("return an error if the entity does not contain the component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)

		err = Remove[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})

	t.Run("successfully removes a component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)

		err = Remove[componentA](&world, entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := Get1[componentA](&world, entity)
		assert.Error(err)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := Get1[componentB](&world, entity)
		assert.NoError(err)
		assert.NotNil(b)
	})

	t.Run("successfully removes a component, no matter which one", func(t *testing.T) {
		assert := assert.New(t)

		world := NewWorld()
		entity, err := Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)

		err = Remove[componentB](&world, entity)
		assert.NoError(err)

		assert.Equal(1, world.CountComponents())
		assert.Equal(1, world.CountEntities())

		// can not fetch componentA, which was removed
		a, err := Get1[componentB](&world, entity)
		assert.Error(err)
		assert.Nil(a)

		// can still fetch componentB, which was not removed
		b, err := Get1[componentA](&world, entity)
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
		world := NewWorld()

		err := Remove2[componentA, componentB](&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if any of the components is not present in the entity, but still removes the other one", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		entity, err := Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		err = Remove2[componentB, componentC](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentA](&world, entity)
		assert.NoError(err)

		entity, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		err = Remove2[componentC, componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentA](&world, entity)
		assert.NoError(err)

		entity, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		err = Remove2[componentA, componentC](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.NoError(err)

		entity, err = Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)
		err = Remove2[componentC, componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.NoError(err)
	})

	t.Run("successfully removes the right components, no matter the order of the given components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		entity, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentA, componentB](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.NoError(err)

		entity, err = Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentB, componentC](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentA, componentC](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentC](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentB, componentA](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.NoError(err)

		entity, err = Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentC, componentB](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)

		entity, err = Spawn(&world, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)
		err = Remove2[componentC, componentA](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentC](&world, entity)
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
		world := NewWorld()

		err := Remove3[componentA, componentB, componentC](&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the right components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		entity, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{}, &componentD{})
		assert.NoError(err)
		err = Remove3[componentB, componentD, componentA](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentD](&world, entity)
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
		world := NewWorld()

		err := Remove4[componentA, componentB, componentC, componentD](&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("successfully removes the right components", func(t *testing.T) {
		assert := assert.New(t)
		world := NewWorld()

		entity, err := Spawn(&world, &componentA{}, &componentB{}, &componentC{}, &componentD{}, &componentE{})
		assert.NoError(err)
		err = Remove4[componentB, componentD, componentA, componentE](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentA](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentB](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentC](&world, entity)
		assert.NoError(err)
		_, err = Get1[componentD](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
		_, err = Get1[componentE](&world, entity)
		assert.ErrorIs(err, ErrComponentNotFound)
	})
}
