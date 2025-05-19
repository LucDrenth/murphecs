package ecs

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComponentStorage(t *testing.T) {
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

func TestComponentStorageInsert(t *testing.T) {
	type componentA struct{ Component }

	t.Run("fails when component is not a pointer", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		_, err = componentStorage.insert(&world, componentA{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)
	})

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(capacity, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		_, err = componentStorage.insert(&world, &componentA{})
		assert.NoError(err)

		assert.Equal(capacity, componentStorage.capacity)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for range 10 {
			_, err = componentStorage.insert(&world, &componentA{})
			assert.NoError(err)
		}

		assert.True(componentStorage.capacity > capacity)
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
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentWithPointers](&world))
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

			componentStorage.insert(&world, &item)
		}

		{
			// assert that component matches what we inserted
			item, err := getComponentFromComponentStorage[componentWithPointers](&componentStorage, 0)
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
			item, err := getComponentFromComponentStorage[componentWithPointers](&componentStorage, 0)
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

func TestComponentStorageSet(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when index is out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentStorage, err := createComponentStorage(2, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		err = componentStorage.set(&componentA{}, 3)
		assert.ErrorIs(err, ErrComponentStorageIndexOutOfBounds)
	})

	t.Run("can set a component at the same index multiple times", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		componentStorage, err := createComponentStorage(2, ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for range 10 {
			err = componentStorage.set(&componentA{}, 1)
			assert.NoError(err)
		}

		assert.Equal(uint(2), componentStorage.capacity)
	})
}

func TestGetComponentFromComponentStorage(t *testing.T) {
	type componentA struct {
		value int
		Component
	}

	t.Run("returns an error if index out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()

		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](&world))
		assert.NoError(err)
		component, err := getComponentFromComponentStorage[componentA](&componentStorage, 5)
		assert.Error(err)
		assert.Nil(component)
	})

	t.Run("gets the correct components when not exceeding capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := DefaultWorld()
		capacity := 4

		componentStorage, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range capacity {
			_, err = componentStorage.insert(&world, &componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range capacity {
			component, err := getComponentFromComponentStorage[componentA](&componentStorage, uint(i))
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

		componentStorage, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](&world))
		assert.NoError(err)

		for i := range numberOfInserts {
			_, err = componentStorage.insert(&world, &componentA{
				value: i,
			})
			assert.NoError(err)
		}

		for i := range numberOfInserts {
			component, err := getComponentFromComponentStorage[componentA](&componentStorage, uint(i))
			assert.NoError(err)
			assert.NotNil(component)
			assert.Equal(i, component.value)
		}
	})
}
