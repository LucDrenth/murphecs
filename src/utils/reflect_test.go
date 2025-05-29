package utils

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

type someInterface interface {
	SomeMethod()
	unexportedMethod()
}

type someStructA struct{}

func (s someStructA) SomeMethod()       {}
func (s someStructA) unexportedMethod() {}

type someStructB struct{}

func (s *someStructB) SomeMethod()       {}
func (s *someStructB) unexportedMethod() {}

func TestMethodHasPointerReceiver(t *testing.T) {
	t.Run("returns error if interface is nil", func(t *testing.T) {
		assert := assert.New(t)

		_, err := MethodHasPointerReceiver(nil, "")
		assert.Error(err)
	})

	t.Run("returns error if input is not a struct", func(t *testing.T) {
		assert := assert.New(t)

		_, err := MethodHasPointerReceiver(10, "")
		assert.Error(err)
	})

	t.Run("returns error if method does not exist", func(t *testing.T) {
		assert := assert.New(t)

		_, err := MethodHasPointerReceiver(someStructA{}, "nonExisting")
		assert.Error(err)
	})

	t.Run("returns error if trying to check for unexported method", func(t *testing.T) {
		assert := assert.New(t)

		_, err := MethodHasPointerReceiver(someStructA{}, "unexportedMethod")
		assert.Error(err)
	})

	t.Run("returns false if method does not have pointer receiver", func(t *testing.T) {
		assert := assert.New(t)

		result, err := MethodHasPointerReceiver(someStructA{}, "SomeMethod")
		assert.NoError(err)
		assert.False(result)

		var input someInterface = someStructA{}
		result, err = MethodHasPointerReceiver(input, "SomeMethod")
		assert.NoError(err)
		assert.False(result)
	})

	t.Run("returns true if method has pointer receiver", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		result, err := MethodHasPointerReceiver(someStructB{}, "SomeMethod")
		require.NoError(err)
		assert.True(result)

		var input someInterface = &someStructB{}
		result, err = MethodHasPointerReceiver(input, "SomeMethod")
		assert.NoError(err)
		assert.True(result)
	})
}
