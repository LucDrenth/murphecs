package ecs

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComponentRegistry(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when using capacity of 0", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		_, err := createComponentStorage(0, ComponentIdFor[componentA](&world))
		assert.ErrorIs(err, ErrInvalidComponentStorageCapacity)
	})

	t.Run("does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		_, err := createComponentStorage(1, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = createComponentStorage(1024, ComponentIdFor[componentA](&world))
		assert.NoError(err)
	})
}

func TestComponentRegistryInsert(t *testing.T) {
	type componentA struct{ Component }

	t.Run("fails when component is not a pointer", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentRegistry, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = componentRegistry.insert(componentA{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)
	})

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentRegistry, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		_, err = componentRegistry.insert(&componentA{})
		assert.NoError(err)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentRegistry, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for range 10 {
			_, err = componentRegistry.insert(&componentA{})
			assert.NoError(err)
		}
	})

	type componentWithPointers struct {
		Component
		name     *string
		aMap     map[string]int
		intSlice []*int
	}

	t.Run("works well with garbage collector", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentRegistry, err := createComponentStorage(4, ComponentIdFor[componentWithPointers](&world))
		assert.NoError(err)

		name := "a name"
		sliceLength := 100

		{
			item := componentWithPointers{
				name: &name,
				aMap: map[string]int{
					"1": 1,
					"3": 3,
				},
				intSlice: []*int{},
			}
			for i := range sliceLength {
				item.intSlice = append(item.intSlice, &i)
			}

			componentRegistry.insert(&item)
		}

		{
			// assert that component matches what we inserted
			item, err := getComponentFromComponentStorage[componentWithPointers](&componentRegistry, 0)
			assert.NoError(err)
			assert.Equal(name, *item.name)
			assert.Equal(
				map[string]int{
					"1": 1,
					"3": 3,
				},
				item.aMap,
			)
			assert.Len(item.intSlice, 100)
			for i := range sliceLength {
				assert.Equal(i, *item.intSlice[i])
			}
		}

		runtime.GC()

		{
			// assert that component matches what we inserted after garbage collection has run
			item, err := getComponentFromComponentStorage[componentWithPointers](&componentRegistry, 0)
			assert.NoError(err)
			assert.Equal(name, *item.name)
			assert.Equal(
				map[string]int{
					"1": 1,
					"3": 3,
				},
				item.aMap,
			)
			assert.Len(item.intSlice, 100)
			for i := range sliceLength {
				assert.Equal(i, *item.intSlice[i])
			}
		}
	})
}

func TestComponentRegistrySet(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when index is out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentRegistry, err := createComponentStorage(2, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		err = componentRegistry.set(&componentA{}, 3)
		assert.ErrorIs(err, ErrComponentRegistryIndexOutOfBounds)
	})

	t.Run("can set a component at the same index multiple times", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentRegistry, err := createComponentStorage(2, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for range 10 {
			err = componentRegistry.set(&componentA{}, 1)
			assert.NoError(err)
		}

		assert.Equal(uint(2), componentRegistry.capacity)
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

		componentRegistry, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		component, err := getComponentFromComponentStorage[componentA](&componentRegistry, 5)
		assert.Error(err)
		assert.Nil(component)
	})

	t.Run("gets the correct components when not exceeding capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := 4

		componentRegistry, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range capacity {
			_, err = componentRegistry.insert(&componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range capacity {
			component, err := getComponentFromComponentStorage[componentA](&componentRegistry, uint(i))
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

		componentRegistry, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range numberOfInserts {
			_, err = componentRegistry.insert(&componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range numberOfInserts {
			component, err := getComponentFromComponentStorage[componentA](&componentRegistry, uint(i))
			assert.NoError(err)
			assert.NotNil(component)
			assert.Equal(i, component.value)
		}
	})
}
