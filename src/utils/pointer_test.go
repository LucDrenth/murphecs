package utils

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestCopyPointerData(t *testing.T) {
	assert := assert.New(t)

	type itemType uint8 // can be any type, result stays the same

	// setup
	var item itemType = 3
	typeOfItem := reflect.TypeOf(item)
	itemSize := AlignedSize(typeOfItem)
	targetItemIndex := 6
	totalItems := 10
	data := reflect.New(reflect.ArrayOf(int(totalItems), typeOfItem)).Elem()
	dataPointer := data.Addr().UnsafePointer() // points to the start of data

	getItem := func(index int) *itemType {
		return (*itemType)(unsafe.Add(dataPointer, uintptr(index)*itemSize))
	}

	// check that the expected item is not there yet
	for i := range totalItems {
		assert.Equal(itemType(0), *getItem(i))
	}

	// copy the item to the right place
	destination := unsafe.Add(
		dataPointer,
		uintptr(targetItemIndex)*itemSize,
	)

	source := unsafe.Pointer(&item)
	CopyPointerData(source, destination, itemSize)

	// check that only the expected item is now set
	for i := range totalItems {
		if i == targetItemIndex {
			assert.Equal(item, *getItem(i))
		} else {
			assert.Equal(itemType(0), *getItem(i))
		}
	}
}

func TestClonePointerValue(t *testing.T) {
	type dataStruct struct{ value int }

	t.Run("panics when passing nil", func(t *testing.T) {
		assert := assert.New(t)

		defer func() {
			r := recover()
			assert.NotNil(r)
		}()

		ClonePointerValue[dataStruct](nil)
	})

	t.Run("updating the cloned value does not alert the original value", func(t *testing.T) {
		assert := assert.New(t)

		expectedValue := 10
		data := &dataStruct{value: expectedValue}
		dataCopy := ClonePointerValue(data)
		dataCopy.value *= 2

		assert.Equal(expectedValue, data.value)
		assert.NotEqual(expectedValue, dataCopy.value)
	})
}
