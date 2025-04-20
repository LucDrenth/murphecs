package utils

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	type someStruct struct{}
	type anotherStruct struct{}

	t.Run("normal type reflect equals result of TypeOf with same type", func(t *testing.T) {
		assert := assert.New(t)

		assert.True(
			reflect.TypeOf(someStruct{}) == TypeOf[someStruct](),
		)
		assert.False(
			reflect.TypeOf(anotherStruct{}) == TypeOf[someStruct](),
		)
	})
}

func TestAlignedSize(t *testing.T) {
	type misaligned struct {
		A int8  // 1 byte
		B int64 // 8 bytes (needs 8-byte alignment)
	}

	// The specific number is not relevant for the test. We could just as well use a different number
	// to demonstrate the same point.
	numberOfElements := 3

	fillElements := func(elementSize uintptr, memory []byte, numberOfElements int) {
		for i := range numberOfElements {
			ptr := unsafe.Pointer(&memory[i*int(elementSize)])
			instance := (*misaligned)(ptr)
			instance.A = int8(i)
			instance.B = int64(i)
		}
	}

	getElement := func(elementSize uintptr, memory []byte, index int) *misaligned {
		ptr := unsafe.Pointer(&memory[index*int(elementSize)])
		return (*misaligned)(ptr)
	}

	t.Run("returns unexpected results if we use unaligned struct size", func(t *testing.T) {
		assert := assert.New(t)

		// If we'd use `reflect.TypeOf(misaligned{}).size()`, some systems will return 9.
		// To keep the tests consistent on all platforms, we'll manually assign 9 here.
		size := uintptr(9)
		mem := make([]byte, size*uintptr(numberOfElements))

		fillElements(size, mem, numberOfElements)

		isAllCorrect := true

		// Some (but not all) of the elements will have an unexpected value
		for i := range numberOfElements {
			element := getElement(size, mem, i)
			if int64(i) != element.B {
				isAllCorrect = false
			}
		}

		assert.False(isAllCorrect)
	})

	t.Run("succeeds if we use aligned size", func(t *testing.T) {
		assert := assert.New(t)

		size := AlignedSize(reflect.TypeOf(misaligned{}))
		assert.Equal(uintptr(16), size)
		mem := make([]byte, size*uintptr(numberOfElements))

		fillElements(size, mem, numberOfElements)

		for i := range numberOfElements {
			element := getElement(size, mem, i)
			assert.Equal(int64(i), element.B)
		}
	})
}
