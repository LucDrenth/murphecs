package ecs

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/lucdrenth/murphecs/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateComponentStorage(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when using capacity of 0", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		_, err := createComponentStorage(0, ComponentIdFor[componentA](world))
		assert.ErrorIs(err, ErrInvalidComponentStorageCapacity)
	})

	t.Run("does not return an error", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		_, err := createComponentStorage(1, ComponentIdFor[componentA](world))
		assert.NoError(err)
		_, err = createComponentStorage(1024, ComponentIdFor[componentA](world))
		assert.NoError(err)
	})
}

func TestGetComponentStoragePointer(t *testing.T) {
	assert := assert.New(t)

	world := NewDefaultWorld()
	componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
	assert.NoError(err)

	_, err = componentStorage.getComponentPointer(3)
	assert.NoError(err)
	_, err = componentStorage.getComponentPointer(4)
	assert.Error(err)
	_, err = componentStorage.getComponentPointer(5)
	assert.Error(err)
}

func TestComponentStorageInsert(t *testing.T) {
	type componentA struct{ Component }

	t.Run("fails when component is not a pointer", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()

		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		_, err = componentStorage.insert(world, componentA{})
		assert.ErrorIs(err, ErrComponentIsNotAPointer)
	})

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(capacity, ComponentIdFor[componentA](world))
		assert.NoError(err)

		_, err = componentStorage.insert(world, &componentA{})
		assert.NoError(err)

		assert.Equal(capacity, componentStorage.capacity)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)

		for range 10 {
			_, err = componentStorage.insert(world, &componentA{})
			assert.NoError(err)
		}

		assert.True(componentStorage.capacity > capacity)
	})

	t.Run("works well with garbage collector", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[ComponentWithPointers](world))
		assert.NoError(err)

		{
			item := CreateComponentWithPointers()
			_, err := componentStorage.insert(world, item)
			assert.NoError(err)
		}

		{
			// assert that component matches what we inserted
			item, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0)
			assert.NoError(err)
			err = item.Validate()
			assert.NoError(err)
		}

		runtime.GC()

		{
			// assert that component matches what we inserted after garbage collection has run
			item, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0)
			assert.NoError(err)
			err = item.Validate()
			assert.NoError(err)
		}
	})
}

func TestComponentStorageInsertValue(t *testing.T) {
	type componentA struct{ Component }

	t.Run("successfully inserts when there is enough capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(capacity, ComponentIdFor[componentA](world))
		assert.NoError(err)

		componentValue := reflect.ValueOf(&componentA{})
		_, err = componentStorage.insertValue(world, &componentValue)
		assert.NoError(err)

		assert.Equal(capacity, componentStorage.capacity)
	})

	t.Run("increases capacity and inserts the component when overstepping capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := uint(4)
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)

		for range 10 {
			componentValue := reflect.ValueOf(&componentA{})
			_, err = componentStorage.insertValue(world, &componentValue)
			assert.NoError(err)
		}

		assert.True(componentStorage.capacity > capacity)
	})

	t.Run("works well with garbage collector", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[ComponentWithPointers](world))
		assert.NoError(err)

		{
			item := reflect.ValueOf(CreateComponentWithPointers())
			_, err := componentStorage.insertValue(world, &item)
			assert.NoError(err)
		}

		{
			// assert that component matches what we inserted
			item, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0)
			assert.NoError(err)
			err = item.Validate()
			assert.NoError(err)
		}

		runtime.GC()

		{
			// assert that component matches what we inserted after garbage collection has run
			item, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0)
			assert.NoError(err)
			err = item.Validate()
			assert.NoError(err)
		}
	})
}

