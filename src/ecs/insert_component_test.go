package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInsertComponentA struct{ Component }
type testInsertComponentB struct{ Component }
type testInsertComponentC struct{ Component }
type testInsertComponentD struct{ Component }

func (c testInsertComponentB) RequiredComponents() []AnyComponent {
	return []AnyComponent{
		&testInsertComponentA{},
		testInsertComponentD{},
	}
}

func (c testInsertComponentC) RequiredComponents() []AnyComponent {
	return []AnyComponent{
		&testInsertComponentB{},
	}
}

func TestInsert(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("no error when passing an empty list of components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn()
		assert.NoError(err)

		err = world.Insert(entity)
		assert.NoError(err)
	})

	t.Run("error when any of the given components are nil", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn()
		assert.NoError(err)

		// only 1 nil
		err = world.Insert(entity, nil)
		assert.ErrorIs(err, ErrComponentIsNil)

		// 1 valid, 1 nil
		err = world.Insert(entity, &componentA{}, nil)
		assert.ErrorIs(err, ErrComponentIsNil)

		// 1 nil, 1 valid
		err = world.Insert(entity, nil, &componentA{})
		assert.ErrorIs(err, ErrComponentIsNil)

		// 1 nil, 1 valid, 1 nil
		err = world.Insert(entity, nil, &componentA{}, nil)
		assert.ErrorIs(err, ErrComponentIsNil)

		assert.Equal(1, world.CountEntities())
		assert.Equal(0, world.CountComponents())
	})

	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.Insert(nonExistingEntity, &componentA{})
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if any of the components are already present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{}, &componentB{})
		assert.NoError(err)

		// one component that is already present
		err = world.Insert(entity, &componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// another component that is already present
		err = world.Insert(entity, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present
		err = world.Insert(entity, &componentA{}, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present - different order
		err = world.Insert(entity, &componentB{}, &componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present
		err = world.Insert(entity, &componentB{}, &componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present - different order
		err = world.Insert(entity, &componentC{}, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)
	})

	t.Run("if any component is already present, still inserts the other components that are not present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentB{})
		assert.NoError(err)

		err = world.Insert(entity, &componentA{}, &componentB{}, &componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts the components, and only to the given entity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entityA, err := world.Spawn()
		assert.NoError(err)
		entityB, err := world.Spawn(&componentB{})
		assert.NoError(err)

		err = world.Insert(entityA, &componentA{}, &componentC{})
		assert.NoError(err)

		a, err := world.Get1[*componentA](entityA)
		assert.NoError(err)
		assert.NotNil(a)
		a, err = world.Get1[*componentA](entityB)
		assert.Error(err)
		assert.Nil(a)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts component and their required components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&testInsertComponentA{})
		assert.NoError(err)

		err = world.Insert(entity, &testInsertComponentC{})
		assert.NoError(err)

		assert.Equal(4, world.CountComponents())
	})
}

func TestInsertOrOverwrite(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }
	type componentD struct{ Component }

	t.Run("no error when passing an empty list of components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn()
		assert.NoError(err)

		err = world.InsertOrOverwrite(entity)
		assert.NoError(err)
	})

	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		assert := assert.New(t)
		world := NewDefaultWorld()

		err := world.InsertOrOverwrite(nonExistingEntity, &componentA{})
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("overwrites component if any of the components is already present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{}, &componentWithValueA{value: 10})
		assert.NoError(err)

		err = world.InsertOrOverwrite(entity, &componentB{}, &componentWithValueA{value: 20})
		assert.NoError(err)
		component, err := world.Get1[*componentWithValueA](entity)
		assert.NoError(err)
		assert.Equal(20, component.value)

		// try again with different component order
		err = world.InsertOrOverwrite(entity, &componentWithValueA{value: 30}, &componentC{})
		assert.NoError(err)
		component, err = world.Get1[*componentWithValueA](entity)
		assert.NoError(err)
		assert.Equal(30, component.value)
	})

	t.Run("if any component is already present, still inserts the other components that are not present", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentB{})
		assert.NoError(err)

		err = world.InsertOrOverwrite(entity, &componentA{}, &componentB{}, &componentC{})
		assert.NoError(err)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("succeeds if some of the components are not passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&componentA{})
		assert.NoError(err)

		err = world.InsertOrOverwrite(entity, &componentB{}, componentC{}, &componentD{})
		assert.NoError(err)

		assert.Equal(4, world.CountComponents())
	})

	t.Run("correctly inserts the components, and only to the given entity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entityA, err := world.Spawn()
		assert.NoError(err)
		entityB, err := world.Spawn(&componentB{})
		assert.NoError(err)

		err = world.InsertOrOverwrite(entityA, &componentA{}, &componentC{})
		assert.NoError(err)

		a, err := world.Get1[*componentA](entityA)
		assert.NoError(err)
		assert.NotNil(a)
		a, err = world.Get1[*componentA](entityB)
		assert.Error(err)
		assert.Nil(a)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts component and their required components", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		entity, err := world.Spawn(&testInsertComponentA{})
		assert.NoError(err)

		err = world.InsertOrOverwrite(entity, &testInsertComponentC{})
		assert.NoError(err)

		assert.Equal(4, world.CountComponents())
	})
}
