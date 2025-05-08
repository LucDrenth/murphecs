package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// componentRegistry stores instances one type of component
type componentRegistry struct {
	pointerToStart unsafe.Pointer // points to the start of data
	data           reflect.Value  // buffer for storing components
	componentSize  uintptr        // the amount of memory that each component takes up in data
	componentType  reflect.Type
	nextItemIndex  uint // the next inserted component will be inserted at this index
	capacity       uint // the number of components that can be stored with the current size of data
}

// createComponentRegistry creates a new instance of createComponentRegistry that can hold [capacity] components of type [componentType].
func createComponentRegistry(capacity uint, componentType reflect.Type) (componentRegistry, error) {
	if capacity == 0 {
		return componentRegistry{}, fmt.Errorf("%w: capacity may not be 0", ErrInvalidComponentStorageCapacity)
	}

	data := reflect.New(reflect.ArrayOf(int(capacity), componentType)).Elem()

	return componentRegistry{
		pointerToStart: data.Addr().UnsafePointer(),
		data:           data,
		componentSize:  utils.AlignedSize(componentType),
		componentType:  componentType,
		nextItemIndex:  0,
		capacity:       capacity,
	}, nil
}

// increaseCapacity creates a new component buffer with more capacity
func (c *componentRegistry) increaseCapacity(extraCapacity uint) {
	newCapacity := c.capacity + extraCapacity
	newData := reflect.New(reflect.ArrayOf(int(newCapacity), c.componentType)).Elem()
	reflect.Copy(newData, c.data)

	c.data = newData
	c.pointerToStart = newData.Addr().UnsafePointer()
	c.capacity = newCapacity
}

// insert returns the index at which the component was inserted.
//
// Returns an ErrComponentIsNotAPointer when component is not passed as a reference (e.g. componentA{}, instead of &componentA{})
func (c *componentRegistry) insert(component IComponent) (uint, error) {
	componentValue := reflect.ValueOf(component)

	if componentValue.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("%w: component %s must be a pointer", ErrComponentIsNotAPointer, toComponentDebugType(component))
	}

	componentPointer := componentValue.UnsafePointer()
	insertIndex := c.nextItemIndex

	if c.capacity == insertIndex {
		// double our current capacity
		c.increaseCapacity(c.capacity)
	}

	destination, err := c.getComponentPointer(insertIndex)
	if err != nil {
		// this error should never happen because we always make sure we have enough capacity.
		return 0, fmt.Errorf("unexpected error when getting component pointer: %w", err)
	}

	utils.CopyPointerData(componentPointer, destination, c.componentSize)
	c.nextItemIndex += 1

	return insertIndex, nil
}

// getComponentPointer returns an unsafe.Pointer to the component at index.
//
// Returns an error if index is out of bounds.
func (c *componentRegistry) getComponentPointer(index uint) (unsafe.Pointer, error) {
	if index >= c.capacity {
		return nil, fmt.Errorf("index %d is out of bounds for componentRegistry of size %d", index, c.nextItemIndex)
	}

	return unsafe.Add(c.pointerToStart, uintptr(index)*c.componentSize), nil
}

// getComponentFromComponentRegistry returns a pointer to the component at index.
//
// Returns an error if index is out of bounds.
func getComponentFromComponentRegistry[T IComponent](c *componentRegistry, index uint) (*T, error) {
	componentPointer, err := c.getComponentPointer(index)
	if err != nil {
		return nil, err
	}

	return (*T)(componentPointer), nil
}