func TestComponentStorageSet(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns an error when index is out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(2, ComponentIdFor[componentA](world))
		assert.NoError(err)

		err = componentStorage.set(&componentA{}, 3)
		assert.ErrorIs(err, ErrComponentStorageIndexOutOfBounds)
	})

	t.Run("can set a component at the same index multiple times", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(2, ComponentIdFor[componentA](world))
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

		world := NewDefaultWorld()

		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		component, err := getComponentFromComponentStorage[componentA](&componentStorage, 5)
		assert.Error(err)
		assert.Nil(component)
	})

	t.Run("gets the correct components when not exceeding capacity", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := 4

		componentStorage, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](world))
		assert.NoError(err)

		for i := range capacity {
			_, err = componentStorage.insert(world, &componentA{
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

		world := NewDefaultWorld()
		capacity := 4
		numberOfInserts := 15

		componentStorage, err := createComponentStorage(uint(capacity), ComponentIdFor[componentA](world))
		assert.NoError(err)

		for i := range numberOfInserts {
			_, err = componentStorage.insert(world, &componentA{
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

func TestRemoveFromComponentStorage(t *testing.T) {
	type componentA struct{ Component }

	t.Run("returns error if removing component at invalid index", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)

		_, err = componentStorage.remove(0)
		assert.ErrorIs(err, ErrComponentStorageIndexOutOfBounds)
	})

	t.Run("successfully removes a component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		index, err := componentStorage.insert(world, &componentA{})
		assert.NoError(err)

		_, err = componentStorage.remove(index)
		assert.NoError(err)
	})

	t.Run("does not move any component if the last inserted component is removed", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(1024, ComponentIdFor[componentA](world))
		assert.NoError(err)

		nrComponents := uint(10)
		for range nrComponents {
			_, err := componentStorage.insert(world, &componentA{})
			assert.NoError(err)
		}

		movedComponent, err := componentStorage.remove(nrComponents - 1)
		assert.NoError(err)
		assert.Nil(movedComponent)
	})

	t.Run("moves the last component to the place of the removed component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(1024, ComponentIdFor[componentA](world))
		assert.NoError(err)

		nrComponents := uint(10)
		for range nrComponents {
			_, err := componentStorage.insert(world, &componentA{})
			assert.NoError(err)
		}

		movedComponent, err := componentStorage.remove(5)
		assert.NoError(err)
		assert.NotNil(movedComponent)
		assert.Equal(nrComponents-1, movedComponent.fromIndex)
		assert.Equal(uint(5), movedComponent.toIndex)
	})

	t.Run("remove makes the component storage reuse memory", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		capacity := uint(8)
		componentStorage, err := createComponentStorage(capacity, ComponentIdFor[componentA](world))
		assert.NoError(err)

		for range capacity {
			_, err := componentStorage.insert(world, &emptyComponentA{})
			assert.NoError(err)
		}

		_, err = componentStorage.remove(5)
		assert.NoError(err)
		_, err = componentStorage.insert(world, &emptyComponentA{})
		assert.NoError(err)

		assert.Equal(capacity, componentStorage.capacity)
	})
}

func TestComponentStorageCopyComponent(t *testing.T) {
	type componentA struct {
		Component
		value int
	}

	t.Run("returns an error if 'from' is out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		_, err = componentStorage.insert(world, &componentA{})
		assert.NoError(err)

		err = componentStorage.copyComponent(1, 0)
		assert.ErrorIs(err, ErrComponentStorageIndexOutOfBounds)
	})

	t.Run("returns an error if 'to' is out of bounds", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		_, err = componentStorage.insert(world, &componentA{})
		assert.NoError(err)

		err = componentStorage.copyComponent(0, 1)
		assert.ErrorIs(err, ErrComponentStorageIndexOutOfBounds)
	})

	t.Run("successfully copies component", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[componentA](world))
		assert.NoError(err)
		_, err = componentStorage.insert(world, &componentA{value: 10})
		assert.NoError(err)
		_, err = componentStorage.insert(world, &componentA{value: 20})
		assert.NoError(err)

		err = componentStorage.copyComponent(1, 0)
		assert.NoError(err)
		component, err := getComponentFromComponentStorage[componentA](&componentStorage, 0) // copied-over component
		assert.NoError(err)
		assert.Equal(20, component.value)
		component, err = getComponentFromComponentStorage[componentA](&componentStorage, 1) // original component
		assert.NoError(err)
		assert.Equal(20, component.value)

		// changing the original should not change the copy
		component.value = 30
		component, err = getComponentFromComponentStorage[componentA](&componentStorage, 0) // copied-over component
		assert.NoError(err)
		assert.Equal(20, component.value)
	})

	t.Run("plays well with the garbage collector", func(t *testing.T) {
		assert := assert.New(t)

		world := NewDefaultWorld()
		componentStorage, err := createComponentStorage(4, ComponentIdFor[ComponentWithPointers](world))
		assert.NoError(err)
		_, err = componentStorage.insert(world, &componentA{})
		assert.NoError(err)
		_, err = componentStorage.insert(world, CreateComponentWithPointers())
		assert.NoError(err)

		err = componentStorage.copyComponent(1, 0)
		assert.NoError(err)

		{
			component, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0) // copied-over component
			assert.NoError(err)
			assert.NoError(component.Validate())
			component, err = getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 1) // original component
			assert.NoError(err)
			assert.NoError(component.Validate())
		}

		runtime.GC()

		{
			component, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0) // copied-over component
			assert.NoError(err)
			assert.NoError(component.Validate())
			component, err = getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 1) // original component
			assert.NoError(err)
			assert.NoError(component.Validate())

			// change component value
			component.name = utils.PointerTo(componentWithPointersName + " and something else")
		}

		runtime.GC()

		{
			component, err := getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 0) // copied-over component
			assert.NoError(err)
			assert.NoError(component.Validate()) // should not have changed

			component, err = getComponentFromComponentStorage[ComponentWithPointers](&componentStorage, 1) // original component
			assert.NoError(err)
			assert.Error(component.Validate()) // now has an error because name changed
		}
	})
}
