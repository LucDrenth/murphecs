package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComponentRegistry(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when using capacity of 0", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		_, err := createComponentRegistry(0, ComponentIdFor[componentA](&world))
		assert.ErrorIs(err, ErrInvalidComponentStorageCapacity)
	})

	t.Run("does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		_, err := createComponentRegistry(1, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = createComponentRegistry(1024, ComponentIdFor[componentA](&world))
		assert.NoError(err)
	})
}

func TestComponentRegistryInsert(t *testing.T) {
	type componentA struct{ Component }

	t.Run("fails when component is not a pointer", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentRegistry, err := createComponentRegistry(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = componentRegistry.insert(componentA{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)
	})

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentRegistry, err := createComponentRegistry(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = componentRegistry.insert(&componentA{})
		assert.NoError(err)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentRegistry, err := createComponentRegistry(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for range 10 {
			_, err = componentRegistry.insert(&componentA{})
			assert.NoError(err)
		}
	})
}

func TestGetComponentFromComponentRegistry(t *testing.T) {
	type componentA struct {
		value int
		Component
	}

	t.Run("returns an error if index out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentRegistry, err := createComponentRegistry(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		component, err := getComponentFromComponentRegistry[componentA](&componentRegistry, 5)
		assert.Error(err)
		assert.Nil(component)
	})

	t.Run("gets the correct components when not exceeding capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := 4

		componentRegistry, err := createComponentRegistry(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range capacity {
			_, err = componentRegistry.insert(&componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range capacity {
			component, err := getComponentFromComponentRegistry[componentA](&componentRegistry, uint(i))
			assert.NoError(err)
			assert.NotNil(component)
			assert.Equal(i, component.value)
		}
	})

	t.Run("gets the correct components when exceeding capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := 4
		numberOfInserts := 15

		componentRegistry, err := createComponentRegistry(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range numberOfInserts {
			_, err = componentRegistry.insert(&componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range numberOfInserts {
			component, err := getComponentFromComponentRegistry[componentA](&componentRegistry, uint(i))
			assert.NoError(err)
			assert.NotNil(component)
			assert.Equal(i, component.value)
		}
	})
}
