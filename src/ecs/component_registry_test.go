package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComponentRegistry(t *testing.T) {
	type componentA struct{ Component }

	// since we are dealing with manually managed memory, lets make sure that we can successfully
	// create a component registry without crashing.
	t.Run("does not crash", func(t *testing.T) {
		createComponentRegistry(1, GetComponentType[componentA]())
		createComponentRegistry(1024, GetComponentType[componentA]())
	})
}

func TestComponentRegistryInsert(t *testing.T) {
	type componentA struct{ Component }

	t.Run("fails when component is not a pointer", func(t *testing.T) {
		assert := assert.New(t)

		componentRegistry := createComponentRegistry(4, GetComponentType[componentA]())
		_, err := componentRegistry.insert(componentA{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)
	})

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		componentRegistry := createComponentRegistry(4, GetComponentType[componentA]())
		_, err := componentRegistry.insert(&componentA{})
		assert.NoError(err)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		componentRegistry := createComponentRegistry(4, GetComponentType[componentA]())

		for range 10 {
			_, err := componentRegistry.insert(&componentA{})
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

		componentRegistry := createComponentRegistry(4, GetComponentType[componentA]())
		component, err := getComponentFromComponentRegistry[componentA](&componentRegistry, 5)
		assert.Error(err)
		assert.Nil(component)
	})

	t.Run("gets the correct components", func(t *testing.T) {
		assert := assert.New(t)

		capacity := 4

		componentRegistry := createComponentRegistry(uint(capacity), GetComponentType[componentA]())

		for i := range capacity {
			_, err := componentRegistry.insert(&componentA{
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

		capacity := 4
		numberOfInserts := 15

		componentRegistry := createComponentRegistry(uint(capacity), GetComponentType[componentA]())

		for i := range numberOfInserts {
			_, err := componentRegistry.insert(&componentA{
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
