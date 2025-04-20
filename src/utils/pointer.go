package utils

import (
	"math"
	"unsafe"
)

// PointerTo returns a pointer to value. Useful for pointing to a literal value.
//
// For example: `stringPointer := PointerTo("my string")`, which can not be done
// like this: `stringPointer := &"my string"`
func PointerTo[T any](value T) *T {
	return &value
}

// CopyPointerData copies what source points to, to where destination points to.
func CopyPointerData(source unsafe.Pointer, destination unsafe.Pointer, itemSize uintptr) {
	// We're casting each unsafe.Pointer to a pointer to a very large byte slice. This doesn't mean we're
	// using that much memory — it's just a trick to give us a big enough buffer to copy over the data.
	destinationSlice := (*[math.MaxInt32]byte)(destination)[:itemSize:itemSize]
	dataSlice := (*[math.MaxInt32]byte)(source)[:itemSize:itemSize]
	copy(destinationSlice, dataSlice)
}
