package utils

import "reflect"

// TypeOf returns the type of a generic.
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// AlignedSize returns the aligned size of t.
//
// Some systems already return the aligned size when calling reflect.Type.Size, but some don't.
// Use this function to ensure getting the aligned size. This is necessary when manually managing
// blocks of memory.
func AlignedSize(t reflect.Type) uintptr {
	size := t.Size()
	align := uintptr(t.Align())
	return (size + (align - 1)) / align * align
}
