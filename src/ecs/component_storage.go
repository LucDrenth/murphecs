package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lucdrenth/murphecs/src/utils"
)

// componentStorage stores instances of one type of component
type componentStorage struct {
	pointerToStart     unsafe.Pointer // points to the start of data
	data               reflect.Value  // buffer for storing components
	componentSize      uintptr        // the amount of memory that each component takes up in data
	componentId        ComponentId
	nextItemIndex      uint // the next inserted component will be inserted at this index
	capacity           uint // the number of components that can be stored with the current size of data
	numberOfComponents uint // the number of components that this storage contains
}

// createComponentStorage creates a new instance of createComponentStorage that can hold [capacity] components of type [ComponentId].
func createComponentStorage(capacity uint, componentId ComponentId) (componentStorage, error) {
	if capacity == 0 {
		return componentStorage{}, fmt.Errorf("%w: capacity may not be 0", ErrInvalidComponentStorageCapacity)
	}

	data := reflect.New(reflect.ArrayOf(int(capacity), componentId.componentType)).Elem()

	return componentStorage{
		pointerToStart: data.Addr().UnsafePointer(),
		data:           data,
		componentSize:  utils.AlignedSize(componentId.componentType),
		componentId:    componentId,
		nextItemIndex:  0,
		capacity:       capacity,
	}, nil
}

// increaseCapacity creates a new component buffer with more capacity
func (storage *componentStorage) increaseCapacity(extraCapacity uint) {
	newCapacity := storage.capacity + extraCapacity
	newData := reflect.New(reflect.ArrayOf(int(newCapacity), storage.componentId.componentType)).Elem()
	reflect.Copy(newData, storage.data)

	storage.data = newData
	storage.pointerToStart = newData.Addr().UnsafePointer()
	storage.capacity = newCapacity
}

// insert returns the index at which the component was inserted.
func (storage *componentStorage) insert(world *World, component IComponent) (uint, error) {
	insertIndex := storage.nextItemIndex

	if storage.capacity == insertIndex {
		extraCapacity := world.componentCapacityGrowthStrategy.GetExtraCapacity(storage.capacity)
		if extraCapacity == 0 {
			return 0, fmt.Errorf("component capacity growth is 0")
		}
		storage.increaseCapacity(extraCapacity)
	}

	err := storage.set(component, insertIndex)
	if err != nil {
		return 0, err
	}

	storage.nextItemIndex += 1
	storage.numberOfComponents += 1

	return insertIndex, nil
}

func (storage *componentStorage) set(component IComponent, index uint) error {
	if index >= storage.capacity {
		return fmt.Errorf("%w: %d", ErrComponentStorageIndexOutOfBounds, index)
	}

	componentValue := reflect.ValueOf(component)

	if componentValue.Kind() != reflect.Pointer {
		ptrValue := reflect.New(componentValue.Type())
		ptrValue.Elem().Set(componentValue)
		componentValue = ptrValue
	}

	componentPointer := componentValue.UnsafePointer()

	destination, err := storage.getComponentPointer(index)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		panic(err)
	}

	utils.CopyPointerData(componentPointer, destination, storage.componentSize)

	return nil
}

// insertRaw returns the index at which the component was inserted.
func (storage *componentStorage) insertRaw(world *World, componentPointer unsafe.Pointer) (uint, error) {
	insertIndex := storage.nextItemIndex

	if storage.capacity == insertIndex {
		extraCapacity := world.componentCapacityGrowthStrategy.GetExtraCapacity(storage.capacity)
		if extraCapacity == 0 {
			return 0, fmt.Errorf("component capacity growth is 0")
		}
		storage.increaseCapacity(extraCapacity)
	}

	destination, err := storage.getComponentPointer(insertIndex)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		panic(err)
	}

	utils.CopyPointerData(componentPointer, destination, storage.componentSize)
	storage.nextItemIndex += 1
	storage.numberOfComponents += 1

	return insertIndex, nil
}

