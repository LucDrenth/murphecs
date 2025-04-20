package ecs

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/lucdrenth/murph_engine/src/utils"
)

// componentRegistry stores components of 1 specific type.
type componentRegistry struct {
	pointerToStart unsafe.Pointer // points to the start of data
	data           reflect.Value  // buffer for storing components
	componentSize  uintptr        // the amount of memory that each component takes up in data
	componentType  reflect.Type
	nextItemIndex  uint // the next inserted component will be inserted at this index
	capacity       uint // the number of components that can be stored with the current size of data
}

// newComponentRegistry creates a new instance of newComponentRegistry that can hold [capacity] components of type [componentType].
func newComponentRegistry(capacity uint, componentType reflect.Type) componentRegistry {
	data := reflect.New(reflect.ArrayOf(int(capacity), componentType)).Elem()

	return componentRegistry{
		pointerToStart: data.Addr().UnsafePointer(),
		data:           data,
		componentSize:  utils.AlignedSize(componentType),
		componentType:  componentType,
		nextItemIndex:  0,
		capacity:       capacity,
	}
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

func (c *componentRegistry) insert(component IComponent) error {
	componentValue := reflect.ValueOf(component)

	if componentValue.Kind() != reflect.Ptr {
		return fmt.Errorf("%w: component %s must be a pointer", ErrComponentIsNotAPointer, toComponentDebugType(component))
	}

	componentPointer := componentValue.UnsafePointer()

	if c.capacity == c.nextItemIndex {
		// double our current capacity
		c.increaseCapacity(c.capacity)
	}

	destination, err := c.getComponentPointer(c.nextItemIndex)
	if err != nil {
		return err
	}

	utils.CopyPointerData(componentPointer, destination, c.componentSize)
	c.nextItemIndex += 1

	return nil
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
