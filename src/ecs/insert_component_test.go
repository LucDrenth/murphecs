package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInsertComponentA struct{ Component }
type testInsertComponentB struct{ Component }
type testInsertComponentC struct{ Component }
type testInsertComponentWithFaultyRequiredComponent struct{ Component }
type testInsertComponentD struct{ Component }

func (c testInsertComponentB) RequiredComponents() []IComponent {
	return []IComponent{
		&testInsertComponentA{},
	}
}

func (c testInsertComponentC) RequiredComponents() []IComponent {
	return []IComponent{
		&testInsertComponentB{},
	}
}

func (c testInsertComponentWithFaultyRequiredComponent) RequiredComponents() []IComponent {
	return []IComponent{
		testInsertComponentD{}, // not passed by reference, should result in error
		&testInsertComponentA{},
	}
}

func TestInsert(t *testing.T) {
	type componentA struct{ Component }
	type componentB struct{ Component }
	type componentC struct{ Component }

	t.Run("no error when passing an empty list of components", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world)
		assert.NoError(err)

		err = Insert(&world, entity)
		assert.NoError(err)
	})

	t.Run("returns an error if the entity is not found", func(t *testing.T) {
		assert := assert.New(t)
		world := DefaultWorld()

		err := Insert(&world, nonExistingEntity)
		assert.ErrorIs(err, ErrEntityNotFound)
	})

	t.Run("returns an error if any of the components are already present", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &componentA{}, &componentB{})
		assert.NoError(err)

		// one component that is already present
		err = Insert(&world, entity, &componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// another component that is already present
		err = Insert(&world, entity, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present
		err = Insert(&world, entity, &componentA{}, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// all components already present - different order
		err = Insert(&world, entity, &componentB{}, &componentA{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present
		err = Insert(&world, entity, &componentB{}, &componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		// one component already and 1 component not present - different order
		err = Insert(&world, entity, &componentC{}, &componentB{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)
	})

	t.Run("if any component is already present, still inserts the other components that are not present", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &componentB{})
		assert.NoError(err)

		err = Insert(&world, entity, &componentA{}, &componentB{}, &componentC{})
		assert.ErrorIs(err, ErrComponentAlreadyPresent)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("if any component is not passed by reference, still inserts the other components that are passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &componentA{})
		assert.NoError(err)

		err = Insert(&world, entity, &componentB{}, componentC{}, &componentD{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts the components, and only to the given entity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entityA, err := Spawn(&world)
		assert.NoError(err)
		entityB, err := Spawn(&world, &componentB{})
		assert.NoError(err)

		err = Insert(&world, entityA, &componentA{}, &componentC{})
		assert.NoError(err)

		a, err := Get1[componentA](&world, entityA)
		assert.NoError(err)
		assert.NotNil(a)
		a, err = Get1[componentA](&world, entityB)
		assert.Error(err)
		assert.Nil(a)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("correctly inserts component and their required components", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world, &testInsertComponentA{})
		assert.NoError(err)

		err = Insert(&world, entity, &testInsertComponentC{})
		assert.NoError(err)

		assert.Equal(3, world.CountComponents())
	})

	t.Run("returns an error if any of the required components is not passed by reference, while still inserting the correctly passed required components", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		entity, err := Spawn(&world)
		assert.NoError(err)

		err = Insert(&world, entity, &testInsertComponentWithFaultyRequiredComponent{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)

		assert.Equal(2, world.CountComponents())
	})
}
