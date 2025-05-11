package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// componentStorage stores instances one type of component
type componentStorage struct {
	pointerToStart     unsafe.Pointer // points to the start of data
	data               reflect.Value  // buffer for storing components
	componentSize      uintptr        // the amount of memory that each component takes up in data
	componentId        ComponentId
	nextItemIndex      uint // the next inserted component will be inserted at this index
	capacity           uint // the number of components that can be stored with the current size of data
	numberOfComponents uint // the number of active components that this storage holds
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
//
// Returns an ErrComponentIsNotAPointer when component is not passed as a reference (e.g. componentA{}, instead of &componentA{})
func (storage *componentStorage) insert(component IComponent) (uint, error) {
	componentValue := reflect.ValueOf(component)

	if componentValue.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("%w: component %s must be a pointer", ErrComponentIsNotAPointer, ComponentDebugStringOf(component))
	}

	insertIndex := storage.nextItemIndex

	if storage.capacity == insertIndex {
		// double our current capacity
		storage.increaseCapacity(storage.capacity)
	}

	componentPointer := componentValue.UnsafePointer()

	destination, err := storage.getComponentPointer(insertIndex)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		return 0, fmt.Errorf("unexpected error when getting component pointer: %w", err)
	}

	utils.CopyPointerData(componentPointer, destination, storage.componentSize)
	storage.nextItemIndex += 1
	storage.numberOfComponents += 1

	return insertIndex, nil
}

// insertRaw returns the index at which the component was inserted.
//
// Returns an ErrComponentIsNotAPointer when component is not passed as a reference (e.g. componentA{}, instead of &componentA{})
func (storage *componentStorage) insertRaw(componentPointer unsafe.Pointer) (uint, error) {
	insertIndex := storage.nextItemIndex

	if storage.capacity == insertIndex {
		// double our current capacity
		storage.increaseCapacity(storage.capacity)
	}

	destination, err := storage.getComponentPointer(insertIndex)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		return 0, fmt.Errorf("unexpected error when getting component pointer: %w", err)
	}

	utils.CopyPointerData(componentPointer, destination, storage.componentSize)
	storage.nextItemIndex += 1
	storage.numberOfComponents += 1

	return insertIndex, nil
}

func (storage *componentStorage) remove(index uint) error {
	if index >= storage.capacity {
		return fmt.Errorf("index %d is out of bounds for componentStorage of size %d", index, storage.nextItemIndex)
	}

	storage.numberOfComponents -= 1

	// TODO on insert, reuse the component at [index] that is now free.

	return nil
}

// getComponentPointer returns an unsafe.Pointer to the component at index.
//
// Returns an error if index is out of bounds.
func (storage *componentStorage) getComponentPointer(index uint) (unsafe.Pointer, error) {
	if index >= storage.capacity {
		return nil, fmt.Errorf("index %d is out of bounds for componentStorage of size %d", index, storage.nextItemIndex)
	}

	return unsafe.Add(storage.pointerToStart, uintptr(index)*storage.componentSize), nil
}

// getComponentFromComponentStorage returns a pointer to the component at index.
//
// Returns an error if index is out of bounds.
func getComponentFromComponentStorage[T IComponent](storage *componentStorage, index uint) (*T, error) {
	componentPointer, err := storage.getComponentPointer(index)
	if err != nil {
		return nil, err
	}

	return (*T)(componentPointer), nil
}