// insertValue returns the index at which the component was inserted.
func (storage *componentStorage) insertValue(world *World, component *reflect.Value) (uint, error) {
	insertIndex := storage.nextItemIndex

	if storage.capacity == insertIndex {
		extraCapacity := world.componentCapacityGrowthStrategy.GetExtraCapacity(storage.capacity)
		if extraCapacity == 0 {
			return 0, fmt.Errorf("component capacity growth is 0")
		}
		storage.increaseCapacity(extraCapacity)
	}

	err := storage.setValue(component, insertIndex)
	if err != nil {
		return 0, err
	}

	storage.nextItemIndex += 1
	storage.numberOfComponents += 1

	return insertIndex, nil
}

func (storage *componentStorage) setValue(component *reflect.Value, index uint) error {
	if index >= storage.capacity {
		return fmt.Errorf("%w: %d", ErrComponentStorageIndexOutOfBounds, index)
	}

	componentPointer := component.UnsafePointer()

	destination, err := storage.getComponentPointer(index)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		panic(err)
	}

	utils.CopyPointerData(componentPointer, destination, storage.componentSize)

	return nil
}

// indicates that a component should be moved to another index to free up space
type movedComponent struct {
	fromIndex uint
	toIndex   uint
}

// remove moves the last item of the component storage to the index of the removed component and returns the indices
// of its old and new spot. This should be used to update the entity data.
func (storage *componentStorage) remove(index uint) (*movedComponent, error) {
	if index >= storage.nextItemIndex {
		return nil, fmt.Errorf("%w: %d", ErrComponentStorageIndexOutOfBounds, index)
	}

	storage.numberOfComponents -= 1

	if storage.nextItemIndex == 0 || index >= (storage.nextItemIndex-1) {
		return nil, nil
	}

	// move the last component in the storage to the index of the removed component to reuse the memory block.

	result := movedComponent{
		fromIndex: storage.nextItemIndex - 1,
		toIndex:   index,
	}

	err := storage.copyComponent(result.fromIndex, result.toIndex)
	if err != nil {
		return &result, fmt.Errorf("failed to move component: %w", err)
	}

	storage.nextItemIndex -= 1

	return &result, nil
}

// copyComponent copies a component from one index in the storage to another. Both indices must already be
// a valid component, you can not use an empty index.
func (storage *componentStorage) copyComponent(fromIndex uint, toIndex uint) error {
	if fromIndex == toIndex {
		return nil
	}

	if fromIndex >= storage.nextItemIndex {
		return fmt.Errorf("%w: fromIndex: %d", ErrComponentStorageIndexOutOfBounds, fromIndex)
	}

	if toIndex >= storage.nextItemIndex {
		return fmt.Errorf("%w: toIndex: %d", ErrComponentStorageIndexOutOfBounds, toIndex)
	}

	source, err := storage.getComponentPointer(fromIndex)
	if err != nil {
		panic(err)
	}

	destination, err := storage.getComponentPointer(toIndex)
	if err != nil {
		panic(err)
	}

	utils.CopyPointerData(source, destination, storage.componentSize)

	return nil
}

// getComponentPointer returns an unsafe.Pointer to the component at index.
//
// Returns an error if index is out of bounds.
func (storage *componentStorage) getComponentPointer(index uint) (unsafe.Pointer, error) {
	// This check might be redundant in most cases but lets be extra safe with unsafe.Pointer
	if index >= storage.capacity {
		return nil, fmt.Errorf("%w: index is %d, componentStorage size is %d", ErrComponentStorageIndexOutOfBounds, index, storage.nextItemIndex)
	}

	return unsafe.Add(storage.pointerToStart, uintptr(index)*storage.componentSize), nil
}

// getComponentFromComponentStorage returns a pointer to the component at index.
//
// Returns an error if index is out of bounds.
func getComponentFromComponentStorage[T IComponent](storage *componentStorage, index uint, isPointer bool) (result T, err error) {
	componentPointer, err := storage.getComponentPointer(index)
	if err != nil {
		return result, err
	}

	// Check if the generic type T is a pointer (e.g., *componentA)
	if isPointer {
		// Treat 'result' as a bucket for a pointer and write the address directly.
		*(*unsafe.Pointer)(unsafe.Pointer(&result)) = componentPointer
	} else {
		// Make 'result' a copy of the data at that address.
		result = *(*T)(componentPointer)
	}

	return result, nil
}
